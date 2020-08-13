/*
 * @Author: guiguan
 * @Date:   2020-08-13T12:26:53+10:00
 * @Last modified by:   guiguan
 * @Last modified time: 2020-08-13T18:00:22+10:00
 */

//go:generate protoc --plugin=protoc-gen-doc=proto/protoc-gen-doc --doc_out=proto --doc_opt=markdown,docs.md -I proto --go_out=plugins=grpc:proto proto/hyperledger.proto

package hyperledger

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/SouthbankSoftware/provendb-hyperledger/chaincode/common"
	pb "github.com/SouthbankSoftware/provendb-hyperledger/pkg/hyperledger/proto"
	"github.com/SouthbankSoftware/provendb-tree/pkg/log"
	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const (
	appName                   = "hyperledger"
	org                       = "Org1MSP"
	orgUser                   = "Admin"
	channelName               = "provendb"
	chancodeName              = channelName
	contractFuncNameEmbedData = "embedData"

	walletAppUser  = "appUser"
	walletCredPath = "test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp"

	sdkConnectionProfile = "test-network/connection-profile/provendb-hyperledger.json"
)

var (
	svc     *service
	svcOnce = new(sync.Once)
)

// ServiceConfig represents the configuration of a service
type ServiceConfig struct {
	HostPort string
}

// Service represents a service
type Service interface {
	// Run runs the service
	Run() error

	pb.HyperledgerServiceServer
}

// Service represents a service instance
type service struct {
	*ServiceConfig

	client   *ledger.Client
	contract *gateway.Contract
}

// NewService creates a new singleton service instance
func NewService(config *ServiceConfig) Service {
	svcOnce.Do(func() {
		svc = &service{
			ServiceConfig: config,
		}
	})

	return svc
}

// Run runs the service
func (s *service) Run() error {
	if s.client != nil {
		return fmt.Errorf("the %s service is already running", appName)
	}

	log.SF().Bg().Info("start", zap.String("hostPort", s.HostPort))

	wallet, err := gateway.NewFileSystemWallet("data/wallet")
	if err != nil {
		return fmt.Errorf("failed to create wallet: %w", err)
	}

	if !wallet.Exists(walletAppUser) {
		err = populateWallet(wallet)
		if err != nil {
			return fmt.Errorf("failed to populate wallet contents: %w", err)
		}
	}

	connProfile := config.FromFile("test-network/connection-profile/provendb-hyperledger.json")

	sdk, err := fabsdk.New(connProfile)
	if err != nil {
		return err
	}
	defer sdk.Close()

	chanProvider := sdk.ChannelContext(channelName, fabsdk.WithUser("Admin"), fabsdk.WithOrg("Org1MSP"))
	client, err := ledger.New(chanProvider)
	if err != nil {
		return err
	}

	s.client = client

	gw, err := gateway.Connect(
		gateway.WithConfig(connProfile),
		gateway.WithIdentity(wallet, walletAppUser),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to gateway: %w", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork(channelName)
	if err != nil {
		return fmt.Errorf("failed to get network: %w", err)
	}

	s.contract = network.GetContract(chancodeName)

	eg, egCTX := errgroup.WithContext(context.Background())

	grpcSRV := grpc.NewServer(
		grpc.UnaryInterceptor(
			logUnaryServerInterceptor(),
		),
		grpc.StreamInterceptor(
			logStreamServerInterceptor(),
		),
	)
	pb.RegisterHyperledgerServiceServer(grpcSRV, s)
	reflection.Register(grpcSRV)

	eg.Go(func() error {
		// start gRPC server
		lis, err := net.Listen("tcp", s.HostPort)
		if err != nil {
			return err
		}

		return grpcSRV.Serve(lis)
	})

	eg.Go(func() error {
		// graceful shutdown
		shutdownCh := make(chan os.Signal, 1)

		signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

		select {
		case sn := <-shutdownCh:
			log.SF().Bg().Info("stop",
				zap.String("signal", sn.String()),
			)
		case <-egCTX.Done():
			log.SF().Bg().Info("stop")
		}

		grpcSRV.Stop()

		return errors.New("stopped")
	})

	err = eg.Wait()
	if err == http.ErrServerClosed {
		return nil
	}
	return nil
}

func populateWallet(wallet *gateway.Wallet) error {
	certPath := filepath.Join(walletCredPath, "signcerts", "cert.pem")
	cert, err := ioutil.ReadFile(certPath)
	if err != nil {
		return err
	}

	keyDir := filepath.Join(walletCredPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return errors.New("keystore folder should have contain one file")
	}

	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity(org, string(cert), string(key))
	err = wallet.Put(walletAppUser, identity)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) EmbedData(ctx context.Context, in *pb.EmbedDataRequest) (
	rp *pb.EmbedDataReply, er error) {
	data := in.GetData()
	if data == "" {
		er = status.Error(codes.InvalidArgument, "`data` must be provided")
		return
	}

	rawReply, err := s.contract.SubmitTransaction(contractFuncNameEmbedData, in.GetData())
	if err != nil {
		er = err
		return
	}

	reply := &common.EmbedDataReply{}
	err = json.Unmarshal(rawReply, reply)
	if err != nil {
		er = err
		return
	}

	bn, err := s.getBlockNumberByTxnID(reply.TxnID)
	if err != nil {
		er = err
		return
	}

	ct, err := ptypes.TimestampProto(reply.CreateTime)
	if err != nil {
		er = err
		return
	}

	rp = &pb.EmbedDataReply{
		TxnId:       reply.TxnID,
		CreateTime:  ct,
		BlockNumber: bn,
	}
	return
}

func (s *service) GetTransactionByID(ctx context.Context, in *pb.GetTransactionByIDRequest) (
	tr *pb.Transaction, er error) {
	txnID := in.GetTxnId()
	if txnID == "" {
		er = status.Error(codes.InvalidArgument, "`txn_id` must be provided")
		return
	}

	ti, err := s.getTxnInfoByTxnID(txnID)
	if err != nil {
		if strings.Contains(err.Error(), "Entry not found in index") {
			er = status.Error(codes.NotFound, "transaction doesn't exist")
			return
		}

		er = status.Error(codes.Internal,
			fmt.Sprintf("transaction doesn't look like a ProvenDB data embedding transaction: %s", err))
		return
	}

	var data string

	if args := ti.Args; len(args) != 2 || args[0] != contractFuncNameEmbedData {
		er = status.Error(codes.Internal,
			fmt.Sprintf("transaction doesn't look like a ProvenDB data embedding transaction, with args: %s", args))
		return
	} else {
		data = args[1]
	}

	bn, err := s.getBlockNumberByTxnID(txnID)
	if err != nil {
		er = err
		return
	}

	tr = &pb.Transaction{
		TxnId:       txnID,
		CreateTime:  ti.CreateTime,
		BlockNumber: bn,
		Data:        data,
	}
	return
}
