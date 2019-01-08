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
	"testing"

	"github.com/mesg-foundation/core/utils/hash"
	"github.com/stretchr/testify/require"
)

func TestServiceNamespace(t *testing.T) {
	service, _ := FromService(&Service{Name: "TestServiceNamespace"})
	namespace := service.namespace()
	require.Equal(t, namespace, []string{service.Hash})
}

func TestDependencyNamespace(t *testing.T) {
	service, _ := FromService(&Service{
		Name: "TestDependencyNamespace",
		Dependencies: []*Dependency{
			{
				Key:   "test",
				Image: "http-server",
			},
		},
	})
	dep := service.Dependencies[0]
	require.Equal(t, dep.namespace(), []string{service.Hash, "test"})
}

func TestEventSubscriptionChannel(t *testing.T) {
	service, _ := FromService(&Service{Name: "TestEventSubscriptionChannel"})
	require.Equal(t, service.EventSubscriptionChannel(), hash.Calculate(append(
		service.namespace(),
		eventChannel,
	)))
}

func TestTaskSubscriptionChannel(t *testing.T) {
	service, _ := FromService(&Service{Name: "TaskSubscriptionChannel"})
	require.Equal(t, service.TaskSubscriptionChannel(), hash.Calculate(append(
		service.namespace(),
		taskChannel,
	)))
}

func TestResultSubscriptionChannel(t *testing.T) {
	service, _ := FromService(&Service{Name: "ResultSubscriptionChannel"})
	require.Equal(t, service.ResultSubscriptionChannel(), hash.Calculate(append(
		service.namespace(),
		resultChannel,
	)))
}
