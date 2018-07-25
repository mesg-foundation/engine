package service

import (
	"github.com/spf13/cobra"
)

// Test a service
var Test = &cobra.Command{
	Use: "test",
	Deprecated: `please use the following commands:
mesg-core service dev
mesg-core service execute`,
	DisableAutoGenTag: true,
}
