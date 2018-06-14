package main

import (
	"fmt"

	"github.com/mesg-foundation/core/cmd"
)

// version of this release. Will be replaced automatically when compiling in CI
var version = "v1.0.0-beta"

func init() {
	cmd.RootCmd.Version = version
	cmd.RootCmd.Short = cmd.RootCmd.Short + " " + version
}

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
