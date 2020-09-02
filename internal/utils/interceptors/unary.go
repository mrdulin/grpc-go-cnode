package interceptors

import (
	"context"
  "go.opencensus.io/trace"
  "time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func NewUnaryInterceptor(logger grpclog.LoggerV2) grpc.UnaryServerInterceptor {
	var Unary = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	  ctx, span := trace.StartSpan(ctx, info.FullMethod)
	  defer span.End()
		start := time.Now()
		h, err := handler(ctx, req)
		logger.Infof("Request - Method:%s\tDuration:%s\tError:%v\n", info.FullMethod, time.Since(start), err)
		return h, err
	}
	return Unary
}
