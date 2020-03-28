package main

import (
	"testing"

	servicemodule "github.com/mesg-foundation/engine/x/service"
	"github.com/stretchr/testify/require"
)

func testAccountSequence(t *testing.T) {
	t.Run("wrong msg", func(t *testing.T) {
		_, err = lcd.BroadcastMsg(servicemodule.MsgCreate{
			// Owner:  engineAddress,
			Sid:    "test-account-seq",
			Name:   "test-account-seq",
			Source: "QmXcPDajWN55n1UPV5VNJDEKE96xJFJMe3X7GwN3qx8p7r",
		})
		require.EqualError(t, err, "transaction returned with invalid code 18: invalid request: Owner is a required field")
	})
	t.Run("good msg", func(t *testing.T) {
		_, err = lcd.BroadcastMsg(servicemodule.MsgCreate{
			Owner:  engineAddress,
			Sid:    "test-account-seq",
			Name:   "test-account-seq",
			Source: "QmXcPDajWN55n1UPV5VNJDEKE96xJFJMe3X7GwN3qx8p7r",
		})
		require.NoError(t, err)
	})
	t.FailNow()
}
