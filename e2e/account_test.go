package main

import (
	"context"
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func TestAccount(t *testing.T) {
	var address string
	t.Run("create", func(t *testing.T) {
		resp, err := client.AccountClient.Create(context.Background(), &pb.CreateAccountRequest{
			Name:     "user",
			Password: "pass",
		})
		require.NoError(t, err)
		address = resp.Address

		acc, err := types.AccAddressFromBech32(resp.Address)
		require.NoError(t, err)

		_, err = testkb.Get("user")
		require.NoError(t, err)

		_, err = testkb.GetByAddress(acc)
		require.NoError(t, err)
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

		accs, err := testkb.List()
		require.NoError(t, err)

		require.Equal(t, len(resp.Accounts), len(accs))
	loop:
		for _, acc := range accs {
			addr, name := acc.GetAddress().String(), acc.GetName()
			for _, acc := range resp.Accounts {
				if acc.Name == name && acc.Address == addr {
					continue loop
				}
			}
			t.Fatal("no account", addr, name)

		}
	})

	t.Run("delete", func(t *testing.T) {
		ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs(
			"credential_username", "user",
			"credential_passphrase", "pass",
		))
		_, err := client.AccountClient.Delete(ctx, &pb.DeleteAccountRequest{})
		require.NoError(t, err)
	})
}
