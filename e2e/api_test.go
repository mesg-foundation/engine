package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/mesg-foundation/engine/hash"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var passmd = metadata.Pairs(
	"credential_username", "dev",
	"credential_passphrase", "pass",
)

type apiclient struct {
	pb.ServiceClient
	pb.EventClient
	pb.ExecutionClient
	pb.AccountClient
	pb.ProcessClient
	pb.InstanceClient
	pb.OwnershipClient
}

var client apiclient

func readCreateServiceRequest(filename string) *pb.CreateServiceRequest {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	var req pb.CreateServiceRequest
	if err = json.Unmarshal(b, &req); err != nil {
		log.Fatal(err)
	}
	return &req
}

var (
	testServiceHash  hash.Hash
	testInstanceHash hash.Hash
)

func TestAPI(t *testing.T) {
	conn, err := grpc.DialContext(context.Background(), "localhost:50052", grpc.WithInsecure())
	require.NoError(t, err)

	client = apiclient{
		pb.NewServiceClient(conn),
		pb.NewEventClient(conn),
		pb.NewExecutionClient(conn),
		pb.NewAccountClient(conn),
		pb.NewProcessClient(conn),
		pb.NewInstanceClient(conn),
		pb.NewOwnershipClient(conn),
	}

	// ping server to test connection
	_, err = client.ServiceClient.List(context.Background(), &pb.ListServiceRequest{})
	require.NoError(t, err)

	time.Sleep(1 * time.Second)

	t.Run("account", testAccount)
	t.Run("service", testService)
	t.Run("instance", testInstance)
	t.Run("execution", testExecution)
}
