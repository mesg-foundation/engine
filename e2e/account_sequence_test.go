package main

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
)

func testAccountSequence(t *testing.T) {
	var (
		randomAddress = sdk.AccAddress(crypto.AddressHash([]byte("hello")))
		err           error
	)
	t.Run("wrong msg", func(t *testing.T) {
		_, err = lcd.BroadcastMsg(bank.NewMsgSend(
			engineAddress,
			randomAddress,
			sdk.NewCoins(sdk.NewInt64Coin("wrong", 1000)),
		))
		require.Error(t, err)
		require.Contains(t, err.Error(), "transaction returned with invalid code")
	})
	t.Run("good msg", func(t *testing.T) {
		_, err = lcd.BroadcastMsg(bank.NewMsgSend(
			engineAddress,
			randomAddress,
			sdk.NewCoins(sdk.NewInt64Coin("atto", 1000)),
		))
		require.NoError(t, err)
	})
}
