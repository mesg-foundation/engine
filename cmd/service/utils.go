package cmdService

import (
	"context"
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/api/client"
	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/service"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var cli client.ClientClient

// Set the default path if needed
func defaultPath(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return "./"
}

func handleError(err error) {
	if err != nil {
		fmt.Println(aurora.Red(err))
		os.Exit(0)
	}
}

func loadService(path string) (importedService *service.Service) {
	importedService, err := service.ImportFromPath(path)
	if err != nil {
		fmt.Println(aurora.Red(err))
		fmt.Println("Run the command 'service validate' to get detailed errors")
		os.Exit(0)
	}
	return
}

func startService(service *service.Service) {
	_, err := cli.StartService(context.Background(), &client.StartServiceRequest{
		Service: service,
	})
	handleError(err)
}

func stopService(service *service.Service) {
	_, err := cli.StopService(context.Background(), &client.StopServiceRequest{
		Service: service,
	})
	handleError(err)
}

func init() {
	connection, err := grpc.Dial(viper.GetString(config.APIClientTarget), grpc.WithInsecure())
	handleError(err)
	cli = client.NewClientClient(connection)
}
