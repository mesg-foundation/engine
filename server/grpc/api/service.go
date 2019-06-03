package api

import (
	"github.com/mesg-foundation/core/protobuf/definition"
	"github.com/mesg-foundation/core/service"
)

func FromProtoService(s *definition.Service) *service.Service {
	return &service.Service{
		Hash:          s.Hash,
		Sid:           s.Sid,
		Name:          s.Name,
		Description:   s.Description,
		Repository:    s.Repository,
		Source:        s.Source,
		Tasks:         FromProtoTasks(s.Tasks),
		Events:        FromProtoEvents(s.Events),
		Configuration: FromProtoConfiguration(s.Configuration),
		Dependencies:  FromProtoDependencies(s.Dependencies),
	}
}

func FromProtoTasks(tasks []*definition.Task) []*service.Task {
	ts := make([]*service.Task, len(tasks))
	for i, task := range tasks {
		ts[i] = &service.Task{
			Key:         task.Key,
			Name:        task.Name,
			Description: task.Description,
			Inputs:      FromProtoParameters(task.Inputs),
			Outputs:     FromProtoParameters(task.Outputs),
		}
	}
	return ts
}

func FromProtoEvents(events []*definition.Event) []*service.Event {
	es := make([]*service.Event, len(events))
	for i, event := range events {
		es[i] = &service.Event{
			Key:         event.Key,
			Name:        event.Name,
			Description: event.Description,
			Data:        FromProtoParameters(event.Data),
		}
	}
	return es
}

func FromProtoParameters(params []*definition.Parameter) []*service.Parameter {
	ps := make([]*service.Parameter, len(params))
	for i, param := range params {
		ps[i] = &service.Parameter{
			Key:         param.Key,
			Name:        param.Name,
			Description: param.Description,
			Type:        param.Type,
			Repeated:    param.Repeated,
			Optional:    param.Optional,
			Object:      FromProtoParameters(param.Object),
		}
	}
	return ps
}

func FromProtoConfiguration(configuration *definition.Configuration) *service.Dependency {
	if configuration == nil {
		return nil
	}
	return &service.Dependency{
		Args:        configuration.Args,
		Command:     configuration.Command,
		Ports:       configuration.Ports,
		Volumes:     configuration.Volumes,
		VolumesFrom: configuration.VolumesFrom,
	}
}

func FromProtoDependency(dep *definition.Dependency) *service.Dependency {
	if dep == nil {
		return nil
	}
	return &service.Dependency{
		Key:         dep.Key,
		Image:       dep.Image,
		Volumes:     dep.Volumes,
		VolumesFrom: dep.VolumesFrom,
		Ports:       dep.Ports,
		Command:     dep.Command,
		Args:        dep.Args,
	}
}

func FromProtoDependencies(deps []*definition.Dependency) []*service.Dependency {
	ds := make([]*service.Dependency, len(deps))
	for i, dep := range deps {
		ds[i] = FromProtoDependency(dep)
	}
	return ds
}
