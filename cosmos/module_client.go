package cosmos

import (
	"context"
	"fmt"

	executionpb "github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/ext/xos"
	"github.com/mesg-foundation/engine/ext/xstrings"
	"github.com/mesg-foundation/engine/hash"
	instancepb "github.com/mesg-foundation/engine/instance"
	ownershippb "github.com/mesg-foundation/engine/ownership"
	processpb "github.com/mesg-foundation/engine/process"
	"github.com/mesg-foundation/engine/protobuf/api"
	runnerpb "github.com/mesg-foundation/engine/runner"
	servicepb "github.com/mesg-foundation/engine/service"
	"github.com/mesg-foundation/engine/x/execution"
	"github.com/mesg-foundation/engine/x/instance"
	"github.com/mesg-foundation/engine/x/ownership"
	"github.com/mesg-foundation/engine/x/process"
	"github.com/mesg-foundation/engine/x/runner"
	"github.com/mesg-foundation/engine/x/service"
)

// ModuleClient handles all communication with every module.
type ModuleClient struct {
	*RPC
}

// NewModuleClient creates new module client.
func NewModuleClient(c *RPC) *ModuleClient {
	return &ModuleClient{RPC: c}
}

func sroutef(format string, args ...interface{}) string {
	return fmt.Sprintf("custom/"+format, args...)
}

// CreateService creates a new service from definition.
func (mc *ModuleClient) CreateService(req *api.CreateServiceRequest) (*servicepb.Service, error) {
	acc, err := mc.GetAccount()
	if err != nil {
		return nil, err
	}
	msg := service.MsgCreate{
		Owner:         acc.GetAddress(),
		Sid:           req.Sid,
		Name:          req.Name,
		Description:   req.Description,
		Configuration: req.Configuration,
		Tasks:         req.Tasks,
		Events:        req.Events,
		Dependencies:  req.Dependencies,
		Repository:    req.Repository,
		Source:        req.Source,
	}
	tx, err := mc.BuildAndBroadcastMsg(msg)
	if err != nil {
		return nil, err
	}
	return mc.GetService(tx.Data)
}

// GetService returns the service that matches given hash.
func (mc *ModuleClient) GetService(hash hash.Hash) (*servicepb.Service, error) {
	var out *servicepb.Service
	route := sroutef("%s/%s/%s", service.QuerierRoute, service.QueryGet, hash)
	return out, mc.QueryJSON(route, nil, &out)
}

// ListService returns all services.
func (mc *ModuleClient) ListService() ([]*servicepb.Service, error) {
	var out []*servicepb.Service
	route := sroutef("%s/%s", service.QuerierRoute, service.QueryList)
	return out, mc.QueryJSON(route, nil, &out)
}

// ExistService returns if a service already exists.
func (mc *ModuleClient) ExistService(hash hash.Hash) (bool, error) {
	var out bool
	route := sroutef("%s/%s/%s", service.QuerierRoute, service.QueryExist, hash)
	return out, mc.QueryJSON(route, nil, &out)
}

// CreateProcess creates a new process.
func (mc *ModuleClient) CreateProcess(req *api.CreateProcessRequest) (*processpb.Process, error) {
	acc, err := mc.GetAccount()
	if err != nil {
		return nil, err
	}
	msg := process.MsgCreate{
		Name:  req.Name,
		Edges: req.Edges,
		Nodes: req.Nodes,
		Owner: acc.GetAddress(),
	}
	tx, err := mc.BuildAndBroadcastMsg(msg)
	if err != nil {
		return nil, err
	}
	return mc.GetProcess(tx.Data)
}

// GetInstance returns the instance that matches given hash.
func (mc *ModuleClient) GetInstance(hash hash.Hash) (*instancepb.Instance, error) {
	var out *instancepb.Instance
	route := sroutef("%s/%s/%s", instance.QuerierRoute, instance.QueryGet, hash)
	return out, mc.QueryJSON(route, nil, &out)
}

