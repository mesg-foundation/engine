package resultsdk

import (
	"errors"
	"fmt"

	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/codec"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/result"
	executionsdk "github.com/mesg-foundation/engine/sdk/execution"
	instancesdk "github.com/mesg-foundation/engine/sdk/instance"
	servicesdk "github.com/mesg-foundation/engine/sdk/service"
	abci "github.com/tendermint/tendermint/abci/types"
)

const backendName = "result"

// Backend is the result backend.
type Backend struct {
	storeKey     *cosmostypes.KVStoreKey
	serviceBack  *servicesdk.Backend
	instanceBack *instancesdk.Backend
	execBack     *executionsdk.Backend
}

// NewBackend returns the backend of the result sdk.
func NewBackend(appFactory *cosmos.AppFactory, serviceBack *servicesdk.Backend, instanceBack *instancesdk.Backend, execBack *executionsdk.Backend) *Backend {
	backend := &Backend{
		storeKey:     cosmostypes.NewKVStoreKey(backendName),
		serviceBack:  serviceBack,
		instanceBack: instanceBack,
		execBack:     execBack,
	}
	appBackendBasic := cosmos.NewAppModuleBasic(backendName)
	appBackend := cosmos.NewAppModule(appBackendBasic, backend.handler, backend.querier)
	appFactory.RegisterModule(appBackend)
	appFactory.RegisterStoreKey(backend.storeKey)

	return backend
}

func (s *Backend) handler(request cosmostypes.Request, msg cosmostypes.Msg) (hash.Hash, error) {
	switch msg := msg.(type) {
	case msgCreateResult:
		exec, err := s.Create(request, msg)
		if err != nil {
			return nil, err
		}
		return exec.Hash, nil
	default:
		errmsg := fmt.Sprintf("Unrecognized result msg type: %v", msg.Type())
		return nil, cosmostypes.ErrUnknownRequest(errmsg)
	}
}

func (s *Backend) querier(request cosmostypes.Request, path []string, req abci.RequestQuery) (interface{}, error) {
	switch path[0] {
	case "get":
		hash, err := hash.Decode(path[1])
		if err != nil {
			return nil, err
		}
		return s.Get(request, hash)
	case "list":
		return s.List(request)
	default:
		return nil, errors.New("unknown result query endpoint" + path[0])
	}
}

// Create creates a new result.
func (s *Backend) Create(request cosmostypes.Request, msg msgCreateResult) (*result.Result, error) {
	exec, err := s.execBack.Get(request, msg.Request.RequestHash)
	if err != nil {
		return nil, err
	}

	ress, err := s.fetchResultsOfRequest(request, exec.Hash)
	if err != nil {
		return nil, err
	}
	if len(ress) > 0 {
		return nil, fmt.Errorf("execution %q has already a result", msg.Request.RequestHash)
	}

	var res *result.Result

	switch reqRes := msg.Request.Result.(type) {
	case *api.CreateResultRequest_Outputs:
		if err := s.validateExecutionOutput(request, exec.InstanceHash, exec.TaskKey, reqRes.Outputs); err != nil {
			res = result.NewWithError(exec.Hash, err.Error())
		} else {
			res = result.NewWithOutputs(exec.Hash, reqRes.Outputs)
		}
	case *api.CreateResultRequest_Error:
		res = result.NewWithError(exec.Hash, reqRes.Error)
	default:
		return nil, errors.New("no result outputs or error supplied")
	}

	value, err := codec.MarshalBinaryBare(res)
	if err != nil {
		return nil, err
	}
	store := request.KVStore(s.storeKey)
	if store.Has(res.Hash) {
		return nil, fmt.Errorf("result %q already exists", res.Hash)
	}
	store.Set(res.Hash, value)
	return res, nil
}

func (s *Backend) validateExecutionOutput(request cosmostypes.Request, instanceHash hash.Hash, taskKey string, outputs *types.Struct) error {
	inst, err := s.instanceBack.Get(request, instanceHash)
	if err != nil {
		return err
	}
	srv, err := s.serviceBack.Get(request, inst.ServiceHash)
	if err != nil {
		return err
	}
	return srv.RequireTaskOutputs(taskKey, outputs)
}

// Get returns the result that matches given hash.
func (s *Backend) Get(request cosmostypes.Request, hash hash.Hash) (*result.Result, error) {
	store := request.KVStore(s.storeKey)
	var res *result.Result
	if !store.Has(hash) {
		return nil, fmt.Errorf("result %q not found", hash)
	}
	return res, codec.UnmarshalBinaryBare(store.Get(hash), &res)
}

// List returns all executions.
func (s *Backend) List(request cosmostypes.Request) ([]*result.Result, error) {
	var (
		execs []*result.Result
		iter  = request.KVStore(s.storeKey).Iterator(nil, nil)
	)
	for iter.Valid() {
		var res *result.Result
		value := iter.Value()
		if err := codec.UnmarshalBinaryBare(value, &res); err != nil {
			return nil, err
		}
		execs = append(execs, res)
		iter.Next()
	}
	iter.Close()
	return execs, nil
}

func (s *Backend) fetchResultsOfRequest(request cosmostypes.Request, requestHash hash.Hash) ([]*result.Result, error) {
	var execsMatch []*result.Result
	execs, err := s.List(request)
	if err != nil {
		return nil, err
	}
	for _, res := range execs {
		if res.RequestHash.Equal(requestHash) {
			execsMatch = append(execsMatch, res)
		}
	}
	return execsMatch, nil
}
