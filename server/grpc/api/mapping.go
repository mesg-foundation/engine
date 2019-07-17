package api

import (
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/service"
)

// fromProtoService converts a the protobuf types to the internal service struct
func fromProtoService(s *types.Service) (*service.Service, error) {
	workflows, err := fromProtoWorkflows(s.Workflows)
	if err != nil {
		return nil, err
	}
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
		Workflows:     workflows,
	}, nil
}

func fromProtoTasks(tasks []*types.Service_Task) []*service.Task {
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

func fromProtoEvents(events []*types.Service_Event) []*service.Event {
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

func fromProtoParameters(params []*types.Service_Parameter) []*service.Parameter {
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

func fromProtoConfiguration(configuration *types.Service_Configuration) *service.Dependency {
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

func fromProtoDependency(dep *types.Service_Dependency) *service.Dependency {
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

func fromProtoDependencies(deps []*types.Service_Dependency) []*service.Dependency {
	ds := make([]*service.Dependency, len(deps))
	for i, dep := range deps {
		ds[i] = fromProtoDependency(dep)
	}
	return ds
}

func fromProtoFilters(filters []*types.Service_Workflow_Trigger_Filter) []*service.WorkflowTriggerFilter {
	fs := make([]*service.WorkflowTriggerFilter, len(filters))
	for i, filter := range filters {
		var predicate service.WorkflowPredicate
		// switch filter.Predicate {
		if filter.Predicate == types.Service_Workflow_Trigger_Filter_EQ {
			predicate = service.EQ
		}
		fs[i] = &service.WorkflowTriggerFilter{
			Key:       filter.Key,
			Predicate: predicate,
			Value:     filter.Value,
		}
	}
	return fs
}

func fromProtoWorkflowTask(task *types.Service_Workflow_Task) (*service.WorkflowTask, error) {
	instanceHash, err := hash.Decode(task.InstanceHash)
	if err != nil {
		return nil, err
	}
	return &service.WorkflowTask{
		InstanceHash: instanceHash,
		TaskKey:      task.TaskKey,
	}, nil
}

func fromProtoWorkflows(workflows []*types.Service_Workflow) ([]*service.Workflow, error) {
	wfs := make([]*service.Workflow, len(workflows))
	for i, wf := range workflows {
		var triggerType service.TriggerType
		switch wf.Trigger.Type {
		case types.Service_Workflow_Trigger_Result:
			triggerType = service.RESULT
		case types.Service_Workflow_Trigger_Event:
			triggerType = service.EVENT
		}
		instanceHash, err := hash.Decode(wf.Trigger.InstanceHash)
		if err != nil {
			return nil, err
		}
		task, err := fromProtoWorkflowTask(wf.Task)
		if err != nil {
			return nil, err
		}
		wfs[i] = &service.Workflow{
			Trigger: &service.WorkflowTrigger{
				Type:         triggerType,
				InstanceHash: instanceHash,
				Key:          wf.Trigger.Key,
				Filters:      fromProtoFilters(wf.Trigger.Filters),
			},
			Task: task,
		}
	}
	return wfs, nil
}

// toProtoService converts an internal service struct to the protobuf types
func toProtoService(s *service.Service) *types.Service {
	return &types.Service{
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
		Workflows:     toProtoWorkflows(s.Workflows),
	}
}

func toProtoTasks(tasks []*service.Task) []*types.Service_Task {
	ts := make([]*types.Service_Task, len(tasks))
	for i, task := range tasks {
		ts[i] = &types.Service_Task{
			Key:         task.Key,
			Name:        task.Name,
			Description: task.Description,
			Inputs:      toProtoParameters(task.Inputs),
			Outputs:     toProtoParameters(task.Outputs),
		}
	}
	return ts
}

func toProtoEvents(events []*service.Event) []*types.Service_Event {
	es := make([]*types.Service_Event, len(events))
	for i, event := range events {
		es[i] = &types.Service_Event{
			Key:         event.Key,
			Name:        event.Name,
			Description: event.Description,
			Data:        toProtoParameters(event.Data),
		}
	}
	return es
}

func toProtoParameters(params []*service.Parameter) []*types.Service_Parameter {
	ps := make([]*types.Service_Parameter, len(params))
	for i, param := range params {
		ps[i] = &types.Service_Parameter{
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

func toProtoConfiguration(configuration *service.Dependency) *types.Service_Configuration {
	return &types.Service_Configuration{
		Args:        configuration.Args,
		Command:     configuration.Command,
		Ports:       configuration.Ports,
		Volumes:     configuration.Volumes,
		VolumesFrom: configuration.VolumesFrom,
		Env:         configuration.Env,
	}
}

func toProtoDependency(dep *service.Dependency) *types.Service_Dependency {
	return &types.Service_Dependency{
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

func toProtoDependencies(deps []*service.Dependency) []*types.Service_Dependency {
	ds := make([]*types.Service_Dependency, len(deps))
	for i, dep := range deps {
		ds[i] = toProtoDependency(dep)
	}
	return ds
}

func toProtoFilters(filters []*service.WorkflowTriggerFilter) []*types.Service_Workflow_Trigger_Filter {
	fs := make([]*types.Service_Workflow_Trigger_Filter, len(filters))
	for i, filter := range filters {
		var predicate types.Service_Workflow_Trigger_Filter_Predicate
		// switch filter.Predicate {
		if filter.Predicate == service.EQ {
			predicate = types.Service_Workflow_Trigger_Filter_EQ
		}
		fs[i] = &types.Service_Workflow_Trigger_Filter{
			Key:       filter.Key,
			Predicate: predicate,
			Value:     filter.Value.(string),
		}
	}
	return fs
}

func toProtoWorkflowTask(task *service.WorkflowTask) *types.Service_Workflow_Task {
	return &types.Service_Workflow_Task{
		InstanceHash: task.InstanceHash.String(),
		TaskKey:      task.TaskKey,
	}
}

func toProtoWorkflows(workflows []*service.Workflow) []*types.Service_Workflow {
	wfs := make([]*types.Service_Workflow, len(workflows))
	for i, wf := range workflows {
		var triggerType types.Service_Workflow_Trigger_TriggerType
		switch wf.Trigger.Type {
		case service.EVENT:
			triggerType = types.Service_Workflow_Trigger_Event
		case service.RESULT:
			triggerType = types.Service_Workflow_Trigger_Result
		}
		wfs[i] = &types.Service_Workflow{
			Trigger: &types.Service_Workflow_Trigger{
				Type:         triggerType,
				InstanceHash: wf.Trigger.InstanceHash.String(),
				Key:          wf.Trigger.Key,
				Filters:      toProtoFilters(wf.Trigger.Filters),
			},
			Task: toProtoWorkflowTask(wf.Task),
		}
	}
	return wfs
}
