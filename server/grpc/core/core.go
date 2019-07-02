package core

import (
	"context"

	"github.com/mesg-foundation/engine/config"
	"github.com/mesg-foundation/engine/protobuf/coreapi"
	"github.com/mesg-foundation/engine/sdk"
	"github.com/mesg-foundation/engine/version"
)

// Server is the type to aggregate all the APIs.
type Server struct {
	sdk *sdk.SDK
}

// NewServer creates a new Server.
func NewServer(sdk *sdk.SDK) *Server {
	return &Server{sdk: sdk}
}

// Info returns all necessary information from the core.
func (s *Server) Info(ctx context.Context, request *coreapi.InfoRequest) (*coreapi.InfoReply, error) {
	c, err := config.Global()
	if err != nil {
		return nil, err
	}
	servicesFromConfig, err := c.Services()
	if err != nil {
		return nil, err
	}
	services := make([]*coreapi.InfoReply_CoreService, len(servicesFromConfig))
	for i, s := range servicesFromConfig {
		services[i] = &coreapi.InfoReply_CoreService{
			Sid:  s.Definition.Sid,
			Hash: s.Instance.Hash.String(),
			Url:  s.Definition.Source,
			Key:  s.Key,
		}
	}
	return &coreapi.InfoReply{
		Version:  version.Version,
		Services: services,
	}, nil
}
