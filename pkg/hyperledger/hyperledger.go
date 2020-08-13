/*
 * @Author: guiguan
 * @Date:   2020-08-13T16:50:20+10:00
 * @Last modified by:   guiguan
 * @Last modified time: 2020-08-13T17:37:34+10:00
 */

package hyperledger

import (
	"fmt"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric/protoutil"
)

type txnInfo struct {
	TxnID      string
	CreateTime *timestamp.Timestamp
	Args       []string
}

func (s *service) getBlockNumberByTxnID(txnID string) (bn uint64, er error) {
	block, err := s.client.QueryBlockByTxID(fab.TransactionID(txnID))
	if err != nil {
		er = err
		return
	}

	bn = block.GetHeader().GetNumber()
	return
}

func (s *service) getTxnInfoByTxnID(txnID string) (ti *txnInfo, er error) {
	ptx, err := s.client.QueryTransaction(fab.TransactionID(txnID))
	if err != nil {
		er = err
		return
	}

	payload, err := protoutil.UnmarshalPayload(ptx.GetTransactionEnvelope().GetPayload())
	if err != nil {
		er = fmt.Errorf("failed to extract `Payload` from `Envelope`: %w", err)
		return
	}

	chHeader, err := protoutil.UnmarshalChannelHeader(payload.GetHeader().GetChannelHeader())
	if err != nil {
		er = fmt.Errorf("failed to extract `ChannelHeader` from `Payload`: %w", err)
		return
	}

	txn, err := protoutil.UnmarshalTransaction(payload.GetData())
	if err != nil {
		er = fmt.Errorf("failed to extract `Transaction` from `Payload`: %w", err)
		return
	}

	if len(txn.Actions) == 0 {
		er = fmt.Errorf("expected non-empty transaction actions")
		return
	}

	caPayload, err := protoutil.UnmarshalChaincodeActionPayload(txn.GetActions()[0].GetPayload())
	if err != nil {
		er = fmt.Errorf("failed to extract `ChaincodeActionPayload` from `Transaction`: %w", err)
		return
	}

	cpPayload, err := protoutil.UnmarshalChaincodeProposalPayload(caPayload.GetChaincodeProposalPayload())
	if err != nil {
		er = fmt.Errorf("failed to extract `ChaincodeProposalPayload` from `ChaincodeActionPayload`: %w", err)
		return
	}

	inSpec, err := protoutil.UnmarshalChaincodeInvocationSpec(cpPayload.GetInput())
	if err != nil {
		er = fmt.Errorf("failed to extract `ChaincodeInvocationSpec` from `ChaincodeProposalPayload`: %w", err)
		return
	}

	args := []string{}

	for _, v := range inSpec.GetChaincodeSpec().GetInput().GetArgs() {
		args = append(args, string(v))
	}

	ti = &txnInfo{
		TxnID:      chHeader.GetTxId(),
		CreateTime: chHeader.GetTimestamp(),
		Args:       args,
	}
	return
}
