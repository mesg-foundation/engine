// Package mesg is a service client for mesg-core.
// For more information please visit https://mesg.com.
package mesg

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"encoding/json"

	"context"

	"github.com/mesg-foundation/core/api/service"
	"google.golang.org/grpc"
)

const (
	tcpEndpointEnv = "MESG_ENDPOINT"
	tokenEnv       = "MESG_TOKEN"
)

var (
	errEndpointNotSet = errors.New("endpoint is not set")
	errTokenNotSet    = errors.New("token is not set")
)

// Service represents a MESG service.
type Service struct {
	// endpoint is the mesg-core endpoint.
	endpoint string

	// token is the service id.
	token string

	// client is the gRPC service client of MESG.
	client service.ServiceClient

	// conn is underlying gRPC conn
	conn *grpc.ClientConn

	// dialOptions passed to grpc dial as options.
	dialOptions []grpc.DialOption

	// callTimeout used to timeout gRPC requests or dial.
	callTimeout time.Duration

	// tasks holds task handlers.
	tasks []Taskable
	mt    sync.RWMutex

	// log is a logger for service.
	log *log.Logger

	// logOutput is the output stream of log.
	logOutput io.Writer
}

// Option is the configuration func of Service.
type Option func(*Service)

// New starts a new Service with options.
func New(options ...Option) (*Service, error) {
	s := &Service{
		endpoint:    os.Getenv(tcpEndpointEnv),
		token:       os.Getenv(tokenEnv),
		callTimeout: time.Second * 10,
		logOutput:   ioutil.Discard,
		dialOptions: []grpc.DialOption{grpc.WithInsecure()},
	}
	for _, option := range options {
		option(s)
	}
	s.log = log.New(s.logOutput, "mesg", log.LstdFlags)
	if s.endpoint == "" {
		return nil, errEndpointNotSet
	}
	if s.token == "" {
		return nil, errTokenNotSet
	}
	return s, s.setupServiceClient()
}

// EndpointOption receives the TCP endpoint of mesg-core.
func EndpointOption(address string) Option {
	return func(s *Service) {
		s.endpoint = address
	}
}

// TokenOption receives token which is the unique id of this service.
func TokenOption(token string) Option {
	return func(s *Service) {
		s.token = token
	}
}

// TimeoutOption receives d to use while dialing mesg-core and making requests.
func TimeoutOption(d time.Duration) Option {
	return func(s *Service) {
		s.callTimeout = d
	}
}

// LogOutputOption uses out as a log destination.
func LogOutputOption(out io.Writer) Option {
	return func(s *Service) {
		s.logOutput = out
	}
}

// DialOption used to mock socket communication for unit testing.
func DialOption(dialer Dialer) Option {
	return func(s *Service) {
		s.dialOptions = append(s.dialOptions, grpc.WithDialer(newGRPCDialer(dialer).Dial))
	}
}

func (s *Service) setupServiceClient() error {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), s.callTimeout)
	defer cancel()
	s.conn, err = grpc.DialContext(ctx, s.endpoint, s.dialOptions...)
	if err != nil {
		return err
	}
	s.client = service.NewServiceClient(s.conn)
	return nil
}

// Listen listens requests for given tasks. It's a blocking call.
func (s *Service) Listen(task Taskable, tasks ...Taskable) error {
	s.mt.Lock()
	if len(s.tasks) > 0 {
		s.mt.Unlock()
		return errors.New("tasks already set")
	}
	s.tasks = append(tasks, task)
	s.mt.Unlock()
	if err := s.validateTasks(); err != nil {
		return err
	}
	return s.listenTasks()
}

// validateTasks checks if the tasks handled exectly desribed in mesg.yaml.
// TODO(ilgooz) use validation handlers of core server to do this?
func (s *Service) validateTasks() error { return nil }

func (s *Service) listenTasks() error {
	stream, err := s.client.ListenTask(context.Background(), &service.ListenTaskRequest{
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
		go s.executeTask(data)
	}
}

func (s *Service) executeTask(data *service.TaskData) {
	s.mt.RLock()
	for _, task := range s.tasks {
		if task.Name() == data.TaskKey {
			s.mt.RUnlock()

			execution := newExecution(s, data)
			if err := execution.reply(task.Execute(execution)); err != nil {
				s.log.Println(err)
			}

			return
		}
	}
	s.mt.RUnlock()
}

// Emit emits a MESG event with given data for name.
func (s *Service) Emit(event string, data Data) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.callTimeout)
	defer cancel()
	_, err = s.client.EmitEvent(ctx, &service.EmitEventRequest{
		Token:     s.token,
		EventKey:  event,
		EventData: string(dataBytes),
	})
	return err
}

// Close gracefully closes underlying connections and stops listening for task requests.
func (s *Service) Close() error {
	return s.conn.Close()
}
