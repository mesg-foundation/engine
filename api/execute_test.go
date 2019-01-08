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

package api

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestNotRunningServiceError(t *testing.T) {
	e := NotRunningServiceError{ServiceID: "test"}
	require.Equal(t, `Service "test" is not running`, e.Error())
}

func TestExecuteFunc(t *testing.T) {
	a, _, closer := newAPIAndDockerTest(t)
	defer closer()
	executor := newTaskExecutor(a)
	s, _ := service.FromService(&service.Service{
		Name: "TestExecuteFunc",
		Tasks: []*service.Task{
			{
				Key: "test",
			},
		},
	}, service.ContainerOption(a.container))
	id, err := executor.execute(s, "xxx", "test", map[string]interface{}{}, []string{})
	require.NoError(t, err)
	require.NotNil(t, id)
}

func TestExecuteFuncInvalidTaskName(t *testing.T) {
	a, _, closer := newAPIAndDockerTest(t)
	defer closer()
	executor := newTaskExecutor(a)
	srv := &service.Service{}
	_, err := executor.execute(srv, "xxx", "test", map[string]interface{}{}, []string{})
	require.Error(t, err)
}

func TestCheckServiceNotRunning(t *testing.T) {
	a, _, closer := newAPIAndDockerTest(t)
	defer closer()
	executor := newTaskExecutor(a)
	err := executor.checkServiceStatus(&service.Service{Name: "TestCheckServiceNotRunning"})
	require.Error(t, err)
	_, notRunningError := err.(*NotRunningServiceError)
	require.True(t, notRunningError)
}

func TestCheckService(t *testing.T) {
	a, _, closer := newAPIAndDockerTest(t)
	defer closer()
	executor := newTaskExecutor(a)
	s, _ := service.FromService(&service.Service{
		Name: "TestCheckService",
		Dependencies: []*service.Dependency{
			{
				Key:   "test",
				Image: "http-server",
			},
		},
	}, service.ContainerOption(a.container))
	s.Start()
	err := executor.checkServiceStatus(s)
	require.NoError(t, err)
}
