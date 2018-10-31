package chunker

import (
	"errors"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogReader(t *testing.T) {
	var (
		chunk1 = []byte{1}
		chunk2 = []byte{2}
	)

	s := NewStream()

	go func() {
		s.Provide(chunk1)
		s.Provide(chunk2)
		s.Close()
	}()

	data, err := ioutil.ReadAll(s)
	require.NoError(t, err)
	require.Len(t, data, 2)
	require.Equal(t, chunk1, []byte{data[0]})
	require.Equal(t, chunk2, []byte{data[1]})
}

func TestLogReaderCloseWithError(t *testing.T) {
	anErr := errors.New("oh an error")

	s := NewStream()
	go s.CloseWithError(anErr)

	_, err := ioutil.ReadAll(s)
	require.Error(t, anErr, err)
}
