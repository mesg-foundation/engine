package instancesdk

import (
	"fmt"

	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/database"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/instance"
	servicesdk "github.com/mesg-foundation/engine/sdk/service"
	"github.com/mesg-foundation/engine/x/xos"
)

// Instance exposes service instance APIs of MESG.
type Instance struct {
	container  container.Container
	service    servicesdk.Service
	instanceDB database.InstanceDB

	port       string
	engineName string
}

// New creates a new Instance SDK with given options.
func New(c container.Container, service servicesdk.Service, instanceDB database.InstanceDB, engineName, port string) *Instance {
	return &Instance{
		container:  c,
		service:    service,
		instanceDB: instanceDB,
		port:       port,
		engineName: engineName,
	}
}

// Get retrieves instance by hash.
func (i *Instance) Get(hash hash.Hash) (*instance.Instance, error) {
	return i.instanceDB.Get(hash)
}

// Filter to apply while listing instances.
type Filter struct {
	ServiceHash  hash.Hash
	InstanceHash hash.Hash
}

// List instances by f filter.
func (i *Instance) List(f *Filter) ([]*instance.Instance, error) {
	instances, err := i.instanceDB.GetAll()
	if err != nil {
		return nil, err
	}
	if f == nil {
		return instances, nil
	}

	ret := make([]*instance.Instance, 0)
	for _, instance := range instances {
		if (f.ServiceHash.IsZero() || instance.ServiceHash.Equal(f.ServiceHash)) &&
			(f.InstanceHash.IsZero() || instance.Hash.Equal(f.InstanceHash)) {
			ret = append(ret, instance)
		}
	}
	return ret, nil
}

// Create creates a new service instance for service with id(sid/hash) and applies given env vars.
func (i *Instance) Create(serviceHash hash.Hash, env []string) (*instance.Instance, error) {
	// get the service from service db.
	srv, err := i.service.Get(serviceHash)
	if err != nil {
		return nil, err
	}

	// build service's Docker image and apply to service.
	imageHash, err := build(i.container, srv)
	if err != nil {
		return nil, err
	}

	// calculate the final env vars by overwriting user defined one's with defaults.
	instanceEnv := xos.EnvMergeSlices(srv.Configuration.Env, env)

	// calculate instance's hash.
	inst := &instance.Instance{
		ServiceHash: srv.Hash,
		EnvHash:     hash.Dump(instanceEnv),
	}
	inst.Hash = hash.Dump(inst)

	// check if instance already exists
	if exist, err := i.instanceDB.Exist(inst.Hash); err != nil {
		return nil, err
	} else if exist {
		return nil, &AlreadyExistsError{Hash: inst.Hash}
	}

	// save & start instance.
	if err := i.instanceDB.Save(inst); err != nil {
		return nil, err
	}

	_, err = start(i.container, srv, inst.Hash, imageHash, instanceEnv, i.engineName, i.port)
	return inst, err
}

// Delete deletes an instance.
// if shouldDeleteData is true, any persistent data that belongs to
// the instance and to its dependencies will also be deleted.
func (i *Instance) Delete(hash hash.Hash, shouldDeleteData bool) error {
	inst, err := i.instanceDB.Get(hash)
	if err != nil {
		return err
	}
	// get the service from service db.
	srv, err := i.service.Get(inst.ServiceHash)
	if err != nil {
		return err
	}
	if err := stop(i.container, hash, srv.Dependencies); err != nil {
		return err
	}
	// delete volumes first before the instance. this way if
	// deleting volumes fails, process can be retried by the user again
	// because instance still will be in the db.
	if shouldDeleteData {
		if err := deleteData(i.container, srv); err != nil {
			return err
		}
	}
	return i.instanceDB.Delete(hash)
}

// AlreadyExistsError is an not found error.
type AlreadyExistsError struct {
	Hash hash.Hash
}

func (e *AlreadyExistsError) Error() string {
	return fmt.Sprintf("instance %q already exists", e.Hash.String())
}
