package result

import (
	"github.com/mesg-foundation/core/pubsub"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/types"
)

func getSubscription(request *types.ListenResultRequest) (subscription chan pubsub.Message) {
	service := service.New(request.Service)

	subscription = pubsub.Subscribe(service.ResultSubscriptionChannel())
	return
}

// Listen for results from the services
func (s *Server) Listen(request *types.ListenResultRequest, stream types.Result_ListenServer) (err error) {
	subscription := getSubscription(request)
	for data := range subscription {
		stream.Send(data.(*types.ResultReply))
	}
	return
}
