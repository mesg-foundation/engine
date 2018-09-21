package api

import (
	"testing"

	"github.com/mesg-foundation/core/event"
	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestValidateEventKey(t *testing.T) {
	a, _ := newAPIAndDockerTest(t)
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
	a, _ := newAPIAndDockerTest(t)
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
