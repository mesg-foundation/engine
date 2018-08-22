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
	"github.com/mesg-foundation/core/service"
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

// prepareService downloads if needed, create the service, build it and inject configuration
func prepareService(path string) *service.Service {
	path, didDownload, err := downloadServiceIfNeeded(path)
	utils.HandleError(err)
	if didDownload {
		defer os.RemoveAll(path)
		fmt.Printf("%s Service downloaded with success\n", aurora.Green("✔"))
	}
	importedService, err := importer.From(path)
	handleValidationError(err)
	utils.HandleError(err)
	imageHash, err := buildDockerImage(path)
	utils.HandleError(err)
	fmt.Printf("%s Image built with success\n", aurora.Green("✔"))
	injectConfigurationInDependencies(importedService, imageHash)
	return importedService
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
