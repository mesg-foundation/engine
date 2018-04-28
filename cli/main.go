package main

import (
	"fmt"

	"github.com/mesg-foundation/core/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
