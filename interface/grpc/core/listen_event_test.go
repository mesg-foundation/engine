package core

import (
	"testing"

	"github.com/mesg-foundation/core/event"
	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestValidateEventKey(t *testing.T) {
	s := &service.Service{
		Events: map[string]*service.Event{
			"test": {},
		},
	}
	require.Nil(t, validateEventKey(s, ""))
	require.Nil(t, validateEventKey(s, "*"))
	require.Nil(t, validateEventKey(s, "test"))
	require.NotNil(t, validateEventKey(s, "xxx"))
}

func TestIsSubscribedEvent(t *testing.T) {
	e := &event.Event{Key: "test"}
	r := &ListenEventRequest{}
	require.True(t, isSubscribedEvent(r, e))

	r = &ListenEventRequest{EventFilter: ""}
	require.True(t, isSubscribedEvent(r, e))

	r = &ListenEventRequest{EventFilter: "*"}
	require.True(t, isSubscribedEvent(r, e))

	r = &ListenEventRequest{EventFilter: "test"}
	require.True(t, isSubscribedEvent(r, e))

	r = &ListenEventRequest{EventFilter: "xxx"}
	require.False(t, isSubscribedEvent(r, e))
}
