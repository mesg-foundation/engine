package queue

// func TestOnMessage(t *testing.T) {
// 	type TestOnMessageType struct {
// 		Foo string
// 		Bar int
// 	}
// 	q := Queue{URL: "amqp://guest:guest@127.0.0.1:5672/"}
// 	channels := []Channel{
// 		Channel{
// 			Kind: Events,
// 			Name: "test",
// 		},
// 	}
// 	go q.Listen("TestOnMessage", channels, func(data interface{}) {
// 		fmt.Println("passe")
// 		assert.Equal(t, true, false)
// 	})

// 	time.Sleep(1)
// 	q.Publish("TestOnMessage", channels, TestOnMessageType{
// 		Foo: "test",
// 		Bar: 1,
// 	})
// }

// func TestOnMessage(t *testing.T) {
// 	isTerminated := make(chan bool)
// 	msgs := make(chan amqp.Delivery)

// 	go onMessage(msgs, func(data interface{}) {
// 		isTerminated <- true
// 	}, isTerminated)

// 	msgs <- amqp.Delivery{
// 		Body: []byte("hello"),
// 	}

// 	<-isTerminated
// }
