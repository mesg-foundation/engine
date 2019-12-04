package processsdk

import (
	"errors"
	"fmt"

	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/codec"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/ownership"
	"github.com/mesg-foundation/engine/process"
	instancesdk "github.com/mesg-foundation/engine/sdk/instance"
	ownershipsdk "github.com/mesg-foundation/engine/sdk/ownership"
	abci "github.com/tendermint/tendermint/abci/types"
)

const backendName = "process"

// Backend is the service backend.
type Backend struct {
	storeKey  *cosmostypes.KVStoreKey
	ownership *ownershipsdk.Backend
	instance  *instancesdk.Backend
}

// NewBackend returns the backend of the service sdk.
func NewBackend(appFactory *cosmos.AppFactory, ownership *ownershipsdk.Backend, instance *instancesdk.Backend) *Backend {
	backend := &Backend{
		storeKey:  cosmostypes.NewKVStoreKey(backendName),
		ownership: ownership,
		instance:  instance,
	}
	appBackendBasic := cosmos.NewAppModuleBasic(backendName)
	appBackend := cosmos.NewAppModule(appBackendBasic, backend.handler, backend.querier)
	appFactory.RegisterModule(appBackend)
	appFactory.RegisterStoreKey(backend.storeKey)
	return backend
}

func (s *Backend) handler(req cosmostypes.Request, msg cosmostypes.Msg) (hash.Hash, error) {
	switch msg := msg.(type) {
	case msgCreateProcess:
		p, err := s.Create(req, &msg)
		if err != nil {
			return nil, err
		}
		return p.Hash, nil
	case msgDeleteProcess:
		return nil, s.Delete(req, &msg)
	default:
		errmsg := fmt.Sprintf("unrecognized process msg type: %v", msg.Type())
		return nil, cosmostypes.ErrUnknownRequest(errmsg)
	}
}

func (s *Backend) querier(req cosmostypes.Request, path []string, _ abci.RequestQuery) (interface{}, error) {
	switch path[0] {
	case "get":
		hash, err := hash.Decode(path[1])
		if err != nil {
			return nil, err
		}
		return s.Get(req, hash)
	case "list":
		return s.List(req)
	default:
		return nil, errors.New("unknown service query endpoint" + path[0])
	}
}

// Create creates a new process.
func (s *Backend) Create(req cosmostypes.Request, msg *msgCreateProcess) (*process.Process, error) {
	store := req.KVStore(s.storeKey)
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
			if _, err := s.instance.Get(req, n.Result.InstanceHash); err != nil {
				return nil, err
			}
		case *process.Process_Node_Event_:
			if _, err := s.instance.Get(req, n.Event.InstanceHash); err != nil {
				return nil, err
			}
		case *process.Process_Node_Task_:
			if _, err := s.instance.Get(req, n.Task.InstanceHash); err != nil {
				return nil, err
			}
		}
	}

	value, err := codec.MarshalBinaryBare(p)
	if err != nil {
		return nil, err
	}

	if _, err := s.ownership.Create(req, msg.Owner, p.Hash, ownership.Ownership_Process); err != nil {
		return nil, err
	}

	store.Set(p.Hash, value)
	return p, nil
}

// Delete deletes a process and realated ownership.
func (s *Backend) Delete(req cosmostypes.Request, msg *msgDeleteProcess) error {
	if err := s.ownership.Delete(req, msg.Owner, msg.Request.Hash); err != nil {
		return err
	}
	req.KVStore(s.storeKey).Delete(msg.Request.Hash)
	return nil
}

// Get returns the service that matches given hash.
func (s *Backend) Get(req cosmostypes.Request, hash hash.Hash) (*process.Process, error) {
	store := req.KVStore(s.storeKey)
	if !store.Has(hash) {
		return nil, fmt.Errorf("process %q not found", hash)
	}

	var p *process.Process
	return p, codec.UnmarshalBinaryBare(store.Get(hash), &p)
}

// List returns all services.
func (s *Backend) List(req cosmostypes.Request) ([]*process.Process, error) {
	var (
		processes []*process.Process
		iter      = req.KVStore(s.storeKey).Iterator(nil, nil)
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
