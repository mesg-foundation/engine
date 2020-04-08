package main

import (
	"context"
	"testing"

	"github.com/mesg-foundation/engine/ext/xos"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/acknowledgement"
	"github.com/mesg-foundation/engine/server/grpc/orchestrator"
	grpcorchestrator "github.com/mesg-foundation/engine/server/grpc/orchestrator"
	"github.com/mesg-foundation/engine/service"
	runnerrest "github.com/mesg-foundation/engine/x/runner/client/rest"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func testComplexService(t *testing.T) {
	var (
		testServiceComplexHash     hash.Hash
		testRunnerComplexHash      hash.Hash
		testInstanceComplexHash    hash.Hash
		testInstanceComplexEnvHash hash.Hash
		testServiceComplexStruct   *service.Service
		testInstanceComplexEnv     []string
		registerPayloadSignature   []byte
		err                        error
	)

	t.Run("create service", func(t *testing.T) {
		testComplexCreateServiceMsg.Owner = cliAddress
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
		require.NoError(t, cont.Build(testServiceComplexStruct))
	})

	t.Run("create msg, sign it and inject into env", func(t *testing.T) {
		value := grpcorchestrator.RunnerRegisterRequest{
			ServiceHash: testServiceComplexHash,
			EnvHash:     testInstanceComplexEnvHash,
		}

		registerPayloadSignature, err = signPayload(value)
		require.NoError(t, err)
	})

	t.Run("start runner", func(t *testing.T) {
		require.NoError(t, cont.Start(testServiceComplexStruct, testInstanceComplexHash, testRunnerComplexHash, testInstanceComplexEnvHash, testInstanceComplexEnv, registerPayloadSignature))
	})

	req := orchestrator.EventStreamRequest{}
	stream, err := client.EventClient.Stream(context.Background(), &req, grpc.PerRPCCredentials(&signCred{req}))
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
			case "service_ready", "read_env_ok", "read_env_override_ok", "access_volumes_ok", "access_volumes_from_ok", "resolve_dependence_ok":
				t.Logf("received event %s ", ev.Key)
			default:
				t.Fatalf("failed on event %s", ev.Key)
			}
		}
	})

	t.Run("delete", func(t *testing.T) {
		req := orchestrator.RunnerDeleteRequest{
			RunnerHash: testRunnerComplexHash,
		}
		_, err := client.RunnerClient.Delete(context.Background(), &req, grpc.PerRPCCredentials(&signCred{req}))
		require.NoError(t, err)

		require.NoError(t, cont.Stop(testServiceComplexStruct, testRunnerComplexHash))
	})
}
