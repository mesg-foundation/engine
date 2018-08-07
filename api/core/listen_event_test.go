package core

import (
	"testing"

	"github.com/mesg-foundation/core/event"
	"github.com/mesg-foundation/core/service"
	"github.com/stvp/assert"
)

func TestValidateEventKey(t *testing.T) {
	s := &service.Service{
		Events: map[string]*service.Event{
			"test": {},
		},
	}
	assert.Nil(t, validateEventKey(s, ""))
	assert.Nil(t, validateEventKey(s, "*"))
	assert.Nil(t, validateEventKey(s, "test"))
	assert.NotNil(t, validateEventKey(s, "xxx"))
}

func TestIsSubscribedEvent(t *testing.T) {
	e := &event.Event{Key: "test"}
	r := &ListenEventRequest{}
	assert.True(t, isSubscribedEvent(r, e))

	r = &ListenEventRequest{EventFilter: ""}
	assert.True(t, isSubscribedEvent(r, e))

	r = &ListenEventRequest{EventFilter: "*"}
	assert.True(t, isSubscribedEvent(r, e))

	r = &ListenEventRequest{EventFilter: "test"}
	assert.True(t, isSubscribedEvent(r, e))

	r = &ListenEventRequest{EventFilter: "xxx"}
	assert.False(t, isSubscribedEvent(r, e))
}
