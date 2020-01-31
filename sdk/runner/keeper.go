package runnersdk

import (
	"errors"
	"fmt"

	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/codec"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/runner"
	instancesdk "github.com/mesg-foundation/engine/sdk/instance"
)

// Keeper holds the logic to read and write data.
type Keeper struct {
	storeKey       *cosmostypes.KVStoreKey
	instanceKeeper *instancesdk.Keeper
}

// NewKeeper initialize a new keeper.
func NewKeeper(storeKey *cosmostypes.KVStoreKey, instanceKeeper *instancesdk.Keeper) *Keeper {
	return &Keeper{
		storeKey:       storeKey,
		instanceKeeper: instanceKeeper,
	}
}

// Create creates a new runner.
func (k *Keeper) Create(request cosmostypes.Request, msg *msgCreateRunner) (*runner.Runner, error) {
	store := request.KVStore(k.storeKey)
	inst, err := k.instanceKeeper.FetchOrCreate(request, msg.ServiceHash, msg.EnvHash)
	if err != nil {
		return nil, err
	}
	run := &runner.Runner{
		Address:      msg.Address.String(),
		InstanceHash: inst.Hash,
	}
	run.Hash = hash.Dump(run)
	if store.Has(run.Hash) {
		return nil, fmt.Errorf("runner %q already exists", run.Hash)
	}
	value, err := codec.MarshalBinaryBare(run)
	if err != nil {
		return nil, err
	}
	store.Set(run.Hash, value)
	return run, nil
}

// Delete deletes a runner.
func (k *Keeper) Delete(request cosmostypes.Request, msg *msgDeleteRunner) error {
	store := request.KVStore(k.storeKey)
	run := runner.Runner{}
	value := store.Get(msg.RunnerHash)
	if err := codec.UnmarshalBinaryBare(value, &run); err != nil {
		return err
	}
	if run.Address != msg.Address.String() {
		return errors.New("only the runner owner can remove itself")
	}
	store.Delete(msg.RunnerHash)
	return nil
}

// Get returns the runner that matches given hash.
func (k *Keeper) Get(request cosmostypes.Request, hash hash.Hash) (*runner.Runner, error) {
	store := request.KVStore(k.storeKey)
	if !store.Has(hash) {
		return nil, fmt.Errorf("runner %q not found", hash)
	}
	value := store.Get(hash)
	var run *runner.Runner
	return run, codec.UnmarshalBinaryBare(value, &run)
}

// List returns all runners.
func (k *Keeper) List(request cosmostypes.Request) ([]*runner.Runner, error) {
	var (
		runners []*runner.Runner
		iter    = request.KVStore(k.storeKey).Iterator(nil, nil)
	)
	for iter.Valid() {
		var run *runner.Runner
		value := iter.Value()
		if err := codec.UnmarshalBinaryBare(value, &run); err != nil {
			return nil, err
		}
		runners = append(runners, run)
		iter.Next()
	}
	iter.Close()
	return runners, nil
}
