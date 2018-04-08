package queue

import "github.com/streadway/amqp"

// EventKind is the kind of event that you can publish
type EventKind string

const (
	// Events is any event related to your technology
	Events EventKind = "event"
	// Tasks is an event type to trigger a task in a service
	Tasks EventKind = "task"
	// Systems is an event type to send some meta informations
	Systems EventKind = "system"
)

// Queue is a struct that contains the connection and the list of channel associated
type Queue struct {
	URL      string
	conn     *amqp.Connection
	channels map[string]*amqp.Channel
}

// Channel is the channel where events will transit
type Channel struct {
	Kind EventKind
	Name string
}

// Interface is interface to publish or listen event to a queue system
type Interface interface {
	Close() (err error)

	Publish(channels []Channel, data interface{}) (err error)
	Listen(channels []Channel, onEvent func(data interface{})) (err error)
}
