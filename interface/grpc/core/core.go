package core

import service "github.com/mesg-foundation/core/service"

func toProtoServices(ss []*service.Service) []*Service {
	services := make([]*Service, 0)
	for _, s := range ss {
		services = append(services, toProtoService(s))
	}
	return services
}

func toProtoService(s *service.Service) *Service {
	sv := &Service{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		Tasks:       []*Task{},
	}

	for _, task := range s.Tasks {
		t := &Task{
			Key:         task.Key,
			Name:        task.Name,
			Description: task.Description,
			Inputs:      toProtoParameters(task.Inputs),
			Outputs:     []*Output{},
		}
		for _, output := range task.Outputs {
			o := &Output{
				Key:         output.Key,
				Name:        output.Name,
				Description: output.Description,
				Data:        toProtoParameters(output.Data),
			}
			t.Outputs = append(t.Outputs, o)
		}
		sv.Tasks = append(sv.Tasks, t)
	}

	return sv
}

func toProtoParameters(params []*service.Parameter) []*Parameter {
	ps := make([]*Parameter, 0)
	for _, param := range params {
		p := &Parameter{
			Name:        param.Name,
			Description: param.Description,
			Type:        param.Type,
			Optional:    param.Optional,
		}
		ps = append(ps, p)
	}
	return ps
}
