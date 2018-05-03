package cmdService

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/service"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
	"gopkg.in/yaml.v2"
)

// Init run the Init command for a service
var Init = &cobra.Command{
	Use:   "init",
	Short: "Initialize a service",
	Long: `Initialize a service by creating a mesg.yml and Dockerfile in a dedicated folder.
	
To get more information, see the page [service file from the documentation](https://docs.mesg.tech/service/service-file.html)`,
	Example: `mesg-cli service init
mesg-cli service init --name NAME --description DESCRIPTION --visibility ALL --publish ALL`,
	Run:               initHandler,
	DisableAutoGenTag: true,
}

func initHandler(cmd *cobra.Command, args []string) {
	res := service.Service{}
	res.Publish = string(service.PublishAll)
	res.Visibility = string(service.VisibilityAll)

	fmt.Printf("%s\n", aurora.Bold("Initialization of a new service"))

	res.Name = cmd.Flag("name").Value.String()
	if res.Name == "" && survey.AskOne(&survey.Input{Message: "Name:"}, &res.Name, nil) != nil {
		os.Exit(0)
	}
	key := strings.Replace(strings.ToLower(res.Name), " ", "-", -1)

	res.Description = cmd.Flag("description").Value.String()
	if res.Description == "" && survey.AskOne(&survey.Input{Message: "Description:"}, &res.Description, nil) != nil {
		os.Exit(0)
	}

	publishStr := cmd.Flag("publish").Value.String()
	if publishStr == "" && survey.AskOne(&survey.Select{
		Message: "Publish (ALL):",
		Options: []string{
			string(service.PublishAll),
			string(service.PublishSource),
			string(service.PublishContainer),
			string(service.PublishNone),
		},
	}, &publishStr, nil) != nil {
		os.Exit(0)
	}
	res.Publish = publishStr

	visibilityStr := cmd.Flag("visibility").Value.String()
	if visibilityStr == "" && survey.AskOne(&survey.Select{
		Message: "Visibility (ALL):",
		Options: []string{
			string(service.VisibilityAll),
			string(service.VisibilityUsers),
			string(service.VisibilityWorkers),
			string(service.VisibilityNone),
		},
	}, &visibilityStr, nil) != nil {
		os.Exit(0)
	}
	res.Visibility = visibilityStr

	res.Dependencies = map[string]*service.Dependency{
		key: &service.Dependency{
			Image: strings.Join([]string{"mesg", key}, "/"),
		},
	}

	out, err := yaml.Marshal(res)
	if err != nil {
		panic(err)
	}
	fmt.Println()
	fmt.Printf("%s\n", aurora.Brown("Summary:").Bold())
	fmt.Printf("%s\n", string(out))
	// fmt.Print(string(out))
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s %s\n", aurora.Bold("Service will be created in:"), aurora.Brown(filepath.Join(dir, key)))

	ok := false
	if survey.AskOne(&survey.Confirm{Message: "Is this correct?", Default: true}, &ok, nil) != nil {
		os.Exit(0)
	}
	if !ok {
		return
	}
	if cmd.Flag("current").Value.String() == "true" {
		err = ioutil.WriteFile("mesg.yml", out, os.ModePerm)
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile("Dockerfile", []byte(""), os.ModePerm)
		if err != nil {
			panic(err)
		}
	} else {
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

	fmt.Printf("%s\n", aurora.Green("Service created with success").Bold())
}

func init() {
	Init.Flags().StringP("name", "n", "", "Name")
	Init.Flags().StringP("description", "d", "", "Description")
	Init.Flags().StringP("visibility", "v", "", "Visibility")
	Init.Flags().StringP("publish", "p", "", "Publish")
	Init.Flags().BoolP("current", "c", false, "Create the service in the current path")
}
