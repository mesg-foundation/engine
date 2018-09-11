package core

import (
	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/utils/chunker"
)

// ServiceLogs gives logs of service with the applied dependency filters.
func (s *Server) ServiceLogs(request *ServiceLogsRequest, stream Core_ServiceLogsServer) error {
	sl, err := s.api.ServiceLogs(request.ServiceID,
		api.ServiceLogsDependenciesFilter(request.Dependencies...))
	if err != nil {
		return err
	}

	var (
		chunks = make(chan chunker.Data, 0)
		errs   = make(chan error, 0)
	)

	for _, l := range sl {
		cstd := chunker.New(l.Standard, chunks, errs, chunker.ValueOption(&chunkMeta{
			Dependency: l.Dependency,
			Type:       LogData_Standard,
		}))
		cerr := chunker.New(l.Error, chunks, errs, chunker.ValueOption(&chunkMeta{
			Dependency: l.Dependency,
			Type:       LogData_Error,
		}))
		defer cstd.Close()
		defer cerr.Close()
		defer l.Standard.Close()
		defer l.Error.Close()
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
			data := &LogData{
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
	Type       LogData_Type
}
