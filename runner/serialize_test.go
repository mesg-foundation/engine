package runner

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
)

var data = &Runner{
	Address:      "hello",
	InstanceHash: sdk.AccAddress(crypto.AddressHash([]byte("10"))),
}

func TestHashSerialize(t *testing.T) {
	require.Equal(t, "2:hello;3:cosmos1ffzdc9fkggz2srlgp6grj32uc9sg9qvzu4jym8;", data.HashSerialize())
	require.Equal(t, "cosmos1fuz7jpafvpc08q7d95jhavzpqs85jcamsfm8nj", sdk.AccAddress(crypto.AddressHash([]byte(data.HashSerialize()))).String())
}

func BenchmarkHashSerialize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data.HashSerialize()
	}
}
