// Copyright 2018 MESG Foundation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package chunker

import (
	"io"
	"sync"
)

// Stream implements io.Reader.
type Stream struct {
	recv chan []byte

	done chan struct{}
	m    sync.Mutex

	data []byte
	i    int64
}

// NewStream returns a new stream.
func NewStream() *Stream {
	return &Stream{
		recv: make(chan []byte),
		done: make(chan struct{}),
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
	s.m.Lock()
	defer s.m.Unlock()
	if s.done == nil {
		return nil
	}
	s.done <- struct{}{}
	s.done = nil
	return nil
}
