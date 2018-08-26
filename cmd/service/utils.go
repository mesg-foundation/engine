package service

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/api/core"
	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/service/importer"
	"github.com/mesg-foundation/core/x/xgit"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// TODO(ilgooz): remove this after service package made Newable.
var defaultContainer *container.Container

// TODO(ilgooz): remove this after service package made Newable.
func init() {
	c, err := container.New()
	if err != nil {
		log.Fatal(err)
	}
	defaultContainer = c
}

func cli() core.CoreClient {
	connection, err := grpc.Dial(viper.GetString(config.APIClientTarget), grpc.WithInsecure())
	utils.HandleError(err)
	return core.NewCoreClient(connection)
}

// Set the default path if needed
func defaultPath(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return "./"
}

func handleValidationError(err error) {
	if _, ok := err.(*importer.ValidationError); ok {
		fmt.Println(aurora.Red(err))
		fmt.Println("Run the command 'service validate' for more details")
		os.Exit(0)
	}
}

func downloadServiceIfNeeded(path string) (newPath string, didDownload bool, err error) {
	if !govalidator.IsURL(path) {
		return path, false, nil
	}
	newPath, err = ioutil.TempDir("", utils.TempDirPrefix)
	if err != nil {
		return "", false, err
	}
	utils.ShowSpinnerForFunc(utils.SpinnerOptions{Text: "Downloading service..."}, func() {
		err = xgit.Clone(path, newPath)
	})
	if err != nil {
		return "", false, err
	}
	return newPath, true, nil
}

func buildDockerImage(path string) (imageHash string, err error) {
	utils.ShowSpinnerForFunc(utils.SpinnerOptions{Text: "Building image..."}, func() {
		imageHash, err = defaultContainer.Build(path)
	})
	return
}
