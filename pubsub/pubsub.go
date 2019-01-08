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

package pubsub

import (
	"sync"
)

var (
	listeners = make(map[string][]chan Message)
	mu        sync.Mutex
)

// Message sends subscribe/publish messages.
type Message interface{}

// Publish publishes a message to a channel.
func Publish(channel string, data Message) {
	mu.Lock()
	defer mu.Unlock()
	for _, listener := range listeners[channel] {
		listener <- data
	}
}

// Subscribe subscribes to the channel and returns listener for it.
func Subscribe(channel string) chan Message {
	mu.Lock()
	defer mu.Unlock()
	listener := make(chan Message)
	if listeners[channel] == nil {
		listeners[channel] = make([]chan Message, 0)
	}
	listeners[channel] = append(listeners[channel], listener)
	return listener
}

// Unsubscribe unsubscribes a listener from listening channel.
func Unsubscribe(channel string, listener chan Message) {
	mu.Lock()
	defer mu.Unlock()
	for i, l := range listeners[channel] {
		if l == listener {
			listeners[channel] = append(listeners[channel][:i], listeners[channel][i+1:]...)
			if len(listeners[channel]) == 0 {
				listeners[channel] = nil
			}
			return
		}
	}
}
