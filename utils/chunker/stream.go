package chunker

import (
	"io"
	"sync"
)

// Stream implements io.Reader.
type Stream struct {
	recv chan []byte

	closeErr error
	me       sync.Mutex

	done chan struct{}
	m    sync.Mutex

	data []byte
	i    int64
}

// NewStream returns a new stream.
func NewStream() *Stream {
	return &Stream{
		closeErr: io.EOF,
		recv:     make(chan []byte),
		done:     make(chan struct{}),
	}
}

// Provide provides data for Read.
func (s *Stream) Provide(data []byte) {
	s.recv <- data
}

// Read puts data into p.
func (s *Stream) Read(p []byte) (n int, err error) {
	if s.i >= int64(len(s.data)) {
		for {
			select {
			case <-s.done:
				s.me.Lock()
				defer s.me.Unlock()
				return 0, s.closeErr

			case data := <-s.recv:
				if err != nil {
					return 0, err
				}
				s.data = data
				s.i = 0
				return s.Read(p)
			}
		}
	}
	n = copy(p, s.data[s.i:])
	s.i += int64(n)
	return n, nil
}

// Close closes reader.
func (s *Stream) Close() error {
	s.m.Lock()
	defer s.m.Unlock()
	if s.done == nil {
		return nil
	}
	s.done <- struct{}{}
	s.done = nil
	return nil
}

// CloseWithError closes reader with an error.
func (s *Stream) CloseWithError(err error) {
	s.me.Lock()
	s.closeErr = err
	s.me.Unlock()
	s.Close()
}
