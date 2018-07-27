// Package applicationtest is a testing package for MESG application.
// Use this package while unit testing your programs.
package applicationtest

import (
	"encoding/json"

	"github.com/mesg-foundation/core/api/core"
)

// Server is a test server.
type Server struct {
	core   *coreServer
	socket *Socket
}

// NewServer a new test server.
func NewServer() *Server {
	return &Server{
		core:   newCoreServer(),
		socket: newSocket(),
	}
}

// Start starts the test server.
func (s *Server) Start() error {
	return s.socket.listen(s.core)
}

// Socket returns a in-memory socket for client application.
func (s *Server) Socket() *Socket {
	return s.socket
}

// LastServiceStart returns the last service start request's info.
func (s *Server) LastServiceStart() *ServiceStart {
	req := <-s.core.serviceStartC
	return &ServiceStart{
		serviceID: req.ServiceID,
	}
}

// LastEventListen returns the last event listen request's info.
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

// EmitEvent emits a new event for serviceID with given data.
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

// LastEventListen returns the last result listen request's info.
func (s *Server) LastResultListen() *ResultListen {
	for {
		select {
		case <-s.core.serviceStartC:
		case req := <-s.core.listenResultC:
			return &ResultListen{
				serviceID: req.ServiceID,
				key:       req.OutputFilter,
				task:      req.TaskFilter,
			}
		}
	}
}

// EmitResult emits a new task result for serviceID with given outputKey and data.
func (s *Server) EmitResult(serviceID, task, outputKey string, data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	rd := &core.ResultData{
		//eexecutionID
		TaskKey:    task,
		OutputKey:  outputKey,
		OutputData: string(bytes),
	}
	s.core.em.Lock()
	if s.core.resultC[serviceID] == nil {
		s.core.resultC[serviceID] = make(chan *core.ResultData, 0)
	}
	s.core.em.Unlock()
	for {
		select {
		case <-s.core.serviceStartC:
		case <-s.core.listenResultC:
		case s.core.resultC[serviceID] <- rd:
			return nil
		}
	}
}

// LastExecute returns the last task execution's info.
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

// MarkServiceAsNonExistent marks a service id as non-exists server.
func (s *Server) MarkServiceAsNonExistent(id string) {
	s.core.nonExistentServices = append(s.core.nonExistentServices, id)
}

// Close closes test server.
func (s *Server) Close() error {
	return s.socket.close()
}
