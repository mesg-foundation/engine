package core

import (
	"errors"

	"github.com/mesg-foundation/core/api"
)

// Server is the type to aggregate all the APIs.
type Server struct {
	api *api.API
}

// Option is a configuration func for Server.
type Option func(*Server)

// NewServer creates a new Server with given options.
func NewServer(options ...Option) (*Server, error) {
	s := &Server{}
	for _, option := range options {
		option(s)
	}
	if s.api == nil {
		return nil, errors.New("api should be provided")
	}
	return s, nil
}

// APIOption sets underlying mesg API.
func APIOption(api *api.API) Option {
	return func(s *Server) {
		s.api = api
	}
}
