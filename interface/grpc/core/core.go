package core

import (
	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/mesg-foundation/core/protobuf/definitions"
	service "github.com/mesg-foundation/core/service"
)

func toProtoServices(ss []*service.Service) []*definitions.Service {
	services := make([]*definitions.Service, len(ss))
	for i, s := range ss {
		services[i] = toProtoService(s)
	}
	return services
}

func toProtoService(s *service.Service) *definitions.Service {
	return &definitions.Service{
		Hash:         s.Hash,
		Sid:          s.Sid,
		Name:         s.Name,
		Description:  s.Description,
		Repository:   s.Repository,
		Tasks:        toProtoTasks(s.Tasks),
		Events:       toProtoEvents(s.Events),
		Dependencies: toProtoDependencies(s.Dependencies),
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

func toProtoTasks(tasks []*service.Task) []*definitions.Task {
	ts := make([]*definitions.Task, len(tasks))
	for i, task := range tasks {
		t := &definitions.Task{
			Key:         task.Key,
			Name:        task.Name,
			Description: task.Description,
			Inputs:      toProtoParameters(task.Inputs),
			Outputs:     []*definitions.Output{},
		}
		for _, output := range task.Outputs {
			o := &definitions.Output{
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

func toProtoEvents(events []*service.Event) []*definitions.Event {
	es := make([]*definitions.Event, len(events))
	for i, event := range events {
		es[i] = &definitions.Event{
			Key:         event.Key,
			Name:        event.Name,
			Description: event.Description,
			Data:        toProtoParameters(event.Data),
		}
	}
	return es
}

func toProtoParameters(params []*service.Parameter) []*definitions.Parameter {
	ps := make([]*definitions.Parameter, len(params))
	for i, param := range params {
		ps[i] = &definitions.Parameter{
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

func toProtoDependency(dep *service.Dependency) *definitions.Dependency {
	if dep == nil {
		return nil
	}
	return &definitions.Dependency{
		Key:         dep.Key,
		Image:       dep.Image,
		Volumes:     dep.Volumes,
		Volumesfrom: dep.VolumesFrom,
		Ports:       dep.Ports,
		Command:     dep.Command,
	}
}

func toProtoDependencies(deps []*service.Dependency) []*definitions.Dependency {
	ds := make([]*definitions.Dependency, len(deps))
	for i, dep := range deps {
		ds[i] = toProtoDependency(dep)
	}
	return ds
}
