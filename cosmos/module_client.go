package cosmos

import (
	"context"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	executionpb "github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/ext/xos"
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
	*Client
}

// NewModuleClient creates new module client.
func NewModuleClient(c *Client) *ModuleClient {
	return &ModuleClient{Client: c}
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
	msg := service.NewMsgCreateService(acc.GetAddress(), req)
	tx, err := mc.BuildAndBroadcastMsg(msg)
	if err != nil {
		return nil, err
	}
	return mc.GetService(tx.Data)
}

// GetService returns the service that matches given hash.
func (mc *ModuleClient) GetService(hash sdk.AccAddress) (*servicepb.Service, error) {
	var out *servicepb.Service
	route := sroutef("%s/%s/%s", service.QuerierRoute, service.QueryGetService, hash)
	return out, mc.QueryJSON(route, nil, &out)
}

// ListService returns all services.
func (mc *ModuleClient) ListService() ([]*servicepb.Service, error) {
	var out []*servicepb.Service
	route := sroutef("%s/%s", service.QuerierRoute, service.QueryListService)
	return out, mc.QueryJSON(route, nil, &out)
}

// ExistService returns if a service already exists.
func (mc *ModuleClient) ExistService(hash sdk.AccAddress) (bool, error) {
	var out bool
	route := sroutef("%s/%s/%s", service.QuerierRoute, service.QueryExistService, hash)
	return out, mc.QueryJSON(route, nil, &out)
}

// HashService returns the calculate hash of a service.
func (mc *ModuleClient) HashService(req *api.CreateServiceRequest) (sdk.AccAddress, error) {
	var out sdk.AccAddress
	route := sroutef("%s/%s", service.QuerierRoute, service.QueryHashService)
	return out, mc.QueryJSON(route, req, &out)
}

// CreateProcess creates a new process.
func (mc *ModuleClient) CreateProcess(req *api.CreateProcessRequest) (*processpb.Process, error) {
	acc, err := mc.GetAccount()
	if err != nil {
		return nil, err
	}
	msg := process.NewMsgCreateProcess(acc.GetAddress(), req)
	tx, err := mc.BuildAndBroadcastMsg(msg)
	if err != nil {
		return nil, err
	}
	return mc.GetProcess(tx.Data)
}

// GetInstance returns the instance that matches given hash.
func (mc *ModuleClient) GetInstance(hash sdk.AccAddress) (*instancepb.Instance, error) {
	var out *instancepb.Instance
	route := sroutef("%s/%s/%s", instance.QuerierRoute, instance.QueryGetInstance, hash)
	return out, mc.QueryJSON(route, nil, &out)
}

// ListInstance returns all instances.
func (mc *ModuleClient) ListInstance(req *api.ListInstanceRequest) ([]*instancepb.Instance, error) {
	var out []*instancepb.Instance
	route := sroutef("%s/%s", instance.QuerierRoute, instance.QueryListInstances)
	return out, mc.QueryJSON(route, req, &out)
}

// ListOwnership returns all ownerships.
func (mc *ModuleClient) ListOwnership() ([]*ownershippb.Ownership, error) {
	var out []*ownershippb.Ownership
	route := sroutef("%s/%s", ownership.QuerierRoute, ownership.QueryListOwnerships)
	return out, mc.QueryJSON(route, nil, &out)
}

// DeleteProcess deletes the process by hash.
func (mc *ModuleClient) DeleteProcess(req *api.DeleteProcessRequest) error {
	acc, err := mc.GetAccount()
	if err != nil {
		return err
	}
	msg := process.NewMsgDeleteProcess(acc.GetAddress(), req)
	_, err = mc.BuildAndBroadcastMsg(msg)
	return err
}

// GetProcess returns the process that matches given hash.
func (mc *ModuleClient) GetProcess(hash sdk.AccAddress) (*processpb.Process, error) {
	var out *processpb.Process
	route := sroutef("%s/%s/%s", process.QuerierRoute, process.QueryGetProcess, hash.String())
	return out, mc.QueryJSON(route, nil, &out)
}

