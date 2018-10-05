package client

import (
	"fmt"
	"os"
	"sync"

	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/mesg-foundation/core/utils/clierrors"
	"google.golang.org/grpc"
)

var _client coreapi.CoreClient
var once sync.Once

// API returns the client necessary to access the API
func API() (coreapi.CoreClient, error) {
	return getClient()
}

func getClient() (cli coreapi.CoreClient, err error) {
	once.Do(func() {
		c, err := config.Global()
		if err != nil {
			fmt.Fprintln(os.Stderr, clierrors.ErrorMessage(err))
			os.Exit(1)
		}
		var connection *grpc.ClientConn
		connection, err = grpc.Dial(c.Client.Address, grpc.WithInsecure())
		if err != nil {
			return
		}
		_client = coreapi.NewCoreClient(connection)
	})
	cli = _client
	return
}
