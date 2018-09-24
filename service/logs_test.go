package service

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServiceLogs(t *testing.T) {
	testDependencyLogs(t, func(s *Service, dependencyKey string) (rstd, rerr io.ReadCloser,
		err error) {
		l, err := s.Logs(dependencyKey)
		require.NoError(t, err)
		require.Len(t, l, 1)
		return l[0].Standard, l[0].Error, nil
	})
}
