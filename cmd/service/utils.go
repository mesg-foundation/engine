package service

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/interface/grpc/core"
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

func gitClone(repoURL string, path string, message string) error {
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

func createTempFolder() (path string, err error) {
	return ioutil.TempDir("", "mesg-")
}
