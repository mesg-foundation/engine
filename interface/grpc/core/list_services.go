package core

import (
	"context"

	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/protobuf/coreapi"
)

// ListServices lists services.
func (s *Server) ListServices(ctx context.Context, request *coreapi.ListServicesRequest) (*coreapi.ListServicesReply, error) {
	var filters []api.ListServicesFilter
	if request.FilterRunning {
		filters = append(filters, api.ListRunningServicesFilter())
	}
	services, err := s.api.ListServices(filters...)
	if err != nil {
		return nil, err
	}
	protoServices := toProtoServices(services)
	if request.FilterRunning {
		for _, service := range protoServices {
			service.IsRunning = true
		}
	}
	return &coreapi.ListServicesReply{Services: protoServices}, nil
}
