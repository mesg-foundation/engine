package grpc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestServerServe(t *testing.T) {
	s := New(nil)
	go func() {
		time.Sleep(500 * time.Millisecond)
		s.Close()
	}()
	require.NoError(t, s.Serve("localhost:50052"))
}
