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
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/service/manager"
	"github.com/mesg-foundation/core/x/xos"
	"github.com/mr-tron/base58"
)

// Instance exposes service instance APIs of MESG.
type Instance struct {
	m          manager.Manager
	container  container.Container
	db         database.ServiceDB
	instanceDB database.InstanceDB
}

// New creates a new Instance SDK with given options.
func New(m manager.Manager, c container.Container, db database.ServiceDB, instanceDB database.InstanceDB) *Instance {
	return &Instance{
		m:          m,
		container:  c,
		db:         db,
		instanceDB: instanceDB,
	}
}

// Create creates a new service instance for service with id(sid/hash) and applies given env vars.
func (i *Instance) Create(id string, env []string) (*service.Instance, error) {
	// get the service from service db.
	srv, err := i.db.Get(id)
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

	// overwrite default env vars with user defined ones.
	instanceEnv := xos.EnvMergeMaps(xos.EnvSliceToMap(srv.Configuration.Env), xos.EnvSliceToMap(env))

	// calculate instance's hash.
	h := sha256.New()
	h.Write([]byte(srv.Hash))
	h.Write([]byte(xos.EnvMapToString(instanceEnv)))
	instanceHash := base58.Encode(h.Sum(nil))

	// check if instance is already running.
	_, err = i.instanceDB.Get(instanceHash)
	if err == nil {
		return nil, errors.New("service's instance is already running")
	}
	if !database.IsErrNotFound(err) {
		return nil, err
	}

	// save & start instance.
	o := &service.Instance{
		Sid:         srv.Sid,
		Hash:        instanceHash,
		ServiceHash: srv.Hash,
		ImageHash:   imageHash,
	}
	if err := i.instanceDB.Save(o); err != nil {
		return nil, err
	}

	srv.Hash = instanceHash
	srv.Configuration.Image = imageHash
	srv.Configuration.Env = xos.EnvMapToSlice(instanceEnv)
	serviceIDs, networkID, err := i.m.Start(srv)
	if err != nil {
		return nil, err
	}

	// TODO(ilgooz) this should be set by Manager.Start(), by the underlying
	// container orchestration tool.
	o.COType = "docker"
	o.COInfo = struct {
		ServiceIDs []string
		NetworkID  string
	}{serviceIDs, networkID}

	return o, i.instanceDB.Save(o)
}
