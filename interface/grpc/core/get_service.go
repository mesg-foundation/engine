package core

import (
	"context"

	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/mesg-foundation/core/service"
)

// GetService returns service serviceID.
func (s *Server) GetService(ctx context.Context, request *coreapi.GetServiceRequest) (*coreapi.GetServiceReply, error) {
	ss, err := s.api.GetService(request.ServiceID)
	if err != nil {
		return nil, err
	}
	protoService := toProtoService(ss)
	status, err := ss.Status()
	if err != nil {
		return nil, err
	}
	protoService.IsRunning = status == service.RUNNING
	return &coreapi.GetServiceReply{Service: protoService}, nil
}
