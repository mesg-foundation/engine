package main

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	creditmodule "github.com/mesg-foundation/engine/x/credit"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
)

func testCredit(t *testing.T) {
	var (
		newAddr = sdk.AccAddress(crypto.AddressHash([]byte("addr")))
	)
	t.Run("account does not exist", func(t *testing.T) {
		var accR *auth.BaseAccount
		require.NoError(t, lcd.Get("auth/accounts/"+newAddr.String(), &accR))
		require.True(t, accR.Address.Empty())
	})
	t.Run("add 1000 credits", func(t *testing.T) {
		t.Run("balance before", func(t *testing.T) {
			var balance sdk.Int
			require.NoError(t, lcd.Get("credit/get/"+newAddr.String(), &balance))
			require.Equal(t, sdk.NewInt(0), balance)
		})
		t.Run("add credits", func(t *testing.T) {
			msg := creditmodule.MsgAdd{
				Signer:  cliAddress,
				Address: newAddr,
				Amount:  sdk.NewInt(1000),
			}
			_, err = lcd.BroadcastMsg(msg)
			require.NoError(t, err)
		})
		t.Run("balance after", func(t *testing.T) {
			var balance sdk.Int
			require.NoError(t, lcd.Get("credit/get/"+newAddr.String(), &balance))
			require.Equal(t, sdk.NewInt(1000), balance)
		})
	})
	t.Run("account exists", func(t *testing.T) {
		var accR *auth.BaseAccount
		require.NoError(t, lcd.Get("auth/accounts/"+newAddr.String(), &accR))
		require.False(t, accR.Address.Empty())
	})
	t.Run("not authorized", func(t *testing.T) {
		msg := creditmodule.MsgAdd{
			Signer:  engineAddress,
			Address: newAddr,
			Amount:  sdk.NewInt(1000),
		}
		_, err = lcdEngine.BroadcastMsg(msg)
		require.EqualError(t, err, "transaction returned with invalid code 4: unauthorized: the signer is not a minter: failed to execute message; message index: 0")
	})
}
