package application

import "context"

type Execution struct {
	ID  string
	Err error
}

type Stream struct {
	Executions chan *Execution
	Err        chan error
	cancel     context.CancelFunc
}

func (s *Stream) Close() error {
	s.cancel()
	return nil
}
