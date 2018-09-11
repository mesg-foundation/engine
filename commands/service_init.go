package commands

import (
	"fmt"
	"os"

	"github.com/mesg-foundation/core/utils/servicetemplate"
	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

type serviceInitCmd struct {
	baseCmd

	name         string
	description  string
	templateURL  string
	templateName string
	dir          string

	e ServiceExecutor
}

// service template select options.
const (
	addMyOwn  = "Add my own"
	customURL = "Enter template URL"
)

func newServiceInitCmd(e ServiceExecutor) *serviceInitCmd {
	c := &serviceInitCmd{e: e}
	c.cmd = newCommand(&cobra.Command{
		Use:   "init",
		Short: "Initialize a service",
		Long: `Initialize a service by creating a mesg.yml and Dockerfile in a dedicated directory.
	
To get more information, see the page [service file from the documentation](https://docs.mesg.com/guide/service/service-file.html)`,
		Example: `mesg-core service init
mesg-core service init --name NAME --description DESCRIPTION
mesg-core service init --current`,
		Args:    cobra.NoArgs,
		PreRunE: c.preRunE,
		RunE:    c.runE,
	})
	c.cmd.Flags().StringVar(&c.dir, "dir", c.dir, "Create the service in the directory")
	c.cmd.Flags().StringVarP(&c.templateURL, "template", "t", c.templateURL, "Specify the template URL to use")
	return c
}

func (c *serviceInitCmd) preRunE(cmd *cobra.Command, args []string) error {
	if err := c.selectOutputDirectory(); err != nil {
		return err
	}

	if c.templateURL != "" {
		c.templateName = c.templateURL
		return nil
	}

	// no template specify - download and select one
	list, err := c.e.ServiceInitTemplateList()
	if err != nil {
		return err
	}

	var result string
	if err := survey.AskOne(&survey.Select{
		Message: "Select a template to use",
		Options: templatesToOptions(list),
	}, &result, nil); err != nil {
		return err
	}

	if result == addMyOwn {
		fmt.Println("You can create and add your own template to this list.")
		fmt.Println("Go to the Awesome Github to see how")
		fmt.Println("https://github.com/mesg-foundation/awesome")
		os.Exit(0)
	}

	if result == customURL {
		if err := survey.AskOne(&survey.Input{
			Message: "Enter template URL",
		}, &c.templateURL, nil); err != nil {
			return err
		}
	}

	for _, l := range list {
		if result == l.String() {
			c.templateURL = l.URL
			c.templateName = l.Name
		}
	}
	return nil
}

func (c *serviceInitCmd) runE(cmd *cobra.Command, args []string) error {
	return c.e.ServiceInitDownloadTemplate(&servicetemplate.Template{
		Name: c.templateName,
		URL:  c.templateURL,
	}, c.dir)
}

func (c *serviceInitCmd) selectOutputDirectory() error {
	if c.dir != "" {
		return nil
	}
	defval := c.name
	if defval == "" {
		defval = "."
	}

	return survey.AskOne(&survey.Input{
		Message: "Enter the output directory",
		Default: defval,
	}, &c.dir, nil)
}

func templatesToOptions(templates []*servicetemplate.Template) []string {
	var options []string
	for _, template := range templates {
		options = append(options, template.String())
	}
	options = append(options, customURL)
	options = append(options, addMyOwn)
	return options
}
