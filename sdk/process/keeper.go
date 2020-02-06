package processsdk

import (
	"fmt"

	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/codec"
	"github.com/mesg-foundation/engine/hash"
	ownershippb "github.com/mesg-foundation/engine/ownership"
	"github.com/mesg-foundation/engine/process"
	instancesdk "github.com/mesg-foundation/engine/sdk/instance"
	"github.com/mesg-foundation/engine/x/ownership"
)

// Keeper is the service keeper.
type Keeper struct {
	storeKey        *cosmostypes.KVStoreKey
	ownershipKeeper ownership.Keeper
	instanceKeeper  *instancesdk.Keeper
}

// NewKeeper returns the keeper of the service sdk.
func NewKeeper(storeKey *cosmostypes.KVStoreKey, ownershipKeeper ownership.Keeper, instanceKeeper *instancesdk.Keeper) *Keeper {
	return &Keeper{
		storeKey:        storeKey,
		ownershipKeeper: ownershipKeeper,
		instanceKeeper:  instanceKeeper,
	}
}

// Create creates a new process.
func (k *Keeper) Create(req cosmostypes.Request, msg *msgCreateProcess) (*process.Process, error) {
	store := req.KVStore(k.storeKey)
	p := &process.Process{
		Name:  msg.Request.Name,
		Nodes: msg.Request.Nodes,
		Edges: msg.Request.Edges,
	}
	p.Hash = hash.Dump(p)
	if store.Has(p.Hash) {
		return nil, fmt.Errorf("process %q already exists", p.Hash)
	}

	for _, node := range p.Nodes {
		switch n := node.Type.(type) {
		case *process.Process_Node_Result_:
			if _, err := k.instanceKeeper.Get(req, n.Result.InstanceHash); err != nil {
				return nil, err
			}
		case *process.Process_Node_Event_:
			if _, err := k.instanceKeeper.Get(req, n.Event.InstanceHash); err != nil {
				return nil, err
			}
		case *process.Process_Node_Task_:
			if _, err := k.instanceKeeper.Get(req, n.Task.InstanceHash); err != nil {
				return nil, err
			}
		}
	}

	value, err := codec.MarshalBinaryBare(p)
	if err != nil {
		return nil, err
	}

	if _, err := k.ownershipKeeper.Set(req, msg.Owner, p.Hash, ownershippb.Ownership_Process); err != nil {
		return nil, err
	}

	store.Set(p.Hash, value)
	return p, nil
}

// Delete deletes a process and realated ownership.
func (k *Keeper) Delete(req cosmostypes.Request, msg *msgDeleteProcess) error {
	if err := k.ownershipKeeper.Delete(req, msg.Owner, msg.Request.Hash); err != nil {
		return err
	}
	req.KVStore(k.storeKey).Delete(msg.Request.Hash)
	return nil
}

// Get returns the service that matches given hash.
func (k *Keeper) Get(req cosmostypes.Request, hash hash.Hash) (*process.Process, error) {
	store := req.KVStore(k.storeKey)
	if !store.Has(hash) {
		return nil, fmt.Errorf("process %q not found", hash)
	}

	var p *process.Process
	return p, codec.UnmarshalBinaryBare(store.Get(hash), &p)
}

// List returns all services.
func (k *Keeper) List(req cosmostypes.Request) ([]*process.Process, error) {
	var (
		processes []*process.Process
		iter      = req.KVStore(k.storeKey).Iterator(nil, nil)
	)
	for iter.Valid() {
		var p *process.Process
		if err := codec.UnmarshalBinaryBare(iter.Value(), &p); err != nil {
			return nil, err
		}
		processes = append(processes, p)
		iter.Next()
	}
	iter.Close()
	return processes, nil
}
