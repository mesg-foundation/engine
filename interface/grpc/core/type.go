package core

import (
	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/systemservices"
)

// Server is the type to aggregate all the APIs.
type Server struct {
	api *api.API
	ss  *systemservices.SystemServices
}

// NewServer creates a new Server.
func NewServer(api *api.API, ss *systemservices.SystemServices) *Server {
	return &Server{api: api, ss: ss}
}
