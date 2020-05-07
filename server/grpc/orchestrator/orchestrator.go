package orchestrator

import (
	"github.com/mesg-foundation/engine/orchestrator"
	"github.com/mesg-foundation/engine/protobuf/acknowledgement"
)

type orchestratorServer struct {
	orch *orchestrator.Orchestrator
	auth *Authorizer
}

// NewOrchestratorServer creates a new Orchestrator Server.
func NewOrchestratorServer(orch *orchestrator.Orchestrator, auth *Authorizer) OrchestratorServer {
	return &orchestratorServer{
		orch: orch,
		auth: auth,
	}
}

// Stream returns stream of events.
func (s *orchestratorServer) Logs(req *OrchestratorLogsRequest, stream Orchestrator_LogsServer) error {
	// check authorization
	if err := s.auth.IsAuthorized(stream.Context(), req); err != nil {
		return err
	}

	// create listener
	topics := make([]string, 0)
	for _, h := range req.ProcessHashes {
		topics = append(topics, h.String())
	}
	if len(topics) == 0 {
		topics = append(topics, orchestrator.AllLogTopic)
	}
	logger := s.orch.NewLogger(topics...)
	go logger.Listen()
	defer logger.Close()

	// send header to notify client that the stream is ready.
	if err := acknowledgement.SetStreamReady(stream); err != nil {
		return err
	}

	for {
		select {
		case log := <-logger.Logs:
			if err := stream.Send(&OrchestratorLogsResponse{
				ProcessHash:     log.ProcessHash,
				NodeKey:         log.NodeKey,
				NodeType:        log.NodeType,
				EventHash:       log.EventHash,
				ParentHash:      log.ParentHash,
				Msg:             log.Msg,
				Error:           log.Error,
				CreatedExecHash: log.CreatedExecHash,
			}); err != nil {
				return err
			}
		case <-stream.Context().Done():
			return stream.Context().Err()
		}
	}
}
