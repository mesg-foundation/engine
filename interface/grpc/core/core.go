package core

import (
	"github.com/mesg-foundation/core/protobuf/core"
	service "github.com/mesg-foundation/core/service"
)

func toProtoServices(ss []*service.Service) []*core.Service {
	services := make([]*core.Service, 0)
	for _, s := range ss {
		services = append(services, toProtoService(s))
	}
	return services
}

func toProtoService(s *service.Service) *core.Service {
	return &core.Service{
		ID:           s.ID,
		Name:         s.Name,
		Description:  s.Description,
		Repository:   s.Repository,
		Tasks:        toProtoTasks(s.Tasks),
		Events:       toProtoEvents(s.Events),
		Dependencies: toProtoDependencies(s.Dependencies),
	}
}

func toProtoTasks(tasks []*service.Task) []*core.Task {
	ts := make([]*core.Task, 0)
	for _, task := range tasks {
		t := &core.Task{
			Key:         task.Key,
			Name:        task.Name,
			Description: task.Description,
			Inputs:      toProtoParameters(task.Inputs),
			Outputs:     []*core.Output{},
		}
		for _, output := range task.Outputs {
			o := &core.Output{
				Key:         output.Key,
				Name:        output.Name,
				Description: output.Description,
				Data:        toProtoParameters(output.Data),
			}
			t.Outputs = append(t.Outputs, o)
		}
		ts = append(ts, t)
	}
	return ts
}

func toProtoEvents(events []*service.Event) []*core.Event {
	es := make([]*core.Event, 0)
	for _, event := range events {
		e := &core.Event{
			Key:         event.Key,
			Name:        event.Name,
			Description: event.Description,
			Data:        toProtoParameters(event.Data),
		}
		es = append(es, e)
	}
	return es
}

func toProtoParameters(params []*service.Parameter) []*core.Parameter {
	ps := make([]*core.Parameter, 0)
	for _, param := range params {
		p := &core.Parameter{
			Key:         param.Key,
			Name:        param.Name,
			Description: param.Description,
			Type:        param.Type,
			Optional:    param.Optional,
		}
		ps = append(ps, p)
	}
	return ps
}

func toProtoDependency(dep *service.Dependency) *core.Dependency {
	if dep == nil {
		return nil
	}
	return &core.Dependency{
		Key:         dep.Key,
		Image:       dep.Image,
		Volumes:     dep.Volumes,
		Volumesfrom: dep.VolumesFrom,
		Ports:       dep.Ports,
		Command:     dep.Command,
	}
}

func toProtoDependencies(deps []*service.Dependency) []*core.Dependency {
	ds := make([]*core.Dependency, 0)
	for _, dep := range deps {
		ds = append(ds, toProtoDependency(dep))
	}
	return ds
}
