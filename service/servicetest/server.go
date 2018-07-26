package servicetest

import (
	"encoding/json"
	"errors"

	"github.com/mesg-foundation/core/api/service"
	uuid "github.com/satori/go.uuid"
)

// Service is a test service.
type Server struct {
	service *serviceServer
	socket  *Socket
}

// New creates a new test service.
func NewServer() *Server {
	return &Server{
		service: newServiceServer(),
		socket:  newSocket(),
	}
}

func (s *Server) Start() error {
	return s.socket.listen(s.service)
}

func (s *Server) Socket() *Socket {
	return s.socket
}

// LastEmit returns last emitted event.
func (s *Server) LastEmit() *Event {
	e := <-s.service.emitC
	return &Event{
		name:  e.EventKey,
		data:  e.EventData,
		token: e.Token,
	}
}

// Execute executes a task with data.
func (s *Server) Execute(task string, data interface{}) (string, *Execution, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", nil, err
	}

	uuidV4, err := uuid.NewV4()
	if err != nil {
		return "", nil, err
	}
	id := uuidV4.String()
	taskData := &service.TaskData{
		ExecutionID: id,
		TaskKey:     task,
		InputData:   string(bytes),
	}
	select {
	case <-s.service.closingC:
		return "", nil, errors.New("connection closed")
	case s.service.taskC <- taskData:
	}

	select {
	case <-s.service.closingC:
		return "", nil, errors.New("connection closed")

	case resp := <-s.service.submitC:
		return id, &Execution{
			id:   resp.ExecutionID,
			key:  resp.OutputKey,
			data: resp.OutputData,
		}, nil
	}
}

// ListenToken returns the token of service that started listening for tasks.
func (s *Server) ListenToken() string {
	return s.service.token
}

// Close ends waiting for task requests.
func (s *Server) Close() error {
	return s.socket.close()
}
