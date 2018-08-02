package main

import (
	"log"

	"github.com/mesg-foundation/core/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	log.Fatalln(doc.GenMarkdownTree(cmd.RootCmd, "./docs/cli"))
}