// ListProcess returns all processes.
func (mc *ModuleClient) ListProcess() ([]*processpb.Process, error) {
	var out []*processpb.Process
	route := sroutef("%s/%s", process.QuerierRoute, process.QueryListProcesses)
	return out, mc.QueryJSON(route, nil, &out)
}

// CreateExecution creates a new execution.
func (mc *ModuleClient) CreateExecution(req *api.CreateExecutionRequest) (*executionpb.Execution, error) {
	acc, err := mc.GetAccount()
	if err != nil {
		return nil, err
	}
	msg := execution.NewMsgCreateExecution(req, acc.GetAddress())
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
	msg := execution.NewMsgUpdateExecution(req, acc.GetAddress())
	tx, err := mc.BuildAndBroadcastMsg(msg)
	if err != nil {
		return nil, err
	}
	return mc.GetExecution(tx.Data)
}

// GetExecution returns the execution that matches given hash.
func (mc *ModuleClient) GetExecution(hash sdk.AccAddress) (*executionpb.Execution, error) {
	var out *executionpb.Execution
	route := sroutef("%s/%s/%s", execution.QuerierRoute, execution.QueryGetExecution, hash)
	return out, mc.QueryJSON(route, nil, &out)
}

// ListExecution returns all executions.
func (mc *ModuleClient) ListExecution() ([]*executionpb.Execution, error) {
	var out []*executionpb.Execution
	route := sroutef("%s/%s", execution.QuerierRoute, execution.QueryListExecution)
	return out, mc.QueryJSON(route, nil, &out)
}

// StreamExecution returns execution that matches given hash.
func (mc *ModuleClient) StreamExecution(ctx context.Context, req *api.StreamExecutionRequest) (chan *executionpb.Execution, chan error, error) {
	if err := req.Filter.Validate(); err != nil {
		return nil, nil, err
	}

	stream, serrC, err := mc.Stream(ctx, EventModuleQuery(execution.ModuleName))
	if err != nil {
		return nil, nil, err
	}

	execC := make(chan *executionpb.Execution)
	errC := make(chan error)
	go func() {
	loop:
		for {
			select {
			case hash := <-stream:
				exec, err := mc.GetExecution(hash)
				if err != nil {
					errC <- err
					break
				}
				if req.Filter.Match(exec) {
					execC <- exec
				}
			case err := <-serrC:
				errC <- err
			case <-ctx.Done():
				break loop
			}
		}
		close(errC)
		close(execC)
	}()
	return execC, errC, nil
}

// CreateRunner creates a new runner.
func (mc *ModuleClient) CreateRunner(req *api.CreateRunnerRequest) (*runnerpb.Runner, error) {
	s, err := mc.GetService(req.ServiceHash)
	if err != nil {
		return nil, err
	}
	envHash := hash.Sum([]byte(strings.Join(xos.EnvMergeSlices(s.Configuration.Env, req.Env), ",")))
	acc, err := mc.GetAccount()
	if err != nil {
		return nil, err
	}

	msg := runner.NewMsgCreateRunner(acc.GetAddress(), req.ServiceHash, envHash)
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
	msg := runner.NewMsgDeleteRunner(acc.GetAddress(), req.Hash)
	_, err = mc.BuildAndBroadcastMsg(msg)
	return err
}

// GetRunner returns the runner that matches given hash.
func (mc *ModuleClient) GetRunner(hash sdk.AccAddress) (*runnerpb.Runner, error) {
	var out *runnerpb.Runner
	route := sroutef("%s/%s/%s", runner.QuerierRoute, runner.QueryGetRunner, hash)
	return out, mc.QueryJSON(route, nil, &out)
}

// FilterRunner to apply while listing runners.
type FilterRunner struct {
	Address      string
	InstanceHash sdk.AccAddress
}

// ListRunner returns all runners.
func (mc *ModuleClient) ListRunner(f *FilterRunner) ([]*runnerpb.Runner, error) {
	var rs []*runnerpb.Runner
	route := sroutef("%s/%s", runner.QuerierRoute, runner.QueryListRunners)
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
		if (f.Address == "" || r.Address == f.Address) &&
			(f.InstanceHash.Empty() || r.InstanceHash.Equals(f.InstanceHash)) {
			out = append(out, r)
		}
	}
	return out, nil
}
