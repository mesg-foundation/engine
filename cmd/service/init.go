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
	Example: `mesg-core service init
mesg-core service init --name NAME --description DESCRIPTION --visibility ALL --publish ALL`,
	Run:               initHandler,
	DisableAutoGenTag: true,
}

func initHandler(cmd *cobra.Command, args []string) {
	fmt.Printf("%s\n", aurora.Bold("Initialization of a new service"))

	res := buildService(cmd)

	out, err := yaml.Marshal(res)
	if err != nil {
		panic(err)
	}
	fmt.Println()
	fmt.Printf("%s\n", aurora.Brown("Summary:").Bold())
	fmt.Printf("%s\n", string(out))

	ok := false
	if survey.AskOne(&survey.Confirm{Message: "Is this correct?", Default: true}, &ok, nil) != nil {
		os.Exit(0)
	}
	if !ok {
		return
	}

	if cmd.Flag("current").Value.String() == "true" {
		err = writeInCurrentFolder(out)
	} else {
		err = writeInFolder(strings.Replace(strings.ToLower(res.Name), " ", "-", -1), out)
	}
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", aurora.Green("Service created with success").Bold())
}

func askOpts(label string, value string, opts []string) string {
	if value == "" && survey.AskOne(&survey.Select{
		Message: label,
		Options: opts,
	}, &value, nil) != nil {
		os.Exit(0)
	}
	return value
}

func ask(label string, value string) string {
	if value != "" {
		return value
	}
	if survey.AskOne(&survey.Input{Message: label}, &value, nil) != nil {
		os.Exit(0)
	}
	return value
}

func buildService(cmd *cobra.Command) (res service.Service) {
	res.Name = ask("Name:", cmd.Flag("name").Value.String())
	res.Publish = string(service.PublishAll)
	res.Visibility = string(service.VisibilityAll)
	res.Description = ask("Description:", cmd.Flag("description").Value.String())
	res.Visibility = askOpts("Visibility (ALL):", cmd.Flag("visibility").Value.String(), []string{
		string(service.VisibilityAll),
		string(service.VisibilityUsers),
		string(service.VisibilityWorkers),
		string(service.VisibilityNone),
	})
	res.Publish = askOpts("Publish (ALL):", cmd.Flag("publish").Value.String(), []string{
		string(service.PublishAll),
		string(service.PublishSource),
		string(service.PublishContainer),
		string(service.PublishNone),
	})
	return
}

func writeInCurrentFolder(content []byte) (err error) {
	return writeInFolder("./", content)
}

func writeInFolder(folder string, content []byte) (err error) {
	if folder != "./" {
		err = os.Mkdir(folder, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	err = ioutil.WriteFile(filepath.Join(folder, "mesg.yml"), content, os.ModePerm)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filepath.Join(folder, "Dockerfile"), []byte(""), os.ModePerm)
	if err != nil {
		panic(err)
	}
	return
}

func init() {
	Init.Flags().StringP("name", "n", "", "Name")
	Init.Flags().StringP("description", "d", "", "Description")
	Init.Flags().StringP("visibility", "v", "", "Visibility")
	Init.Flags().StringP("publish", "p", "", "Publish")
	Init.Flags().BoolP("current", "c", false, "Create the service in the current path")
}
