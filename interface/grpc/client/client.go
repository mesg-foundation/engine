package client

import (
	"sync"

	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/interface/grpc/core"
	"google.golang.org/grpc"
)

var _client core.CoreClient
var once sync.Once

// API returns the client necessary to access the API
func API() (core.CoreClient, error) {
	return getClient()
}

func getClient() (cli core.CoreClient, err error) {
	once.Do(func() {
		apiAddress, err := config.APIAddress().GetValue()
		if err != nil {
			return
		}
		apiPort, err := config.APIPort().GetValue()
		if err != nil {
			return
		}
		utils.HandleError(err)

		var connection *grpc.ClientConn
		connection, err = grpc.Dial(apiAddress+":"+apiPort, grpc.WithInsecure())
		if err != nil {
			return
		}
		_client = core.NewCoreClient(connection)
	})
	cli = _client
	return
}
