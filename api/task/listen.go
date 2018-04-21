package task

import (
	"github.com/mesg-foundation/application/service"
	"github.com/mesg-foundation/application/types"
)

// Listen for tasks
func (s *Server) Listen(request *types.ListenTaskRequest, stream types.Task_ListenServer) (err error) {
	service := service.New(request.Service)

	onMessage := register(service)

	for x := range onMessage {
		stream.Send(x)
	}

	return
}
