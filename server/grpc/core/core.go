package core

import (
	"context"

	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/mesg-foundation/core/sdk"
	"github.com/mesg-foundation/core/version"
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
	services := make([]*coreapi.InfoReply_CoreService, len(c.Services()))
	for i, s := range c.Services() {
		services[i] = &coreapi.InfoReply_CoreService{
			Sid:  s.Sid,
			Hash: s.Hash.String(),
			Url:  s.URL,
			Key:  s.Key,
		}
	}
	return &coreapi.InfoReply{
		Version:  version.Version,
		Services: services,
	}, nil
}
