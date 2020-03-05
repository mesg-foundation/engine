package instance

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/hash"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
)

var data = &Instance{
	ServiceHash: sdk.AccAddress(crypto.AddressHash([]byte("10"))),
	EnvHash:     hash.Int(1),
}

func TestHashSerialize(t *testing.T) {
	require.Equal(t, "2:cosmos1ffzdc9fkggz2srlgp6grj32uc9sg9qvzu4jym8;3:4uQeVj5tqViQh7yWWGStvkEG1Zmhx6uasJtWCJziofM;", data.HashSerialize())
	require.Equal(t, "cosmos17zm8re7hr6m8e96n0vwf6guvffrl49z5td6gg0", sdk.AccAddress(crypto.AddressHash([]byte(data.HashSerialize()))).String())
}

func BenchmarkHashSerialize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data.HashSerialize()
	}
}
