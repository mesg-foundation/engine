package client

import (
	"context"

	"github.com/mesg-foundation/core/api/core"
)

func startServices(wf *Workflow) {
	if wf.OnEvent != nil {
		startService(wf.OnEvent.Service)
	}
	if wf.OnResult != nil {
		startService(wf.OnResult.Service)
	}
	startService(wf.Execute.Service)
}

func stopServices(wf *Workflow) {
	if wf.OnEvent != nil {
		stopService(wf.OnEvent.Service)
	}
	if wf.OnResult != nil {
		stopService(wf.OnResult.Service)
	}
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
