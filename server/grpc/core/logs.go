package core

import (
	"github.com/mesg-foundation/core/protobuf/acknowledgement"
	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/mesg-foundation/core/sdk"
	"github.com/mesg-foundation/core/utils/chunker"
	"github.com/mr-tron/base58"
)

// ServiceLogs gives logs of service with the applied dependency filters.
func (s *Server) ServiceLogs(request *coreapi.ServiceLogsRequest, stream coreapi.Core_ServiceLogsServer) error {
	hash, err := base58.Decode(request.ServiceID)
	if err != nil {
		return err
	}

	sl, err := s.sdk.ServiceLogs(hash,
		sdk.ServiceLogsDependenciesFilter(request.Dependencies...))
	if err != nil {
		return err
	}

	var (
		chunks = make(chan chunker.Data)
		errs   = make(chan error)
	)

	for _, l := range sl {
		cstd := chunker.New(l.Standard, chunks, errs, chunker.ValueOption(&chunkMeta{
			Dependency: l.Dependency,
			Type:       coreapi.LogData_Standard,
		}))
		cerr := chunker.New(l.Error, chunks, errs, chunker.ValueOption(&chunkMeta{
			Dependency: l.Dependency,
			Type:       coreapi.LogData_Error,
		}))
		defer cstd.Close()
		defer cerr.Close()
		defer l.Standard.Close()
		defer l.Error.Close()
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
