package service

import (
	"regexp"

	"github.com/streadway/amqp"
)

const (
	queueEndpoint      = "amqp://guest:guest@localhost:5672/"
	channelType        = "topic"
	queueName          = "events"
	eventRoutingRegexp = "^event."
	systemClose        = "system.close"
	systemInit         = "system.init"
)

// ListenEvents listen all the events that matches the event parameter from the services
// event parameter can match all using the "*" value
func ListenEvents(service *Service, eventFilter string, onEvent func(event amqp.Delivery)) (err error) {
	if eventFilter == "" {
		eventFilter = "*"
	}
	conn, err := amqp.Dial(queueEndpoint)
	if err != nil {
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(service.Namespace(), channelType, true, false, false, false, nil)
	if err != nil {
		return
	}

	q, err := ch.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		return
	}

	ch.QueueBind(q.Name, "event."+eventFilter, service.Namespace(), false, nil)
	ch.QueueBind(q.Name, "system.*", service.Namespace(), false, nil)

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return
	}

	isTerminated := make(chan bool)
	go handleMessages(msgs, onEvent, isTerminated)
	<-isTerminated

	return
}

func handleMessages(msgs <-chan amqp.Delivery, onEvent func(event amqp.Delivery), isTerminated chan bool) {
	isEvent, err := regexp.Compile(eventRoutingRegexp)
	if err != nil {
		panic(err)
	}
	for d := range msgs {
		switch {
		case d.RoutingKey == systemClose:
			isTerminated <- true
		case isEvent.Match([]byte(d.RoutingKey)):
			onEvent(d)
		}
	}
}
