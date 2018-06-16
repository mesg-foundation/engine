package main

import (
	"github.com/mesg-foundation/core/cmd"
	"github.com/mesg-foundation/core/cmd/utils"
)

// version of this release. Will be replaced automatically when compiling in CI
var version = "vX.X.X"

func init() {
	cmd.RootCmd.Version = version
	cmd.RootCmd.Short = cmd.RootCmd.Short + " " + version
}

func main() {
	err := cmd.RootCmd.Execute()
	cmdUtils.HandleError(err)
}
