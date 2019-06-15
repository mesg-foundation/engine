package api

import (
	"github.com/mesg-foundation/core/protobuf/definition"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/workflow"
)

func FromProtoService(s *definition.Service) *service.Service {
	return &service.Service{
		Hash:          s.Hash,
		Sid:           s.Sid,
		Name:          s.Name,
		Description:   s.Description,
		Repository:    s.Repository,
		Source:        s.Source,
		Tasks:         fromProtoTasks(s.Tasks),
		Events:        fromProtoEvents(s.Events),
		Configuration: fromProtoConfiguration(s.Configuration),
		Dependencies:  fromProtoDependencies(s.Dependencies),
		Workflows:     fromProtoWorkflows(s.Workflows),
	}
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
	}
}

func fromProtoDependencies(deps []*definition.Dependency) []*service.Dependency {
	ds := make([]*service.Dependency, len(deps))
	for i, dep := range deps {
		ds[i] = fromProtoDependency(dep)
	}
	return ds
}

func fromProtoWorkflow(w definition.Workflow) workflow.Workflow {
	ww := workflow.Workflow{Key: w.Key}
	if w.Trigger != nil {
		ww.Trigger.InstanceHash = w.Trigger.InstanceHash
		ww.Trigger.EventKey = w.Trigger.EventKey
		if w.Trigger.Filter != nil {
			ww.Trigger.Filter.TaskKey = w.Trigger.Filter.TaskKey
		}
	}
	for _, t := range w.Tasks {
		ww.Tasks = append(ww.Tasks, workflow.Task{
			InstanceHash: t.InstanceHash,
			Key:          t.Key,
		})
	}
	return ww
}

func fromProtoWorkflows(workflows []*definition.Workflow) []workflow.Workflow {
	ws := make([]workflow.Workflow, len(workflows))
	for i, w := range workflows {
		ws[i] = fromProtoWorkflow(*w)
	}
	return ws
}
