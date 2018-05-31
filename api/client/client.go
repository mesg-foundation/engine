package client

import (
	"sync"

	"github.com/mesg-foundation/core/api/core"
	"github.com/mesg-foundation/core/config"
	"github.com/spf13/viper"
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
		var connection *grpc.ClientConn
		connection, err = grpc.Dial(viper.GetString(config.APIServerAddress), grpc.WithInsecure())
		if err != nil {
			return
		}
		_client = core.NewCoreClient(connection)
	})
	cli = _client
	return
}
