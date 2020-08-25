package interceptors

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func NewUnaryInterceptor(logger grpclog.LoggerV2) grpc.UnaryServerInterceptor {
	var Unary = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()
		h, err := handler(ctx, req)
		logger.Infof("Request - Method:%s\tDuration:%s\tError:%v\n", info.FullMethod, time.Since(start), err)
		return h, err
	}
	return Unary
}
