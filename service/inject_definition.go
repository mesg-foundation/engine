// Copyright 2018 MESG Foundation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"sort"

	"github.com/mesg-foundation/core/service/importer"
	"github.com/mesg-foundation/core/x/xstrings"
)

// injectDefinition applies service definition to Service type.
func (s *Service) injectDefinition(def *importer.ServiceDefinition) {
	s.Name = def.Name
	s.Sid = def.Sid
	s.Description = def.Description
	s.Repository = def.Repository
	s.Events = s.defEventsToService(def.Events)
	s.Tasks = s.defTasksToService(def.Tasks)
	s.Dependencies = s.defDependenciesToService(def.Dependencies)

	configuration := &Dependency{
		Key: importer.ConfigurationDependencyKey,
	}
	if def.Configuration != nil {
		configuration.Command = def.Configuration.Command
		configuration.Args = def.Configuration.Args
		configuration.Ports = def.Configuration.Ports
		configuration.Volumes = def.Configuration.Volumes
		configuration.VolumesFrom = def.Configuration.VolumesFrom
		configuration.Env = def.Configuration.Env
	}
	s.Dependencies = append(s.Dependencies, configuration)
}

func (s *Service) defTasksToService(tasks map[string]*importer.Task) []*Task {
	var (
		keys []string
		ts   = make([]*Task, len(tasks))
	)

	for key := range tasks {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for key, task := range tasks {
		i := xstrings.SliceIndex(keys, key)
		ts[i] = &Task{
			Key:         key,
			Name:        task.Name,
			Description: task.Description,
			Inputs:      s.defParametersToService(task.Inputs),
			Outputs:     s.defOutputsToService(task.Outputs),
		}
	}
	return ts
}

func (s *Service) defOutputsToService(outputs map[string]*importer.Output) []*Output {
	var (
		keys []string
		ots  = make([]*Output, len(outputs))
	)

	for key := range outputs {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for key, output := range outputs {
		i := xstrings.SliceIndex(keys, key)
		ots[i] = &Output{
			Key:         key,
			Name:        output.Name,
			Description: output.Description,
			Data:        s.defParametersToService(output.Data),
		}
	}
	return ots
}

func (s *Service) defEventsToService(events map[string]*importer.Event) []*Event {
	var (
		keys []string
		es   = make([]*Event, len(events))
	)

	for key := range events {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for key, event := range events {
		i := xstrings.SliceIndex(keys, key)
		es[i] = &Event{
			Key:         key,
			Name:        event.Name,
			Description: event.Description,
			Data:        s.defParametersToService(event.Data),
		}
	}
	return es
}

func (s *Service) defDependenciesToService(dependencies map[string]*importer.Dependency) []*Dependency {
	var (
		keys []string
		deps = make([]*Dependency, len(dependencies))
	)

	for key := range dependencies {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for key, dep := range dependencies {
		i := xstrings.SliceIndex(keys, key)
		deps[i] = &Dependency{
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

func (s *Service) defParametersToService(params map[string]*importer.Parameter) []*Parameter {
	var (
		keys []string
		ps   = make([]*Parameter, len(params))
	)

	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for key, param := range params {
		i := xstrings.SliceIndex(keys, key)
		ps[i] = &Parameter{
			Key:         key,
			Name:        param.Name,
			Description: param.Description,
			Type:        param.Type,
			Optional:    param.Optional,
			Repeated:    param.Repeated,
		}
	}
	return ps
}
