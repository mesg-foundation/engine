package main

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	cliDoc := "./docs/cli"
	os.Mkdir(cliDoc, os.ModePerm)
	if err := doc.GenMarkdownTree(cmd.RootCmd, cliDoc); err != nil {
		fmt.Println(aurora.Red(err))
	}
}
