package core

import (
	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/mesg-foundation/core/utils/chunker"
)

// WorkflowLogs gives the logs stream of a workflow.
func (s *Server) WorkflowLogs(request *coreapi.WorkflowLogsRequest, stream coreapi.Core_WorkflowLogsServer) error {
	stdLogs, errLogs, err := s.ss.Workflow().Logs(request.ID)
	if err != nil {
		return err
	}

	var (
		chunks = make(chan chunker.Data)
		errs   = make(chan error)
	)

	cstd := chunker.New(stdLogs, chunks, errs, chunker.ValueOption(&workflowChunkMeta{
		Type: coreapi.WorkflowLogData_Standard,
	}))
	cerr := chunker.New(errLogs, chunks, errs, chunker.ValueOption(&workflowChunkMeta{
		Type: coreapi.WorkflowLogData_Error,
	}))
	defer cstd.Close()
	defer cerr.Close()
	defer stdLogs.Close()
	defer errLogs.Close()

	ctx := stream.Context()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case err := <-errs:
			return err

		case chunk := <-chunks:
			meta := chunk.Value.(*workflowChunkMeta)
			data := &coreapi.WorkflowLogData{
				Type: meta.Type,
				Data: chunk.Data,
			}
			if err := stream.Send(data); err != nil {
				return err
			}
		}
	}
}

// workflowChunkMeta keeps the meta data for workflow log chunks.
type workflowChunkMeta struct {
	Type coreapi.WorkflowLogData_Type
}
