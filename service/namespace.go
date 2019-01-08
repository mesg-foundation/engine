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
	"github.com/mesg-foundation/core/utils/hash"
)

// NAMESPACE is the namespace used for the docker services.
const eventChannel string = "Event"
const taskChannel string = "Task"
const resultChannel string = "Result"

// namespace returns the namespace of the service.
func (service *Service) namespace() []string {
	return []string{service.Hash}
}

// namespace returns the namespace of a dependency.
func (d *Dependency) namespace() []string {
	return append(d.service.namespace(), d.Key)
}

// EventSubscriptionChannel returns the channel to listen for events from this service.
func (service *Service) EventSubscriptionChannel() string {
	return hash.Calculate(append(
		service.namespace(),
		eventChannel,
	))
}

// TaskSubscriptionChannel returns the channel to listen for tasks from this service.
func (service *Service) TaskSubscriptionChannel() string {
	return hash.Calculate(append(
		service.namespace(),
		taskChannel,
	))
}

// ResultSubscriptionChannel returns the channel to listen for tasks from this service.
func (service *Service) ResultSubscriptionChannel() string {
	return hash.Calculate(append(
		service.namespace(),
		resultChannel,
	))
}
