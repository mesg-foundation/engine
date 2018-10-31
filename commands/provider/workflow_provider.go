package provider

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/mesg-foundation/core/utils/chunker"
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
		File: fileData,
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

// WorkflowLog keeps workflow logs' standard and error streams.
type WorkflowLog struct {
	Standard, Error *chunker.Stream
}

// WorkflowLogs returns workflow log streams.
func (p *WorkflowProvider) WorkflowLogs(id string) (log *WorkflowLog, close func(), err error) {
	ctx, cancel := context.WithCancel(context.Background())
	stream, err := p.client.WorkflowLogs(ctx, &coreapi.WorkflowLogsRequest{
		ID: id,
	})
	if err != nil {
		cancel()
		return nil, nil, err
	}

	log = &WorkflowLog{
		Standard: chunker.NewStream(),
		Error:    chunker.NewStream(),
	}

	closer := func() {
		cancel()
		log.Standard.Close()
		log.Error.Close()
	}

	// TODO(ilgooz):
	// * this one actually should be handled by Close() (io.Closer) called from chunker.Stream.
	//   chunkler.Stream should accept a closeFunc handler in the initialization time to call it later.
	// * listening gRPC stream shouldn't be done by listenWorkflowLogs() as well.
	//   chunker.Stream can accept a readFunc that will be called when Read() (io.Reader) called from it.
	//   content of listenWorkflowLogs() should be implemented inside readFunc handler.
	// * to make these changes possible, chunker Streams should be initialized from a parent Stream.
	//   otherwise reading source will be duplicated which is error prone.
	//   and similar thing will happen for the closeFunc as well.
	// * when this change is made, do not forget updating other places that uses pkg chunker.
	go func() {
		p.listenWorkflowLogs(stream, log)
		closer()
	}()

	return log, closer, nil
}

// listenWorkflowLogs listens gRPC stream to get workflow logs.
func (p *WorkflowProvider) listenWorkflowLogs(stream coreapi.Core_WorkflowLogsClient, log *WorkflowLog) {
	for {
		data, err := stream.Recv()
		if err != nil {
			log.Standard.CloseWithError(err)
			return
		}

		var out *chunker.Stream
		switch data.Type {
		case coreapi.WorkflowLogData_Standard:
			out = log.Standard
		case coreapi.WorkflowLogData_Error:
			out = log.Error
		}
		out.Provide(data.Data)
	}
}
