package cmdWorkflow

import (
	"fmt"
	"time"

	"github.com/kyokomi/emoji"
	"github.com/mesg-foundation/application/cmd/utils"

	"github.com/spf13/cobra"
)

// Test run the test command for a workflow
var Test = &cobra.Command{
	Use:               "test",
	Short:             "Test a workflow",
	Args:              cobra.MinimumNArgs(1),
	Example:           "mesg-cli workflow test workflow.yml",
	Run:               testHandler,
	DisableAutoGenTag: true,
}

func testHandler(cmd *cobra.Command, args []string) {
	s := cmdUtils.StartSpinner(cmdUtils.SpinnerOptions{Text: "Preparing testing environment..."})
	time.Sleep(2 * time.Second)
	s.Stop()
	s = cmdUtils.StartSpinner(cmdUtils.SpinnerOptions{Text: "Task #1: processing..."})
	time.Sleep(time.Second)
	s.Stop()
	fmt.Println(emoji.Sprint(":white_check_mark: Task #1: onSuccess(foo = 12, bar = 23)"))
	s = cmdUtils.StartSpinner(cmdUtils.SpinnerOptions{Text: "Task #2: processing..."})
	time.Sleep(time.Second)
	s.Stop()
	fmt.Println(emoji.Sprint(":white_check_mark: Task #2: onSent(id = 123)"))
	// TODO test the workflow
}

func init() {
	Test.Flags().BoolP("live", "l", false, "Use live events")
	Test.Flags().StringP("event", "e", "", "Path to the event file")
}
