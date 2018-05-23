package daemon

import (
	// "bytes"
	// "context"
	// "fmt"
	// "os"
	// "time"

	// godocker "github.com/fsouza/go-dockerclient"
	// "github.com/logrusorgru/aurora"
	// "github.com/mesg-foundation/core/daemon"
	// "github.com/mesg-foundation/core/docker"
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
	// isRunning, err := daemon.IsRunning()
	// if err != nil {
	// 	fmt.Println(aurora.Red(err))
	// 	return
	// }
	// if isRunning {
	// 	client, err := docker.Client()
	// 	if err != nil {
	// 		fmt.Println(aurora.Red(err))
	// 		return
	// 	}

	// 	var stream bytes.Buffer
	// 	go func() {
	// 		err = client.GetServiceLogs(godocker.LogsServiceOptions{
	// 			Context:      context.Background(),
	// 			Service:      service.ID,
	// 			Follow:       true,
	// 			Stdout:       true,
	// 			Stderr:       true,
	// 			Timestamps:   false,
	// 			OutputStream: &stream,
	// 			ErrorStream:  &stream,
	// 		})
	// 		if err != nil {
	// 			fmt.Println(aurora.Red(err))
	// 			os.Exit(1)
	// 		}
	// 	}()

	// 	buf := make([]byte, 1024)
	// 	for {
	// 		n, _ := stream.Read(buf)
	// 		if n != 0 {
	// 			fmt.Print(string(buf[:n]))
	// 		}
	// 		time.Sleep(500 * time.Millisecond)
	// 	}
	// }
}
