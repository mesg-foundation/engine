package queue

import "github.com/streadway/amqp"

func (queue *Queue) connect() (err error) {
	if queue.conn != nil {
		return
	}
	queue.conn, err = amqp.Dial(queue.URL)
	return
}

func (queue *Queue) disconnect() (err error) {
	err = queue.conn.Close()
	return
}
