package service

// Task describes a service task.
type Task struct {
	// Key is the key of task.
	Key string `hash:"-" yaml:"-"`

	// Name is the name of task.
	Name string `hash:"name:1" yaml:"name"`

	// Description is the description of task.
	Description string `hash:"name:2" yaml:"description"`

	// ServiceName is the service name of task.
	// TODO(ilgooz) remove this or replace with Service type in next PRs.
	ServiceName string `hash:"-" yaml:"-"`

	// Inputs are the definition of the execution inputs of task.
	Inputs map[string]*Parameter `hash:"name:3" yaml:"inputs"`

	// Outputs are the definition of the execution results of task.
	Outputs map[string]*Output `hash:"name:4" yaml:"outputs"`
}

// Output describes task output.
type Output struct {
	// Key is the key of output.
	Key string `hash:"-" yaml:"-"`

	// Name is the name of task output.
	Name string `hash:"name:1" yaml:"name"`

	// Description is the description of task output.
	Description string `hash:"name:2" yaml:"description"`

	// TaskKey is the task key of the output.
	// TODO(ilgooz) remove this or replace with Task type in next PRs.
	TaskKey string `hash:"-" yaml:"-"`

	// ServiceName is the service name of tasj output.
	// TODO(ilgooz) remove this or replace with Service type in next PRs.
	ServiceName string `hash:"-" yaml:"-"`

	// Data holds the output parameters of a task output.
	Data map[string]*Parameter `hash:"3" yaml:"data"`
}

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
