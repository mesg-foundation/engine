package applicationtest

import (
	"encoding/json"

	"github.com/mesg-foundation/core/api/core"
)

// Service is a test service.
type Server struct {
	core   *coreServer
	socket *Socket
}

// New creates a new test service.
func NewServer() *Server {
	return &Server{
		core:   newCoreServer(),
		socket: newSocket(),
	}
}

func (s *Server) Start() error {
	return s.socket.listen(s.core)
}

func (s *Server) Socket() *Socket {
	return s.socket
}

func (s *Server) LastEventListen() *EventListen {
	for {
		select {
		case <-s.core.serviceStartC:
		case req := <-s.core.listenEventC:
			return &EventListen{
				serviceID: req.ServiceID,
				event:     req.EventFilter,
			}
		}
	}
}

func (s *Server) LastServiceStart() *ServiceStart {
	req := <-s.core.serviceStartC
	return &ServiceStart{
		serviceID: req.ServiceID,
	}
}

func (s *Server) EmitEvent(serviceID, event string, data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	ed := &core.EventData{
		EventKey:  event,
		EventData: string(bytes),
	}
	s.core.em.Lock()
	if s.core.eventC[serviceID] == nil {
		s.core.eventC[serviceID] = make(chan *core.EventData, 0)
	}
	s.core.em.Unlock()
	for {
		select {
		case <-s.core.serviceStartC:
		case <-s.core.listenEventC:
		case s.core.eventC[serviceID] <- ed:
			return nil
		}
	}
}

func (s *Server) EmitResult() {

}

func (s *Server) LastExecute() *Execute {
	for {
		select {
		case <-s.core.serviceStartC:
		case <-s.core.listenEventC:
		case req := <-s.core.executeC:
			return &Execute{
				serviceID: req.ServiceID,
				task:      req.TaskKey,
				data:      req.InputData,
			}
		}
	}
}

func (s *Server) MarkServiceAsNonExistent(id string) {
	s.core.nonExistentServices = append(s.core.nonExistentServices, id)
}

// Close ends waiting for task requests.
func (s *Server) Close() error {
	return s.socket.close()
}
