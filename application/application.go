// Package application is an application client for mesg-core.
// For more information please visit https://mesg.com.
package application

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

	"github.com/mesg-foundation/core/api/core"
	"google.golang.org/grpc"
)

const (
	endpointEnv     = "MESG_ENDPOINT_TCP"
	defaultEndpoint = "localhost:50052"
)

// Application represents is a MESG application.
type Application struct {
	endpoint string

	// Client is the gRPC core client of MESG.
	client core.CoreClient
	conn   *grpc.ClientConn

	callTimeout time.Duration

	dialOptions []grpc.DialOption

	log       *log.Logger
	logOutput io.Writer
}

// Option is the configuration func of Application.
type Option func(*Application)

// New returns a new Application with options.
func New(options ...Option) (*Application, error) {
	a := &Application{
		endpoint:    os.Getenv(endpointEnv),
		callTimeout: time.Second * 10,
		logOutput:   os.Stdout,
		dialOptions: []grpc.DialOption{grpc.WithInsecure()},
	}
	for _, option := range options {
		option(a)
	}
	a.log = log.New(a.logOutput, "mesg", log.LstdFlags)
	if a.endpoint == "" {
		a.endpoint = defaultEndpoint
	}
	return a, a.setupCoreClient()
}

// EndpointOption receives the endpoint of MESG.
func EndpointOption(address string) Option {
	return func(a *Application) {
		a.endpoint = address
	}
}

// LogOutputOption uses out as a log destination.
func LogOutputOption(out io.Writer) Option {
	return func(app *Application) {
		app.logOutput = out
	}
}

func DialOption(dialer Dialer) Option {
	return func(a *Application) {
		a.dialOptions = append(a.dialOptions, grpc.WithDialer(newGRPCDialer(dialer).Dial))
	}
}

func (a *Application) setupCoreClient() error {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), a.callTimeout)
	defer cancel()
	a.conn, err = grpc.DialContext(ctx, a.endpoint, a.dialOptions...)
	if err != nil {
		return err
	}
	a.client = core.NewCoreClient(a.conn)
	return nil
}

// Execute executes a task for serviceID with given data. Task results set into
// out, if out set to nil Execute will not block to receive a result from task.
// Use WhenResult to dynamically wait for task results.
func (a *Application) Execute(serviceID, task string, data Data) (executionID string, err error) {
	inputDataBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), a.callTimeout)
	defer cancel()
	resp, err := a.client.ExecuteTask(ctx, &core.ExecuteTaskRequest{
		ServiceID: serviceID,
		TaskKey:   task,
		InputData: string(inputDataBytes),
	})
	if err != nil {
		return "", err
	}
	return resp.ExecutionID, nil
}

func (a *Application) startServices(ids ...string) error {
	idsLen := len(ids)
	if idsLen == 0 {
		return nil
	}

	errC := make(chan error, 0)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, id := range ids {
		go a.startService(ctx, id, errC)
	}

	for i := 0; i < idsLen; i++ {
		if err := <-errC; err != nil {
			return err
		}
	}
	return nil
}

func (a *Application) startService(ctx context.Context, id string, errC chan error) {
	_, err := a.client.StartService(ctx, &core.StartServiceRequest{
		ServiceID: id,
	})
	errC <- err
}

// Close gracefully closes underlying connections and stops listening for events.
func (a *Application) Close() error {
	return a.conn.Close()
}
