package runner

import (
	"testing"

	"github.com/mesg-foundation/engine/hash"
	"github.com/stretchr/testify/require"
)

var data = &Runner{
	Address:      "hello",
	InstanceHash: hash.Int(10),
}

func TestHashSerialize(t *testing.T) {
	require.Equal(t, "2:hello;3:g35TxFqwMx95vCk63fTxGTHb6ei4W24qg5t2x6xD3cT;", data.HashSerialize())
	require.Equal(t, "8kS6ayPnvzjqNFTY9RYyEwhD3D4HMq9zZnDg7iFcUQYi", hash.Dump(data).String())
}

func BenchmarkHashSerialize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data.HashSerialize()
	}
}
