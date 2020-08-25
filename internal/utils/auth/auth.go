package auth

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const jwt = "Bearer 123"

type Authentication struct {
	Authorization string
}

func (a *Authentication) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{"authorization": a.Authorization}, nil
}

func (a *Authentication) RequireTransportSecurity() bool {
	return false
}

func (a *Authentication) Auth(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.InvalidArgument, "missing credentials")
	}
	var (
		token string
	)
	if val, ok := md["authorization"]; ok {
		token = val[0]
	}
	if token != jwt {
		return status.Errorf(codes.Unauthenticated, "invalid token")
	}
	return nil
}
