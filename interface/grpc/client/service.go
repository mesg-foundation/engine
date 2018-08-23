package client

import (
	"context"

	"github.com/mesg-foundation/core/interface/grpc/core"
)

func (wf *Workflow) services() (services []string) {
	presence := make(map[string]bool)
	if wf.OnEvent != nil && !presence[wf.OnEvent.ServiceID] {
		services = append(services, wf.OnEvent.ServiceID)
		presence[wf.OnEvent.ServiceID] = true
	}
	if wf.OnResult != nil && !presence[wf.OnResult.ServiceID] {
		services = append(services, wf.OnResult.ServiceID)
		presence[wf.OnResult.ServiceID] = true
	}
	if wf.Execute != nil && !presence[wf.Execute.ServiceID] {
		services = append(services, wf.Execute.ServiceID)
		presence[wf.Execute.ServiceID] = true
	}
	return
}

func iterateService(wf *Workflow, action func(string) error) (err error) {
	for _, ID := range wf.services() {
		err = action(ID)
		if err != nil {
			break
		}
	}
	return
}

func startServices(wf *Workflow) error {
	return iterateService(wf, func(ID string) (err error) {
		_, err = wf.client.StartService(context.Background(), &core.StartServiceRequest{
			ServiceID: ID,
		})
		return
	})
}

func stopServices(wf *Workflow) error {
	return iterateService(wf, func(ID string) (err error) {
		_, err = wf.client.StopService(context.Background(), &core.StopServiceRequest{
			ServiceID: ID,
		})
		return
	})
}
