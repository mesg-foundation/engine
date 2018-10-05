package main

import (
	"log"

	"github.com/mesg-foundation/core/commands"
	"github.com/spf13/cobra/doc"
)

func main() {
	log.Fatalln(doc.GenMarkdownTree(commands.Build(nil), "./docs/cli"))
}
