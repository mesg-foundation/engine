package service

import (
	"io"
	"testing"

	"github.com/mesg-foundation/core/container"
	"github.com/stretchr/testify/require"
)

func TestServiceLogs(t *testing.T) {
	testDependencyLogs(t, func(s *Service, c container.Container, dependencyKey string) (rstd, rerr io.ReadCloser,
		err error) {
		l, err := s.Logs(c, dependencyKey)
		require.NoError(t, err)
		require.Len(t, l, 1)
		return l[0].Standard, l[0].Error, nil
	})
}
