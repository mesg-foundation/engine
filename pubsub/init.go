package pubsub

import (
	"sync"
)

var listeners map[string][]chan Message
var mu sync.Mutex

// Message is the interface to send subscribe/publish messages
type Message interface {
}

func init() {
	listeners = make(map[string][]chan Message)
}
