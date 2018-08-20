package services

import (
	"sync"
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestDb(t *testing.T) {
	db, err := open()
	defer close()
	require.Nil(t, err)
	require.NotNil(t, db)
}

// Test to stress the database with concurrency access
// BUG: https://github.com/mesg-foundation/core/issues/163
func TestConcurrency(t *testing.T) {
	var wg sync.WaitGroup
	service := &service.Service{
		Name: "TestConcurrency",
	}
	hash, _ := Save(service)
	defer Delete(hash)
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			s, err := Get(hash)
			require.Nil(t, err)
			require.Equal(t, s.Name, service.Name)
			wg.Done()
		}()
	}
	wg.Wait()
}
