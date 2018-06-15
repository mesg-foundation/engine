package main

import (
	"github.com/mesg-foundation/core/cmd"
	"github.com/mesg-foundation/core/cmd/utils"
)

func main() {
	err := cmd.RootCmd.Execute()
	cmdUtils.HandleError(err)
}
