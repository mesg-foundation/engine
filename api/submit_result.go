package api

import (
	"encoding/json"
	"fmt"

	"github.com/mesg-foundation/core/database"
	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/pubsub"
)

// SubmitResult submits results for executionID.
func (api *API) SubmitResult(executionID string, outputKey string, outputData []byte) error {
	exec, stateChanged, err := api.processExecution(executionID, outputKey, outputData)
	if stateChanged {
		// only publish to listeners when the execution's state changed.
		go pubsub.Publish(exec.Service.ResultSubscriptionChannel(), exec)
	}
	// always return any error to the service.
	return err
}

// processExecution processes execution and marks it as complated or failed.
func (api *API) processExecution(executionID string, outputKey string, outputData []byte) (exec *execution.Execution, stateChanged bool, err error) {
	tx, err := api.execDB.OpenTransaction()
	if err != nil {
		return nil, false, err
	}

	exec, err = tx.Find(executionID)
	if err != nil {
		tx.Discard()
		return nil, false, err
	}

	var outputDataMap map[string]interface{}
	if err := json.Unmarshal(outputData, &outputDataMap); err != nil {
		return saveExecution(tx, exec, fmt.Errorf("invalid output data error: %s", err))
	}

	err = exec.Service.ValidateTaskOutput(exec.TaskKey, outputKey, outputDataMap)
	if err != nil {
		return saveExecution(tx, exec, err)
	}

	err = exec.Complete(outputKey, outputDataMap)
	return saveExecution(tx, exec, err)
}

func saveExecution(tx database.ExecutionTransaction, exec *execution.Execution, err error) (execOut *execution.Execution, stateChanged bool, errOut error) {
	if err != nil {
		if errFailed := exec.Failed(err); errFailed != nil {
			tx.Discard()
			return exec, false, errFailed
		}
	}
	if errSave := tx.Save(exec); errSave != nil {
		tx.Discard()
		return exec, true, errSave
	}
	return exec, true, tx.Commit()
}
