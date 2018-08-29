package commands

import (
	"github.com/spf13/cobra"
)

type serviceTestCmd struct {
	baseCmd
}

func newServiceTestCmd() *serviceTestCmd {
	c := &serviceTestCmd{}
	c.cmd = newCommand(&cobra.Command{
		Use: "test",
		Deprecated: `please use the following commands:
mesg-core service dev
mesg-core service execute`,
	})
	return c
}
