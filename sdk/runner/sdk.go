package runnersdk

import (
	"fmt"

	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/instance"
	"github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/runner"
	accountsdk "github.com/mesg-foundation/engine/sdk/account"
	instancesdk "github.com/mesg-foundation/engine/sdk/instance"
	servicesdk "github.com/mesg-foundation/engine/sdk/service"
	"github.com/mesg-foundation/engine/x/xos"
	"github.com/tendermint/tendermint/mempool"
)

// SDK is the runner sdk.
type SDK struct {
	accountSDK   *accountsdk.SDK
	serviceSDK   *servicesdk.SDK
	instanceSDK  *instancesdk.SDK
	client       *cosmos.Client
	container    container.Container
	port         string
	engineName   string
	ipfsEndpoint string
}

// Filter to apply while listing runners.
type Filter struct {
	Address      string
	InstanceHash hash.Hash
}

// New returns the runner sdk.
func New(client *cosmos.Client, accountSDK *accountsdk.SDK, serviceSDK *servicesdk.SDK, instanceSDK *instancesdk.SDK, container container.Container, engineName, port, ipfsEndpoint string) *SDK {
	sdk := &SDK{
		container:    container,
		accountSDK:   accountSDK,
		serviceSDK:   serviceSDK,
		instanceSDK:  instanceSDK,
		client:       client,
		port:         port,
		engineName:   engineName,
		ipfsEndpoint: ipfsEndpoint,
	}
	return sdk
}

// Create creates a new runner.
func (s *SDK) Create(req *api.CreateRunnerRequest, accountName, accountPassword string) (*runner.Runner, error) {
	account, err := s.accountSDK.Get(accountName)
	if err != nil {
		return nil, err
	}
	// TODO: pass account by parameters
	user, err := cosmostypes.AccAddressFromBech32(account.Address)
	if err != nil {
		return nil, err
	}

	// calculate instance's hash.
	// TODO: this should be merged with the same logic currently in instance sdk
	srv, err := s.serviceSDK.Get(req.ServiceHash)
	if err != nil {
		return nil, err
	}
	instanceEnv := xos.EnvMergeSlices(srv.Configuration.Env, req.Env)
	envHash := hash.Dump(instanceEnv)
	// TODO: should be done by instance
	instanceHash := hash.Dump(&instance.Instance{
		ServiceHash: srv.Hash,
		EnvHash:     envHash,
	})

	// start the container
	imageHash, err := build(s.container, srv, s.ipfsEndpoint)
	if err != nil {
		return nil, err
	}
	_, err = start(s.container, srv, instanceHash, imageHash, instanceEnv, s.engineName, s.port)
	if err != nil {
		return nil, err
	}
	onError := func() {
		stop(s.container, instanceHash, srv.Dependencies)
	}

	msg := newMsgCreateRunner(user, req.ServiceHash, envHash)
	tx, err := s.client.BuildAndBroadcastMsg(msg, accountName, accountPassword)
	if err != nil {
		defer onError()
		if err == mempool.ErrTxInCache {
			return nil, fmt.Errorf("runner already exists: %w", err)
		}
		return nil, err
	}
	return s.Get(tx.Data)
}

// Delete deletes an existing runner.
func (s *SDK) Delete(req *api.DeleteRunnerRequest, accountName, accountPassword string) error {
	account, err := s.accountSDK.Get(accountName)
	if err != nil {
		return err
	}
	// TODO: pass account by parameters
	user, err := cosmostypes.AccAddressFromBech32(account.Address)
	if err != nil {
		return err
	}

	// get runner before deleting it
	runner, err := s.Get(req.Hash)
	if err != nil {
		return err
	}

	msg := newMsgDeleteRunner(user, req.Hash)
	_, err = s.client.BuildAndBroadcastMsg(msg, accountName, accountPassword)
	if err != nil {
		if err == mempool.ErrTxInCache {
			return fmt.Errorf("runner already deleted: %w", err)
		}
		return err
	}

	// get instance and service
	inst, err := s.instanceSDK.Get(runner.InstanceHash)
	if err != nil {
		return err
	}
	srv, err := s.serviceSDK.Get(inst.ServiceHash)
	if err != nil {
		return err
	}

	// stop the local running service
	if err := stop(s.container, inst.Hash, srv.Dependencies); err != nil {
		return err
	}

	// remove local volume
	if req.DeleteData {
		if err := deleteData(s.container, srv); err != nil {
			return err
		}
	}

	return nil
}

// Get returns the runner that matches given hash.
func (s *SDK) Get(hash hash.Hash) (*runner.Runner, error) {
	var runner runner.Runner
	if err := s.client.Query("custom/"+backendName+"/get/"+hash.String(), nil, &runner); err != nil {
		return nil, err
	}
	return &runner, nil
}

// List returns all runners.
func (s *SDK) List(f *Filter) ([]*runner.Runner, error) {
	var runners []*runner.Runner
	if err := s.client.Query("custom/"+backendName+"/list", nil, &runners); err != nil {
		return nil, err
	}
	// no filter, returns
	if f == nil {
		return runners, nil
	}
	// filter results
	ret := make([]*runner.Runner, 0)
	for _, runner := range runners {
		if (f.Address == "" || runner.Address == f.Address) &&
			(f.InstanceHash.IsZero() || runner.Hash.Equal(f.InstanceHash)) {
			ret = append(ret, runner)
		}
	}
	return ret, nil
}

// Exists returns if a runner already exists.
func (s *SDK) Exists(hash hash.Hash) (bool, error) {
	var exists bool
	if err := s.client.Query("custom/"+backendName+"/exists/"+hash.String(), nil, &exists); err != nil {
		return false, err
	}
	return exists, nil
}
