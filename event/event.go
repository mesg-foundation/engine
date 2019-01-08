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

package event

import (
	"time"

	"github.com/mesg-foundation/core/pubsub"
	"github.com/mesg-foundation/core/service"
)

// Event stores all informations about Events.
type Event struct {
	Service   *service.Service
	Key       string
	Data      interface{}
	CreatedAt time.Time
}

// Create creates an event eventKey with eventData for service s.
func Create(s *service.Service, eventKey string, eventData map[string]interface{}) (*Event, error) {
	event, err := s.GetEvent(eventKey)
	if err != nil {
		return nil, err
	}
	if err := event.RequireData(eventData); err != nil {
		return nil, err
	}
	return &Event{
		Service:   s,
		Key:       eventKey,
		Data:      eventData,
		CreatedAt: time.Now(),
	}, nil
}

// Publish publishes an event for every listener.
func (event *Event) Publish() {
	channel := event.Service.EventSubscriptionChannel()
	go pubsub.Publish(channel, event)
}
