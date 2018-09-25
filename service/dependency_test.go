package service

import (
	"io"
	"sync"
	"testing"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/stretchr/testify/require"
)

func TestDependencyLogs(t *testing.T) {
	testDependencyLogs(t, func(s *Service, dependencyKey string) (rstd, rerr io.ReadCloser,
		err error) {
		dep, err := s.getDependency(dependencyKey)
		require.NoError(t, err)
		return dep.Logs()
	})
}

func testDependencyLogs(t *testing.T, do func(s *Service, dependencyKey string) (rstd, rerr io.ReadCloser,
	err error)) {
	var (
		dependencyKey = "1"
		stdData       = []byte{1, 2}
		errData       = []byte{3, 4}
	)

	rp, wp := io.Pipe()
	wstd := stdcopy.NewStdWriter(wp, stdcopy.Stdout)
	werr := stdcopy.NewStdWriter(wp, stdcopy.Stderr)

	go wstd.Write(stdData)
	go werr.Write(errData)

	s, dt := newFromServiceAndDockerTest(t, &Service{
		Dependencies: []*Dependency{
			{Key: dependencyKey},
		},
	})
	dt.ProvideServiceLogs(rp, nil)

	rstd, rerr, err := do(s, dependencyKey)
	require.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		buf := make([]byte, 2)
		_, err := rstd.Read(buf)
		require.NoError(t, err)
		require.Equal(t, stdData, buf)
	}()

	go func() {
		defer wg.Done()
		buf := make([]byte, 2)
		_, err = rerr.Read(buf)
		require.NoError(t, err)
		require.Equal(t, errData, buf)
	}()

	wg.Wait()
}
