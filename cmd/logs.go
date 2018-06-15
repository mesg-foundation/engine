package cmd

import (
	"fmt"
	"os"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/daemon"
	"github.com/spf13/cobra"
)

// Logs of the core
var Logs = &cobra.Command{
	Use:               "logs",
	Short:             "Show the MESG Core's logs",
	Run:               logsHandler,
	DisableAutoGenTag: true,
}

func init() {
	RootCmd.AddCommand(Logs)
}

func logsHandler(cmd *cobra.Command, args []string) {
	status, err := daemon.Status()
	if err != nil {
		fmt.Println(aurora.Red(err))
		return
	}
	if status == container.STOPPED {
		fmt.Println(aurora.Brown("MESG Core is stopped"))
		return
	}
	reader, err := daemon.Logs()
	defer reader.Close()
	if err != nil {
		fmt.Println(aurora.Red(err))
		return
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, reader)
}
