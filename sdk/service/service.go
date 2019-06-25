package servicesdk

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/cskr/pubsub"
	"github.com/docker/docker/pkg/archive"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/database"
	"github.com/mesg-foundation/core/hash"
	"github.com/mesg-foundation/core/hash/dirhash"
	"github.com/mesg-foundation/core/service"
)

// Service exposes service APIs of MESG.
type Service struct {
	ps *pubsub.PubSub

	container container.Container
	serviceDB database.ServiceDB
}

// New creates a new Service SDK with given options.
func New(c container.Container, serviceDB database.ServiceDB) *Service {
	return &Service{
		ps:        pubsub.New(0),
		container: c,
		serviceDB: serviceDB,
	}
}

// Create creates a new service from definition.
func (s *Service) Create(srv *service.Service) (*service.Service, error) {
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

	// calculate and apply hash to service.
	dh := dirhash.New(path)
	h, err := dh.Sum(hash.Dump(srv))
	if err != nil {
		return nil, err
	}
	srv.Hash = hash.Hash(h)

	// check if service is already deployed.
	if _, err := s.serviceDB.Get(srv.Hash); err == nil {
		return nil, errors.New("service is already deployed")
	}

	// build service's Docker image.
	_, err = s.container.Build(path)
	if err != nil {
		return nil, err
	}
	// TODO: the following test should be moved in New function
	if srv.Sid == "" {
		// make sure that sid doesn't have the same length with id.
		srv.Sid = "_" + srv.Hash.String()
	}

	return srv, s.serviceDB.Save(srv)
}

// Delete deletes the service by hash.
func (s *Service) Delete(hash hash.Hash) error {
	return s.serviceDB.Delete(hash)
}

// Get returns the service that matches given hash.
func (s *Service) Get(hash hash.Hash) (*service.Service, error) {
	return s.serviceDB.Get(hash)
}

// List returns all services.
func (s *Service) List() ([]*service.Service, error) {
	return s.serviceDB.All()
}
