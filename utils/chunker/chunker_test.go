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
