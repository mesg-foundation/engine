package cmdWorkflow

import (
	"fmt"

	"github.com/mesg-foundation/application/cmd/utils"

	"github.com/spf13/cobra"
)

// Log workflow executions
var Log = &cobra.Command{
	Use:               "log ID",
	Short:             "Log details of a workflow",
	Args:              cobra.MinimumNArgs(1),
	Example:           "mesg-cli workflow log XX",
	Run:               logHandler,
	DisableAutoGenTag: true,
}

func logHandler(cmd *cobra.Command, args []string) {
	account := cmdUtils.AccountFromFlagOrAsk(cmd, "Select an account you want to log with")

	execution := cmd.Flag("execution").Value.String()
	task := cmd.Flag("task").Value.String()
	from := cmd.Flag("from").Value.String()
	to := cmd.Flag("to").Value.String()

	fmt.Printf("Loggin results with account %s of workflow %s, execution: %s, task: %s, from: %s, to: %s", account, args[0], execution, task, from, to)
}

func init() {
	Log.Flags().StringP("execution", "e", "", "Log a specific execution of the workflow")
	Log.Flags().StringP("task", "t", "", "Log a specific task of an execution of the workflow")
	Log.Flags().StringP("from", "", "", "From date in ISO format")
	Log.Flags().StringP("to", "", "", "To date in ISO format")
	cmdUtils.Accountable(Log)
}
