package event

import (
	"github.com/mesg-foundation/core/pubsub"
)

// Publish publishes an event for every listener.
func (event *Event) Publish() {
	channel := event.Service.EventSubscriptionChannel()

	go pubsub.Publish(channel, event)
}
