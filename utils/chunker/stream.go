package chunker

import "io"

// Stream implements io.Reader.
type Stream struct {
	recv chan []byte
	done chan struct{}

	data []byte
	i    int64
}

// NewStream returns a new stream.
func NewStream() *Stream {
	return &Stream{
		recv: make(chan []byte, 0),
		done: make(chan struct{}, 0),
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
				return 0, io.EOF

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
	if s.done == nil {
		return nil
	}
	close(s.done)
	s.done = nil
	return nil
}
