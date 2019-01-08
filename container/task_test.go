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

package container

import (
	"errors"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/mesg-foundation/core/container/dockertest"
	"github.com/stretchr/testify/require"
)

func TestListTasks(t *testing.T) {
	namespace := []string{"namespace"}
	tasks := []swarm.Task{
		{ID: "1"},
		{ID: "2"},
	}

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	dt.ProvideTaskList(tasks, nil)

	tasks1, err := c.ListTasks(namespace)
	require.NoError(t, err)
	require.Equal(t, tasks, tasks1)
	require.Equal(t, len(tasks), len(tasks1))

	require.Equal(t, types.TaskListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: "com.docker.stack.namespace=" + c.Namespace(namespace),
		}),
	}, (<-dt.LastTaskList()).Options)
}

var errTaskList = errors.New("task list")

func TestListTasksError(t *testing.T) {
	namespace := []string{"namespace"}

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	dt.ProvideTaskList(nil, errTaskList)

	_, err := c.ListTasks(namespace)
	require.Equal(t, errTaskList, err)
}

func TestTasksError(t *testing.T) {
	namespace := []string{"namespace"}
	tasks := []swarm.Task{
		{
			ID:     "1",
			Status: swarm.TaskStatus{Err: "1-err"},
		},
		{
			ID:     "1",
			Status: swarm.TaskStatus{Err: "2-err"},
		},
	}

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	dt.ProvideTaskList(tasks, nil)

	errors, err := c.TasksError(namespace)
	require.NoError(t, err)
	require.Equal(t, len(tasks), len(errors))
	require.Equal(t, tasks[0].Status.Err, errors[0])
	require.Equal(t, tasks[1].Status.Err, errors[1])
}
