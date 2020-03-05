package ownership

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
)

var data = &Ownership{
	Owner:        "hello",
	ResourceHash: sdk.AccAddress(crypto.AddressHash([]byte("5"))),
	Resource:     Ownership_Process,
}

func TestHashSerialize(t *testing.T) {
	require.Equal(t, "2:hello;3:cosmos1auk3yl0r0w2zh2ksv9z72jcvvxdp7g3jseahyw;", data.HashSerialize())
	require.Equal(t, "cosmos1qzqzt8t9hyz2kc7sv4md4skztxj8wew7a5l6c5", sdk.AccAddress(crypto.AddressHash([]byte(data.HashSerialize()))).String())
}

func BenchmarkHashSerialize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data.HashSerialize()
	}
}
