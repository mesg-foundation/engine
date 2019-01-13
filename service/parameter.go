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
	Type string `hash:"name:4"`

	// Object  keeps definition of object type.
	Object map[string]*Parameter `hash:"name:7"`

	// Optional indicates if parameter is optional.
	Optional bool `hash:"name:5"`

	// Repeated is to have an array of this parameter
	Repeated bool `hash:"name:6"`
}
