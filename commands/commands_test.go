package commands

import (
	"bufio"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

type outputStream struct {
	t *testing.T
	w io.Writer
	r *bufio.Reader
}

func newoutputStream(t *testing.T) *outputStream {
	r, w := io.Pipe()
	return &outputStream{
		t: t,
		w: w,
		r: bufio.NewReader(r),
	}
}

func (o *outputStream) Write(b []byte) (n int, err error) {
	return o.w.Write(b)
}

func (o *outputStream) ReadLine() []byte {
	line, _, err := o.r.ReadLine()
	require.NoError(o.t, err)
	return line
}

func (o *outputStream) ReadAll() []byte {
	bytes, err := ioutil.ReadAll(o.r)
	require.NoError(o.t, err)
	return bytes
}
