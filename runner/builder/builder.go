package builder

import (
	"errors"
	"fmt"

	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/ext/xos"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/hash/serializer"
	instancepb "github.com/mesg-foundation/engine/instance"
	"github.com/mesg-foundation/engine/protobuf/api"
	runnerpb "github.com/mesg-foundation/engine/runner"
)

// Builder is the runner builder.
type Builder struct {
	mc           *cosmos.ModuleClient
	container    container.Container
	engineName   string
	port         string
	ipfsEndpoint string
}

// New returns the runner sdk.
func New(mc *cosmos.ModuleClient, container container.Container, engineName, port, ipfsEndpoint string) *Builder {
	sdk := &Builder{
		container:    container,
		mc:           mc,
		engineName:   engineName,
		port:         port,
		ipfsEndpoint: ipfsEndpoint,
	}
	return sdk
}

// Create creates a new runner.
func (b *Builder) Create(req *api.CreateRunnerRequest) (*runnerpb.Runner, error) {
	// calculate instance's hash.
	// TODO: this should be merged with the same logic currently in instance sdk
	srv, err := b.mc.GetService(req.ServiceHash)
	if err != nil {
		return nil, err
	}

	instanceEnv := xos.EnvMergeSlices(srv.Configuration.Env, req.Env)
	envHash := hash.Dump(serializer.StringSlice(instanceEnv))
	// TODO: should be done by instance or runner
	instanceHash := hash.Dump(&instancepb.Instance{
		ServiceHash: srv.Hash,
		EnvHash:     envHash,
	})
	acc, err := b.mc.GetAccount()
	if err != nil {
		return nil, err
	}
	expRunnerHash := hash.Dump(&runnerpb.Runner{
		Address:      acc.GetAddress().String(),
		InstanceHash: instanceHash,
	})

	if runExisting, _ := b.mc.GetRunner(expRunnerHash); runExisting != nil {
		return nil, fmt.Errorf("runner %q already exists", runExisting.Hash)
	}

	// start the container
	imageHash, err := build(b.container, srv, b.ipfsEndpoint)
	if err != nil {
		return nil, err
	}
	_, err = start(b.container, srv, instanceHash, expRunnerHash, imageHash, instanceEnv, b.engineName, b.port)
	if err != nil {
		return nil, err
	}

	run, err := b.mc.CreateRunner(req)
	if err != nil {
		stop(b.container, expRunnerHash, srv.Dependencies)
		return nil, err
	}

	if !run.Hash.Equal(expRunnerHash) {
		stop(b.container, expRunnerHash, srv.Dependencies)
		return nil, errors.New("calculated runner hash is not the same")
	}
	return run, nil
}

// Delete deletes an existing runner.
func (b *Builder) Delete(req *api.DeleteRunnerRequest) error {
	// get runner before deleting it
	r, err := b.mc.GetRunner(req.Hash)
	if err != nil {
		return err
	}

	if err := b.mc.DeleteRunner(req); err != nil {
		return err
	}

	inst, err := b.mc.GetInstance(r.InstanceHash)
	if err != nil {
		return err
	}

	srv, err := b.mc.GetService(inst.ServiceHash)
	if err != nil {
		return err
	}

	// stop the local running service
	if err := stop(b.container, r.Hash, srv.Dependencies); err != nil {
		return err
	}

	// remove local volume
	if req.DeleteData {
		if err := deleteData(b.container, srv); err != nil {
			return err
		}
	}

	return nil
}
