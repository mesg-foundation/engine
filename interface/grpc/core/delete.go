package core

import (
	"context"

	"github.com/mesg-foundation/core/protobuf/core"
)

// DeleteService stops and deletes service serviceID.
func (s *Server) DeleteService(ctx context.Context, request *core.DeleteServiceRequest) (*core.DeleteServiceReply, error) {
	return &core.DeleteServiceReply{}, s.api.DeleteService(request.ServiceID)
}
