package runnersdk

import (
	"errors"
	"fmt"

	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/ext/xos"
	"github.com/mesg-foundation/engine/hash"
	instancepb "github.com/mesg-foundation/engine/instance"
	"github.com/mesg-foundation/engine/protobuf/api"
	runnerpb "github.com/mesg-foundation/engine/runner"
	servicesdk "github.com/mesg-foundation/engine/sdk/service"
	"github.com/mesg-foundation/engine/x/instance"
	"github.com/mesg-foundation/engine/x/runner"
)

// SDK is the runner sdk.
type SDK struct {
	serviceSDK   *servicesdk.SDK
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
func New(client *cosmos.Client, serviceSDK *servicesdk.SDK, container container.Container, engineName, port, ipfsEndpoint string) *SDK {
	sdk := &SDK{
		container:    container,
		serviceSDK:   serviceSDK,
		client:       client,
		port:         port,
		engineName:   engineName,
		ipfsEndpoint: ipfsEndpoint,
	}
	return sdk
}

// Create creates a new runner.
func (s *SDK) Create(req *api.CreateRunnerRequest) (*runnerpb.Runner, error) {
	// calculate instance's hash.
	// TODO: this should be merged with the same logic currently in instance sdk
	srv, err := s.serviceSDK.Get(req.ServiceHash)
	if err != nil {
		return nil, err
	}
	instanceEnv := xos.EnvMergeSlices(srv.Configuration.Env, req.Env)
	envHash := hash.Dump(instanceEnv)
	// TODO: should be done by instance or runner
	instanceHash := hash.Dump(&instancepb.Instance{
		ServiceHash: srv.Hash,
		EnvHash:     envHash,
	})
	acc, err := s.client.GetAccount()
	if err != nil {
		return nil, err
	}
	expRunnerHash := hash.Dump(&runnerpb.Runner{
		Address:      acc.GetAddress().String(),
		InstanceHash: instanceHash,
	})

	if runExisting, _ := s.Get(expRunnerHash); runExisting != nil {
		return nil, fmt.Errorf("runner %q already exists", runExisting.Hash)
	}

	// start the container
	imageHash, err := build(s.container, srv, s.ipfsEndpoint)
	if err != nil {
		return nil, err
	}
	_, err = start(s.container, srv, instanceHash, expRunnerHash, imageHash, instanceEnv, s.engineName, s.port)
	if err != nil {
		return nil, err
	}
	onError := func() {
		stop(s.container, expRunnerHash, srv.Dependencies)
	}

	msg := runner.NewMsgCreateRunner(acc.GetAddress(), req.ServiceHash, envHash)
	tx, err := s.client.BuildAndBroadcastMsg(msg)
	if err != nil {
		onError()
		return nil, err
	}
	run, err := s.Get(tx.Data)
	if err != nil {
		onError()
		return nil, err
	}
	if !run.Hash.Equal(expRunnerHash) {
		onError()
		return nil, errors.New("calculated runner hash is not the same")
	}
	return run, nil
}

// Delete deletes an existing runner.
func (s *SDK) Delete(req *api.DeleteRunnerRequest) error {
	// get runner before deleting it
	r, err := s.Get(req.Hash)
	if err != nil {
		return err
	}
	acc, err := s.client.GetAccount()
	if err != nil {
		return err
	}
	msg := runner.NewMsgDeleteRunner(acc.GetAddress(), req.Hash)
	_, err = s.client.BuildAndBroadcastMsg(msg)
	if err != nil {
		return err
	}

	var inst instancepb.Instance
	if err := s.client.QueryJSON(fmt.Sprintf("custom/%s/%s/%s", instance.QuerierRoute, instance.QueryGetInstance, r.InstanceHash), nil, &inst); err != nil {
		return err
	}

	srv, err := s.serviceSDK.Get(inst.ServiceHash)
	if err != nil {
		return err
	}

	// stop the local running service
	if err := stop(s.container, r.Hash, srv.Dependencies); err != nil {
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
func (s *SDK) Get(hash hash.Hash) (*runnerpb.Runner, error) {
	var r runnerpb.Runner
	if err := s.client.QueryJSON(fmt.Sprintf("custom/%s/%s/%s", runner.ModuleName, runner.QueryGetRunner, hash), nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// List returns all runners.
func (s *SDK) List(f *Filter) ([]*runnerpb.Runner, error) {
	var runners []*runnerpb.Runner
	if err := s.client.QueryJSON(fmt.Sprintf("custom/%s/%s", runner.ModuleName, runner.QueryListRunners), nil, &runners); err != nil {
		return nil, err
	}
	// no filter, returns
	if f == nil {
		return runners, nil
	}
	// filter results
	ret := make([]*runnerpb.Runner, 0)
	for _, r := range runners {
		if (f.Address == "" || r.Address == f.Address) &&
			(f.InstanceHash.IsZero() || r.InstanceHash.Equal(f.InstanceHash)) {
			ret = append(ret, r)
		}
	}
	return ret, nil
}
