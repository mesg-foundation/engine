package workflow

type WorkflowDefinition struct {
	Name        string
	Description string
	Services    []ServiceDefinition
	Configs     []ConfigDefinition
	Events      []EventDefinition
}

type ConfigDefinition struct {
	Key   string
	Value interface{}
}

type ServiceDefinition struct {
	Name string
	ID   string
}

type EventDefinition struct {
	ServiceName string
	EventKey    string
	Map         []MapDefinition
	Execute     ExecuteDefinition
}

type MapDefinition struct {
	Key   string
	Value interface{}
}

type ExecuteDefinition struct {
	ServiceName string
	TaskKey     string
}
