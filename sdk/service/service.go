package servicesdk

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cskr/pubsub"
	"github.com/docker/docker/pkg/archive"
	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/database"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/hash/dirhash"
	"github.com/mesg-foundation/engine/service"
	"github.com/mesg-foundation/engine/service/validator"
	abci "github.com/tendermint/tendermint/abci/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

// Service exposes service APIs of MESG.
type Service struct {
	AppModuleBasic

	ps *pubsub.PubSub

	container container.Container
	serviceDB database.ServiceDB

	keeper     Keeper
	abciClient rpcclient.Client
	cdc        *codec.Codec
}

const GetQueryPath = "get"

type GetQuery struct {
	Service service.Service `json:"service"`
}

func (s *Service) SetKeeper(keeper Keeper) {
	s.keeper = keeper
}
func (s *Service) SetClient(abciClient rpcclient.Client) {
	s.abciClient = abciClient
}
func (s *Service) SetCodec(cdc *codec.Codec) {
	s.cdc = cdc
}

// NewQuerier is the module level router for state queries.
func (s Service) NewQuerierHandler() sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case GetQueryPath:
			service, err := s.querierGet(ctx, path[1:])
			if err != nil {
				return nil, sdk.ErrInternal(err.Error())
			}
			res, err := s.keeper.cdc.MarshalJSON(GetQuery{service})
			if err != nil {
				return nil, sdk.ErrInternal(err.Error())
			}
			return res, nil
		default:
			return nil, sdk.ErrUnknownRequest("unknown nameservice query endpoint")
		}
	}
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

	// check if service already exists.
	if _, err := s.serviceDB.Get(srv.Hash); err == nil {
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

	return srv, s.serviceDB.Save(srv)
}

// Delete deletes the service by hash.
func (s *Service) Delete(hash hash.Hash) error {
	return s.serviceDB.Delete(hash)
}

// Get returns the service that matches given hash.
func (s *Service) Get(hash hash.Hash) (*service.Service, error) {
	result, err := s.abciClient.ABCIQuery("custom/service/get/"+hash.String(), nil)
	if err != nil {
		return nil, err
	}

	resp := result.Response
	if !resp.IsOK() {
		return nil, errors.New(resp.Log)
	}

	var out GetQuery
	if err := s.cdc.UnmarshalJSON(resp.Value, &out); err != nil {
		return nil, err
	}

	return &out.Service, err
}

func (s *Service) querierGet(ctx sdk.Context, path []string) (service.Service, error) {
	srv := s.keeper.Get(ctx, path[0])
	if srv.Hash.IsZero() {
		return service.Service{}, sdk.ErrUnknownRequest("could not resolve service")
	}
	return srv, nil
}

// func (s *Service) querierGet(hash hash.Hash) (*service.Service, error) {
// 	return s.serviceDB.Get(hash)
// }

// List returns all services.
func (s *Service) List() ([]*service.Service, error) {
	return s.serviceDB.All()
}

// AlreadyExistsError is an not found error.
type AlreadyExistsError struct {
	Hash hash.Hash
}

func (e *AlreadyExistsError) Error() string {
	return fmt.Sprintf("service %q already exists", e.Hash.String())
}
