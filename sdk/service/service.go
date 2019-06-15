package servicesdk

import (
	"crypto/sha256"
	"errors"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/cnf/structhash"
	"github.com/cskr/pubsub"
	"github.com/docker/docker/pkg/archive"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/database"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/service/manager"
	"github.com/mesg-foundation/core/utils/dirhash"
	"github.com/mr-tron/base58"
)

// Service exposes service APIs of MESG.
type Service struct {
	ps *pubsub.PubSub

	m         manager.Manager
	container container.Container
	db        database.ServiceDB
	execDB    database.ExecutionDB
}

// New creates a new Service SDK with given options.
func New(m manager.Manager, c container.Container, db database.ServiceDB, execDB database.ExecutionDB) *Service {
	return &Service{
		ps:        pubsub.New(0),
		m:         m,
		container: c,
		db:        db,
		execDB:    execDB,
	}
}

// Create creates a new service from definition.
func (s *Service) Create(srv *service.Service) error {
	// calculate and apply hash to service.
	h := sha256.New()
	h.Write(structhash.Sha1(srv, 1))

	var path string
	if srv.Source != "" {
		// download and untar service context into path.
		var err error
		path, err = ioutil.TempDir("", "mesg")
		if err != nil {
			return err
		}
		defer os.RemoveAll(path)

		resp, err := http.Get("http://ipfs.app.mesg.com:8080/ipfs/" + srv.Source)
		if err != nil {
			return err
		}
		if resp.StatusCode != 200 {
			return errors.New("service's source code is not reachable")
		}
		defer resp.Body.Close()

		if err := archive.Untar(resp.Body, path, nil); err != nil {
			return err
		}

		dh := dirhash.New(path)
		ph, err := dh.Sum(nil)
		if err != nil {
			return err
		}
		h.Write(ph)
	}

	srv.Hash = base58.Encode(h.Sum(nil))

	// check if service is already deployed.
	if _, err := s.db.Get(srv.Hash); err == nil {
		return errors.New("service is already deployed: " + srv.Hash)
	}

	if path != "" {
		// build service's Docker image and apply to service.
		imageHash, err := s.container.Build(path)
		if err != nil {
			return err
		}
		srv.Configuration.Image = imageHash
	}

	// TODO: the following test should be moved in New function
	if srv.Sid == "" {
		// make sure that sid doesn't have the same length with id.
		srv.Sid = "_" + srv.Hash
	}

	return s.db.Save(srv)
}
