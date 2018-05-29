package daemon

import (
	"bytes"
	"fmt"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/container"
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
	service, err := getService()
	if err != nil {
		fmt.Println(aurora.Red(err))
		return
	}
	if service != nil {
		var stream bytes.Buffer
		err = container.ServiceLogs([]string{name}, &stream)
		if err != nil {
			fmt.Println(aurora.Red(err))
			return
		}
		buf := make([]byte, 1024)
		for {
			n, _ := stream.Read(buf)
			if n != 0 {
				fmt.Print(string(buf[:n]))
			}
			time.Sleep(500 * time.Millisecond)
		}
	}
}
