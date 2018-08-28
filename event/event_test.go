package event

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	s := service.Service{
		Name: "TestCreate",
		Events: map[string]*service.Event{
			"test": {},
		},
	}
	var data map[string]interface{}
	exec, err := Create(&s, "test", data)
	require.Nil(t, err)
	require.Equal(t, &s, exec.Service)
	require.Equal(t, data, exec.Data)
	require.Equal(t, "test", exec.Key)
	require.NotNil(t, exec.CreatedAt)
}
