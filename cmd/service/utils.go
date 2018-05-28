package cmdService

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fsouza/go-dockerclient"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/api/core"
	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/mesg-foundation/core/config"
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

func buildDockerImage(path string, dockerImage string) (imageHash string, err error) {
	s := cmdUtils.StartSpinner(cmdUtils.SpinnerOptions{Text: "Building image..."})
	defer s.Stop()
	dockerclient, err := docker.NewClientFromEnv()
	var stream bytes.Buffer
	go func() {
		err = dockerclient.BuildImage(docker.BuildImageOptions{
			Context: context.Background(),
			Name: strings.Join([]string{
				"mesg",
				strings.Replace(strings.ToLower(dockerImage), " ", "-", -1),
			}, "/"),
			RmTmpContainer:      true,
			ForceRmTmpContainer: true,
			ContextDir:          path,
			OutputStream:        &stream,
			SuppressOutput:      true,
		})
		if err != nil {
			fmt.Println(aurora.Red(err))
			os.Exit(1)
		}
	}()

	buf := make([]byte, 1024)
	for {
		n, _ := stream.Read(buf)
		imageHash = strings.Join([]string{
			imageHash,
			string(buf[:n]),
		}, "")
		if len(imageHash) > 0 {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	return
}

func init() {
	connection, err := grpc.Dial(viper.GetString(config.APIClientTarget), grpc.WithInsecure())
	handleError(err)
	cli = core.NewCoreClient(connection)
}
