package api

import (
	"github.com/mesg-foundation/core/protobuf/definition"
	"github.com/mesg-foundation/core/service"
	"github.com/mr-tron/base58"
)

// FromProtoService converts a the protobuf definition to the internal service struct
// TODO: should not be public. Need to move server/grpc/service.go to server/grpc/api/service.go
func FromProtoService(s *definition.Service) (*service.Service, error) {
	hash, err := base58.Decode(s.Hash)
	if err != nil {
		return nil, err
	}

	return &service.Service{
		Hash:          hash,
		Sid:           s.Sid,
		Name:          s.Name,
		Description:   s.Description,
		Repository:    s.Repository,
		Source:        s.Source,
		Tasks:         fromProtoTasks(s.Tasks),
		Events:        fromProtoEvents(s.Events),
		Configuration: fromProtoConfiguration(s.Configuration),
		Dependencies:  fromProtoDependencies(s.Dependencies),
	}, nil
}

func fromProtoTasks(tasks []*definition.Task) []*service.Task {
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

func fromProtoEvents(events []*definition.Event) []*service.Event {
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

func fromProtoParameters(params []*definition.Parameter) []*service.Parameter {
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

func fromProtoConfiguration(configuration *definition.Configuration) *service.Dependency {
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

func fromProtoDependency(dep *definition.Dependency) *service.Dependency {
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
		Env:         dep.Env,
	}
}

func fromProtoDependencies(deps []*definition.Dependency) []*service.Dependency {
	ds := make([]*service.Dependency, len(deps))
	for i, dep := range deps {
		ds[i] = fromProtoDependency(dep)
	}
	return ds
}

// ToProtoServices converts internal services struct to their protobuf definition
// TODO: should not be public. Need to move server/grpc/service.go to server/grpc/api/service.go and delete server/grpc/core package
func ToProtoServices(ss []*service.Service) []*definition.Service {
	services := make([]*definition.Service, len(ss))
	for i, s := range ss {
		services[i] = ToProtoService(s)
	}
	return services
}

// ToProtoService converts an internal service struct to the protobuf definition
// TODO: should not be public. Need to move server/grpc/service.go to server/grpc/api/service.go and delete server/grpc/core package
func ToProtoService(s *service.Service) *definition.Service {
	return &definition.Service{
		Hash:          s.Hash,
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

func toProtoTasks(tasks []*service.Task) []*definition.Task {
	ts := make([]*definition.Task, len(tasks))
	for i, task := range tasks {
		ts[i] = &definition.Task{
			Key:         task.Key,
			Name:        task.Name,
			Description: task.Description,
			Inputs:      toProtoParameters(task.Inputs),
			Outputs:     toProtoParameters(task.Outputs),
		}
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
		Env:         configuration.Env,
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
		Env:         dep.Env,
	}
}

func toProtoDependencies(deps []*service.Dependency) []*definition.Dependency {
	ds := make([]*definition.Dependency, len(deps))
	for i, dep := range deps {
		ds[i] = toProtoDependency(dep)
	}
	return ds
}
