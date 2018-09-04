package service

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/interface/grpc/core"
	"github.com/mesg-foundation/core/service/importer"
	"google.golang.org/grpc"
)

func cli() core.CoreClient {
	apiAddress, err := config.APIAddress().GetValue()
	utils.HandleError(err)
	apiPort, err := config.APIPort().GetValue()
	utils.HandleError(err)
	connection, err := grpc.Dial(apiAddress+":"+apiPort, grpc.WithInsecure())
	utils.HandleError(err)
	return core.NewCoreClient(connection)
}

// Set the default path if needed
func defaultPath(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return "./"
}

func handleValidationError(err error) {
	if _, ok := err.(*importer.ValidationError); ok {
		fmt.Println(aurora.Red(err))
		fmt.Println("Run the command 'service validate' for more details")
		os.Exit(0)
	}
}
