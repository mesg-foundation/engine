package cmdService

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/api/core"
	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/service"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var cli core.CoreClient

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

func buildDockerImage(path string, name string) (imageHash string) {
	s := cmdUtils.StartSpinner(cmdUtils.SpinnerOptions{Text: "Building image..."})
	tag, err := container.Build(path, []string{name})
	s.Stop()
	if err != nil {
		handleError(err)
		return
	}
	fmt.Println(aurora.Green("Image built with success. Tagged: " + tag))
	return
}

func init() {
	connection, err := grpc.Dial(viper.GetString(config.APIClientTarget), grpc.WithInsecure())
	handleError(err)
	cli = core.NewCoreClient(connection)
}
