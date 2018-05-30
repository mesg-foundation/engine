package client

import (
	"context"

	"github.com/mesg-foundation/core/api/core"
)

func startServices(wf *Workflow) {
	if wf.OnEvent != nil {
		startService(wf.OnEvent.ServiceID)
	}
	if wf.OnResult != nil {
		startService(wf.OnResult.ServiceID)
	}
	startService(wf.Execute.ServiceID)
}

func stopServices(wf *Workflow) {
	if wf.OnEvent != nil {
		stopService(wf.OnEvent.ServiceID)
	}
	if wf.OnResult != nil {
		stopService(wf.OnResult.ServiceID)
	}
	stopService(wf.Execute.ServiceID)
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
