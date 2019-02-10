package core

import (
	"context"

	"github.com/mesg-foundation/core/protobuf/coreapi"
)

// ListServices lists services.
func (s *Server) ListServices(ctx context.Context, request *coreapi.ListServicesRequest) (*coreapi.ListServicesReply, error) {
	services, err := s.api.ListServices()
	if err != nil {
		return nil, err
	}

	var protoServices []*coreapi.Service
	for _, s := range services {
		protoServices = append(protoServices, toProtoService(s))
	}

	return &coreapi.ListServicesReply{Services: protoServices}, nil
}
