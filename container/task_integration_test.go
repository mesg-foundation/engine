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

// +build integration

package container

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestIntegrationListTasks(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	namespace := []string{"TestListTasks"}
	startTestService(namespace)
	defer c.StopService(namespace)
	tasks, err := c.ListTasks(namespace)
	require.NoError(t, err)
	require.NotNil(t, tasks)
	require.Equal(t, 1, len(tasks))
}

func TestIntegrationTasksError(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	namespace := []string{"TestTasksError"}
	c.StartService(ServiceOptions{
		Image:     "fiifioewifewiewfifewijopwjeokpfeo",
		Namespace: namespace,
	})
	defer c.StopService(namespace)
	var errors []string
	for {
		errors, err = c.TasksError(namespace)
		if err != nil || len(errors) > 0 {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	require.NoError(t, err)
	require.NotNil(t, errors)
	require.True(t, len(errors) > 0)
	require.Equal(t, "No such image: fiifioewifewiewfifewijopwjeokpfeo:latest", errors[0])
}
