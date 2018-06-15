package cmdService

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/asaskevich/govalidator"
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

// prepareService downloads if needed, create the service, build it and inject configuration
func prepareService(path string) (importedService *service.Service) {
	path, didDownload, err := downloadServiceIfNeeded(path)
	cmdUtils.HandleError(err)
	if didDownload {
		fmt.Println(aurora.Green("Service downloaded with success"))
		fmt.Println("Temp folder: " + path)
	}
	defer removeTempFolder(path, didDownload)
	importedService, err = service.ImportFromPath(path)
	if err != nil {
		fmt.Println("Run the command 'service validate' to get detailed errors")
		cmdUtils.HandleError(err)
	}
	imageHash, err := buildDockerImage(path)
	cmdUtils.HandleError(err)
	fmt.Println(aurora.Green("Image built with success"))
	fmt.Println("Image hash:", imageHash)
	injectConfigurationInDependencies(importedService, imageHash)
	return
}

func downloadServiceIfNeeded(path string) (newPath string, didDownload bool, err error) {
	didDownload = false
	newPath = path
	if govalidator.IsURL(path) {
		newPath, err = createTempFolder()
		if err != nil {
			return
		}
		err = gitClone(path, newPath)
		didDownload = true
	}
	return
}

func gitClone(url string, path string) (err error) {
	cmdUtils.ShowSpinnerForFunc(cmdUtils.SpinnerOptions{Text: "Downloading service..."}, func() {
		_, err = git.PlainClone(path, false, &git.CloneOptions{
			URL: url,
		})
	})
	return
}

func createTempFolder() (path string, err error) {
	path, err = ioutil.TempDir("", "mesg-")
	return
}

func removeTempFolder(path string, didDownload bool) (err error) {
	if didDownload {
		err = os.RemoveAll(path)
	}
	return
}

func buildDockerImage(path string) (imageHash string, err error) {
	cmdUtils.ShowSpinnerForFunc(cmdUtils.SpinnerOptions{Text: "Building image..."}, func() {
		imageHash, err = container.Build(path)
	})
	return
}

func injectConfigurationInDependencies(s *service.Service, imageHash string) {
	config := s.Configuration
	if config == nil {
		config = &service.Dependency{}
	}
	dependency := &service.Dependency{
		Command:     config.Command,
		Ports:       config.Ports,
		Volumes:     config.Volumes,
		Volumesfrom: config.Volumesfrom,
		Image:       imageHash,
	}
	if s.Dependencies == nil {
		s.Dependencies = make(map[string]*service.Dependency)
	}
	s.Dependencies["service"] = dependency
}

func init() {
	connection, err := grpc.Dial(viper.GetString(config.APIClientTarget), grpc.WithInsecure())
	cmdUtils.HandleError(err)
	cli = core.NewCoreClient(connection)
}
