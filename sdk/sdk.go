package sdk

import (
	"github.com/cskr/pubsub"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/database"
	eventsdk "github.com/mesg-foundation/core/sdk/event"
	executionsdk "github.com/mesg-foundation/core/sdk/execution"
	instancesdk "github.com/mesg-foundation/core/sdk/instance"
	servicesdk "github.com/mesg-foundation/core/sdk/service"
)

// SDK exposes all functionalities of MESG core.
type SDK struct {
	Service   *servicesdk.Service
	Instance  *instancesdk.Instance
	Execution *executionsdk.Execution
	Event     *eventsdk.Event

	ps *pubsub.PubSub

	container container.Container
	db        database.ServiceDB
	execDB    database.ExecutionDB
}

// New creates a new SDK with given options.
func New(c container.Container, serviceDB database.ServiceDB, instanceDB database.InstanceDB, execDB database.ExecutionDB) *SDK {
	ps := pubsub.New(0)
	return &SDK{
		Service:   servicesdk.New(c, db, execDB),
		Instance:  instancesdk.New(c, db, instanceDB),
		Execution: executionsdk.New(ps, db, execDB, instanceDB),
		Event:     eventsdk.New(ps, db),
		ps:        ps,
		container: c,
		db:        db,
		execDB:    execDB,
	}
}
