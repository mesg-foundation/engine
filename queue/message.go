package queue

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

func message(data interface{}) (msg amqp.Publishing, err error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}
	msg = amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonData,
	}
	return
}
