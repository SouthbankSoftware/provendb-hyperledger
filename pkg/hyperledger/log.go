/*
 * @Author: guiguan
 * @Date:   2020-01-22T10:42:09+11:00
 * @Last modified by:   guiguan
 * @Last modified time: 2020-08-13T16:25:32+10:00
 */

package hyperledger

import (
	"context"
	"errors"
	"fmt"
	"path"

	"github.com/SouthbankSoftware/provendb-tree/pkg/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func handleGRPCError(ctx context.Context, fullMethod string, err error) error {
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return err
		}

		if errors.Is(err, context.Canceled) {
			return status.Error(codes.Aborted, "client canceled")
		}

		// automatically report
		return status.Error(codes.Internal, log.KillBugStr(ctx,
			fmt.Sprintf("failed to invoke `%s`", path.Base(fullMethod)),
			zap.Error(err),
		))
	}

	return nil
}

func logUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		resp, err := handler(ctx, req)

		err = handleGRPCError(ctx, info.FullMethod, err)

		return resp, err
	}
}

func logStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		err := handler(srv, stream)

		err = handleGRPCError(stream.Context(), info.FullMethod, err)

		return err
	}
}
