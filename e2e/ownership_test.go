package main

import (
	"context"
	"testing"

	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/ownership"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/stretchr/testify/require"
)

func testOwnership(t *testing.T) {
	t.Run("list", func(t *testing.T) {
		ownerships, err := client.OwnershipClient.List(context.Background(), &pb.ListOwnershipRequest{})
		require.NoError(t, err)

		acc, err := client.AccountClient.Get(context.Background(), &pb.GetAccountRequest{Name: "engine"})
		require.NoError(t, err)

		require.Len(t, ownerships.Ownerships, 1)
		require.Equal(t,
			hash.Dump(&ownership.Ownership{
				Owner:        acc.Address,
				Resource:     ownership.Ownership_Service,
				ResourceHash: testServiceHash,
			}),
			ownerships.Ownerships[0].Hash,
		)
	})
}
