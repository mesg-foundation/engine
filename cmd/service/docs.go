package service

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/cmd/service/assets"
	"github.com/mesg-foundation/core/cmd/utils"
	s "github.com/mesg-foundation/core/service"
	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// Docs a service to the marketplace
var Docs = &cobra.Command{
	Use:   "gen-doc",
	Short: "Generate the documentation for the service in a README.md file",
	Example: `mesg-core service gen-doc
mesg-core service gen-doc ./PATH_TO_SERVICE`,
	Run:               genDocHandler,
	DisableAutoGenTag: true,
}

func genDocHandler(cmd *cobra.Command, args []string) {
	path := defaultPath(args)
	readmePath := filepath.Join(path, "README.md")
	if _, err := os.Stat(readmePath); err == nil {
		var value bool
		if survey.AskOne(&survey.Confirm{Message: "The file README.md already exists. Do you want to overwrite it?"}, &value, nil) != nil {
			return
		}
		if !value {
			return
		}
	}

	service, err := s.ImportFromPath(path)
	utils.HandleError(err)

	f, err := os.OpenFile(readmePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	defer f.Close()

	readmeTemplate, err := assets.Asset("cmd/service/assets/readmeTemplate.md")
	utils.HandleError(err)

	tmpl, err := template.New("doc").Parse(string(readmeTemplate))
	utils.HandleError(err)
	err = tmpl.Execute(f, service)
	utils.HandleError(err)

	fmt.Println(aurora.Green("File README.md generated with success"))
}
