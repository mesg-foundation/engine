package service

// Task describes a service task.
type Task struct {
	// Key is the key of task.
	Key string `hash:"name:1"`

	// Name is the name of task.
	Name string `hash:"name:2"`

	// Description is the description of task.
	Description string `hash:"name:3"`

	// Inputs are the definition of the execution inputs of task.
	Inputs []*Parameter `hash:"name:4"`

	// Outputs are the definition of the execution results of task.
	Outputs []*Parameter `hash:"name:5"`

	// serviceName is the task's service's name.
	serviceName string `hash:"-"`
}

// GetTask returns task taskKey of service.
func (s *Service) GetTask(taskKey string) (*Task, error) {
	for _, task := range s.Tasks {
		if task.Key == taskKey {
			task.serviceName = s.Name
			return task, nil
		}
	}
	return nil, &TaskNotFoundError{
		TaskKey:     taskKey,
		ServiceName: s.Name,
	}
}

// GetInputParameter returns input inputKey parameter of task.
func (t *Task) GetInputParameter(inputKey string) (*Parameter, error) {
	for _, input := range t.Inputs {
		if input.Key == inputKey {
			return input, nil
		}
	}
	return nil, &TaskInputNotFoundError{
		TaskKey:      t.Key,
		TaskInputKey: inputKey,
		ServiceName:  t.serviceName,
	}
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
			ServiceName: t.serviceName,
			Warnings:    warnings,
		}
	}
	return nil
}

// ValidateOutputs produces warnings for task outputs that doesn't satisfy their parameter schemas.
func (t *Task) ValidateOutputs(data map[string]interface{}) []*ParameterWarning {
	return validateParametersSchema(t.Outputs, data)
}

// RequireOutputData requires task outputs to be matched with parameter schemas.
func (t *Task) RequireOutputData(data map[string]interface{}) error {
	warnings := t.ValidateOutputs(data)
	if len(warnings) > 0 {
		return &InvalidTaskOutputError{
			TaskKey:     t.Key,
			ServiceName: t.serviceName,
			Warnings:    warnings,
		}
	}
	return nil
}
