package service

import (
	"testing"

	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/stretchr/testify/require"
)

func TestRequire(t *testing.T) {
	sp := []*Service_Parameter{
		{
			Key:      "key",
			Type:     "Number",
			Optional: true,
		},
	}

	invalid := &types.Struct{
		Fields: map[string]*types.Value{
			"key": {
				Kind: &types.Value_StringValue{},
			},
		},
	}

	s := &Service{
		Events: []*Service_Event{
			{
				Key:  "event",
				Data: sp,
			},
		},
		Tasks: []*Service_Task{
			{
				Key:     "task",
				Inputs:  sp,
				Outputs: sp,
			},
		},
	}

	t.Run("RequireTaskInputs", func(t *testing.T) {
		require.NoError(t, s.RequireTaskInputs("task", nil))
		require.Error(t, s.RequireTaskInputs("task", invalid))
		require.Error(t, s.RequireTaskInputs("-", nil))
	})

	t.Run("RequireTaskOutputs", func(t *testing.T) {
		require.NoError(t, s.RequireTaskOutputs("task", nil))
		require.Error(t, s.RequireTaskOutputs("task", invalid))
		require.Error(t, s.RequireTaskOutputs("-", nil))
	})

	t.Run("RequireEventData", func(t *testing.T) {
		require.NoError(t, s.RequireEventData("event", nil))
		require.Error(t, s.RequireEventData("event", invalid))
		require.Error(t, s.RequireEventData("-", nil))
	})
}
