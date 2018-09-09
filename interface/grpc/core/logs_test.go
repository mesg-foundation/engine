package core

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogPiper(t *testing.T) {
	var (
		typ       = LogData_Data_Standard
		dep       = "1"
		dataChunk = []byte{1, 2}
		stream    = errorCloser{bytes.NewReader(dataChunk)}
		chunks    = make(chan logChunk)
	)

	newLogPiper(typ, dep, stream, chunks)

	require.Equal(t, logChunk{
		Dependency: dep,
		Type:       typ,
		Data:       dataChunk,
	}, <-chunks)

	stream.Close()
	require.Equal(t, logChunk{Err: io.EOF}, <-chunks)
}

type errorCloser struct {
	io.Reader
}

func (c errorCloser) Close() error {
	return io.EOF
}
