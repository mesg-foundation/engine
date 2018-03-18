package main

import (
	"fmt"

	"github.com/mesg-foundation/application/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
	if err := doc.GenMarkdownTree(cmd.RootCmd, "./docs/cli"); err != nil {
		fmt.Println(err)
	}

}
