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
	sv := &service.Service{
		ID:            s.ID,
		Name:          s.Name,
		Description:   s.Description,
		Repository:    s.Repository,
		Tasks:         map[string]*service.Task{},
		Events:        map[string]*service.Event{},
		Dependencies:  map[string]*service.Dependency{},
		Configuration: nil,
	}

	if s.Configuration != nil {
		sv.Configuration = &service.Dependency{
			Image:       s.Configuration.Image,
			Volumes:     s.Configuration.Volumes,
			VolumesFrom: s.Configuration.Volumesfrom,
			Ports:       s.Configuration.Ports,
			Command:     s.Configuration.Command,
		}
	}

	for eventKey, event := range s.Events {
		sv.Events[eventKey] = &service.Event{
			Key:         event.Key,
			Name:        event.Name,
			Description: event.Description,
			ServiceName: event.ServiceName,
			Data:        toParameters(event.Data),
		}
	}

	for dependencyKey, dependency := range s.Dependencies {
		sv.Dependencies[dependencyKey] = &service.Dependency{
			Image:       dependency.Image,
			Volumes:     dependency.Volumes,
			VolumesFrom: dependency.Volumesfrom,
			Ports:       dependency.Ports,
			Command:     dependency.Command,
		}
	}

	for taskKey, task := range s.Tasks {
		t := &service.Task{
			Key:         task.Key,
			Name:        task.Name,
			Description: task.Description,
			ServiceName: task.ServiceName,
			Inputs:      toParameters(task.Inputs),
			Outputs:     map[string]*service.Output{},
		}
		for outputKey, output := range task.Outputs {
			t.Outputs[outputKey] = &service.Output{
				Key:         output.Key,
				Name:        output.Name,
				Description: output.Description,
				TaskKey:     output.TaskKey,
				ServiceName: output.ServiceName,
				Data:        toParameters(output.Data),
			}
		}
		sv.Tasks[taskKey] = t
	}

	return sv
}

// TODO(ilgooz) rm this when we stop using internal methods of service in cmd.
func toParameters(params map[string]*core.Parameter) map[string]*service.Parameter {
	gParams := make(map[string]*service.Parameter, 0)
	for key, param := range params {
		gParams[key] = &service.Parameter{
			Name:        param.Name,
			Description: param.Description,
			Type:        param.Type,
			Optional:    param.Optional,
		}
	}
	return gParams
}
