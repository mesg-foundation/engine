// Package mesg is a service client for mesg-core.
// For more information please visit https://mesg.com.
package mesg

import (
	"errors"
	"fmt"
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

	// dialOptions holds dial options of gRPC.
	dialOptions []grpc.DialOption

	// callTimeout used to timeout gRPC requests or dial.
	callTimeout time.Duration

	// cancel stops receiving from gRPC task stream.
	cancel context.CancelFunc
	mc     sync.Mutex // protects cancel and conn.Close().

	// isListening set true after the first call to Listen().
	isListening bool
	ml          sync.Mutex

	// gracefulWait will be in the done state when all processing
	// task executions are done.
	gracefulWait *sync.WaitGroup

	// taskables holds task handlers.
	taskables []Taskable

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
		endpoint:     os.Getenv(tcpEndpointEnv),
		token:        os.Getenv(tokenEnv),
		callTimeout:  time.Second * 10,
		gracefulWait: &sync.WaitGroup{},
		logOutput:    ioutil.Discard,
		dialOptions:  []grpc.DialOption{grpc.WithInsecure()},
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
	s.ml.Lock()
	if s.isListening {
		s.ml.Unlock()
		return errAlreadyListening{}
	}
	s.isListening = true
	s.ml.Unlock()
	s.taskables = append(tasks, task)
	if err := s.validateTasks(); err != nil {
		return err
	}
	return s.listenTasks()
}

// validateTasks checks if the tasks handled exectly desribed in mesg.yaml.
// TODO(ilgooz) use validation handlers of core server to do this?
func (s *Service) validateTasks() error { return nil }

func (s *Service) listenTasks() error {
	var ctx context.Context
	s.mc.Lock()
	ctx, s.cancel = context.WithCancel(context.Background())
	s.mc.Unlock()
	stream, err := s.client.ListenTask(ctx, &service.ListenTaskRequest{
		Token: s.token,
	})
	if err != nil {
		return err
	}
	for {
		s.gracefulWait.Add(1)
		data, err := stream.Recv()
		if err != nil {
			s.gracefulWait.Done()
			return err
		}
		go s.executeTask(data)
	}
}

func (s *Service) getTaskableByName(key string) Taskable {
	for _, taskable := range s.taskables {
		if taskable.Key() == key {
			return taskable
		}
	}
	return nil
}

func (s *Service) executeTask(data *service.TaskData) {
	defer s.gracefulWait.Done()
	taskable := s.getTaskableByName(data.TaskKey)
	if taskable == nil {
		s.log.Println(errNonExistentTask{data.TaskKey})
		return
	}
	execution := newExecution(s, data)
	if err := execution.reply(taskable.Execute(execution)); err != nil {
		s.log.Println(err)
	}
}

// Emit emits a MESG event eventKey with given eventData.
func (s *Service) Emit(eventKey string, eventData Data) error {
	dataBytes, err := json.Marshal(eventData)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.callTimeout)
	defer cancel()
	_, err = s.client.EmitEvent(ctx, &service.EmitEventRequest{
		Token:     s.token,
		EventKey:  eventKey,
		EventData: string(dataBytes),
	})
	return err
}

// Close gracefully stops listening for future task execution requests and waits
// current ones to complete before closing underlying connection.
func (s *Service) Close() error {
	s.mc.Lock()
	defer s.mc.Unlock()
	s.cancel()
	s.gracefulWait.Wait()
	return s.conn.Close()
}

type errNonExistentTask struct {
	name string
}

func (e errNonExistentTask) Error() string {
	return fmt.Sprintf("task %q does not exists", e.name)
}

type errAlreadyListening struct{}

func (e errAlreadyListening) Error() string {
	return "already listening for tasks"
}
