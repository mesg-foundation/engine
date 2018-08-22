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

	s := &service.Service{
		Events: map[string]*service.Event{
			"test": {},
		},
	}
	require.Nil(t, ln.validateEventKey(s, ""))
	require.Nil(t, ln.validateEventKey(s, "*"))
	require.Nil(t, ln.validateEventKey(s, "test"))
	require.NotNil(t, ln.validateEventKey(s, "xxx"))
}

func TestIsSubscribedEvent(t *testing.T) {
	a, _ := newAPIAndDockerTest(t)
	ln := newEventListener(a)

	e := &event.Event{Key: "test"}
	require.True(t, ln.isSubscribedEvent("", e))
	require.True(t, ln.isSubscribedEvent("*", e))
	require.True(t, ln.isSubscribedEvent("test", e))
	require.False(t, ln.isSubscribedEvent("xxx", e))
}
