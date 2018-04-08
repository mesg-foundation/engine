package queue

import (
	"encoding/json"
	"strings"
)

// Listen see "./interface.go"
func (queue *Queue) Listen(namespace string, channels []Channel, onEvent func(data interface{})) (err error) {
	if queue.conn == nil {
		err = queue.connect()
	}
	if err != nil {
		return
	}
	ch, err := queue.conn.Channel()
	if err != nil {
		return
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(namespace, "topic", true, false, false, false, nil)
	if err != nil {
		return
	}

	q, err := ch.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		return
	}

	for _, channel := range channels {
		ch.QueueBind(q.Name, strings.Join([]string{
			string(channel.Kind),
			channel.Name,
		}, "."), namespace, false, nil)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return
	}

	isTerminated := make(chan bool)
	go func() {
		for d := range msgs {
			var res interface{}
			err := json.Unmarshal(d.Body, &res)
			if err != nil {
				panic(err)
			}
			onEvent(res)
		}
	}()
	<-isTerminated

	return
}
