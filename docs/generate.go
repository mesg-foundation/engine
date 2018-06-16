package main

import (
	"os"

	"github.com/mesg-foundation/core/cmd"
	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/spf13/cobra/doc"
)

func main() {
	cliDoc := "./docs/cli"
	os.Mkdir(cliDoc, os.ModePerm)
	err := doc.GenMarkdownTree(cmd.RootCmd, cliDoc)
	utils.HandleError(err)
}
