package core

import (
	"github.com/mesg-foundation/core/protobuf/coreapi"
	service "github.com/mesg-foundation/core/service"
)

func toProtoService(s *service.Service) *coreapi.Service {
	return &coreapi.Service{
		Hash:          s.Hash,
		Sid:           s.Sid,
		Name:          s.Name,
		Description:   s.Description,
		Repository:    s.Repository,
		Tasks:         toProtoTasks(s.Tasks),
		Events:        toProtoEvents(s.Events),
		Status:        toProtoServiceStatusType(s.Status),
		Configuration: toProtoDependency(s.Configuration),
		Dependencies:  toProtoDependencies(s.Dependencies),
	}
}

func toProtoServiceStatusType(s service.Status) coreapi.Service_Status {
	switch s {
	case service.StatusStopped:
		return coreapi.Service_STOPPED
	case service.StatusStarting:
		return coreapi.Service_STARTING
	case service.StatusPartial:
		return coreapi.Service_PARTIAL
	case service.StatusRunning:
		return coreapi.Service_RUNNING
	default:
		return coreapi.Service_UNKNOWN
	}
}

func toProtoTasks(tasks []*service.Task) []*coreapi.Task {
	ts := make([]*coreapi.Task, len(tasks))
	for i, task := range tasks {
		ts[i] = &coreapi.Task{
			Name:        task.Name,
			Description: task.Description,
			Inputs:      toProtoParameters(task.Inputs),
			Outputs:     make([]*coreapi.Output, 0),
		}
		for j, output := range task.Outputs {
			ts[i].Outputs[j] = &coreapi.Output{
				Name:        output.Name,
				Description: output.Description,
				Data:        toProtoParameters(output.Data),
			}
		}
	}
	return ts
}

func toProtoEvents(events []*service.Event) []*coreapi.Event {
	es := make([]*coreapi.Event, len(events))
	for i, event := range events {
		es[i] = &coreapi.Event{
			Name:        event.Name,
			Description: event.Description,
			Data:        toProtoParameters(event.Data),
		}
	}
	return es
}

func toProtoParameters(params []*service.Parameter) []*coreapi.Parameter {
	ps := make([]*coreapi.Parameter, len(params))
	for i, param := range params {
		ps[i] = &coreapi.Parameter{
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

func toProtoDependency(dep *service.Dependency) *coreapi.Dependency {
	if dep == nil {
		return nil
	}
	return &coreapi.Dependency{
		Image:       dep.Image,
		Volumes:     dep.Volumes,
		Volumesfrom: dep.VolumesFrom,
		Ports:       dep.Ports,
		Command:     dep.Command,
	}
}

func toProtoDependencies(deps []*service.Dependency) []*coreapi.Dependency {
	ds := make([]*coreapi.Dependency, len(deps))
	for i, dep := range deps {
		ds[i] = toProtoDependency(dep)
	}
	return ds
}
