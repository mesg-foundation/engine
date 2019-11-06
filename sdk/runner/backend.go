package runnersdk

import (
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/runner"
	instancesdk "github.com/mesg-foundation/engine/sdk/instance"
	abci "github.com/tendermint/tendermint/abci/types"
)

const backendName = "runner"

// Backend is the runner backend.
type Backend struct {
	cdc          *codec.Codec
	storeKey     *cosmostypes.KVStoreKey
	instanceBack *instancesdk.Backend
}

// NewBackend returns the backend of the runner sdk.
func NewBackend(appFactory *cosmos.AppFactory, instanceBack *instancesdk.Backend) *Backend {
	backend := &Backend{
		cdc:          appFactory.Cdc(),
		storeKey:     cosmostypes.NewKVStoreKey(backendName),
		instanceBack: instanceBack,
	}
	appBackendBasic := cosmos.NewAppModuleBasic(backendName)
	appBackend := cosmos.NewAppModule(appBackendBasic, backend.cdc, backend.handler, backend.querier)
	appFactory.RegisterModule(appBackend)
	appFactory.RegisterStoreKey(backend.storeKey)

	backend.cdc.RegisterConcrete(msgCreateRunner{}, "runner/create", nil)
	backend.cdc.RegisterConcrete(msgDeleteRunner{}, "runner/delete", nil)

	return backend
}

func (s *Backend) handler(request cosmostypes.Request, msg cosmostypes.Msg) cosmostypes.Result {
	switch msg := msg.(type) {
	case msgCreateRunner:
		run, err := s.Create(request, &msg)
		if err != nil {
			return cosmostypes.ErrInternal(err.Error()).Result()
		}
		return cosmostypes.Result{
			Data: run.Hash,
		}
	case msgDeleteRunner:
		if err := s.Delete(request, &msg); err != nil {
			return cosmostypes.ErrInternal(err.Error()).Result()
		}
		return cosmostypes.Result{}
	default:
		errmsg := fmt.Sprintf("Unrecognized runner Msg type: %v", msg.Type())
		return cosmostypes.ErrUnknownRequest(errmsg).Result()
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
	case "exists":
		hash, err := hash.Decode(path[1])
		if err != nil {
			return nil, err
		}
		return s.Exists(request, hash)
	default:
		return nil, errors.New("unknown runner query endpoint" + path[0])
	}
}

// Create creates a new runner.
func (s *Backend) Create(request cosmostypes.Request, msg *msgCreateRunner) (*runner.Runner, error) {
	store := cosmos.NewDB(request.KVStore(s.storeKey), s.cdc)

	// get instance and create it if needed
	inst, err := s.instanceBack.FetchOrCreate(request, msg.ServiceHash, msg.EnvHash)
	if err != nil {
		return nil, err
	}

	// create runner
	run := &runner.Runner{
		Address:      msg.Address.String(),
		InstanceHash: inst.Hash,
	}
	run.Hash = hash.Dump(run)

	// check if runner already exists.
	if exist, err := store.Has(run.Hash); err != nil {
		return nil, err
	} else if exist {
		return nil, errors.New("runner %q already exists" + run.Hash.String())
	}

	// save runner
	if err := store.Save(run.Hash, run); err != nil {
		return nil, err
	}
	return run, nil
}

// Delete deletes a runner.
func (s *Backend) Delete(request cosmostypes.Request, msg *msgDeleteRunner) error {
	store := cosmos.NewDB(request.KVStore(s.storeKey), s.cdc)
	runner := runner.Runner{}
	if err := store.Get(msg.RunnerHash, &runner); err != nil {
		return err
	}
	if runner.Address != msg.Address.String() {
		return errors.New("only the runner owner can remove itself")
	}
	return store.Delete(msg.RunnerHash)
}

// Get returns the runner that matches given hash.
func (s *Backend) Get(request cosmostypes.Request, hash hash.Hash) (*runner.Runner, error) {
	runner := runner.Runner{}
	if err := cosmos.NewDB(request.KVStore(s.storeKey), s.cdc).Get(hash, &runner); err != nil {
		return nil, err
	}
	return &runner, nil
}

// Exists returns true if a specific set of data exists in the database, false otherwise
func (s *Backend) Exists(request cosmostypes.Request, hash hash.Hash) (bool, error) {
	return cosmos.NewDB(request.KVStore(s.storeKey), s.cdc).Has(hash)
}

// List returns all runners.
func (s *Backend) List(request cosmostypes.Request) ([]*runner.Runner, error) {
	var (
		runners []*runner.Runner
		runner  *runner.Runner
		iter    = cosmos.NewDB(request.KVStore(s.storeKey), s.cdc).NewIterator()
	)
	for iter.Next() {
		if err := iter.Value(runner); err != nil {
			return nil, err
		}
		runners = append(runners, runner)
	}
	iter.Release()
	return runners, iter.Error()
}
