package ownership

import (
	"testing"

	"github.com/mesg-foundation/engine/hash"
	"github.com/stretchr/testify/require"
)

var data = &Ownership{
	Owner:        "hello",
	ResourceHash: hash.Int(5),
	Resource:     Ownership_Process,
}

func TestHashSerialize(t *testing.T) {
	require.Equal(t, "2:hello;3:LX3EUdRUBUa3TbsYXLEUdj9J3prXkWXvLYSWyYyc2Jj;", data.HashSerialize())
	require.Equal(t, "GzZBiyQWRkDAAwnpckM9VvZxtZZEZii7UUsjjxycZJ8N", hash.Dump(data).String())
}

func BenchmarkHashSerialize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data.HashSerialize()
	}
}
