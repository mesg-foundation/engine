package main

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/ownership"
	"github.com/mesg-foundation/engine/process"
	processmodule "github.com/mesg-foundation/engine/x/process"
	processrest "github.com/mesg-foundation/engine/x/process/client/rest"
	"github.com/stretchr/testify/require"
)

func testProcess(t *testing.T) {
	var (
		err            error
		processHash    hash.Hash
		processAddress sdk.AccAddress
		msg            = &processmodule.MsgCreate{
			Owner: cliAddress,
			Name:  "test-process",
			Nodes: []*process.Process_Node{
				{
					Key: "n0",
					Type: &process.Process_Node_Event_{
						Event: &process.Process_Node_Event{
							InstanceHash: testInstanceHash,
							EventKey:     "service_ready",
						},
					},
				},
				{
					Key: "n1",
					Type: &process.Process_Node_Task_{
						Task: &process.Process_Node_Task{
							InstanceHash: testInstanceHash,
							TaskKey:      "service_ready",
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
		processHash, err = lcd.BroadcastMsg(msg)
		require.NoError(t, err)
	})

	t.Run("get", func(t *testing.T) {
		var p *process.Process
		require.NoError(t, lcd.Get("process/get/"+processHash.String(), &p))
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
		require.NoError(t, lcd.Get("ownership/list", &ownerships))
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
		require.NoError(t, lcd.Get("bank/balances/"+processAddress.String(), &coins))
		require.True(t, coins.IsEqual(processInitialBalance), coins)
	})

	t.Run("list", func(t *testing.T) {
		ps := make([]*process.Process, 0)
		require.NoError(t, lcd.Get("process/list", &ps))
		require.Len(t, ps, 1)
	})

	t.Run("hash", func(t *testing.T) {
		msg := processrest.HashRequest{
			Name:  msg.Name,
			Nodes: msg.Nodes,
			Edges: msg.Edges,
		}
		var hash hash.Hash
		require.NoError(t, lcd.Post("process/hash", msg, &hash))
		require.Equal(t, processHash, hash)
	})

	t.Run("delete", func(t *testing.T) {
		_, err := lcd.BroadcastMsg(processmodule.MsgDelete{
			Owner: cliAddress,
			Hash:  processHash,
		})
		require.NoError(t, err)
	})

	t.Run("check ownership deletion", func(t *testing.T) {
		ownerships := make([]*ownership.Ownership, 0)
		require.NoError(t, lcd.Get("ownership/list", &ownerships))
		require.Len(t, ownerships, 2)
	})

	t.Run("check coins on process", func(t *testing.T) {
		var coins sdk.Coins
		require.NoError(t, lcd.Get("bank/balances/"+processAddress.String(), &coins))
		require.True(t, coins.IsZero(), coins)
	})
}
