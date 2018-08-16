package core

import (
	"errors"

	"github.com/mesg-foundation/core/mesg"
)

// Server is the type to aggregate all the APIs.
type Server struct {
	mesg *mesg.MESG
}

// Option is a configuration func for Server.
type Option func(*Server)

// NewServer creates a new Server with given options.
func NewServer(options ...Option) (*Server, error) {
	s := &Server{}
	for _, option := range options {
		option(s)
	}
	if s.mesg == nil {
		return nil, errors.New("mesg should be provided")
	}
	return s, nil
}

func MESGOption(mesg *mesg.MESG) Option {
	return func(s *Server) {
		s.mesg = mesg
	}
}
