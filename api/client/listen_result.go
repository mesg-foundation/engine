package client

import (
	"github.com/mesg-foundation/core/pubsub"
)

// Listen for results from the services
func (s *Server) ListenResult(request *ListenResultRequest, stream Client_ListenResultServer) (err error) {
	subscription := pubsub.Subscribe(request.Service.ResultSubscriptionChannel())
	for data := range subscription {
		stream.Send(data.(*ResultData))
	}
	return
}
