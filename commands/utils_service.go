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
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		Tasks:       map[string]*service.Task{},
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
