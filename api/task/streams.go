package task

import (
	"sync"

	"github.com/mesg-foundation/application/service"
	"github.com/mesg-foundation/application/types"
)

var listeners map[string][]chan *types.TaskReply
var mu sync.Mutex

func register(service *service.Service) (res chan *types.TaskReply) {
	res = make(chan *types.TaskReply)
	mu.Lock()
	defer mu.Unlock()
	if listeners[service.Name] == nil {
		listeners[service.Name] = make([]chan *types.TaskReply, 0)
	}
	listeners[service.Name] = append(listeners[service.Name], res)
	return
}

func emit(service *service.Service, data *types.TaskReply) {
	mu.Lock()
	defer mu.Unlock()
	for _, event := range listeners[service.Name] {
		if event != nil {
			event <- data
		}
	}
}

func init() {
	listeners = make(map[string][]chan *types.TaskReply)
}
