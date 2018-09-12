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
		ID:            s.ID,
		Name:          s.Name,
		Description:   s.Description,
		Repository:    s.Repository,
		Tasks:         toProtoTasks(s.Tasks),
		Events:        toProtoEvents(s.Events),
		Dependencies:  toProtoDependencies(s.Dependencies),
		Configuration: toProtoDependency(s.Configuration),
	}
}

func toProtoTasks(tasks map[string]*service.Task) map[string]*Task {
	ts := make(map[string]*Task, 0)
	for taskKey, task := range tasks {
		t := &Task{
			Key:         task.Key,
			Name:        task.Name,
			Description: task.Description,
			ServiceName: task.ServiceName,
			Inputs:      toProtoParameters(task.Inputs),
			Outputs:     map[string]*Output{},
		}
		for outputKey, output := range task.Outputs {
			t.Outputs[outputKey] = &Output{
				Key:         output.Key,
				Name:        output.Name,
				Description: output.Description,
				TaskKey:     output.TaskKey,
				ServiceName: output.ServiceName,
				Data:        toProtoParameters(output.Data),
			}
		}
		ts[taskKey] = t
	}
	return ts
}

func toProtoEvents(events map[string]*service.Event) map[string]*Event {
	es := make(map[string]*Event, 0)
	for eventKey, event := range events {
		es[eventKey] = &Event{
			Key:         event.Key,
			Name:        event.Name,
			Description: event.Description,
			ServiceName: event.ServiceName,
			Data:        toProtoParameters(event.Data),
		}
	}
	return es
}

func toProtoParameters(params map[string]*service.Parameter) map[string]*Parameter {
	ps := make(map[string]*Parameter, 0)
	for key, param := range params {
		ps[key] = &Parameter{
			Name:        param.Name,
			Description: param.Description,
			Type:        param.Type,
			Optional:    param.Optional,
		}
	}
	return ps
}

func toProtoDependency(dep *service.Dependency) *Dependency {
	if dep == nil {
		return nil
	}
	return &Dependency{
		Image:       dep.Image,
		Volumes:     dep.Volumes,
		Volumesfrom: dep.VolumesFrom,
		Ports:       dep.Ports,
		Command:     dep.Command,
	}
}

func toProtoDependencies(deps map[string]*service.Dependency) map[string]*Dependency {
	ds := make(map[string]*Dependency, 0)
	for key, dep := range deps {
		ds[key] = toProtoDependency(dep)
	}
	return ds
}
