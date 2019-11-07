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
	default:
		return nil, errors.New("unknown runner query endpoint" + path[0])
	}
}

// Create creates a new runner.
func (s *Backend) Create(request cosmostypes.Request, msg *msgCreateRunner) (*runner.Runner, error) {
	store := request.KVStore(s.storeKey)
	inst, err := s.instanceBack.FetchOrCreate(request, msg.ServiceHash, msg.EnvHash)
	if err != nil {
		return nil, err
	}
	run := &runner.Runner{
		Address:      msg.Address.String(),
		InstanceHash: inst.Hash,
	}
	run.Hash = hash.Dump(run)
	if store.Has(run.Hash) {
		return nil, errors.New("runner %q already exists" + run.Hash.String())
	}
	value, err := s.cdc.MarshalBinaryBare(run)
	if err != nil {
		return nil, err
	}
	store.Set(run.Hash, value)
	return run, nil
}

// Delete deletes a runner.
func (s *Backend) Delete(request cosmostypes.Request, msg *msgDeleteRunner) error {
	store := request.KVStore(s.storeKey)
	run := runner.Runner{}
	value := store.Get(msg.RunnerHash)
	if err := s.cdc.UnmarshalBinaryBare(value, &run); err != nil {
		return err
	}
	if run.Address != msg.Address.String() {
		return errors.New("only the runner owner can remove itself")
	}
	store.Delete(msg.RunnerHash)
	return nil
}

// Get returns the runner that matches given hash.
func (s *Backend) Get(request cosmostypes.Request, hash hash.Hash) (*runner.Runner, error) {
	var run *runner.Runner
	value := request.KVStore(s.storeKey).Get(hash)
	return run, s.cdc.UnmarshalBinaryBare(value, run)
}

// List returns all runners.
func (s *Backend) List(request cosmostypes.Request) ([]*runner.Runner, error) {
	var (
		runners []*runner.Runner
		iter    = request.KVStore(s.storeKey).Iterator(nil, nil)
	)
	for iter.Valid() {
		var run *runner.Runner
		value := iter.Value()
		if err := s.cdc.UnmarshalBinaryBare(value, run); err != nil {
			return nil, err
		}
		runners = append(runners, run)
		iter.Next()
	}
	iter.Close()
	return runners, nil
}
