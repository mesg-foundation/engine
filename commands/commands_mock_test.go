package commands

import (
	"io"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/mesg-foundation/core/commands/provider"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/mesg-foundation/core/utils/servicetemplate"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	_ RootExecutor    = (*mockRootExecutor)(nil)
	_ ServiceExecutor = (*mockServiceExecutor)(nil)
)

// captureStd is helper function that captures Stdout and Stderr and returns function
// that returns standard output and standard error as string.
func captureStd(t *testing.T) func() (stdout string, stderr string) {
	var (
		bufout strings.Builder
		buferr strings.Builder
		wg     sync.WaitGroup

		stdout = os.Stdout
		stderr = os.Stderr
	)

	or, ow, err := os.Pipe()
	require.NoError(t, err)

	er, ew, err := os.Pipe()
	require.NoError(t, err)

	os.Stdout = ow
	os.Stderr = ew

	wg.Add(1)
	// copy out and err to buffers
	go func() {
		_, err := io.Copy(&bufout, or)
		require.NoError(t, err)
		or.Close()

		_, err = io.Copy(&buferr, er)
		require.NoError(t, err)
		er.Close()

		wg.Done()
	}()

	return func() (string, string) {
		// close writers and wait for copy to finish
		ow.Close()
		ew.Close()
		wg.Wait()

		// set back oginal stdout and stderr
		os.Stdout = stdout
		os.Stderr = stderr

		// return stdout and stderr
		return bufout.String(), buferr.String()
	}
}

type mockExecutor struct {
	*mock.Mock
	*mockRootExecutor
	*mockServiceExecutor
	*mockWorkflowExecutor
}

func newMockExecutor() *mockExecutor {
	m := &mock.Mock{}
	return &mockExecutor{
		Mock:                 m,
		mockRootExecutor:     &mockRootExecutor{m},
		mockServiceExecutor:  &mockServiceExecutor{m},
		mockWorkflowExecutor: &mockWorkflowExecutor{m},
	}
}

type mockRootExecutor struct {
	*mock.Mock
}

func (m *mockRootExecutor) Start() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockRootExecutor) Stop() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockRootExecutor) Status() (container.StatusType, error) {
	args := m.Called()
	return args.Get(0).(container.StatusType), args.Error(1)
}

func (m *mockRootExecutor) Logs() (io.ReadCloser, error) {
	args := m.Called()
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

type mockServiceExecutor struct {
	*mock.Mock
}

func (m *mockServiceExecutor) ServiceByID(id string) (*coreapi.Service, error) {
	args := m.Called()
	return args.Get(0).(*coreapi.Service), args.Error(1)
}

func (m *mockServiceExecutor) ServiceDeleteAll() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockServiceExecutor) ServiceDelete(ids ...string) error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockServiceExecutor) ServiceDeploy(path string, statuses chan provider.DeployStatus) (id string, validationError, err error) {
	args := m.Called()
	return args.String(0), args.Error(1), args.Error(2)
}

func (m *mockServiceExecutor) ServiceListenEvents(id, taskFilter string) (chan *coreapi.EventData, chan error, error) {
	args := m.Called()
	return args.Get(0).(chan *coreapi.EventData), args.Get(1).(chan error), args.Error(2)
}

func (m *mockServiceExecutor) ServiceListenResults(id, taskFilter, outputFilter string, tagFilters []string) (chan *coreapi.ResultData, chan error, error) {
	args := m.Called()
	return args.Get(0).(chan *coreapi.ResultData), args.Get(1).(chan error), args.Error(2)
}

func (m *mockServiceExecutor) ServiceLogs(id string, dependencies ...string) (logs []*provider.Log, close func(), err error) {
	args := m.Called()
	return args.Get(0).([]*provider.Log), args.Get(1).(func()), args.Error(2)
}

func (m *mockServiceExecutor) ServiceExecuteTask(id, taskKey, inputData string, tags []string) error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockServiceExecutor) ServiceStart(id string) error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockServiceExecutor) ServiceStop(id string) error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockServiceExecutor) ServiceValidate(path string) (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *mockServiceExecutor) ServiceGenerateDocs(path string) error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockServiceExecutor) ServiceList() ([]*coreapi.Service, error) {
	args := m.Called()
	return args.Get(0).([]*coreapi.Service), args.Error(1)
}

func (m *mockServiceExecutor) ServiceInit(name, description, templateURL string, currentDir bool) error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockServiceExecutor) ServiceInitTemplateList() ([]*servicetemplate.Template, error) {
	args := m.Called()
	return args.Get(0).([]*servicetemplate.Template), args.Error(1)
}

func (m *mockServiceExecutor) ServiceInitDownloadTemplate(t *servicetemplate.Template, dst string) error {
	args := m.Called()
	return args.Error(0)
}

type mockWorkflowExecutor struct {
	*mock.Mock
}

func (m *mockWorkflowExecutor) CreateWorkflow(filePath string, name string) (id string, err error) {
	args := m.Called(filePath, name)
	return args.Get(0).(string), args.Error(1)
}

func (m *mockWorkflowExecutor) DeleteWorkflow(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
