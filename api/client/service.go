package client

import (
	"context"

	"github.com/mesg-foundation/core/api/core"
)

func (wf *Workflow) services() []string {
	var (
		services []string
		presence = make(map[string]bool)
	)
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
	return services
}

func startServices(wf *Workflow) error {
	for _, ID := range wf.services() {
		if _, err := wf.client.StartService(context.Background(),
			&core.StartServiceRequest{ServiceID: ID}); err != nil {
			return err
		}
	}
	return nil
}

func stopServices(wf *Workflow) error {
	for _, ID := range wf.services() {
		if _, err := wf.client.StopService(context.Background(),
			&core.StopServiceRequest{ServiceID: ID}); err != nil {
			return err
		}
	}
	return nil
}
