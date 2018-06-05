package cmdService

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/api/core"
	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/service"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	git "gopkg.in/src-d/go-git.v4"
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

func gitClone(url string) (path string, err error) {
	tmpFile := "/tmp/mesg-templates"
	err = os.MkdirAll(tmpFile, os.ModePerm)
	if err != nil {
		return
	}
	path, err = ioutil.TempDir(tmpFile, string(time.Now().UnixNano()))
	if err != nil {
		return
	}
	s := cmdUtils.StartSpinner(cmdUtils.SpinnerOptions{Text: "Fetching service..."})
	defer s.Stop()
	_, err = git.PlainClone(path, false, &git.CloneOptions{
		URL: url,
	})
	return
}

func loadService(path string) (importedService *service.Service) {
	var err error
	if _, err = url.ParseRequestURI(path); err == nil {
		path, err = gitClone(path)
		handleError(err)
	}

	importedService, err = service.ImportFromPath(path)
	if err != nil {
		fmt.Println(aurora.Red(err))
		fmt.Println("Run the command 'service validate' to get detailed errors")
		os.Exit(0)
	}
	return
}

func buildDockerImage(path string, name string) (imageHash string) {
	s := cmdUtils.StartSpinner(cmdUtils.SpinnerOptions{Text: "Building image..."})
	tag, err := container.Build(path, []string{name})
	s.Stop()
	if err != nil {
		handleError(err)
		return
	}
	fmt.Println(aurora.Green("Image built with success. Tagged: " + tag))
	return
}

func init() {
	connection, err := grpc.Dial(viper.GetString(config.APIClientTarget), grpc.WithInsecure())
	handleError(err)
	cli = core.NewCoreClient(connection)
}
