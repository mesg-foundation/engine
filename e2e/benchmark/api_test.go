package main

import (
	"context"
	"testing"

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

var client apiclient

func initClient() error {
	conn, err := grpc.DialContext(context.Background(), "localhost:50052", grpc.WithInsecure())
	if err != nil {
		return err
	}
	client = apiclient{
		pb.NewServiceClient(conn),
		pb.NewEventClient(conn),
		pb.NewExecutionClient(conn),
		pb.NewProcessClient(conn),
		pb.NewInstanceClient(conn),
		pb.NewOwnershipClient(conn),
		pb.NewRunnerClient(conn),
	}
	return nil
}

func BenchmarkAPI(b *testing.B) {
	if testing.Short() {
		b.Skip()
	}

	require.NoError(b, initClient())

	// ping server to test connection
	_, err := client.ServiceClient.List(context.Background(), &pb.ListServiceRequest{})
	require.NoError(b, err)

	// benchmark tests
	b.Run("execution", benchmarkExecution)
}
