package main

import (
	"context"
	"testing"

	"github.com/mesg-foundation/engine/ext/xos"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/acknowledgement"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/runner/builder"
	grpcrunner "github.com/mesg-foundation/engine/server/grpc/runner"
	"github.com/mesg-foundation/engine/service"
	runnermodule "github.com/mesg-foundation/engine/x/runner"
	runnerrest "github.com/mesg-foundation/engine/x/runner/client/rest"
	"github.com/stretchr/testify/require"
)

func testComplexService(t *testing.T) {
	var (
		testServiceComplexHash      hash.Hash
		testRunnerComplexHash       hash.Hash
		testInstanceComplexHash     hash.Hash
		testInstanceComplexEnvHash  hash.Hash
		testServiceComplexStruct    *service.Service
		testServiceComplexImageHash string
		testInstanceComplexEnv      []string
		err                         error
	)

	t.Run("create service", func(t *testing.T) {
		testComplexCreateServiceMsg.Owner = engineAddress
		testServiceComplexHash, err = lcd.BroadcastMsg(testComplexCreateServiceMsg)
		require.NoError(t, err)
	})

	t.Run("get", func(t *testing.T) {
		require.NoError(t, lcd.Get("service/get/"+testServiceComplexHash.String(), &testServiceComplexStruct))
		require.Equal(t, testServiceComplexHash, testServiceComplexStruct.Hash)
	})
	testInstanceComplexEnv = xos.EnvMergeSlices(testServiceComplexStruct.Configuration.Env, []string{"ENVB=is_override"})

	t.Run("get runner hashes", func(t *testing.T) {
		var res runnerrest.HashResponse
		err := lcd.Post("runner/hash", &runnerrest.HashRequest{
			ServiceHash: testServiceComplexHash,
			Address:     engineAddress.String(),
			Env:         testInstanceComplexEnv,
		}, &res)
		require.NoError(t, err)
		testRunnerComplexHash = res.RunnerHash
		testInstanceComplexHash = res.InstanceHash
		testInstanceComplexEnvHash = res.EnvHash
	})

	t.Run("build service image", func(t *testing.T) {
		var err error
		testServiceComplexImageHash, err = builder.Build(cont, testServiceComplexStruct, ipfsEndpoint)
		require.NoError(t, err)
	})

	t.Run("create msg, sign it and inject into env", func(t *testing.T) {
		value := grpcrunner.RegisterRequestPayload_Value{
			ServiceHash: testServiceComplexHash,
			EnvHash:     testInstanceComplexEnvHash,
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

		testInstanceComplexEnv = append(testInstanceComplexEnv, "MESG_REGISTER_PAYLOAD="+string(encodedPayload))
	})

	t.Run("start runner", func(t *testing.T) {
		require.NoError(t, builder.Start(cont, testServiceComplexStruct, testInstanceComplexHash, testRunnerComplexHash, testServiceComplexImageHash, testInstanceComplexEnv, engineName, enginePort))
	})

	stream, err := client.EventClient.Stream(context.Background(), &pb.StreamEventRequest{})
	require.NoError(t, err)
	acknowledgement.WaitForStreamToBeReady(stream)

	t.Run("check events", func(t *testing.T) {
		okEventsNo := 6
		for i := 0; i < okEventsNo; {
			ev, err := stream.Recv()
			require.NoError(t, err)

			if !ev.InstanceHash.Equal(testInstanceComplexHash) {
				continue
			}
			i++

			switch ev.Key {
			case "test_service_ready", "read_env_ok", "read_env_override_ok", "access_volumes_ok", "access_volumes_from_ok", "resolve_dependence_ok":
				t.Logf("received event %s ", ev.Key)
			default:
				t.Fatalf("failed on event %s", ev.Key)
			}
		}
	})

	t.Run("delete", func(t *testing.T) {
		t.Skip("FIXME: this test timeout on CIRCLE CI. works well on local computer")
		msg := runnermodule.MsgDelete{
			Owner: engineAddress,
			Hash:  testRunnerComplexHash,
		}

		_, err := lcd.BroadcastMsg(msg)
		require.NoError(t, err)

		require.NoError(t, builder.Stop(cont, testRunnerComplexHash, testServiceComplexStruct.Dependencies))
	})
}
