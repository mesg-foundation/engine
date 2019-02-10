package core

import (
	"io"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/mesg-foundation/core/protobuf/acknowledgement"
	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/mesg-foundation/core/utils/chunker"
)

// ServiceLogs gives logs of service with the applied dependency filters.
func (s *Server) ServiceLogs(request *coreapi.ServiceLogsRequest, stream coreapi.Core_ServiceLogsServer) error {
	logs, err := s.api.ServiceLogs(request.ServiceID, request.Dependencies)
	if err != nil {
		return err
	}

	var (
		chunks = make(chan chunker.Data)
		errs   = make(chan error)
	)

	for _, log := range logs {
		stdout, dstdout := io.Pipe()
		stderr, dstderr := io.Pipe()
		go func(r io.Reader) {
			stdcopy.StdCopy(dstdout, dstderr, r)
		}(log.Reader())

		cstd := chunker.New(stdout, chunks, errs, chunker.ValueOption(&chunkMeta{
			Dependency: log.Dependency(),
			Type:       coreapi.LogData_Standard,
		}))
		cerr := chunker.New(stderr, chunks, errs, chunker.ValueOption(&chunkMeta{
			Dependency: log.Dependency(),
			Type:       coreapi.LogData_Error,
		}))
		defer log.Close()
		defer cstd.Close()
		defer cerr.Close()
	}

	// send header to notify client that the stream is ready.
	if err := acknowledgement.SetStreamReady(stream); err != nil {
		return err
	}

	ctx := stream.Context()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case err := <-errs:
			return err

		case chunk := <-chunks:
			meta := chunk.Value.(*chunkMeta)
			data := &coreapi.LogData{
				Dependency: meta.Dependency,
				Type:       meta.Type,
				Data:       chunk.Data,
			}
			if err := stream.Send(data); err != nil {
				return err
			}
		}
	}
}

// chunkMeta is a meta data for chunks.
type chunkMeta struct {
	Dependency string
	Type       coreapi.LogData_Type
}
