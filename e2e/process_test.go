package main

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/ownership"
	"github.com/mesg-foundation/engine/process"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	processmodule "github.com/mesg-foundation/engine/x/process"
	"github.com/stretchr/testify/require"
)

func testProcess(t *testing.T) {
	var (
		processHash    hash.Hash
		processAddress sdk.AccAddress
		msg            = &processmodule.MsgCreate{
			Owner: engineAddress,
			Name:  "test-process",
			Nodes: []*process.Process_Node{
				{
					Key: "n0",
					Type: &process.Process_Node_Event_{
						Event: &process.Process_Node_Event{
							InstanceHash: testInstanceHash,
							EventKey:     "test_service_ready",
						},
					},
				},
				{
					Key: "n1",
					Type: &process.Process_Node_Task_{
						Task: &process.Process_Node_Task{
							InstanceHash: testInstanceHash,
							TaskKey:      "test_service_ready",
						},
					},
				},
			},
			Edges: []*process.Process_Edge{
				{
					Src: "n0",
					Dst: "n1",
				},
			},
		}
	)

	t.Run("create", func(t *testing.T) {
		res, err := cclient.BuildAndBroadcastMsg(msg)
		require.NoError(t, err)
		processHash = res.Data
	})

	t.Run("get", func(t *testing.T) {
		var p *process.Process
		lcdGet(t, "process/get/"+processHash.String(), &p)
		require.True(t, p.Equal(&process.Process{
			Hash:    p.Hash,
			Address: p.Address,
			Name:    msg.Name,
			Nodes:   msg.Nodes,
			Edges:   msg.Edges,
		}))
		processAddress = p.Address
	})

	t.Run("check ownership creation", func(t *testing.T) {
		ownerships := make([]*ownership.Ownership, 0)
		lcdGet(t, "ownership/list", &ownerships)
		owners := make([]*ownership.Ownership, 0)
		for _, o := range ownerships {
			if o.ResourceHash.Equal(processHash) && o.Resource == ownership.Ownership_Process && o.Owner != "" {
				owners = append(owners, o)
			}
		}
		require.Len(t, owners, 1)
	})

	t.Run("check coins on process", func(t *testing.T) {
		var coins sdk.Coins
		param := bank.NewQueryBalanceParams(processAddress)
		require.NoError(t, cclient.QueryJSON("custom/bank/balances", param, &coins))
		require.True(t, coins.IsEqual(processInitialBalance), coins)
	})

	t.Run("list", func(t *testing.T) {
		ps := make([]*process.Process, 0)
		lcdGet(t, "process/list", &ps)
		require.Len(t, ps, 1)
	})

	t.Run("hash", func(t *testing.T) {
		var hash hash.Hash
		lcdPost(t, "process/hash", msg, &hash)
		require.Equal(t, processHash, hash)
	})

	t.Run("delete", func(t *testing.T) {
		_, err := client.ProcessClient.Delete(context.Background(), &pb.DeleteProcessRequest{Hash: processHash})
		require.NoError(t, err)
	})

	t.Run("check ownership deletion", func(t *testing.T) {
		ownerships := make([]*ownership.Ownership, 0)
		lcdGet(t, "ownership/list", &ownerships)
		require.Len(t, ownerships, 2)
	})

	t.Run("check coins on process", func(t *testing.T) {
		var coins sdk.Coins
		param := bank.NewQueryBalanceParams(processAddress)
		require.NoError(t, cclient.QueryJSON("custom/bank/balances", param, &coins))
		require.True(t, coins.IsZero(), coins)
	})
}
