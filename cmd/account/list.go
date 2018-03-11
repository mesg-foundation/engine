package cmdAccount

import (
	"fmt"

	"github.com/spf13/cobra"
)

// List all the accounts
var List = &cobra.Command{
	Use:     "list",
	Short:   "List all the accounts on this computer",
	Example: "mesg-cli account list",
	Run:     listHandler,
}

func listHandler(cmd *cobra.Command, args []string) {
	// TODO add real listing
	fmt.Println("0x0000000000000000000000000000000000000000")
	fmt.Println("0x0000000000000000000000000000000000000001")
}
