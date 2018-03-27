package cmdService

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Status command returns started and paused services
var Status = &cobra.Command{
	Use:               "status",
	Short:             "List started and paused services",
	Example:           "mesg-cli service status",
	Run:               statusHandler,
	DisableAutoGenTag: true,
}

func statusHandler(cmd *cobra.Command, args []string) {
	// TODO: List started and paused services
	fmt.Println("Started services")
	fmt.Println("- ethereum")
	fmt.Println("- web-server")
	fmt.Println()
	fmt.Println("Paused services")
	fmt.Println("- bitcoin")
}
