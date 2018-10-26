package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/mesg-foundation/core/commands/provider"
	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/mesg-foundation/core/x/xsignal"
	"github.com/spf13/cobra"
)

type workflowLogsCmd struct {
	baseCmd

	e WorkflowExecutor
}

func newWorkflowLogsCmd(e WorkflowExecutor) *workflowLogsCmd {
	c := &workflowLogsCmd{
		e: e,
	}

	c.cmd = newCommand(&cobra.Command{
		Use:     "logs",
		Short:   "Show logs of a workflow",
		Example: `mesg-core workflow logs WORKFLOW`,
		Args:    cobra.ExactArgs(1),
		RunE:    c.runE,
	})

	return c
}

func (c *workflowLogsCmd) runE(cmd *cobra.Command, args []string) error {
	closer, err := c.showLogs(c.e, args[0])
	if err != nil {
		return err
	}
	defer closer()

	<-xsignal.WaitForInterrupt()
	return nil
}

// WorkflowLog keeps workflow logs.
type WorkflowLog struct {
	RunStart                      bool   `json:"runStart"`
	WorkflowID                    string `json:"workflowID"`
	WorkflowName                  string `json:"workflowName"`
	WorkflowDefinitionName        string `json:"workflowDefinitionName"`
	WorkflowDefinitionDescription string `json:"workflowDefinitionDescription"`
}

// EventLog keeps event logs.
type EventLog struct {
	WorkflowID    string      `json:"workflowID"`
	WorkflowName  string      `json:"workflowName"`
	ServiceName   string      `json:"serviceName"`
	EventKey      string      `json:"eventKey"`
	ExecutionData interface{} `json:"executionData"`
}

// ExecutionLog keeps execution logs.
type ExecutionLog struct {
	WorkflowID   string `json:"workflowID"`
	WorkflowName string `json:"workflowName"`
	ServiceName  string `json:"serviceName"`
	TaskKey      string `json:"taskKey"`
}

// logLine represents a log line received from log stream.
type logLine struct {
	Workflow  *WorkflowLog  `json:"workflow"`
	Event     *EventLog     `json:"event"`
	Execution *ExecutionLog `json:"execution"`
}

// showLogs prints logs for workflowID.
func (c *workflowLogsCmd) showLogs(e WorkflowExecutor, workflowID string) (closer func(), err error) {
	var (
		log *provider.WorkflowLog
	)
	pretty.Progress("Loading logs...", func() {
		log, closer, err = e.WorkflowLogs(workflowID)
	})
	if err != nil {
		return nil, err
	}

	go c.printLog(os.Stdout, log.Standard)
	go c.printLog(os.Stderr, log.Error)

	return closer, nil
}

// printLog prints logs from workflow's r stream by its log types.
func (c *workflowLogsCmd) printLog(out io.Writer, r io.Reader) {
	dc := json.NewDecoder(r)

	for {
		var line *logLine
		if err := dc.Decode(&line); err != nil {
			fmt.Println(err)
			return
		}
		switch {
		case line.Workflow != nil:
			c.printWorkflowLog(out, line.Workflow)
		case line.Event != nil:
			c.printEventLog(out, line.Event)
		case line.Execution != nil:
			c.printExecutionLog(out, line.Execution)
		}
	}
}

// coloring for logs.
var (
	colorGreen  = color.New(color.FgGreen)
	colorYellow = color.New(color.FgYellow)
	colorBold   = color.New(color.Bold)

	colorAttention = color.New(color.FgYellow, color.Bold)
	colorInfo      = color.New(color.FgBlue)
	colorError     = color.New(color.FgRed, color.Bold)
)

// printWorkflowLog prints logs related directly with workflow.
func (c *workflowLogsCmd) printWorkflowLog(out io.Writer, workflow *WorkflowLog) {
	fmt.Println(colorGreen.Sprintf("✔ %s workflow started", colorBold.Sprintf("%s", workflow.WorkflowDefinitionName)))
	fmt.Println(colorYellow.Sprintf("%s", workflow.WorkflowDefinitionDescription))
}

// printWorkflowLog prints logs related with workflow's events.
func (c *workflowLogsCmd) printEventLog(out io.Writer, event *EventLog) error {
	data, err := json.MarshalIndent(event.ExecutionData, "", "  ")
	if err != nil {
		return err
	}

	fmt.Printf(">> event %s received on %s service, execution data will be: %s\n",
		colorAttention.Sprintf("%s", event.EventKey),
		colorAttention.Sprintf("%s", event.ServiceName),
		colorInfo.Sprintf(" %+v", string(data)))

	return nil
}

// printWorkflowLog prints logs related with workflow's executions.
func (c *workflowLogsCmd) printExecutionLog(out io.Writer, execution *ExecutionLog) {
	fmt.Printf("<< execution successfully made for %s task on %s service\n",
		colorAttention.Sprintf("%s", execution.TaskKey),
		colorAttention.Sprintf("%s", execution.ServiceName))
}
