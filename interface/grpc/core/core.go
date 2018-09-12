package core

import (
	service "github.com/mesg-foundation/core/service"
)

func toProtoServices(ss []*service.Service) []*Service {
	services := make([]*Service, 0)
	for _, s := range ss {
		services = append(services, toProtoService(s))
	}
	return services
}

func toProtoService(s *service.Service) *Service {
	return &Service{
		ID:           s.ID,
		Name:         s.Name,
		Description:  s.Description,
		Repository:   s.Repository,
		Tasks:        toProtoTasks(s.Tasks),
		Events:       toProtoEvents(s.Events),
		Dependencies: toProtoDependencies(s.Dependencies),
	}
}

func toProtoTasks(tasks []*service.Task) []*Task {
	ts := make([]*Task, 0)
	for _, task := range tasks {
		t := &Task{
			Key:         task.Key,
			Name:        task.Name,
			Description: task.Description,
			Inputs:      toProtoParameters(task.Inputs),
			Outputs:     []*Output{},
		}
		for _, output := range task.Outputs {
			o := &Output{
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

func toProtoEvents(events []*service.Event) []*Event {
	es := make([]*Event, 0)
	for _, event := range events {
		e := &Event{
			Key:         event.Key,
			Name:        event.Name,
			Description: event.Description,
			Data:        toProtoParameters(event.Data),
		}
		es = append(es, e)
	}
	return es
}

func toProtoParameters(params []*service.Parameter) []*Parameter {
	ps := make([]*Parameter, 0)
	for _, param := range params {
		p := &Parameter{
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

func toProtoDependency(dep *service.Dependency) *Dependency {
	if dep == nil {
		return nil
	}
	return &Dependency{
		Key:         dep.Key,
		Image:       dep.Image,
		Volumes:     dep.Volumes,
		Volumesfrom: dep.VolumesFrom,
		Ports:       dep.Ports,
		Command:     dep.Command,
	}
}

func toProtoDependencies(deps []*service.Dependency) []*Dependency {
	ds := make([]*Dependency, 0)
	for _, dep := range deps {
		ds = append(ds, toProtoDependency(dep))
	}
	return ds
}
