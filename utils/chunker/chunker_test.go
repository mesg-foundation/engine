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
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	var (
		chunk  = []byte{1, 2}
		r      = errorCloser{bytes.NewReader(chunk)}
		chunks = make(chan Data)
		errs   = make(chan error)
		value  = "1"
	)

	New(r, chunks, errs, ValueOption(value))

	require.Equal(t, Data{
		Value: value,
		Data:  chunk,
	}, <-chunks)

	r.Close()
	require.Equal(t, io.EOF, <-errs)
}

type errorCloser struct {
	io.Reader
}

func (c errorCloser) Close() error {
	return io.EOF
}
