package main

import (
	"github.com/mesg-foundation/core/cmd"
	"github.com/mesg-foundation/core/cmd/utils"
)

// Version of this release. Will be replaced automatically when compiling in CI
var Version = "vX.X.X"

func init() {
	cmd.RootCmd.Version = Version
	cmd.RootCmd.Short = cmd.RootCmd.Short + " " + Version
}

func main() {
	err := cmd.RootCmd.Execute()
	utils.HandleError(err)
}
