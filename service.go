// Package mesg is a service and application client for mesg-core.
// For more information please visit https://mesg.com.
package mesg

import (
	"errors"
	"os"
	"sync"
	"time"

	"encoding/json"

	"context"

	"github.com/mesg-foundation/core/api/service"
	"google.golang.org/grpc"
)

const tcpEndpointEnv = "MESG_ENDPOINT_TCP"
const tokenEnv = "MESG_TOKEN"

var defaultService *Service
var once sync.Once

// Service is a MESG service.
type Service struct {
	endpoint string
	token    string

	serviceClient service.ServiceClient
	conn          *grpc.ClientConn

	// testing is used to know if we're in the `go test` mode to disable grpc conn operations.
	// this needed because we can't mock grpc.ClientConn since it isn't an interface.
	// we're also not able to create our own interface for grpc.ClientConn because it has
	// methods that accepts paramaters from it's `internal` package. We're not able to access
	// packages named as internal from outside of their root directory.
	testing bool

	callTimeout time.Duration

	tasks map[string]Task
	mt    sync.RWMutex
}

// ServiceOption is the configuration function for Service.
type ServiceOption func(*Service)

// NewService starts a new Service with options.
func NewService(options ...ServiceOption) (*Service, error) {
	s := &Service{
		endpoint:    os.Getenv(tcpEndpointEnv),
		token:       os.Getenv(tokenEnv),
		tasks:       map[string]Task{},
		callTimeout: time.Second * 10,
	}
	for _, option := range options {
		option(s)
	}
	if s.endpoint == "" {
		return nil, errors.New("endpoint is not set")
	}
	if s.token == "" {
		return nil, errors.New("token is not set")
	}
	if !s.testing {
		s.setupServiceClient()
	}
	return s, nil
}

// GetService starts the default service. Calling this more than one ise safe.
func GetService() (*Service, error) {
	var err error
	once.Do(func() {
		defaultService, err = NewService()
	})
	return defaultService, err
}

// ServiceEndpointOption receives the TCP endpoint of MESG.
func ServiceEndpointOption(address string) ServiceOption {
	return func(s *Service) {
		s.endpoint = address
	}
}

// ServiceTokenOption receives the service id.
func ServiceTokenOption(token string) ServiceOption {
	return func(s *Service) {
		s.token = token
	}
}

// ServiceTimeoutOption receives d to use while dialing mesg-core and making requests.
func ServiceTimeoutOption(d time.Duration) ServiceOption {
	return func(s *Service) {
		s.callTimeout = d
	}
}

func serviceTestingOption(testing bool) ServiceOption {
	return func(s *Service) {
		s.testing = testing
	}
}

func (s *Service) setupServiceClient() error {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), s.callTimeout)
	defer cancel()
	s.conn, err = grpc.DialContext(ctx, s.endpoint, grpc.WithInsecure())
	if err != nil {
		return err
	}
	s.serviceClient = service.NewServiceClient(s.conn)
	return nil
}

// ListenTasks blocks while listening for tasks.
func (s *Service) ListenTasks(task Task, tasks ...Task) error {
	s.mt.Lock()
	if len(s.tasks) > 0 {
		s.mt.Unlock()
		return errors.New("tasks already set")
	}
	s.tasks[task.name] = task
	for _, task := range tasks {
		s.tasks[task.name] = task
	}
	s.mt.Unlock()
	if err := s.validateTasks(); err != nil {
		return err
	}
	return s.listenTasks()
}

func (s *Service) validateTasks() error { return nil }

func (s *Service) listenTasks() error {
	stream, err := s.serviceClient.ListenTask(context.Background(), &service.ListenTaskRequest{
		Token: s.token,
	})
	if err != nil {
		return err
	}
	for {
		data, err := stream.Recv()
		if err != nil {
			return err
		}
		s.executeTask(data)
	}
}

func (s *Service) executeTask(data *service.TaskData) {
	s.mt.RLock()
	fn, ok := s.tasks[data.TaskKey]
	s.mt.RUnlock()
	if !ok {
		return
	}
	req := &Request{
		executionID: data.ExecutionID,
		key:         data.TaskKey,
		data:        data.InputData,
		service:     s,
	}
	go fn.handler(req)
}

// EmitEvent emits a MESG event with given data for name.
func (s *Service) EmitEvent(name string, data interface{}) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.callTimeout)
	defer cancel()
	_, err = s.serviceClient.EmitEvent(ctx, &service.EmitEventRequest{
		Token:     s.token,
		EventKey:  name,
		EventData: string(dataBytes),
	})
	return err
}

// Close gracefully closes underlying conections and frees blocking calls.
func (s *Service) Close() error {
	if !s.testing {
		return s.conn.Close()
	}
	return nil
}
