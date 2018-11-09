package yaml

import (
	"io/ioutil"
	"os"
)

// YAMLDefinition TODO.
type YAMLDefinition struct{}

// Parse parses workflow yaml into definition object.
func Parse(yaml []byte) (WorkflowDefinition, error) {
	return WorkflowDefinition{}, nil
}

// ParseFromFile parses workflow.yaml file info definition object.
func ParseFromFile(path string) (WorkflowDefinition, error) {
	file, err := os.Open(path)
	if err != nil {
		return WorkflowDefinition{}, err
	}
	defer file.Close()

	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		return WorkflowDefinition{}, err
	}

	return Parse(fileData)
}
