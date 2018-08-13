package main

import (
	"github.com/mesg-foundation/core/cmd"
	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/mesg-foundation/core/version"
)

func init() {
	cmd.RootCmd.Version = version.Version
	cmd.RootCmd.Short = cmd.RootCmd.Short + " " + version.Version
}

func main() {
	err := cmd.RootCmd.Execute()
	utils.HandleError(err)
}
