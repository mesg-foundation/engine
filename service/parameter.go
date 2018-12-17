package service

// Parameter describes task input parameters, output parameters of a task
// output and input parameters of an event.
type Parameter struct {
	// Key is the key of parameter.
	Key string `hash:"name:1"`

	// Name is the name of parameter.
	Name string `hash:"name:2"`

	// Description is the description of parameter.
	Description string `hash:"name:3"`

	// Type is the data type of parameter.
	// It indicates basic types(String, Number, Boolean, Any).
	// It's filled when there are no child parameters.
	Type string `hash:"name:4"`

	// Parameters are the child parameters of parameter to keep nested data.
	Parameters []*Parameter `hash:"name:6"`

	// Repeated shows that if type or child parameters are repeated.
	Repeated bool `hash:"name:7"`

	// Optional indicates if parameter is optional.
	Optional bool `hash:"name:5"`
}
