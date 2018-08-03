package service

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/api/core"
	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/service/importer"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

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

func downloadServiceIfNeeded(path string) (string, bool, error) {
	if !govalidator.IsURL(path) {
		return path, false, nil
	}
	newPath, err := createTempFolder()
	if err != nil {
		return "", false, err
	}
	if err := gitClone(path, newPath, "Downloading service..."); err != nil {
		return "", false, err
	}
	return newPath, true, nil
}

func gitClone(repoURL string, path string, message string) (err error) {
	u, err := url.Parse(repoURL)
	if err != nil {
		return err
	}
	if u.Scheme == "" {
		u.Scheme = "https"
	}
	options := &git.CloneOptions{}
	if u.Fragment != "" {
		options.ReferenceName = plumbing.ReferenceName("refs/heads/" + u.Fragment)
		u.Fragment = ""
	}
	options.URL = u.String()
	utils.ShowSpinnerForFunc(utils.SpinnerOptions{Text: message}, func() {
		_, err = git.PlainClone(path, false, options)
	})
	return err
}

func createTempFolder() (string, error) {
	return ioutil.TempDir("", "mesg-")
}

func buildDockerImage(path string) (imageHash string, err error) {
	utils.ShowSpinnerForFunc(utils.SpinnerOptions{Text: "Building image..."}, func() {
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
