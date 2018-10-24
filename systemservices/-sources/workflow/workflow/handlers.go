package workflow

// WorkflowDocument combines a workflow with additional info.
type WorkflowDocument struct {
	// ID is the unique id for workflow.
	ID string

	// Name is the optionally set unique name for workflow.
	Name string

	// Definition of workflow.
	Definition WorkflowDefinition
}

// errorOutput is the error output data.
type errorOutput struct {
	Message string `json:"message"`
}
