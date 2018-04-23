package task

import (
	"github.com/mesg-foundation/application/service"
	"github.com/mesg-foundation/application/types"
)

// Listen for tasks
func (s *Server) Listen(request *types.ListenTaskRequest, stream types.Task_ListenServer) (err error) {
	service := service.New(request.Service)

	subscription := subscribe(service)

	for data := range subscription {
		stream.Send(data)
	}

	return
}
