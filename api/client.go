package api

import (
	"github.com/mesg-foundation/application/api/service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50052"
)

var connectionInstance *grpc.ClientConn

func getConnection() (connection *grpc.ClientConn, err error) {
	if connectionInstance == nil {
		connectionInstance, err = grpc.Dial(address, grpc.WithInsecure())
	}
	connection = connectionInstance
	return
}

// CloseClient closes the connection (if exist)
func CloseClient() {
	if connectionInstance != nil {
		connectionInstance.Close()
		connectionInstance = nil
	}
}

// ServiceClient returns a Service Client
func ServiceClient() (client apiService.ServiceClient, ctx context.Context, err error) {
	conn, err := getConnection()
	if err != nil {
		return
	}
	client = apiService.NewServiceClient(conn)
	ctx = context.Background()
	return
}
