package mesg

import "github.com/mesg-foundation/core/container"

// MESG gives all functionalies of MESG core.
type MESG struct {
	container *container.Container
}

// Option is a configuration func for MESG.
type Option func(*MESG)

// New creates a new MESG with given options.
func New(options ...Option) (*MESG, error) {
	m := &MESG{}
	for _, option := range options {
		option(m)
	}
	var err error
	if m.container == nil {
		m.container, err = container.New()
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}

func DockerClientOption(container *container.Container) Option {
	return func(m *MESG) {
		m.container = container
	}
}
