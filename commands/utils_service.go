package commands

import (
	"github.com/mesg-foundation/core/interface/grpc/core"
	"github.com/mesg-foundation/core/service"
)

// TODO(ilgooz) rm this when we stop using internal methods of service in cmd.
func toServices(ss []*core.Service) []*service.Service {
	services := make([]*service.Service, 0)
	for _, s := range ss {
		services = append(services, toService(s))
	}
	return services
}

// TODO(ilgooz) rm this when we stop using internal methods of service in cmd.
func toService(s *core.Service) *service.Service {
	return &service.Service{
		ID:           s.ID,
		Name:         s.Name,
		Description:  s.Description,
		Repository:   s.Repository,
		Tasks:        toTasks(s.Tasks),
		Events:       toEvents(s.Events),
		Dependencies: toDependencies(s.Dependencies),
	}
}

func toTasks(tasks []*core.Task) []*service.Task {
	ts := make([]*service.Task, 0, len(tasks))
	for _, task := range tasks {
		t := &service.Task{
			Key:         task.Key,
			Name:        task.Name,
			Description: task.Description,
			Inputs:      toParameters(task.Inputs),
			Outputs:     []*service.Output{},
		}
		for _, output := range task.Outputs {
			o := &service.Output{
				Key:         output.Key,
				Name:        output.Name,
				Description: output.Description,
				Data:        toParameters(output.Data),
			}
			t.Outputs = append(t.Outputs, o)
		}
		ts = append(ts, t)
	}
	return ts
}

func toEvents(events []*core.Event) []*service.Event {
	es := make([]*service.Event, 0, len(events))
	for eventKey, event := range events {
		es[eventKey] = &service.Event{
			Key:         event.Key,
			Name:        event.Name,
			Description: event.Description,
			Data:        toParameters(event.Data),
		}
	}
	return es
}

// TODO(ilgooz) rm this when we stop using internal methods of service in cmd.
func toParameters(params []*core.Parameter) []*service.Parameter {
	gParams := make([]*service.Parameter, 0, len(params))
	for _, param := range params {
		p := &service.Parameter{
			Name:        param.Name,
			Description: param.Description,
			Type:        param.Type,
			Optional:    param.Optional,
		}
		gParams = append(gParams, p)
	}
	return gParams
}

func toDependency(dep *core.Dependency) *service.Dependency {
	if dep == nil {
		return nil
	}
	return &service.Dependency{
		Image:       dep.Image,
		Volumes:     dep.Volumes,
		VolumesFrom: dep.Volumesfrom,
		Ports:       dep.Ports,
		Command:     dep.Command,
	}
}

func toDependencies(deps []*core.Dependency) []*service.Dependency {
	ds := make([]*service.Dependency, len(deps))
	for key, dep := range deps {
		ds[key] = toDependency(dep)
	}
	return ds
}
