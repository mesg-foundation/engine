package core

import (
	"context"

	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/mesg-foundation/core/version"
)

// Info returns all necessary information from the core.
func (s *Server) Info(ctx context.Context, request *coreapi.InfoRequest) (*coreapi.InfoReply, error) {
	c, err := config.Global()
	if err != nil {
		return nil, err
	}
	services := make(map[string]*coreapi.InfoReply_CoreService)
	for i, s := range c.Services() {
		services[i] = &coreapi.InfoReply_CoreService{
			Sid:  s.Sid,
			Hash: s.Hash,
			Url:  s.URL,
		}
	}
	return &coreapi.InfoReply{
		Address:  c.Client.Address,
		Image:    c.Core.Image,
		Version:  version.Version,
		Services: services,
	}, nil
}
