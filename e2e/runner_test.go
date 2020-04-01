package main

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/ext/xos"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/acknowledgement"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/runner"
	"github.com/mesg-foundation/engine/runner/builder"
	grpcrunner "github.com/mesg-foundation/engine/server/grpc/runner"
	runnermodule "github.com/mesg-foundation/engine/x/runner"
	runnerrest "github.com/mesg-foundation/engine/x/runner/client/rest"
	"github.com/stretchr/testify/require"
)

var (
	testRunnerHash       hash.Hash
	testInstanceEnvHash  hash.Hash
	testRunnerAddress    sdk.AccAddress
	testServiceImageHash string
)

func testRunner(t *testing.T) {
	var (
		testInstanceEnv = xos.EnvMergeSlices(testServiceStruct.Configuration.Env, []string{"BAR=3", "REQUIRED=4"})
	)
	t.Run("hash", func(t *testing.T) {
		var res runnerrest.HashResponse
		err := lcd.Post("runner/hash", &runnerrest.HashRequest{
			ServiceHash: testServiceHash,
			Address:     engineAddress.String(),
			Env:         testInstanceEnv,
		}, &res)
		require.NoError(t, err)
		testRunnerHash = res.RunnerHash
		testInstanceHash = res.InstanceHash
		testInstanceEnvHash = res.EnvHash
	})

	t.Run("build service image", func(t *testing.T) {
		var err error
		testServiceImageHash, err = builder.Build(cont, testServiceStruct, ipfsEndpoint)
		require.NoError(t, err)
	})

	t.Run("create msg, sign it and inject into env", func(t *testing.T) {
		value := grpcrunner.RegisterRequestPayload_Value{
			ServiceHash: testServiceHash,
			EnvHash:     testInstanceEnvHash,
		}

		encodedValue, err := cdc.MarshalJSON(value)
		require.NoError(t, err)

		signature, _, err := kb.Sign(engineAccountName, engineAccountPassword, encodedValue)
		require.NoError(t, err)

		payload := grpcrunner.RegisterRequestPayload{
			Signature: signature,
			Value:     value,
		}
		encodedPayload, err := cdc.MarshalJSON(payload)
		require.NoError(t, err)

		testInstanceEnv = append(testInstanceEnv, "MESG_REGISTER_PAYLOAD="+string(encodedPayload))
	})

	t.Run("wait for service to be ready", func(t *testing.T) {
		stream, err := client.EventClient.Stream(context.Background(), &pb.StreamEventRequest{
			Filter: &pb.StreamEventRequest_Filter{
				Key: "test_service_ready",
			},
		})
		require.NoError(t, err)
		acknowledgement.WaitForStreamToBeReady(stream)

		t.Run("start", func(t *testing.T) {
			require.NoError(t, builder.Start(cont, testServiceStruct, testInstanceHash, testRunnerHash, testServiceImageHash, testInstanceEnv, engineName, enginePort))
		})

		// wait for service to be ready
		_, err = stream.Recv()
		require.NoError(t, err)
	})

	t.Run("get", func(t *testing.T) {
		var run *runner.Runner
		require.NoError(t, lcd.Get("runner/get/"+testRunnerHash.String(), &run))
		require.Equal(t, testRunnerHash, run.Hash)
		testRunnerAddress = run.Address
	})

	t.Run("list", func(t *testing.T) {
		rs := make([]*runner.Runner, 0)
		require.NoError(t, lcd.Get("runner/list", &rs))
		require.Len(t, rs, 1)
		require.Equal(t, testInstanceHash, rs[0].InstanceHash)
		require.Equal(t, testRunnerHash, rs[0].Hash)
	})
}

func testDeleteRunner(t *testing.T) {
	msg := runnermodule.MsgDelete{
		Owner: engineAddress,
		Hash:  testRunnerHash,
	}
	_, err := lcd.BroadcastMsg(msg)
	require.NoError(t, err)

	require.NoError(t, builder.Stop(cont, testRunnerHash, testServiceStruct.Dependencies))

	t.Run("check deletion", func(t *testing.T) {
		rs := make([]*runner.Runner, 0)
		require.NoError(t, lcd.Get("runner/list", &rs))
		require.Len(t, rs, 0)
	})

	t.Run("check coins on runner", func(t *testing.T) {
		var coins sdk.Coins
		require.NoError(t, lcd.Get("bank/balances/"+testRunnerAddress.String(), &coins))
		require.True(t, coins.IsZero(), coins)
	})
}
