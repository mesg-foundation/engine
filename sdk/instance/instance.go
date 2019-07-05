package instancesdk

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/docker/docker/pkg/archive"
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
	service    *servicesdk.Service
	instanceDB database.InstanceDB
}

// New creates a new Instance SDK with given options.
func New(c container.Container, service *servicesdk.Service, instanceDB database.InstanceDB) *Instance {
	return &Instance{
		container:  c,
		service:    service,
		instanceDB: instanceDB,
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

	// download and untar service context into path.
	path, err := ioutil.TempDir("", "mesg")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(path)

	resp, err := http.Get("http://ipfs.app.mesg.com:8080/ipfs/" + srv.Source)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("service's source code is not reachable")
	}
	defer resp.Body.Close()

	if err := archive.Untar(resp.Body, path, nil); err != nil {
		return nil, err
	}

	// build service's Docker image and apply to service.
	imageHash, err := i.container.Build(path)
	if err != nil {
		return nil, err
	}

	// calculate the final env vars by overwriting user defined one's with defaults.
	instanceEnv := xos.EnvMergeMaps(xos.EnvSliceToMap(srv.Configuration.Env), xos.EnvSliceToMap(env))

	// calculate instance's hash.
	h := hash.New()
	h.Write(srv.Hash)
	h.Write([]byte(xos.EnvMapToString(instanceEnv)))
	instanceHash := h.Sum(nil)

	// check if instance already exists
	_, err = i.instanceDB.Get(instanceHash)
	if err == nil {
		return nil, &AlreadyExistsError{Hash: instanceHash}
	}
	if !database.IsErrNotFound(err) {
		return nil, err
	}

	// save & start instance.
	inst := &instance.Instance{
		Hash:        instanceHash,
		ServiceHash: srv.Hash,
	}
	if err := i.instanceDB.Save(inst); err != nil {
		return nil, err
	}

	_, err = i.start(inst, imageHash, xos.EnvMapToSlice(instanceEnv))
	return inst, err
}

// Delete deletes an instance.
// if deleteData is enabled, any persistent data that belongs to
// the instance and to its dependencies will also be deleted.
func (i *Instance) Delete(hash hash.Hash, deleteData bool) error {
	inst, err := i.instanceDB.Get(hash)
	if err != nil {
		return err
	}
	if err := i.stop(inst); err != nil {
		return err
	}
	// delete volumes first before the instance. this way if
	// deleting volumes fails, process can be retried by the user again
	// because instance still will be in the db.
	if deleteData {
		if err := i.deleteData(inst); err != nil {
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
