package service

import (
	"github.com/mesg-foundation/core/api"
)

// Server binds all api functions.
type Server struct {
	api *api.API
}

// NewServer creates a new Server.
func NewServer(api *api.API) *Server {
	return &Server{api: api}
}
