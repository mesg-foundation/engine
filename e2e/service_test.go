package main

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/ownership"
	"github.com/mesg-foundation/engine/service"
	servicerest "github.com/mesg-foundation/engine/x/service/client/rest"
	"github.com/stretchr/testify/require"
)

var (
	testServiceHash    hash.Hash
	testServiceAddress sdk.AccAddress
)

func testService(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		testCreateServiceMsg.Owner = engineAddress
		testServiceHash = lcdBroadcastMsg(testCreateServiceMsg)
	})

	t.Run("get", func(t *testing.T) {
		var s *service.Service
		lcdGet("service/get/"+testServiceHash.String(), &s)
		require.Equal(t, testServiceHash, s.Hash)
		testServiceAddress = s.Address
	})

	t.Run("list", func(t *testing.T) {
		ss := make([]*service.Service, 0)
		lcdGet("service/list", &ss)
		require.Len(t, ss, 1)
	})

	t.Run("exists", func(t *testing.T) {
		var exist bool
		lcdGet("service/exist/"+testServiceHash.String(), &exist)
		require.True(t, exist)
		lcdGet("service/exist/"+hash.Int(1).String(), &exist)
		require.False(t, exist)
	})

	t.Run("hash", func(t *testing.T) {
		msg := servicerest.HashRequest{
			Sid:           testCreateServiceMsg.Sid,
			Name:          testCreateServiceMsg.Name,
			Description:   testCreateServiceMsg.Description,
			Configuration: testCreateServiceMsg.Configuration,
			Tasks:         testCreateServiceMsg.Tasks,
			Events:        testCreateServiceMsg.Events,
			Dependencies:  testCreateServiceMsg.Dependencies,
			Repository:    testCreateServiceMsg.Repository,
			Source:        testCreateServiceMsg.Source,
		}
		var hash hash.Hash
		lcdPost("service/hash", msg, &hash)
		require.Equal(t, testServiceHash, hash)
	})

	t.Run("check ownership creation", func(t *testing.T) {
		ownerships := make([]*ownership.Ownership, 0)
		lcdGet("ownership/list", &ownerships)
		require.Len(t, ownerships, 1)
		require.NotEmpty(t, ownerships[0].Owner)
		require.Equal(t, ownership.Ownership_Service, ownerships[0].Resource)
		require.Equal(t, testServiceHash, ownerships[0].ResourceHash)
	})
}
