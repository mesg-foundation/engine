package queue

import (
	"strings"

	"github.com/streadway/amqp"
)

// Publish see "./interface.go"
func (queue *Queue) Publish(namespace string, channels []Channel, data interface{}) (err error) {
	for _, ch := range channels {
		err = queue.publish(namespace, ch, data)
		if err != nil {
			break
		}
	}
	return
}

func (queue *Queue) publish(namespace string, channel Channel, data interface{}) (err error) {
	if queue.conn == nil {
		err = queue.connect()
	}
	if err != nil {
		return
	}
	ch := queue.channels[namespace]
	if ch == nil {
		ch, err = queue.createInternalChannel(channel)
		if err != nil {
			return
		}
	}
	msg, err := message(data)
	if err != nil {
		return
	}
	err = ch.Publish(namespace, strings.Join([]string{
		string(channel.Kind),
		channel.Name,
	}, "."), false, false, msg)
	return
}

func (queue *Queue) createInternalChannel(channel Channel) (ch *amqp.Channel, err error) {
	ch, err = queue.conn.Channel()
	if err != nil {
		return
	}
	err = ch.ExchangeDeclare(channel.namespace(), "topic", true, false, false, false, nil)
	if err != nil {
		return
	}
	if queue.channels == nil {
		queue.channels = make(map[string]*amqp.Channel)
	}
	queue.channels[channel.namespace()] = ch
	return
}
