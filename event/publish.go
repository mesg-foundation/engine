package event

import (
	"github.com/mesg-foundation/core/pubsub"
)

// Publish an event for everyone listening
func (event *Event) Publish() {
	channel := event.Service.EventSubscriptionChannel()

	go pubsub.Publish(channel, event)
}
