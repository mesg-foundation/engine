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
	"testing"

	"github.com/stretchr/testify/require"
)

type messageStructTest struct {
	test string
}

func TestPublish(t *testing.T) {
	key := "TestPublish"
	data := messageStructTest{test: "TestPublish"}

	res := Subscribe(key)
	go Publish(key, data)
	x := <-res
	require.Equal(t, x, data)
}

func TestPublishMultipleListeners(t *testing.T) {
	key := "TestPublishMultipleListeners"
	data := messageStructTest{test: "TestPublishMultipleListeners"}
	res1 := Subscribe(key)
	res2 := Subscribe(key)
	go Publish(key, data)
	x := <-res1
	y := <-res2
	require.Equal(t, x, data)
	require.Equal(t, y, data)
}

func TestSubscribe(t *testing.T) {
	key := "TestSubscribe"
	res := Subscribe(key)
	require.NotNil(t, res)
	require.Len(t, listeners[key], 1)
}

func TestSubscribeMultipleTimes(t *testing.T) {
	key := "TestSubscribeMultipleTimes"
	Subscribe(key)
	Subscribe(key)
	require.Len(t, listeners[key], 2)
}

func TestUnsubscribe(t *testing.T) {
	key := "TestUnsubscribe"
	channel := Subscribe(key)
	channel1 := Subscribe(key)
	Unsubscribe(key, channel)
	require.Len(t, listeners[key], 1)
	require.NotNil(t, listeners[key])
	Unsubscribe(key, channel1)
	require.Len(t, listeners[key], 0)
	require.Nil(t, listeners[key])
}
