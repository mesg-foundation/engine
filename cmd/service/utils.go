package cmdService

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/service"
)

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
	_, err := service.Start()
	if err != nil {
		fmt.Println(aurora.Red(err))
		os.Exit(0)
	}
	fmt.Println(aurora.Green("Service started"))
}
