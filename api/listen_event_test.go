package api

import (
	"testing"

	"github.com/mesg-foundation/core/event"
	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestValidateEventKey(t *testing.T) {
	var (
		l = &EventListener{}
		s = &service.Service{
			Events: []*service.Event{
				{
					Key: "test",
				},
			},
		}
	)

	l.eventKey = ""
	require.Nil(t, l.validateEventKey(s))

	l.eventKey = "*"
	require.Nil(t, l.validateEventKey(s))

	l.eventKey = "test"
	require.Nil(t, l.validateEventKey(s))

	l.eventKey = "xxx"
	require.NotNil(t, l.validateEventKey(s))
}

func TestIsSubscribedEvent(t *testing.T) {
	var (
		l = &EventListener{}
		e = &event.Event{Key: "test"}
	)

	l.eventKey = ""
	require.True(t, l.isSubscribedEvent(e))

	l.eventKey = "*"
	require.True(t, l.isSubscribedEvent(e))

	l.eventKey = "test"
	require.True(t, l.isSubscribedEvent(e))

	l.eventKey = "xxx"
	require.False(t, l.isSubscribedEvent(e))
}
