package processesdk

import (
	"fmt"

	"github.com/mesg-foundation/engine/database"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/process"
	instancesdk "github.com/mesg-foundation/engine/sdk/instance"
)

// Process exposes process APIs of MESG.
type Process struct {
	processDB database.ProcessDB
	instance  *instancesdk.Instance
}

// New creates a new Process SDK with given options.
func New(instance *instancesdk.Instance, processDB database.ProcessDB) *Process {
	return &Process{
		processDB: processDB,
		instance:  instance,
	}
}

// Create creates a new service from definition.
func (w *Process) Create(wf *process.Process) (*process.Process, error) {
	wf.Hash = hash.Dump(wf)

	for _, node := range wf.Nodes {
		switch n := node.Type.(type) {
		case *process.Process_Node_Result_:
			if _, err := w.instance.Get(n.Result.InstanceHash); err != nil {
				return nil, err
			}
		case *process.Process_Node_Event_:
			if _, err := w.instance.Get(n.Event.InstanceHash); err != nil {
				return nil, err
			}
		case *process.Process_Node_Task_:
			if _, err := w.instance.Get(n.Task.InstanceHash); err != nil {
				return nil, err
			}
		}
	}

	// check if process already exists.
	if _, err := w.processDB.Get(wf.Hash); err == nil {
		return nil, &AlreadyExistsError{Hash: wf.Hash}
	}

	if err := wf.Validate(); err != nil {
		return nil, err
	}
	return wf, w.processDB.Save(wf)
}

// Delete deletes the process by hash.
func (w *Process) Delete(hash hash.Hash) error {
	return w.processDB.Delete(hash)
}

// Get returns the process that matches given hash.
func (w *Process) Get(hash hash.Hash) (*process.Process, error) {
	return w.processDB.Get(hash)
}

// List returns all processes.
func (w *Process) List() ([]*process.Process, error) {
	return w.processDB.All()
}

// AlreadyExistsError is an not found error.
type AlreadyExistsError struct {
	Hash hash.Hash
}

func (e *AlreadyExistsError) Error() string {
	return fmt.Sprintf("process %q already exists", e.Hash.String())
}
