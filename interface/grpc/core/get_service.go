package core

import (
	"context"

	"github.com/mesg-foundation/core/protobuf/coreapi"
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
	protoService.Status = toProtoServiceStatusType(status)
	return &coreapi.GetServiceReply{Service: protoService}, nil
}
