package sdk

import (
	"errors"

	"github.com/cskr/pubsub"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/database"
	eventsdk "github.com/mesg-foundation/core/sdk/event"
	executionsdk "github.com/mesg-foundation/core/sdk/execution"
	instancesdk "github.com/mesg-foundation/core/sdk/instance"
	servicesdk "github.com/mesg-foundation/core/sdk/service"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/service/manager"
)

// SDK exposes all functionalities of MESG core.
type SDK struct {
	Service   *servicesdk.Service
	Instance  *instancesdk.Instance
	Execution *executionsdk.Execution
	Event     *eventsdk.Event

	ps *pubsub.PubSub

	m         manager.Manager
	container container.Container
	db        database.ServiceDB
	execDB    database.ExecutionDB
}

// New creates a new SDK with given options.
func New(m manager.Manager, c container.Container, db database.ServiceDB, instanceDB database.InstanceDB, execDB database.ExecutionDB) *SDK {
	ps := pubsub.New(0)
	return &SDK{
		Service:   servicesdk.New(m, c, db, execDB),
		Instance:  instancesdk.New(c, db, instanceDB),
		Execution: executionsdk.New(m, ps, db, execDB),
		Event:     eventsdk.New(ps, db),
		ps:        ps,
		m:         m,
		container: c,
		db:        db,
		execDB:    execDB,
	}
}

// GetService returns service serviceID.
func (sdk *SDK) GetService(serviceID string) (*service.Service, error) {
	return sdk.db.Get(serviceID)
}

// ListServices returns all services.
func (sdk *SDK) ListServices() ([]*service.Service, error) {
	return sdk.db.All()
}

// Status returns the status of a service
func (sdk *SDK) Status(service *service.Service) (service.StatusType, error) {
	return sdk.m.Status(service)
}

// StartService starts service serviceID.
func (sdk *SDK) StartService(serviceID string) error {
	s, err := sdk.db.Get(serviceID)
	if err != nil {
		return err
	}
	_, err = sdk.m.Start(s)
	return err
}

// StopService stops service serviceID.
func (sdk *SDK) StopService(serviceID string) error {
	s, err := sdk.db.Get(serviceID)
	if err != nil {
		return err
	}
	return sdk.m.Stop(s)
}

// DeleteService stops and deletes service serviceID.
// when deleteData is enabled, any persistent data that belongs to
// the service and to its dependencies also will be deleted.
func (sdk *SDK) DeleteService(serviceID string, deleteData bool) error {
	s, err := sdk.db.Get(serviceID)
	if err != nil {
		return err
	}

	// check if service has any running instances.
	instances, err := sdk.Instance.GetAllByService(s.Hash)
	if err != nil {
		return err
	}
	if len(instances) > 0 {
		return errors.New("service has running instances. in order to delete the service, stop its instances first")
	}

	if err := sdk.m.Stop(s); err != nil {
		return err
	}
	// delete volumes first before the service. this way if
	// deleting volumes fails, process can be retried by the user again
	// because service still will be in the db.
	if deleteData {
		if err := sdk.m.Delete(s); err != nil {
			return err
		}
	}
	return sdk.db.Delete(serviceID)
}
