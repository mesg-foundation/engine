package main

import (
	"context"
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/core"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestAccount(t *testing.T) {
	genTestData(t)
	defer os.RemoveAll("testdata")

	kb, closer := core.Start("./testdata/testuser-0/config.yml")
	defer closer()
	conn, err := grpc.DialContext(context.Background(), "localhost:50052", grpc.WithInsecure())
	require.NoError(t, err)
	c := pb.NewAccountClient(conn)

	var address string

	t.Run("create", func(t *testing.T) {
		resp, err := c.Create(context.Background(), &pb.CreateAccountRequest{
			Name:     "user",
			Password: "pass",
		})
		require.NoError(t, err)
		address = resp.Address

		acc, err := types.AccAddressFromBech32(resp.Address)
		require.NoError(t, err)

		_, err = kb.Get("user")
		require.NoError(t, err)

		_, err = kb.GetByAddress(acc)
		require.NoError(t, err)
	})

	t.Run("get", func(t *testing.T) {
		acc, err := c.Get(context.Background(), &pb.GetAccountRequest{Name: "user"})
		require.NoError(t, err)

		require.Equal(t, "user", acc.Name)
		require.Equal(t, address, acc.Address)
	})

	t.Run("list", func(t *testing.T) {
		resp, err := c.List(context.Background(), &pb.ListAccountRequest{})
		require.NoError(t, err)

		accs, err := kb.List()
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
		md := metadata.Pairs(
			"credential_username", "user",
			"credential_passphrase", "pass",
		)
		ctx := metadata.NewOutgoingContext(context.Background(), md)
		_, err := c.Delete(ctx, &pb.DeleteAccountRequest{})
		require.NoError(t, err)
	})
}
