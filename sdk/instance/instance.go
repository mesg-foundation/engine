package instancesdk

import (
	"crypto/sha256"
	"errors"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/docker/docker/pkg/archive"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/database"
	"github.com/mesg-foundation/core/instance"
	"github.com/mesg-foundation/core/x/xos"
)

// Instance exposes service instance APIs of MESG.
type Instance struct {
	container  container.Container
	serviceDB  database.ServiceDB
	instanceDB database.InstanceDB
}

// New creates a new Instance SDK with given options.
func New(c container.Container, serviceDB database.ServiceDB, instanceDB database.InstanceDB) *Instance {
	return &Instance{
		container:  c,
		serviceDB:  serviceDB,
		instanceDB: instanceDB,
	}
}

// Filter to apply while listing instances.
type Filter struct {
	ServiceHash []byte
}

// List instances by f filter.
func (i *Instance) List(f *Filter) ([]*instance.Instance, error) {
	if f != nil && len(f.ServiceHash) > 0 {
		return i.instanceDB.GetAllByService(f.ServiceHash)
	}
	return i.instanceDB.GetAll()
}

// Create creates a new service instance for service with id(sid/hash) and applies given env vars.
func (i *Instance) Create(serviceHash []byte, env []string) (*instance.Instance, error) {
	// get the service from service db.
	srv, err := i.serviceDB.Get(serviceHash)
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
	_, err = i.container.Build(path)
	if err != nil {
		return nil, err
	}

	// calculate the final env vars by overwriting user defined one's with defaults.
	instanceEnv := xos.EnvMergeMaps(xos.EnvSliceToMap(srv.Configuration.Env), xos.EnvSliceToMap(env))

	// calculate instance's hash.
	h := sha256.New()
	h.Write(srv.Hash)
	h.Write([]byte(xos.EnvMapToString(instanceEnv)))
	instanceHash := h.Sum(nil)

	// check if instance is already running.
	_, err = i.instanceDB.Get(instanceHash)
	if err == nil {
		return nil, errors.New("service's instance is already running")
	}
	if !database.IsErrNotFound(err) {
		return nil, err
	}

	// save & start instance.
	o := &instance.Instance{
		Hash:        instanceHash,
		ServiceHash: srv.Hash,
	}
	if err := i.instanceDB.Save(o); err != nil {
		return nil, err
	}

	_, err = i.start(o, xos.EnvMapToSlice(instanceEnv))
	return o, err
}

// GetAllByService retrives all instances of service by service's hash.
func (i *Instance) GetAllByService(serviceHash []byte) ([]*instance.Instance, error) {
	return i.instanceDB.GetAllByService(serviceHash)
}

// Delete an instance
func (i *Instance) Delete(hash []byte) error {
	inst, err := i.instanceDB.Get(hash)
	if err != nil {
		return err
	}
	if err := i.stop(inst); err != nil {
		return err
	}
	return i.instanceDB.Delete(hash)
}
