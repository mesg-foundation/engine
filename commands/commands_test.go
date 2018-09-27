package commands

import (
	"bufio"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

// outputStream is a testing utility to capture io writes and read
// written data with different ways.
type outputStream struct {
	t *testing.T
	w io.Writer
	r *bufio.Reader
}

// newOutputStream receives t to assert errors generated during
// the writes and reads.
func newOutputStream(t *testing.T) *outputStream {
	r, w := io.Pipe()
	return &outputStream{
		t: t,
		w: w,
		r: bufio.NewReader(r),
	}
}

// Write implements io.Writer.
func (o *outputStream) Write(b []byte) (n int, err error) {
	return o.w.Write(b)
}

// ReadLine returns the next line in the stream.
func (o *outputStream) ReadLine() []byte {
	line, _, err := o.r.ReadLine()
	require.NoError(o.t, err)
	return line
}

// ReadAll returns all the data in thee stream until it's closed.
func (o *outputStream) ReadAll() []byte {
	bytes, err := ioutil.ReadAll(o.r)
	require.NoError(o.t, err)
	return bytes
}
