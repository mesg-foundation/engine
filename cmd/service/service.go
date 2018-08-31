package service

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
		Tasks:       []*service.Task{},
	}

	for _, task := range s.Tasks {
		t := &service.Task{
			Key:         task.Key,
			Name:        task.Name,
			Description: task.Description,
			Inputs:      toParameters(task.Inputs),
			Outputs:     []*service.Output{},
		}
		for _, output := range task.Outputs {
			o := &service.Output{
				Key:         output.Key,
				Name:        output.Name,
				Description: output.Description,
				Data:        toParameters(output.Data),
			}
			t.Outputs = append(t.Outputs, o)
		}
		sv.Tasks = append(sv.Tasks, t)
	}

	return sv
}

// TODO(ilgooz) rm this when we stop using internal methods of service in cmd.
func toParameters(params []*core.Parameter) []*service.Parameter {
	gParams := make([]*service.Parameter, 0)
	for _, param := range params {
		p := &service.Parameter{
			Name:        param.Name,
			Description: param.Description,
			Type:        param.Type,
			Optional:    param.Optional,
		}
		gParams = append(gParams, p)
	}
	return gParams
}
