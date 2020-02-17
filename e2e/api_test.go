package main

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/mesg-foundation/engine/app"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

type apiclient struct {
	pb.ServiceClient
	pb.EventClient
	pb.ExecutionClient
	pb.ProcessClient
	pb.InstanceClient
	pb.OwnershipClient
	pb.RunnerClient
}

var (
	client apiclient
	cdc    = app.MakeCodec()
)

const (
	lcdEndpoint        = "http://127.0.0.1:1317/"
	lcdPostContentType = "application/json"
)

func lcdGet(t *testing.T, path string, ptr interface{}) {
	resp, err := http.Get(lcdEndpoint + path)
	require.NoError(t, err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	cosResp := rest.ResponseWithHeight{}
	require.NoError(t, cdc.UnmarshalJSON(body, &cosResp))
	if len(cosResp.Result) > 0 {
		require.NoError(t, cdc.UnmarshalJSON(cosResp.Result, ptr))
	}
}

func lcdPost(t *testing.T, path string, req interface{}, ptr interface{}) {
	reqBody, err := cdc.MarshalJSON(req)
	require.NoError(t, err)
	resp, err := http.Post(lcdEndpoint+path, lcdPostContentType, bytes.NewReader(reqBody))
	require.NoError(t, err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	cosResp := rest.ResponseWithHeight{}
	require.NoError(t, cdc.UnmarshalJSON(body, &cosResp))
	if len(cosResp.Result) > 0 {
		require.NoError(t, cdc.UnmarshalJSON(cosResp.Result, ptr))
	}
}

func TestAPI(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	conn, err := grpc.DialContext(context.Background(), "localhost:50052", grpc.WithInsecure())
	require.NoError(t, err)

	client = apiclient{
		pb.NewServiceClient(conn),
		pb.NewEventClient(conn),
		pb.NewExecutionClient(conn),
		pb.NewProcessClient(conn),
		pb.NewInstanceClient(conn),
		pb.NewOwnershipClient(conn),
		pb.NewRunnerClient(conn),
	}

	// ping server to test connection
	_, err = client.ServiceClient.List(context.Background(), &pb.ListServiceRequest{})
	require.NoError(t, err)

	// basic tests
	t.Run("service", testService)
	t.Run("runner", testRunner)
	t.Run("process", testProcess)
	t.Run("instance", testInstance)
	t.Run("event", testEvent)
	t.Run("execution", testExecution)
	t.Run("orchestrator", testOrchestrator)
	t.Run("runner/delete", testDeleteRunner)
	t.Run("complex-service", testComplexService)
}
