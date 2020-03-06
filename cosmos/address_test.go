package cosmos

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
)

func TestServiceAddress(t *testing.T) {
	address := ServiceAddress(crypto.AddressHash([]byte("hello")))
	fmt.Println(address.String())
	json, err := sdk.AccAddress(address).MarshalJSON()
	require.NoError(t, err)
	fmt.Println(string(json))
}
