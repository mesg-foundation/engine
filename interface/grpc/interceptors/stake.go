package interceptors

import (
	"context"
	"regexp"

	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/config"
	"google.golang.org/grpc"
)

const (
	unaryPattern  = "^/api.Core/(StopService|StartService|ExecuteTask|CreateWorkflow)$"
	streamPattern = "^/api.Core/(ListenEvent|ListenResult|ServiceLogs)$"
)

// UnaryStakeInterceptor ...
func UnaryStakeInterceptor(api *api.API) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if err := checkAuthorization(api, unaryPattern, info.FullMethod); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

// StreamStakeInterceptor ...
func StreamStakeInterceptor(api *api.API) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if err := checkAuthorization(api, streamPattern, info.FullMethod); err != nil {
			return err
		}
		return handler(srv, ss)
	}
}

func checkAuthorization(api *api.API, pattern string, method string) error {
	ok, _ := regexp.MatchString(pattern, method)
	if !ok {
		return nil
	}
	c, err := config.Global()
	if err != nil {
		return err
	}
	k, err := c.GetKey()
	if err != nil {
		return err
	}
	return api.RequireStake(k.Address)
}
