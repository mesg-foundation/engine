package main

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/ownership"
	"github.com/mesg-foundation/engine/service"
	"github.com/stretchr/testify/require"
)

var (
	testServiceHash    hash.Hash
	testServiceAddress sdk.AccAddress
)

func testService(t *testing.T) {
	msg := newTestCreateServiceMsg()

	t.Run("create", func(t *testing.T) {
		msg.Owner = engineAddress
		testServiceHash = lcdBroadcastMsg(t, msg)
	})

	t.Run("get", func(t *testing.T) {
		var s *service.Service
		lcdGet(t, "service/get/"+testServiceHash.String(), &s)
		require.Equal(t, testServiceHash, s.Hash)
		testServiceAddress = s.Address
	})

	t.Run("list", func(t *testing.T) {
		ss := make([]*service.Service, 0)
		lcdGet(t, "service/list", &ss)
		require.Len(t, ss, 1)
	})

	t.Run("exists", func(t *testing.T) {
		var exist bool
		lcdGet(t, "service/exist/"+testServiceHash.String(), &exist)
		require.True(t, exist)
		lcdGet(t, "service/exist/"+hash.Int(1).String(), &exist)
		require.False(t, exist)
	})

	t.Run("hash", func(t *testing.T) {
		var hash hash.Hash
		lcdPost(t, "service/hash", msg, &hash)
		require.Equal(t, testServiceHash, hash)
	})

	t.Run("check ownership creation", func(t *testing.T) {
		ownerships := make([]*ownership.Ownership, 0)
		lcdGet(t, "ownership/list", &ownerships)
		require.Len(t, ownerships, 1)
		require.NotEmpty(t, ownerships[0].Owner)
		require.Equal(t, ownership.Ownership_Service, ownerships[0].Resource)
		require.Equal(t, testServiceHash, ownerships[0].ResourceHash)
	})
}
