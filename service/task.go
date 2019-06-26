package service

// Task describes a service task.
type Task struct {
	// Key is the key of task.
	Key string `hash:"name:1" validate:"printascii"`

	// Name is the name of task.
	Name string `hash:"name:2" validate:"printascii"`

	// Description is the description of task.
	Description string `hash:"name:3" validate:"printascii"`

	// Inputs are the definition of the execution inputs of task.
	Inputs []*Parameter `hash:"name:4" validate:"dive,required"`

	// Outputs are the definition of the execution results of task.
	Outputs []*Parameter `hash:"name:5" validate:"dive,required"`
}

// GetTask returns task taskKey of service.
func (s *Service) GetTask(taskKey string) (*Task, error) {
	for _, task := range s.Tasks {
		if task.Key == taskKey {
			return task, nil
		}
	}
	return nil, &TaskNotFoundError{
		TaskKey:     taskKey,
		ServiceName: s.Name,
	}
}

// ValidateTaskInputs produces warnings for task inputs that doesn't satisfy their parameter schemas.
func (s *Service) ValidateTaskInputs(taskKey string, taskInputs map[string]interface{}) ([]*ParameterWarning, error) {
	t, err := s.GetTask(taskKey)
	if err != nil {
		return nil, err
	}
	return validateParametersSchema(t.Inputs, taskInputs), nil
}

// ValidateTaskOutputs produces warnings for task outputs that doesn't satisfy their parameter schemas.
func (s *Service) ValidateTaskOutputs(taskKey string, taskOutputs map[string]interface{}) ([]*ParameterWarning, error) {
	t, err := s.GetTask(taskKey)
	if err != nil {
		return nil, err
	}
	return validateParametersSchema(t.Outputs, taskOutputs), nil
}

// RequireTaskInputs requires task inputs to match with parameter schemas.
func (s *Service) RequireTaskInputs(taskKey string, taskInputs map[string]interface{}) error {
	warnings, err := s.ValidateTaskInputs(taskKey, taskInputs)
	if err != nil {
		return err
	}
	if len(warnings) > 0 {
		return &InvalidTaskInputError{
			TaskKey:     taskKey,
			ServiceName: s.Name,
			Warnings:    warnings,
		}
	}
	return nil
}

// RequireTaskOutputs requires task outputs to match with parameter schemas.
func (s *Service) RequireTaskOutputs(taskKey string, taskOutputs map[string]interface{}) error {
	warnings, err := s.ValidateTaskOutputs(taskKey, taskOutputs)
	if err != nil {
		return err
	}
	if len(warnings) > 0 {
		return &InvalidTaskOutputError{
			TaskKey:     taskKey,
			ServiceName: s.Name,
			Warnings:    warnings,
		}
	}
	return nil
}
