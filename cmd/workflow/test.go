package workflow

import (
	"fmt"
	"time"

	"github.com/kyokomi/emoji"
	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/spf13/cobra"
)

// Test run the test command for a workflow
var Test = &cobra.Command{
	Use:   "test ./PATH_TO_WORKFLOW_FILE",
	Short: "Test a workflow",
	Long: `Test a workflow locally

To get more information, see the [test page from the documentation](https://docs.mesg.tech/workflow/test.html)`,
	Args: cobra.MinimumNArgs(1),
	Example: `mesg-core workflow test ./PATH_TO_WORKFLOW_FILE.yml
mesg-core workflow test ./PATH_TO_WORKFLOW_FILE.yml --live
mesg-core workflow test ./PATH_TO_WORKFLOW_FILE.yml --task TASK_ID --live
mesg-core workflow test ./PATH_TO_WORKFLOW_FILE.yml --task TASK_ID --event ./PATH_TO_EVENT_DATA_FILE.yml
mesg-core workflow test ./PATH_TO_WORKFLOW_FILE.yml --live --keep-alive`,
	Run:               testHandler,
	DisableAutoGenTag: true,
}

func testHandler(cmd *cobra.Command, args []string) {
	s := utils.StartSpinner(utils.SpinnerOptions{Text: "Preparing testing environment..."})
	time.Sleep(2 * time.Second)
	s.Stop()
	s = utils.StartSpinner(utils.SpinnerOptions{Text: "Task #1: processing..."})
	time.Sleep(time.Second)
	s.Stop()
	fmt.Println(emoji.Sprint(":white_check_mark: Task #1: onSuccess(foo = 12, bar = 23)"))
	if cmd.Flag("task").Value.String() == "" {
		s = utils.StartSpinner(utils.SpinnerOptions{Text: "Task #2: processing..."})
		time.Sleep(time.Second)
		s.Stop()
		fmt.Println(emoji.Sprint(":white_check_mark: Task #2: onSent(id = 123)"))
	} else {
		fmt.Println("Bypass other tasks")
	}
	// TODO test the workflow
}

func init() {
	Test.Flags().BoolP("live", "l", false, "Use live events")
	Test.Flags().StringP("event", "e", "", "Path to the event data file")
	Test.Flags().StringP("task", "t", "", "Run the test on a specific task of the workflow")
	Test.Flags().BoolP("keep-alive", "k", false, "Keep the services alive (re-run without the option to stop them)")
}
