package api

import (
	"sort"

	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/service/importer"
	"github.com/mesg-foundation/core/x/xstrings"
)

// injectDefinition applies service definition to Service type.
func injectDefinition(s *service.Service, def *importer.ServiceDefinition) {
	s.Name = def.Name
	s.Sid = def.Sid
	s.Description = def.Description
	s.Repository = def.Repository
	s.Events = defEventsToService(def.Events)
	s.Tasks = defTasksToService(def.Tasks)
	s.Dependencies = defDependenciesToService(def.Dependencies)

	s.Configuration = &service.Dependency{
		Key: service.MainServiceKey,
	}
	if def.Configuration != nil {
		s.Configuration.Command = def.Configuration.Command
		s.Configuration.Args = def.Configuration.Args
		s.Configuration.Ports = def.Configuration.Ports
		s.Configuration.Volumes = def.Configuration.Volumes
		s.Configuration.VolumesFrom = def.Configuration.VolumesFrom
		s.Configuration.Env = def.Configuration.Env
	}
}

func defTasksToService(tasks map[string]*importer.Task) []*service.Task {
	var (
		keys []string
		ts   = make([]*service.Task, len(tasks))
	)

	for key := range tasks {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for key, task := range tasks {
		i := xstrings.SliceIndex(keys, key)
		ts[i] = &service.Task{
			Key:         key,
			Name:        task.Name,
			Description: task.Description,
			Inputs:      defParametersToService(task.Inputs),
			Outputs:     defOutputsToService(task.Outputs),
		}
	}
	return ts
}

func defOutputsToService(outputs map[string]*importer.Output) []*service.Output {
	var (
		keys []string
		ots  = make([]*service.Output, len(outputs))
	)

	for key := range outputs {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for key, output := range outputs {
		i := xstrings.SliceIndex(keys, key)
		ots[i] = &service.Output{
			Key:         key,
			Name:        output.Name,
			Description: output.Description,
			Data:        defParametersToService(output.Data),
		}
	}
	return ots
}

func defEventsToService(events map[string]*importer.Event) []*service.Event {
	var (
		keys []string
		es   = make([]*service.Event, len(events))
	)

	for key := range events {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for key, event := range events {
		i := xstrings.SliceIndex(keys, key)
		es[i] = &service.Event{
			Key:         key,
			Name:        event.Name,
			Description: event.Description,
			Data:        defParametersToService(event.Data),
		}
	}
	return es
}

func defDependenciesToService(dependencies map[string]*importer.Dependency) []*service.Dependency {
	var (
		keys []string
		deps = make([]*service.Dependency, len(dependencies))
	)

	for key := range dependencies {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for key, dep := range dependencies {
		i := xstrings.SliceIndex(keys, key)
		deps[i] = &service.Dependency{
			Key:         key,
			Image:       dep.Image,
			Volumes:     dep.Volumes,
			VolumesFrom: dep.VolumesFrom,
			Ports:       dep.Ports,
			Command:     dep.Command,
			Args:        dep.Args,
			Env:         dep.Env,
		}
	}
	return deps
}

func defParametersToService(params map[string]*importer.Parameter) []*service.Parameter {
	var (
		keys []string
		ps   = make([]*service.Parameter, len(params))
	)

	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for key, param := range params {
		i := xstrings.SliceIndex(keys, key)
		ps[i] = &service.Parameter{
			Key:         key,
			Name:        param.Name,
			Description: param.Description,
			Type:        param.Type,
			Optional:    param.Optional,
			Repeated:    param.Repeated,
			Object:      defParametersToService(param.Object),
		}
	}
	return ps
}
