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

import "fmt"

// EventNotFoundError is an error returned when corresponding event cannot be found in service.
type EventNotFoundError struct {
	EventKey    string
	ServiceName string
}

func (e *EventNotFoundError) Error() string {
	return fmt.Sprintf("Event %q not found in service %q", e.EventKey, e.ServiceName)
}

// TaskNotFoundError is an error returned when corresponding task cannot be found in service.
type TaskNotFoundError struct {
	TaskKey     string
	ServiceName string
}

func (e *TaskNotFoundError) Error() string {
	return fmt.Sprintf("Task %q not found in service %q", e.TaskKey, e.ServiceName)
}

// TaskInputNotFoundError is an error returned when service doesn't contain corresponding input.
type TaskInputNotFoundError struct {
	TaskKey      string
	TaskInputKey string
	ServiceName  string
}

func (e *TaskInputNotFoundError) Error() string {
	return fmt.Sprintf("Input %q of task %q not found in service %q", e.TaskInputKey, e.TaskKey,
		e.ServiceName)
}

// TaskOutputNotFoundError is an error returned when service doesn't contain corresponding output.
type TaskOutputNotFoundError struct {
	TaskKey       string
	TaskOutputKey string
	ServiceName   string
}

func (e *TaskOutputNotFoundError) Error() string {
	return fmt.Sprintf("Output %q of task %q not found in service %q", e.TaskOutputKey, e.TaskKey,
		e.ServiceName)
}

// InvalidEventDataError is an error returned when the data of corresponding event is not valid.
type InvalidEventDataError struct {
	EventKey    string
	ServiceName string
	Warnings    []*ParameterWarning
}

func (e *InvalidEventDataError) Error() string {
	s := fmt.Sprintf("Data of event %q is invalid in service %q", e.EventKey, e.ServiceName)
	for _, warning := range e.Warnings {
		s = fmt.Sprintf("%s. %s", s, warning)
	}
	return s
}

// InvalidTaskInputError is an error returned when the inputs of corresponding task are not valid.
type InvalidTaskInputError struct {
	TaskKey     string
	ServiceName string
	Warnings    []*ParameterWarning
}

func (e *InvalidTaskInputError) Error() string {
	s := fmt.Sprintf("Inputs of task %q are invalid in service %q", e.TaskKey, e.ServiceName)
	for _, warning := range e.Warnings {
		s = fmt.Sprintf("%s. %s", s, warning)
	}
	return s
}

// InvalidTaskOutputError is an error returned when the outputs of corresponding task are not valid.
type InvalidTaskOutputError struct {
	TaskKey       string
	TaskOutputKey string
	ServiceName   string
	Warnings      []*ParameterWarning
}

func (e *InvalidTaskOutputError) Error() string {
	s := fmt.Sprintf("Outputs %q of task %q are invalid in service %q", e.TaskOutputKey, e.TaskKey,
		e.ServiceName)
	for _, warning := range e.Warnings {
		s = fmt.Sprintf("%s. %s", s, warning)
	}
	return s
}
