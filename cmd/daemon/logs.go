package daemon

import (
	"fmt"
	"os"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/daemon"
	"github.com/spf13/cobra"
)

// Logs the daemon
var Logs = &cobra.Command{
	Use:               "logs",
	Short:             "Show the daemon's logs",
	Run:               logsHandler,
	DisableAutoGenTag: true,
}

func logsHandler(cmd *cobra.Command, args []string) {
	isRunning, err := daemon.IsRunning()
	if err != nil {
		fmt.Println(aurora.Red(err))
		return
	}
	if !isRunning {
		fmt.Println(aurora.Brown("Daemon is stopped"))
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
