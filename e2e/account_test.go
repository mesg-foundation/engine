package main

import (
	"context"
	"testing"

	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/stretchr/testify/require"
)

func testAccount(t *testing.T) {
	var address string
	t.Run("create", func(t *testing.T) {
		resp, err := client.AccountClient.Create(context.Background(), &pb.CreateAccountRequest{
			Name:     "user",
			Password: "pass",
		})
		require.NoError(t, err)
		address = resp.Address
	})

	t.Run("get", func(t *testing.T) {
		acc, err := client.AccountClient.Get(context.Background(), &pb.GetAccountRequest{Name: "user"})
		require.NoError(t, err)

		require.Equal(t, "user", acc.Name)
		require.Equal(t, address, acc.Address)
	})

	t.Run("list", func(t *testing.T) {
		resp, err := client.AccountClient.List(context.Background(), &pb.ListAccountRequest{})
		require.NoError(t, err)

		require.Equal(t, len(resp.Accounts), 2)
		require.Equal(t, resp.Accounts[0].Name, "engine")
		require.Equal(t, resp.Accounts[1].Name, "user")
	})

	// t.Run("delete", func(t *testing.T) {
	// 	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs(
	// 		"credential_username", "engine",
	// 		"credential_passphrase", "pass",
	// 	))
	// 	_, err := client.AccountClient.Delete(ctx, &pb.DeleteAccountRequest{})
	// 	require.NoError(t, err)

	// 	resp, err := client.AccountClient.List(context.Background(), &pb.ListAccountRequest{})
	// 	require.NoError(t, err)

	// 	require.Equal(t, len(resp.Accounts), 1)
	// 	require.Equal(t, resp.Accounts[0].Name, "user")

	// 	// recreate engine account to make sure rest of tests are working.
	// 	// to delete when proper account management is done.
	// 	_, err = client.AccountClient.Create(context.Background(), &pb.CreateAccountRequest{
	// 		Name:     "engine",
	// 		Password: "pass",
	// 	})
	// 	require.NoError(t, err)
	// })
}