// ListInstance returns all instances.
func (mc *ModuleClient) ListInstance(req *api.ListInstanceRequest) ([]*instancepb.Instance, error) {
	var out []*instancepb.Instance
	route := sroutef("%s/%s", instance.QuerierRoute, instance.QueryList)
	return out, mc.QueryJSON(route, req, &out)
}

// ListOwnership returns all ownerships.
func (mc *ModuleClient) ListOwnership() ([]*ownershippb.Ownership, error) {
	var out []*ownershippb.Ownership
	route := sroutef("%s/%s", ownership.QuerierRoute, ownership.QueryList)
	return out, mc.QueryJSON(route, nil, &out)
}

// DeleteProcess deletes the process by hash.
func (mc *ModuleClient) DeleteProcess(req *api.DeleteProcessRequest) error {
	acc, err := mc.GetAccount()
	if err != nil {
		return err
	}
	msg := process.MsgDelete{
		Hash:  req.Hash,
		Owner: acc.GetAddress(),
	}
	_, err = mc.BuildAndBroadcastMsg(msg)
	return err
}

// GetProcess returns the process that matches given hash.
func (mc *ModuleClient) GetProcess(hash hash.Hash) (*processpb.Process, error) {
	var out *processpb.Process
	route := sroutef("%s/%s/%s", process.QuerierRoute, process.QueryGet, hash.String())
	return out, mc.QueryJSON(route, nil, &out)
}

// ListProcess returns all processes.
func (mc *ModuleClient) ListProcess() ([]*processpb.Process, error) {
	var out []*processpb.Process
	route := sroutef("%s/%s", process.QuerierRoute, process.QueryList)
	return out, mc.QueryJSON(route, nil, &out)
}

// CreateExecution creates a new execution.
func (mc *ModuleClient) CreateExecution(req *api.CreateExecutionRequest) (*executionpb.Execution, error) {
	acc, err := mc.GetAccount()
	if err != nil {
		return nil, err
	}
	msg := execution.MsgCreate{
		Signer:       acc.GetAddress(),
		EventHash:    req.EventHash,
		ExecutorHash: req.ExecutorHash,
		Inputs:       req.Inputs,
		NodeKey:      req.NodeKey,
		ParentHash:   req.ParentHash,
		Price:        req.Price,
		ProcessHash:  req.ProcessHash,
		Tags:         req.Tags,
		TaskKey:      req.TaskKey,
	}
	tx, err := mc.BuildAndBroadcastMsg(msg)
	if err != nil {
		return nil, err
	}
	return mc.GetExecution(tx.Data)
}

// UpdateExecution updates a execution.
func (mc *ModuleClient) UpdateExecution(req *api.UpdateExecutionRequest) (*executionpb.Execution, error) {
	acc, err := mc.GetAccount()
	if err != nil {
		return nil, err
	}
	msg := execution.MsgUpdate{
		Executor: acc.GetAddress(),
		Hash:     req.Hash,
	}
	switch result := req.Result.(type) {
	case *api.UpdateExecutionRequest_Outputs:
		msg.Result = &execution.MsgUpdateOutputs{
			Outputs: result.Outputs,
		}
	case *api.UpdateExecutionRequest_Error:
		msg.Result = &execution.MsgUpdateError{
			Error: result.Error,
		}
	}
	tx, err := mc.BuildAndBroadcastMsg(msg)
	if err != nil {
		return nil, err
	}
	return mc.GetExecution(tx.Data)
}

// GetExecution returns the execution that matches given hash.
func (mc *ModuleClient) GetExecution(hash hash.Hash) (*executionpb.Execution, error) {
	var out *executionpb.Execution
	route := sroutef("%s/%s/%s", execution.QuerierRoute, execution.QueryGet, hash)
	return out, mc.QueryJSON(route, nil, &out)
}

// ListExecution returns all executions.
func (mc *ModuleClient) ListExecution() ([]*executionpb.Execution, error) {
	var out []*executionpb.Execution
	route := sroutef("%s/%s", execution.QuerierRoute, execution.QueryList)
	return out, mc.QueryJSON(route, nil, &out)
}

