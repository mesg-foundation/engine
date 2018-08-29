package service

// GetTask returns task taskKey of service.
func (s *Service) GetTask(taskKey string) (*Task, error) {
	task, ok := s.Tasks[taskKey]
	if !ok {
		return nil, &TaskNotFoundError{
			TaskKey:     taskKey,
			ServiceName: s.Name,
		}
	}
	return task, nil
}

// ValidateInputs produces warnings for task inputs that doesn't satisfy their parameter schemas.
func (t *Task) ValidateInputs(taskInputs map[string]interface{}) []*ParameterWarning {
	return validateParametersSchema(t.Inputs, taskInputs)
}

// RequireInputs requires task inputs to be matched with parameter schemas.
func (t *Task) RequireInputs(taskInputs map[string]interface{}) error {
	warnings := t.ValidateInputs(taskInputs)
	if len(warnings) > 0 {
		return &InvalidTaskInputError{
			TaskKey:     t.Key,
			ServiceName: t.ServiceName,
			Warnings:    warnings,
		}
	}
	return nil
}

// GetOutput returns output outputKey of task.
func (t *Task) GetOutput(outputKey string) (*Output, error) {
	output, ok := t.Outputs[outputKey]
	if !ok {
		return nil, &TaskOutputNotFoundError{
			TaskKey:       t.Key,
			TaskOutputKey: outputKey,
			ServiceName:   t.ServiceName,
		}
	}
	return output, nil
}

// ValidateData produces warnings for task outputs that doesn't satisfy their parameter schemas.
func (o *Output) ValidateData(outputData map[string]interface{}) []*ParameterWarning {
	return validateParametersSchema(o.Data, outputData)
}

// RequireData requires task outputs to be matched with parameter schemas.
func (o *Output) RequireData(outputData map[string]interface{}) error {
	warnings := o.ValidateData(outputData)
	if len(warnings) > 0 {
		return &InvalidTaskOutputError{
			TaskKey:       o.TaskKey,
			TaskOutputKey: o.Key,
			ServiceName:   o.ServiceName,
			Warnings:      warnings,
		}
	}
	return nil
}
