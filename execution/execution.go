package execution

import (
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/types"
)

// NewRequest returns a new execution request
func NewRequest(processHash, instanceHash, parentHash, eventHash hash.Hash, stepID string, taskKey string, inputs *types.Struct, tags []string, executorHash hash.Hash) *ExecutionRequest {
	exec := &ExecutionRequest{
		ProcessHash:  processHash,
		EventHash:    eventHash,
		InstanceHash: instanceHash,
		ParentHash:   parentHash,
		Inputs:       inputs,
		TaskKey:      taskKey,
		StepID:       stepID,
		Tags:         tags,
		ExecutorHash: executorHash,
	}
	exec.Hash = hash.Dump(exec)
	return exec
}

// NewResultWithOutputs returns a new execution result
func NewResultWithOutputs(requestHash hash.Hash, outputs *types.Struct) *ExecutionResult {
	execResult := &ExecutionResult{
		RequestHash: requestHash,
		Result:      &ExecutionResult_Outputs{outputs},
	}
	execResult.Hash = hash.Dump(execResult)
	return execResult
}

// NewResultWithError returns a new execution result
func NewResultWithError(requestHash hash.Hash, err string) *ExecutionResult {
	execResult := &ExecutionResult{
		RequestHash: requestHash,
		Result:      &ExecutionResult_Error{err},
	}
	execResult.Hash = hash.Dump(execResult)
	return execResult
}

func ToExecution(request *ExecutionRequest, result *ExecutionResult) *Execution {
	status := Status_InProgress
	exec := &Execution{
		ProcessHash:  request.ProcessHash,
		EventHash:    request.EventHash,
		InstanceHash: request.InstanceHash,
		ParentHash:   request.ParentHash,
		Inputs:       request.Inputs,
		TaskKey:      request.TaskKey,
		StepID:       request.StepID,
		Tags:         request.Tags,
		ExecutorHash: request.ExecutorHash,
	}
	if result != nil {
		status = Status_Completed
		exec.Outputs = result.GetOutputs()
		exec.Error = result.GetError()
		if _, ok := result.GetResult().(*ExecutionResult_Error); ok {
			status = Status_Failed
		}
	}
	exec.Status = status
	exec.Hash = hash.Dump(exec)
	return exec
}
