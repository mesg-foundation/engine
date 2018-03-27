package cmdService

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/mesg-foundation/application/service"

	"github.com/spf13/cobra"
)

// Init run the Init command for a service
var Init = &cobra.Command{
	Use:               "init NAME",
	Short:             "Initialize a service",
	Example:           "mesg-cli service init",
	Run:               initHandler,
	DisableAutoGenTag: true,
}

func initHandler(cmd *cobra.Command, args []string) {
	res := service.Service{}
	res.Publish = service.PublishAll
	res.Visibility = service.VisibilityAll

	res.Name = cmd.Flag("name").Value.String()
	if res.Name == "" && survey.AskOne(&survey.Input{Message: "Name of your service"}, &res.Name, nil) != nil {
		os.Exit(0)
	}
	key := strings.Replace(strings.ToLower(res.Name), " ", "-", -1)

	res.Description = cmd.Flag("description").Value.String()
	if res.Description == "" && survey.AskOne(&survey.Input{Message: "Description of your service"}, &res.Description, nil) != nil {
		os.Exit(0)
	}

	publishStr := cmd.Flag("publish").Value.String()
	if publishStr == "" && survey.AskOne(&survey.Select{
		Message: "Publish details of the service (ALL)",
		Options: []string{
			string(service.PublishAll),
			string(service.PublishSource),
			string(service.PublishContainer),
			string(service.PublishNone),
		},
	}, &publishStr, nil) != nil {
		os.Exit(0)
	}
	res.Publish = service.Publish(publishStr)

	visibilityStr := cmd.Flag("visibility").Value.String()
	if visibilityStr == "" && survey.AskOne(&survey.Select{
		Message: "Visibility of the service (ALL)",
		Options: []string{
			string(service.VisibilityAll),
			string(service.VisibilityUsers),
			string(service.VisibilityWorkers),
			string(service.VisibilityNone),
		},
	}, &visibilityStr, nil) != nil {
		os.Exit(0)
	}
	res.Visibility = service.Visibility(visibilityStr)

	res.Dependencies = service.Dependencies{
		key: service.Dependency{
			Image: strings.Join([]string{"mesg", key}, "/"),
		},
	}

	out, err := yaml.Marshal(res)
	if err != nil {
		panic(err)
	}
	fmt.Println()
	fmt.Println(string(out))
	ok := false
	if survey.AskOne(&survey.Confirm{Message: "Is this ok ?"}, &ok, nil) != nil {
		os.Exit(0)
	}
	if !ok {
		return
	}
	err = os.Mkdir(key, os.ModePerm)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filepath.Join(key, "mesg.yml"), out, os.ModePerm)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filepath.Join(key, "Dockerfile"), []byte(""), os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func init() {
	Init.Flags().StringP("name", "n", "", "Name of the service")
	Init.Flags().StringP("description", "d", "", "Description of the service")
	Init.Flags().StringP("visibility", "v", "", "Visibility of the service")
	Init.Flags().StringP("publish", "p", "", "Publish details of the service")
}
