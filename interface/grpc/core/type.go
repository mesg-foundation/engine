package core

import (
	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/container"
)

// Server is the type to aggregate all the APIs.
type Server struct {
	container container.Container
	api       *api.API
}

// NewServer creates a new Server.
func NewServer(c container.Container, api *api.API) *Server {
	return &Server{container: c, api: api}
}
