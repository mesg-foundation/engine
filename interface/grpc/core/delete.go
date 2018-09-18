package core

import (
	"context"

	"github.com/mesg-foundation/core/protobuf/coreapi"
)

// DeleteService stops and deletes service serviceID.
func (s *Server) DeleteService(ctx context.Context, request *coreapi.DeleteServiceRequest) (*coreapi.DeleteServiceReply, error) {
	return &coreapi.DeleteServiceReply{}, s.api.DeleteService(request.ServiceID)
}
