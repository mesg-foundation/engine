package servicesdk

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/docker/docker/pkg/archive"
	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/hash/dirhash"
	"github.com/mesg-foundation/engine/service"
	"github.com/mesg-foundation/engine/service/validator"
)

type Default struct {
// TODO: rename to deprected
// TODO: split it to deprected and utils
	container     container.Container
	keeperFactory KeeperFactory
}

func NewDefault(c container.Container, keeperFactory KeeperFactory) *Default {
	return &Default{
		container:     c,
		keeperFactory: keeperFactory,
	}
}

// Create creates a new service from definition.
func (s *Default) Create(ctx context.Context, srv *service.Service) (*service.Service, error) {
	if srv.Configuration == nil {
		srv.Configuration = &service.Configuration{}
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

	// calculate and apply hash to service.
	dh := dirhash.New(path)
	h, err := dh.Sum(hash.Dump(srv))
	if err != nil {
		return nil, err
	}
	srv.Hash = hash.Hash(h)

	keeper, err := s.keeperFactory(ctx)
	if err != nil {
		return nil, err
	}

	// check if service already exists.
	if _, err := keeper.Get(srv.Hash); err == nil {
		return nil, &AlreadyExistsError{Hash: srv.Hash}
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

	if err := validator.ValidateService(srv); err != nil {
		return nil, err
	}
	return srv, keeper.Save(srv)
}

// Delete deletes the service by hash.
func (s *Default) Delete(ctx context.Context, hash hash.Hash) error {
	keeper, err := s.keeperFactory(ctx)
	if err != nil {
		return err
	}
	return keeper.Delete(hash)
}

// Get returns the service that matches given hash.
func (s *Default) Get(ctx context.Context, hash hash.Hash) (*service.Service, error) {
	keeper, err := s.keeperFactory(ctx)
	if err != nil {
		return nil, err
	}
	return keeper.Get(hash)
}

// List returns all services.
func (s *Default) List(ctx context.Context) ([]*service.Service, error) {
	keeper, err := s.keeperFactory(ctx)
	if err != nil {
		return nil, err
	}
	return keeper.All()
}
