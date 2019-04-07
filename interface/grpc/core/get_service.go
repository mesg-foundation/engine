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
	status, err := ss.Status()
	if err != nil {
		return nil, err
	}
	details := &coreapi.Service{
		Definition: toProtoService(ss),
		Status:     toProtoServiceStatusType(status),
	}
	return &coreapi.GetServiceReply{Service: details}, nil
}
