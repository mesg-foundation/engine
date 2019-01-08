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

	"github.com/mesg-foundation/core/event"
	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestValidateEventKey(t *testing.T) {
	a, _, closer := newAPIAndDockerTest(t)
	defer closer()
	ln := newEventListener(a)

	s, _ := service.FromService(&service.Service{
		Events: []*service.Event{
			{
				Key: "test",
			},
		},
	}, service.ContainerOption(a.container))

	ln.eventKey = ""
	require.Nil(t, ln.validateEventKey(s))

	ln.eventKey = "*"
	require.Nil(t, ln.validateEventKey(s))

	ln.eventKey = "test"
	require.Nil(t, ln.validateEventKey(s))

	ln.eventKey = "xxx"
	require.NotNil(t, ln.validateEventKey(s))
}

func TestIsSubscribedEvent(t *testing.T) {
	a, _, closer := newAPIAndDockerTest(t)
	defer closer()
	ln := newEventListener(a)

	e := &event.Event{Key: "test"}

	ln.eventKey = ""
	require.True(t, ln.isSubscribedEvent(e))

	ln.eventKey = "*"
	require.True(t, ln.isSubscribedEvent(e))

	ln.eventKey = "test"
	require.True(t, ln.isSubscribedEvent(e))

	ln.eventKey = "xxx"
	require.False(t, ln.isSubscribedEvent(e))
}
