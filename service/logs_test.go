package service

import (
	"io"
	"sync"
	"testing"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/mesg-foundation/core/container/mocks"
	"github.com/stretchr/testify/require"
)

func TestServiceLogs(t *testing.T) {
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

	var (
		s = &Service{
			Hash: "1",
			Dependencies: []*Dependency{
				{Key: dependencyKey},
			},
		}
		mc = &mocks.Container{}
	)

	d, _ := s.getDependency(dependencyKey)
	mc.On("ServiceLogs", d.namespace(s.namespace())).Once().Return(rp, nil)

	l, err := s.Logs(mc, dependencyKey)
	require.NoError(t, err)
	require.Len(t, l, 1)
	rstd, rerr := l[0].Standard, l[0].Error

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
	mc.AssertExpectations(t)
}
