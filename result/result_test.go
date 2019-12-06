package result

import (
	"testing"

	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	var (
		execHash = hash.Int(1)
	)
	t.Run("NewWithOutputs", func(t *testing.T) {
		res := NewWithOutputs(execHash, &types.Struct{
			Fields: map[string]*types.Value{
				"test": {Kind: &types.Value_StringValue{StringValue: "hello"}},
			},
		})
		require.True(t, res.Hash.Valid())
		require.True(t, execHash.Equal(res.ExecutionHash))
	})
	t.Run("NewWithError", func(t *testing.T) {
		resErr := NewWithError(execHash, "error string")
		require.True(t, resErr.Hash.Valid())
		require.True(t, execHash.Equal(resErr.ExecutionHash))
	})
}
