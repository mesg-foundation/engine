package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	types "github.com/mesg-foundation/engine/protobuf/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	// env variables for configure mesg client.
	envMesgEndpoint     = "MESG_ENDPOINT"
	envMesgInstanceHash = "MESG_INSTANCE_HASH"
	envMesgRunnerHash   = "MESG_RUNNER_HASH"
)

// Client is a client to connect to all mesg exposed API.
type Client struct {
	// all clients registered by mesg server.
	pb.EventClient
	pb.ExecutionClient

	// instance hash that could be used in api calls.
	InstanceHash hash.Hash

	// runner hash that could be used in api calls.
	RunnerHash hash.Hash
}

// New creates a new client from env variables supplied by mesg engine.
func New() (*Client, error) {
	endpoint := os.Getenv(envMesgEndpoint)
	if endpoint == "" {
		return nil, fmt.Errorf("client: mesg server address env(%s) is empty", envMesgEndpoint)
	}

	instanceHash, err := hash.Decode(os.Getenv(envMesgInstanceHash))
	if err != nil {
		return nil, fmt.Errorf("client: error with mesg's instance hash env(%s): %s", envMesgInstanceHash, err.Error())
	}

	runnerHash, err := hash.Decode(os.Getenv(envMesgRunnerHash))
	if err != nil {
		return nil, fmt.Errorf("client: error with mesg's runner hash env(%s): %s", envMesgRunnerHash, err.Error())
	}

	dialoptions := []grpc.DialOption{
		// Keep alive prevents Docker network to drop TCP idle connections after 15 minutes.
		// See: https://forum.mesg.com/t/solution-summary-for-docker-dropping-connections-after-15-min/246
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time: 5 * time.Minute, // 5 minutes is the minimun time of gRPC enforcement policy.
		}),
		grpc.WithTimeout(10 * time.Second),
		grpc.WithInsecure(),
	}

	conn, err := grpc.DialContext(context.Background(), endpoint, dialoptions...)
	if err != nil {
		return nil, fmt.Errorf("client: connection error: %s", err)
	}

	return &Client{
		ExecutionClient: pb.NewExecutionClient(conn),
		EventClient:     pb.NewEventClient(conn),
		InstanceHash:    instanceHash,
		RunnerHash:      runnerHash,
	}, nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	client, err := New()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("connected to %s\n", os.Getenv(envMesgEndpoint))

	stream, err := client.ExecutionClient.Stream(context.Background(), &pb.StreamExecutionRequest{
		Filter: &pb.StreamExecutionRequest_Filter{
			Statuses:     []execution.Status{execution.Status_InProgress},
			ExecutorHash: client.RunnerHash,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("created execution stream\n")

	if _, err := client.EventClient.Create(context.Background(), &pb.CreateEventRequest{
		InstanceHash: client.InstanceHash,
		Key:          "test_service_ready",
	}); err != nil {
		log.Fatal(err)
	}
	log.Printf("emitted test_service_ready event\n")

	for {
		exec, err := stream.Recv()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("received execution %s %s\n", exec.TaskKey, exec.Hash)

		req := &pb.UpdateExecutionRequest{
			Hash: exec.Hash,
		}

		if exec.TaskKey == "task1" || exec.TaskKey == "task2" {
			req.Result = &pb.UpdateExecutionRequest_Outputs{
				Outputs: &types.Struct{
					Fields: map[string]*types.Value{
						"msg": {
							Kind: &types.Value_StringValue{
								StringValue: exec.Inputs.Fields["msg"].GetStringValue(),
							},
						},
						"timestamp": {
							Kind: &types.Value_NumberValue{
								NumberValue: float64(time.Now().Unix()),
							},
						},
					},
				},
			}
		} else if exec.TaskKey == "task_complex" {
			req.Result = &pb.UpdateExecutionRequest_Outputs{
				Outputs: &types.Struct{
					Fields: map[string]*types.Value{
						"msg": {
							Kind: &types.Value_StructValue{
								StructValue: &types.Struct{
									Fields: map[string]*types.Value{
										"msg": {
											Kind: &types.Value_StringValue{
												StringValue: exec.Inputs.Fields["msg"].GetStructValue().Fields["msg"].GetStringValue(),
											},
										},
										"timestamp": {
											Kind: &types.Value_NumberValue{
												NumberValue: float64(time.Now().Unix()),
											},
										},
										"array": {
											Kind: &types.Value_ListValue{
												ListValue: exec.Inputs.Fields["msg"].GetStructValue().Fields["array"].GetListValue(),
											},
										},
									},
								},
							},
						},
					},
				},
			}
		} else {
			log.Fatalf("Taskkey %q not implemented", exec.TaskKey)
		}

		if _, err := client.ExecutionClient.Update(context.Background(), req); err != nil {
			log.Fatal(err)
		}
		log.Printf("execution result submitted\n")

		if _, err := client.EventClient.Create(context.Background(), &pb.CreateEventRequest{
			InstanceHash: client.InstanceHash,
			Key:          "event_after_task",
			Data: &types.Struct{
				Fields: map[string]*types.Value{
					"task_key": {
						Kind: &types.Value_StringValue{
							StringValue: exec.TaskKey,
						},
					},
					"timestamp": {
						Kind: &types.Value_NumberValue{
							NumberValue: float64(time.Now().Unix()),
						},
					},
				},
			},
		}); err != nil {
			log.Fatal(err)
		}
		log.Printf("emitted event_after_task event\n")
	}
}
