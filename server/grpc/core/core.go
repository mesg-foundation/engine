package core

import (
	"context"

	"github.com/mesg-foundation/engine/protobuf/coreapi"
	"github.com/mesg-foundation/engine/version"
)

// Server is the type to aggregate all the APIs.
type Server struct{}

// NewServer creates a new Server.
func NewServer() *Server {
	return &Server{}
}

// Info returns all necessary information from the core.
func (s *Server) Info(ctx context.Context, request *coreapi.InfoRequest) (*coreapi.InfoReply, error) {
	return &coreapi.InfoReply{Version: version.Version}, nil
}
