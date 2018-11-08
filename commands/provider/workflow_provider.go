package provider

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/mesg-foundation/core/protobuf/coreapi"
)

// WorkflowProvider is a struct that provides all methods required by workflow command.
type WorkflowProvider struct {
	client coreapi.CoreClient
}

// NewWorkflowProvider creates new WorkflowProvider.
func NewWorkflowProvider(c coreapi.CoreClient) *WorkflowProvider {
	return &WorkflowProvider{client: c}
}

// CreateWorkflow creates and runs a new workflow with given yaml file and
// optionally given unique name.
func (p *WorkflowProvider) CreateWorkflow(filePath string, name string) (id string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	reply, err := p.client.CreateWorkflow(context.Background(), &coreapi.CreateWorkflowRequest{
		YAML: fileData,
		Name: name,
	})
	if err != nil {
		return "", err
	}
	return reply.ID, nil
}

// DeleteWorkflow stops and deletes workflow with id.
func (p *WorkflowProvider) DeleteWorkflow(id string) error {
	_, err := p.client.DeleteWorkflow(context.Background(), &coreapi.DeleteWorkflowRequest{ID: id})
	return err
}
