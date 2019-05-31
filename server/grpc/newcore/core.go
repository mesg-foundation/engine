package newcore

import (
	"context"

	"github.com/mesg-foundation/core/protobuf/service"
	servicesdk "github.com/mesg-foundation/core/sdk/service"
)

// Server is the type to aggregate all the APIs.
type Server struct {
	sdk *servicesdk.ServiceSDK
}

// NewServer creates a new Server.
func NewServer(sdk *servicesdk.ServiceSDK) *Server {
	return &Server{sdk: sdk}
}

// Create creates a new service from definition.
func (s *Server) Create(ctx context.Context, request *service.CreateRequest) (*service.CreateResponse, error) {
	srv := fromProtoService(request.Definition)
	if err := s.sdk.Create(srv); err != nil {
		return nil, err
	}
	return &service.CreateResponse{Sid: srv.Sid, Hash: srv.Hash}, nil
}