// StreamExecution returns execution that matches given hash.
func (mc *ModuleClient) StreamExecution(ctx context.Context, req *api.StreamExecutionRequest) (chan *executionpb.Execution, chan error, error) {
	if err := req.Filter.Validate(); err != nil {
		return nil, nil, err
	}

	subscriber := xstrings.RandASCIILetters(8)
	query := fmt.Sprintf("%s.%s EXISTS", execution.EventType, execution.AttributeKeyHash)
	eventStream, err := mc.Subscribe(ctx, subscriber, query, 0)
	if err != nil {
		return nil, nil, err
	}

	execC := make(chan *executionpb.Execution)
	errC := make(chan error)
	go func() {
	loop:
		for {
			select {
			case event := <-eventStream:
				attrHash := fmt.Sprintf("%s.%s", execution.EventType, execution.AttributeKeyHash)
				attrs := event.Events[attrHash]
				for _, attr := range attrs {
					hash, err := hash.Decode(attr)
					if err != nil {
						errC <- err
						continue
					}
					exec, err := mc.GetExecution(hash)
					if err != nil {
						errC <- err
						continue
					}
					if req.Filter.Match(exec) {
						execC <- exec
					}
				}
			case <-ctx.Done():
				break loop
			}
		}
		if err := mc.Unsubscribe(context.Background(), subscriber, query); err != nil {
			errC <- err
		}
		close(execC)
		close(errC)
	}()
	return execC, errC, nil
}

// CreateRunner creates a new runner.
func (mc *ModuleClient) CreateRunner(req *api.CreateRunnerRequest) (*runnerpb.Runner, error) {
	s, err := mc.GetService(req.ServiceHash)
	if err != nil {
		return nil, err
	}
	envHash := hash.Dump(xos.EnvMergeSlices(s.Configuration.Env, req.Env))
	acc, err := mc.GetAccount()
	if err != nil {
		return nil, err
	}
	msg := runner.MsgCreate{
		Owner:       acc.GetAddress(),
		ServiceHash: req.ServiceHash,
		EnvHash:     envHash,
	}
	tx, err := mc.BuildAndBroadcastMsg(msg)
	if err != nil {
		return nil, err
	}
	return mc.GetRunner(tx.Data)
}

// DeleteRunner deletes an existing runner.
func (mc *ModuleClient) DeleteRunner(req *api.DeleteRunnerRequest) error {
	acc, err := mc.GetAccount()
	if err != nil {
		return err
	}
	msg := runner.MsgDelete{
		Owner: acc.GetAddress(),
		Hash:  req.Hash,
	}
	_, err = mc.BuildAndBroadcastMsg(msg)
	return err
}

// GetRunner returns the runner that matches given hash.
func (mc *ModuleClient) GetRunner(hash hash.Hash) (*runnerpb.Runner, error) {
	var out *runnerpb.Runner
	route := sroutef("%s/%s/%s", runner.QuerierRoute, runner.QueryGet, hash)
	return out, mc.QueryJSON(route, nil, &out)
}

// FilterRunner to apply while listing runners.
type FilterRunner struct {
	Owner        string
	InstanceHash hash.Hash
}

// ListRunner returns all runners.
func (mc *ModuleClient) ListRunner(f *FilterRunner) ([]*runnerpb.Runner, error) {
	var rs []*runnerpb.Runner
	route := sroutef("%s/%s", runner.QuerierRoute, runner.QueryList)
	if err := mc.QueryJSON(route, nil, &rs); err != nil {
		return nil, err
	}

	// no filter, returns
	if f == nil {
		return rs, nil
	}

	// filter results
	out := make([]*runnerpb.Runner, 0)
	for _, r := range rs {
		if (f.Owner == "" || r.Owner == f.Owner) &&
			(f.InstanceHash.IsZero() || r.InstanceHash.Equal(f.InstanceHash)) {
			out = append(out, r)
		}
	}
	return out, nil
}
