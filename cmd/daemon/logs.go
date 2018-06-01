package daemon

import (
	"fmt"
	"io"
	"os"

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
	if isRunning {
		reader, err := daemon.Logs()
		if err != nil {
			fmt.Println(aurora.Red(err))
			return
		}
		_, err = io.Copy(os.Stdout, reader)
		if err != nil && err != io.EOF {
			fmt.Println(aurora.Red(err))
		}
	} else {
		fmt.Println(aurora.Brown("Daemon is stopped"))
	}
}
