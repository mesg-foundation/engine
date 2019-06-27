package api

import (
	"github.com/mesg-foundation/core/protobuf/definition"
	"github.com/mesg-foundation/core/service"
)

// fromProtoService converts a the protobuf definition to the internal service struct
func fromProtoService(s *definition.Service) *service.Service {
	return &service.Service{
		Sid:           s.Sid,
		Name:          s.Name,
		Description:   s.Description,
		Repository:    s.Repository,
		Source:        s.Source,
		Tasks:         fromProtoTasks(s.Tasks),
		Events:        fromProtoEvents(s.Events),
		Configuration: fromProtoConfiguration(s.Configuration),
		Dependencies:  fromProtoDependencies(s.Dependencies),
	}
}

func fromProtoTasks(tasks []*definition.Service_Task) []*service.Task {
	ts := make([]*service.Task, len(tasks))
	for i, task := range tasks {
		ts[i] = &service.Task{
			Key:         task.Key,
			Name:        task.Name,
			Description: task.Description,
			Inputs:      fromProtoParameters(task.Inputs),
			Outputs:     fromProtoParameters(task.Outputs),
		}
	}
	return ts
}

func fromProtoEvents(events []*definition.Service_Event) []*service.Event {
	es := make([]*service.Event, len(events))
	for i, event := range events {
		es[i] = &service.Event{
			Key:         event.Key,
			Name:        event.Name,
			Description: event.Description,
			Data:        fromProtoParameters(event.Data),
		}
	}
	return es
}

func fromProtoParameters(params []*definition.Service_Parameter) []*service.Parameter {
	ps := make([]*service.Parameter, len(params))
	for i, param := range params {
		ps[i] = &service.Parameter{
			Key:         param.Key,
			Name:        param.Name,
			Description: param.Description,
			Type:        param.Type,
			Repeated:    param.Repeated,
			Optional:    param.Optional,
			Object:      fromProtoParameters(param.Object),
		}
	}
	return ps
}

func fromProtoConfiguration(configuration *definition.Service_Configuration) *service.Dependency {
	if configuration == nil {
		return nil
	}
	return &service.Dependency{
		Key:         service.MainServiceKey,
		Args:        configuration.Args,
		Command:     configuration.Command,
		Ports:       configuration.Ports,
		Volumes:     configuration.Volumes,
		VolumesFrom: configuration.VolumesFrom,
		Env:         configuration.Env,
	}
}

func fromProtoDependency(dep *definition.Service_Dependency) *service.Dependency {
	return &service.Dependency{
		Key:         dep.Key,
		Image:       dep.Image,
		Volumes:     dep.Volumes,
		VolumesFrom: dep.VolumesFrom,
		Ports:       dep.Ports,
		Command:     dep.Command,
		Args:        dep.Args,
		Env:         dep.Env,
	}
}

func fromProtoDependencies(deps []*definition.Service_Dependency) []*service.Dependency {
	ds := make([]*service.Dependency, len(deps))
	for i, dep := range deps {
		ds[i] = fromProtoDependency(dep)
	}
	return ds
}

// toProtoService converts an internal service struct to the protobuf definition
func toProtoService(s *service.Service) *definition.Service {
	return &definition.Service{
		Hash:          s.Hash.String(),
		Sid:           s.Sid,
		Name:          s.Name,
		Description:   s.Description,
		Repository:    s.Repository,
		Source:        s.Source,
		Tasks:         toProtoTasks(s.Tasks),
		Events:        toProtoEvents(s.Events),
		Configuration: toProtoConfiguration(s.Configuration),
		Dependencies:  toProtoDependencies(s.Dependencies),
	}
}

func toProtoTasks(tasks []*service.Task) []*definition.Service_Task {
	ts := make([]*definition.Service_Task, len(tasks))
	for i, task := range tasks {
		ts[i] = &definition.Service_Task{
			Key:         task.Key,
			Name:        task.Name,
			Description: task.Description,
			Inputs:      toProtoParameters(task.Inputs),
			Outputs:     toProtoParameters(task.Outputs),
		}
	}
	return ts
}

func toProtoEvents(events []*service.Event) []*definition.Service_Event {
	es := make([]*definition.Service_Event, len(events))
	for i, event := range events {
		es[i] = &definition.Service_Event{
			Key:         event.Key,
			Name:        event.Name,
			Description: event.Description,
			Data:        toProtoParameters(event.Data),
		}
	}
	return es
}

func toProtoParameters(params []*service.Parameter) []*definition.Service_Parameter {
	ps := make([]*definition.Service_Parameter, len(params))
	for i, param := range params {
		ps[i] = &definition.Service_Parameter{
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

func toProtoConfiguration(configuration *service.Dependency) *definition.Service_Configuration {
	return &definition.Service_Configuration{
		Args:        configuration.Args,
		Command:     configuration.Command,
		Ports:       configuration.Ports,
		Volumes:     configuration.Volumes,
		VolumesFrom: configuration.VolumesFrom,
		Env:         configuration.Env,
	}
}

func toProtoDependency(dep *service.Dependency) *definition.Service_Dependency {
	return &definition.Service_Dependency{
		Key:         dep.Key,
		Image:       dep.Image,
		Volumes:     dep.Volumes,
		VolumesFrom: dep.VolumesFrom,
		Ports:       dep.Ports,
		Command:     dep.Command,
		Args:        dep.Args,
		Env:         dep.Env,
	}
}

func toProtoDependencies(deps []*service.Dependency) []*definition.Service_Dependency {
	ds := make([]*definition.Service_Dependency, len(deps))
	for i, dep := range deps {
		ds[i] = toProtoDependency(dep)
	}
	return ds
}
