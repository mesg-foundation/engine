package instance

import (
	"testing"

	"github.com/mesg-foundation/engine/hash"
	"github.com/stretchr/testify/require"
)

var data = &Instance{
	ServiceHash: hash.Int(10),
	EnvHash:     hash.Int(5),
}

func TestHashSerialize(t *testing.T) {
	require.Equal(t, "2:g35TxFqwMx95vCk63fTxGTHb6ei4W24qg5t2x6xD3cT;3:LX3EUdRUBUa3TbsYXLEUdj9J3prXkWXvLYSWyYyc2Jj;", data.HashSerialize())
	require.Equal(t, "BwWnWRgpPfB9SPmSRKYZp8Dq1LEpCpwHAGZFsJJEg1nd", hash.Dump(data).String())
}

func BenchmarkHashSerialize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data.HashSerialize()
	}
}
