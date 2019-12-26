package runnersdk

import (
	"errors"
	"fmt"

	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/instance"
	"github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/runner"
	instancesdk "github.com/mesg-foundation/engine/sdk/instance"
	servicesdk "github.com/mesg-foundation/engine/sdk/service"
	"github.com/mesg-foundation/engine/x/xos"
)

// SDK is the runner sdk.
type SDK struct {
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
func New(client *cosmos.Client, serviceSDK *servicesdk.SDK, instanceSDK *instancesdk.SDK, container container.Container, engineName, port, ipfsEndpoint string) *SDK {
	sdk := &SDK{
		container:    container,
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
func (s *SDK) Create(req *api.CreateRunnerRequest) (*runner.Runner, error) {
	// calculate instance's hash.
	// TODO: this should be merged with the same logic currently in instance sdk
	srv, err := s.serviceSDK.Get(req.ServiceHash)
	if err != nil {
		return nil, err
	}
	instanceEnv := xos.EnvMergeSlices(srv.Configuration.Env, req.Env)
	envHash := hash.Dump(instanceEnv)
	// TODO: should be done by instance or runner
	instanceHash := hash.Dump(&instance.Instance{
		ServiceHash: srv.Hash,
		EnvHash:     envHash,
	})
	acc, err := s.client.GetAccount()
	if err != nil {
		return nil, err
	}
	expRunnerHash := hash.Dump(&runner.Runner{
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

	msg := newMsgCreateRunner(acc.GetAddress(), req.ServiceHash, envHash)
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
	runner, err := s.Get(req.Hash)
	if err != nil {
		return err
	}
	acc, err := s.client.GetAccount()
	if err != nil {
		return err
	}
	msg := newMsgDeleteRunner(acc.GetAddress(), req.Hash)
	_, err = s.client.BuildAndBroadcastMsg(msg)
	if err != nil {
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
	if err := stop(s.container, runner.Hash, srv.Dependencies); err != nil {
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
			(f.InstanceHash.IsZero() || runner.InstanceHash.Equal(f.InstanceHash)) {
			ret = append(ret, runner)
		}
	}
	return ret, nil
}
