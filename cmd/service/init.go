package service

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/mesg-foundation/core/service"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

const templateText = `name: "{{.Name}}"
description: "{{.Description}}"
tasks:
  foo:
    name: "Foo"
    inputs:
      inputA:
        type: String
      inputB:
        type: Number
    outputs:
      outputX:
        data:
          resultX:
            type: String
events:
  eventX:
    data:
      dataA:
        type: String
`

// Init run the Init command for a service
var Init = &cobra.Command{
	Use:   "init",
	Short: "Initialize a service",
	Long: `Initialize a service by creating a mesg.yml and Dockerfile in a dedicated folder.
	
To get more information, see the page [service file from the documentation](https://docs.mesg.com/service/service-file.html)`,
	Example: `mesg-core service init
mesg-core service init --name NAME --description DESCRIPTION
mesg-core service init --current`,
	Run:               initHandler,
	DisableAutoGenTag: true,
}

func initHandler(cmd *cobra.Command, args []string) {
	fmt.Printf("%s\n", aurora.Bold("Initialization of a new service"))

	res := buildService(cmd)

	mesgFile, err := generateMesgFile(res)
	utils.HandleError(err)

	ok := false
	fmt.Println()
	fmt.Println(string(mesgFile))
	if survey.AskOne(&survey.Confirm{Message: "Is this correct?", Default: true}, &ok, nil) != nil {
		os.Exit(0)
	}
	if !ok {
		return
	}

	folder := strings.Replace(strings.ToLower(res.Name), " ", "-", -1)
	if cmd.Flag("current").Value.String() == "true" {
		folder = "./"
	}
	err = writeInFolder(folder, mesgFile)
	utils.HandleError(err)
	fmt.Printf("%s\n", aurora.Green("Service created with success in folder '"+folder+"'").Bold())
}

func generateMesgFile(service *service.Service) (res []byte, err error) {
	var doc bytes.Buffer
	tmpl, err := template.New("service-init").Parse(templateText)
	if err != nil {
		return
	}
	err = tmpl.Execute(&doc, service)
	res = doc.Bytes()
	return
}

func ask(label string, value string, validator survey.Validator) string {
	if value != "" {
		return value
	}
	if survey.AskOne(&survey.Input{Message: label}, &value, validator) != nil {
		os.Exit(0)
	}
	return value
}

func buildService(cmd *cobra.Command) (res *service.Service) {
	res = &service.Service{}
	res.Name = ask("Name:", cmd.Flag("name").Value.String(), survey.Required)
	res.Description = ask("Description:", cmd.Flag("description").Value.String(), nil)
	return
}

func writeInFolder(folder string, content []byte) (err error) {
	if folder != "./" {
		err = os.Mkdir(folder, os.ModePerm)
		if err != nil {
			return
		}
	}
	err = ioutil.WriteFile(filepath.Join(folder, "mesg.yml"), content, os.ModePerm)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(filepath.Join(folder, "Dockerfile"), []byte(""), os.ModePerm)
	if err != nil {
		return
	}
	return
}

func init() {
	Init.Flags().StringP("name", "n", "", "Name")
	Init.Flags().StringP("description", "d", "", "Description")
	Init.Flags().BoolP("current", "c", false, "Create the service in the current path")
}
