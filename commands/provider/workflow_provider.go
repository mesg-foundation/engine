package provider

import (
	"context"

	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/mesg-foundation/core/utils/workflowparser/yaml"
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
	definition, err := yaml.ParseFromFile(filePath)
	if err != nil {
		return "", err
	}
	reply, err := p.client.CreateWorkflow(context.Background(), &coreapi.CreateWorkflowRequest{
		Definition: p.toProtoWorkflowDefinition(definition),
		Name:       name,
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

func (p *WorkflowProvider) toProtoWorkflowDefinition(definition yaml.WorkflowDefinition) *coreapi.CreateWorkflowRequest_WorkflowDefinition {
	return &coreapi.CreateWorkflowRequest_WorkflowDefinition{}
}
