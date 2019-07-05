package client

import (
	"context"

	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/protobuf/types"
)

// TaskFn is a task function handler prototype.
type TaskFn func(inputs []byte) ([]byte, error)

// TaskRunner handles running task in a loop.
type TaskRunner struct {
	client *Client
	defs   map[string]TaskFn
}

// Add sets new task handler. It overwrites previous one.
func (r *TaskRunner) Add(name string, fn TaskFn) {
	r.defs[name] = fn
}

// Run recives executions and runs dedicated task for them.
func (r *TaskRunner) Run() error {
	stream, err := r.client.ExecutionClient.Stream(context.Background(), &pb.StreamExecutionRequest{
		Filter: &pb.StreamExecutionRequest_Filter{
			Statuses:     []types.Status{types.Status_InProgress},
			InstanceHash: r.client.InstanceHash,
		},
	})
	if err != nil {
		return err
	}

	for {
		exec, err := stream.Recv()
		if err != nil {
			return err
		}

		if _, ok := r.defs[exec.TaskKey]; !ok {
			return err
		}

		output, err := r.defs[exec.TaskKey]([]byte(exec.Inputs))
		req := &pb.UpdateExecutionRequest{
			Hash: exec.Hash,
			Result: &pb.UpdateExecutionRequest_Outputs{
				Outputs: string(output),
			},
		}
		if err != nil {
			req.Result = &pb.UpdateExecutionRequest_Error{
				Error: err.Error(),
			}
		}
		if _, err := r.client.ExecutionClient.Update(context.Background(), req); err != nil {
			return err
		}
	}
}
