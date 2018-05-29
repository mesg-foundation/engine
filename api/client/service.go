package client

import (
	"context"

	"github.com/mesg-foundation/core/api/core"
)

func startServices(wf *Workflow) {
	startService(wf.OnEvent.Service)
	startService(wf.Execute.Service)
}

func stopServices(wf *Workflow) {
	stopService(wf.OnEvent.Service)
	stopService(wf.Execute.Service)
}

func startService(ID string) {
	_, err := getClient().StartService(context.Background(), &core.StartServiceRequest{
		ServiceID: ID,
	})
	if err != nil {
		panic(err)
	}
}

func stopService(ID string) {
	getClient().StopService(context.Background(), &core.StopServiceRequest{
		ServiceID: ID,
	})
}
