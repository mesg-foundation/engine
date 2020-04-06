package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	types "github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/server/grpc/orchestrator"
	"github.com/mesg-foundation/engine/server/grpc/runner"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	// env variables for configure mesg client.
	envMesgEndpoint          = "MESG_ENDPOINT"
	envMesgServiceHash       = "MESG_SERVICE_HASH"
	envMesgEnvHash           = "MESG_ENV_HASH"
	envMesgRegisterSignature = "MESG_REGISTER_SIGNATURE"
)

// register
func register() (string, error) {
	endpoint := os.Getenv(envMesgEndpoint)
	if endpoint == "" {
		return "", fmt.Errorf("env %q is empty", envMesgEndpoint)
	}

	conn, err := grpc.DialContext(context.Background(), endpoint, grpc.WithInsecure())
	if err != nil {
		return "", err
	}
	client := orchestrator.NewRunnerClient(conn)

	serviceHash, err := hash.Decode(os.Getenv(envMesgServiceHash))
	if err != nil {
		return "", err
	}
	envHash, err := hash.Decode(os.Getenv(envMesgEnvHash))
	if err != nil {
		return "", err
	}
	signature := os.Getenv(envMesgRegisterSignature)

	resp, err := client.Register(context.Background(), &orchestrator.RunnerRegisterRequest{
		ServiceHash: serviceHash,
		EnvHash:     envHash,
	}, grpc.PerRPCCredentials(&signCred{signature}))
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

func newClient(token string) (runner.RunnerClient, error) {
	endpoint := os.Getenv(envMesgEndpoint)
	if endpoint == "" {
		return nil, fmt.Errorf("env %q is empty", envMesgEndpoint)
	}

	// runner client
	dialoptions := []grpc.DialOption{
		// Keep alive prevents Docker network to drop TCP idle connections after 15 minutes.
		// See: https://forum.mesg.com/t/solution-summary-for-docker-dropping-connections-after-15-min/246
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time: 5 * time.Minute, // 5 minutes is the minimun time of gRPC enforcement policy.
		}),
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(&tokenCred{token}),
	}
	conn, err := grpc.DialContext(context.Background(), endpoint, dialoptions...)
	if err != nil {
		return nil, fmt.Errorf("connection error: %s", err)
	}

	return runner.NewRunnerClient(conn), nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	token, err := register()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("registered with token %s\n", token)

	client, err := newClient(token)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("connected to %s\n", os.Getenv(envMesgEndpoint))

	stream, err := client.Execution(context.Background(), &runner.ExecutionRequest{})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("created execution stream\n")

	if _, err := client.Event(context.Background(), &runner.EventRequest{
		Key: "test_service_ready",
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

		go processExec(client, token, exec)
	}
}

func processExec(client runner.RunnerClient, token string, exec *execution.Execution) {
	req := &runner.ResultRequest{
		ExecutionHash: exec.Hash,
	}

	if exec.TaskKey == "task1" || exec.TaskKey == "task2" {
		req.Result = &runner.ResultRequest_Outputs{
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
		var fields = map[string]*types.Value{
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
		}
		if exec.Inputs.Fields["msg"].GetStructValue().Fields["array"].GetListValue() != nil {
			fields["array"] = &types.Value{
				Kind: &types.Value_ListValue{
					ListValue: exec.Inputs.Fields["msg"].GetStructValue().Fields["array"].GetListValue(),
				},
			}
		}
		req.Result = &runner.ResultRequest_Outputs{
			Outputs: &types.Struct{
				Fields: map[string]*types.Value{
					"msg": {
						Kind: &types.Value_StructValue{
							StructValue: &types.Struct{
								Fields: fields,
							},
						},
					},
				},
			},
		}
	} else {
		log.Fatalf("Taskkey %q not implemented", exec.TaskKey)
	}

	if _, err := client.Result(context.Background(), req); err != nil {
		log.Fatal(err)
	}
	log.Printf("execution result submitted\n")

	if _, err := client.Event(context.Background(), &runner.EventRequest{
		Key: "event_after_task",
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

// tokenCred is a structure that manage a token.
type tokenCred struct {
	token string
}

// GetRequestMetadata returns the metadata for the request.
func (c *tokenCred) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		runner.CredentialToken: c.token,
	}, nil
}

// RequireTransportSecurity tells if the transport should be secured.
func (c *tokenCred) RequireTransportSecurity() bool {
	return false
}

// signCred is a structure that manage a token.
type signCred struct {
	signature string
}

// GetRequestMetadata returns the metadata for the request.
func (c *signCred) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		orchestrator.RequestSignature: c.signature,
	}, nil
}

// RequireTransportSecurity tells if the transport should be secured.
func (c *signCred) RequireTransportSecurity() bool {
	return false
}
