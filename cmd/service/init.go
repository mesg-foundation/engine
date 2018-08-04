package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

const templatesURL = "https://raw.githubusercontent.com/mesg-foundation/awesome/master/templates.json"
const addMyOwn = "Add my own"
const custom = "Enter template URL"

type templateStruct struct {
	Name string
	URL  string
}

// Init run the Init command for a service
var Init = &cobra.Command{
	Use:   "init",
	Short: "Initialize a service",
	Long: `Initialize a service by creating a mesg.yml and Dockerfile in a dedicated folder.
	
To get more information, see the page [service file from the documentation](https://docs.mesg.com/guide/service/service-file.html)`,
	Example: `mesg-core service init
mesg-core service init --name NAME --description DESCRIPTION
mesg-core service init --current`,
	Run:               initHandler,
	DisableAutoGenTag: true,
}

func init() {
	Init.Flags().StringP("name", "n", "", "Name")
	Init.Flags().StringP("description", "d", "", "Description")
	Init.Flags().BoolP("current", "c", false, "Create the service in the current path")
	Init.Flags().StringP("template", "t", "", "Specify the template URL to use")
}

func initHandler(cmd *cobra.Command, args []string) {
	fmt.Printf("%s\n", aurora.Bold("Initialization of a new service"))

	tmpl := &templateStruct{
		URL:  cmd.Flag("template").Value.String(),
		Name: cmd.Flag("template").Value.String(),
	}
	if tmpl.URL == "" {
		templates, err := getTemplates(templatesURL)
		utils.HandleError(err)
		tmpl, err = selectTemplate(templates)
		utils.HandleError(err)
		if tmpl == nil {
			os.Exit(0)
		}
	}
	path, err := downloadTemplate(tmpl)
	utils.HandleError(err)
	defer os.RemoveAll(path)
	replacements := askReplacements(cmd)
	folder := strings.Replace(strings.ToLower(replacements["name"]), " ", "-", -1)
	if cmd.Flag("current").Value.String() == "true" {
		folder = "./"
	}
	err = copyDir(path+"/template", folder, replacements)
	utils.HandleError(err)
	fmt.Println(aurora.Green("Service created in folder: " + folder))
}

func getTemplates(url string) ([]*templateStruct, error) {
	client := http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var templates []*templateStruct
	return templates, json.Unmarshal(body, &templates)
}

func selectTemplate(templates []*templateStruct) (*templateStruct, error) {
	var result string
	if survey.AskOne(&survey.Select{
		Message: "Select a template to use",
		Options: templatesToOption(templates),
	}, &result, nil) != nil {
		os.Exit(0)
	}
	return getTemplateResult(result, templates), nil
}

func templatesToOption(templates []*templateStruct) (options []string) {
	options = []string{}
	for _, template := range templates {
		options = append(options, template.Name+" ("+template.URL+")")
	}
	options = append(options, custom)
	options = append(options, addMyOwn)
	return
}

func getTemplateResult(result string, templates []*templateStruct) (tmpl *templateStruct) {
	if result == addMyOwn {
		fmt.Println(aurora.Green("You can create and add your own template to this list. Go to the Awesome Github to see how: https://github.com/mesg-foundation/awesome"))
		return
	}
	if result == custom {
		var url string
		if survey.AskOne(&survey.Input{Message: "Enter template URL"}, &url, nil) != nil {
			os.Exit(0)
		}
		tmpl = &templateStruct{
			URL:  url,
			Name: url,
		}
	}
	for _, template := range templates {
		if template.Name+" ("+template.URL+")" == result {
			tmpl = template
			break
		}
	}
	return
}

func downloadTemplate(tmpl *templateStruct) (path string, err error) {
	path, err = createTempFolder()
	if err != nil {
		return "", err
	}

	return path, gitClone(tmpl.URL, path, "Downloading template "+tmpl.Name+"...")
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

func askReplacements(cmd *cobra.Command) (replacement map[string]string) {
	replacement = make(map[string]string)
	replacement["name"] = ask("Name:", cmd.Flag("name").Value.String(), survey.Required)
	replacement["description"] = ask("Description:", cmd.Flag("description").Value.String(), nil)
	return
}

func copyDir(src string, dst string, replacement map[string]string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)
	err = os.MkdirAll(dst, os.ModePerm)
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = copyDir(srcPath, dstPath, replacement)
			if err != nil {
				break
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = copyFile(srcPath, dstPath, replacement)
			if err != nil {
				break
			}
		}
	}

	return

}

func copyFile(src, dst string, replacement map[string]string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	return transform(dst, src, replacement)
}

func transform(dest string, source string, replacement map[string]string) error {
	body, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}
	res := string(body)
	for key, value := range replacement {
		res = strings.Replace(res, "{{"+key+"}}", value, -1)
	}
	si, err := os.Stat(source)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(dest, []byte(res), si.Mode())
}
