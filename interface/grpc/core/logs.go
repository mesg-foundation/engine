package core

import (
	"io"

	"github.com/mesg-foundation/core/api"
)

// ServiceLogs gives logs of service with the applied dependency filters.
func (s *Server) ServiceLogs(request *ServiceLogsRequest, stream Core_ServiceLogsServer) error {
	sl, err := s.api.ServiceLogs(request.ServiceID,
		api.ServiceLogsDependenciesFilter(request.Dependencies...))
	if err != nil {
		return err
	}

	// send dependency list as a header data.
	var dependencies []string
	for _, l := range sl {
		dependencies = append(dependencies, l.Dependency)
	}
	if err := stream.Send(&LogData{
		Depedencies: dependencies,
	}); err != nil {
		return err
	}

	results := make(chan logChunk)

	for _, l := range sl {
		rstd := newLogPiper(LogData_Data_Standard, l.Dependency, l.Standard, results)
		rerr := newLogPiper(LogData_Data_Error, l.Dependency, l.Error, results)
		defer rstd.Close()
		defer rerr.Close()
	}

	ctx := stream.Context()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case result := <-results:
			if result.Err != nil {
				return err
			}

			data := &LogData_Data{
				Dependency: result.Dependency,
				Type:       result.Type,
				Data:       result.Data,
			}
			if err := stream.Send(&LogData{Data: data}); err != nil {
				return err
			}
		}
	}
}

// logChunk keeps the information about log data chunk and its owner.
type logChunk struct {
	Dependency string
	Type       LogData_Data_Type
	Data       []byte
	Err        error
}

// logPiper reads logs from given reader for corresponding dependency and log type.
type logPiper struct {
	Type       LogData_Data_Type
	Dependency string
	Stream     io.ReadCloser
	Chunks     chan logChunk
	closing    chan struct{}
}

// newLogPiper creates a new log piper which starts reading from reader and pipes data
// chunks to chunks chan.
func newLogPiper(typ LogData_Data_Type, dependency string, stream io.ReadCloser,
	results chan logChunk) *logPiper {
	r := &logPiper{
		Type:       typ,
		Dependency: dependency,
		Stream:     stream,
		Chunks:     results,
		closing:    make(chan struct{}, 0),
	}
	go r.run()
	return r
}

// run reads log data from reader and sends it to gRPC's stream send queue.
func (r *logPiper) run() {
	buf := make([]byte, 1024)
	for {
		n, err := r.Stream.Read(buf)
		if err != nil {
			select {
			case <-r.closing:
				return

			case r.Chunks <- logChunk{Err: err}:
			}
		}

		select {
		case <-r.closing:
			return

		case r.Chunks <- logChunk{
			Dependency: r.Dependency,
			Type:       r.Type,
			Data:       buf[:n],
		}:
		}
	}
}

func (r *logPiper) Close() error {
	r.Stream.Close()
	close(r.closing)
	return nil
}
