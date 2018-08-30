package service

// Parameter describes task input parameters, output parameters of a task
// output and input parameters of an event.
type Parameter struct {
	// Name is the name of parameter.
	Name string `hash:"name:1" yaml:"name"`

	// Description is the description of parameter.
	Description string `hash:"name:2" yaml:"description"`

	// Type is the data type of parameter.
	// Type DataType `hash:"3"`
	Type string `hash:"name:3" yaml:"type"`

	// Optional indicates if parameter is optional.
	Optional bool `hash:"name:4" yaml:"optiona"`
}
