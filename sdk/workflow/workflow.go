package workflowsdk

import (
	"fmt"

	"github.com/mesg-foundation/engine/database"
	"github.com/mesg-foundation/engine/hash"
	instancesdk "github.com/mesg-foundation/engine/sdk/instance"
	"github.com/mesg-foundation/engine/workflow"
	validator "gopkg.in/go-playground/validator.v9"
)

// Workflow exposes workflow APIs of MESG.
type Workflow struct {
	workflowDB database.WorkflowDB
	instance   *instancesdk.Instance
}

// New creates a new Workflow SDK with given options.
func New(instance *instancesdk.Instance, workflowDB database.WorkflowDB) *Workflow {
	return &Workflow{
		workflowDB: workflowDB,
		instance:   instance,
	}
}

// Create creates a new service from definition.
func (w *Workflow) Create(wf *workflow.Workflow) (*workflow.Workflow, error) {
	wf.Hash = hash.Dump(wf)

	if _, err := w.instance.Get(wf.Trigger.InstanceHash); err != nil {
		return nil, err
	}
	for _, node := range wf.Nodes {
		if _, err := w.instance.Get(node.InstanceHash); err != nil {
			return nil, err
		}
	}

	// check if workflow already exists.
	if _, err := w.workflowDB.Get(wf.Hash); err == nil {
		return nil, &AlreadyExistsError{Hash: wf.Hash}
	}

	if err := validator.New().Struct(wf); err != nil {
		return nil, err
	}
	return wf, w.workflowDB.Save(wf)
}

// Delete deletes the workflow by hash.
func (w *Workflow) Delete(hash hash.Hash) error {
	return w.workflowDB.Delete(hash)
}

// Get returns the workflow that matches given hash.
func (w *Workflow) Get(hash hash.Hash) (*workflow.Workflow, error) {
	return w.workflowDB.Get(hash)
}

// List returns all workflows.
func (w *Workflow) List() ([]*workflow.Workflow, error) {
	return w.workflowDB.All()
}

// AlreadyExistsError is an not found error.
type AlreadyExistsError struct {
	Hash hash.Hash
}

func (e *AlreadyExistsError) Error() string {
	return fmt.Sprintf("workflow %q already exists", e.Hash.String())
}
