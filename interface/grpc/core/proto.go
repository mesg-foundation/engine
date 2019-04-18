package core

import (
	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/mesg-foundation/core/protobuf/definition"
	service "github.com/mesg-foundation/core/service"
)

func toProtoServices(ss []*service.Service) []*definition.Service {
	services := make([]*definition.Service, len(ss))
	for i, s := range ss {
		services[i] = toProtoService(s)
	}
	return services
}

func toProtoService(s *service.Service) *definition.Service {
	return &definition.Service{
		Hash:          s.Hash,
		Sid:           s.Sid,
		Name:          s.Name,
		HashVersion:   s.HashVersion,
		Description:   s.Description,
		Repository:    s.Repository,
		Tasks:         toProtoTasks(s.Tasks),
		Events:        toProtoEvents(s.Events),
		Configuration: toProtoConfiguration(s.Configuration),
		Dependencies:  toProtoDependencies(s.Dependencies),
	}
}

func toProtoServiceStatusType(s service.StatusType) coreapi.Service_Status {
	switch s {
	default:
		return coreapi.Service_UNKNOWN
	case service.STOPPED:
		return coreapi.Service_STOPPED
	case service.STARTING:
		return coreapi.Service_STARTING
	case service.PARTIAL:
		return coreapi.Service_PARTIAL
	case service.RUNNING:
		return coreapi.Service_RUNNING
	}
}

func toProtoTasks(tasks []*service.Task) []*definition.Task {
	ts := make([]*definition.Task, len(tasks))
	for i, task := range tasks {
		t := &definition.Task{
			Key:         task.Key,
			Name:        task.Name,
			Description: task.Description,
			Inputs:      toProtoParameters(task.Inputs),
			Outputs:     []*definition.Output{},
		}
		for _, output := range task.Outputs {
			o := &definition.Output{
				Key:         output.Key,
				Name:        output.Name,
				Description: output.Description,
				Data:        toProtoParameters(output.Data),
			}
			t.Outputs = append(t.Outputs, o)
		}
		ts[i] = t
	}
	return ts
}

func toProtoEvents(events []*service.Event) []*definition.Event {
	es := make([]*definition.Event, len(events))
	for i, event := range events {
		es[i] = &definition.Event{
			Key:         event.Key,
			Name:        event.Name,
			Description: event.Description,
			Data:        toProtoParameters(event.Data),
		}
	}
	return es
}

func toProtoParameters(params []*service.Parameter) []*definition.Parameter {
	ps := make([]*definition.Parameter, len(params))
	for i, param := range params {
		ps[i] = &definition.Parameter{
			Key:         param.Key,
			Name:        param.Name,
			Description: param.Description,
			Type:        param.Type,
			Repeated:    param.Repeated,
			Optional:    param.Optional,
			Object:      toProtoParameters(param.Object),
		}
	}
	return ps
}

func toProtoConfiguration(configuration *service.Dependency) *definition.Configuration {
	if configuration == nil {
		return nil
	}
	return &definition.Configuration{
		Args:        configuration.Args,
		Command:     configuration.Command,
		Ports:       configuration.Ports,
		Volumes:     configuration.Volumes,
		VolumesFrom: configuration.VolumesFrom,
	}
}

func toProtoDependency(dep *service.Dependency) *definition.Dependency {
	if dep == nil {
		return nil
	}
	return &definition.Dependency{
		Key:         dep.Key,
		Image:       dep.Image,
		Volumes:     dep.Volumes,
		VolumesFrom: dep.VolumesFrom,
		Ports:       dep.Ports,
		Command:     dep.Command,
		Args:        dep.Args,
	}
}

func toProtoDependencies(deps []*service.Dependency) []*definition.Dependency {
	ds := make([]*definition.Dependency, len(deps))
	for i, dep := range deps {
		ds[i] = toProtoDependency(dep)
	}
	return ds
}
