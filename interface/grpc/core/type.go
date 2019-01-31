package core

import (
	"github.com/mesg-foundation/core/api"
)

// Server is the type to aggregate all the APIs.
type Server struct {
	api *api.API
}

// NewServer creates a new Server.
func NewServer(api *api.API) *Server {
	return &Server{api: api}
}
