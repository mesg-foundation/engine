package commands

import (
	"io"

	"github.com/mesg-foundation/core/commands/provider"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/utils/servicetemplate"
	"github.com/stretchr/testify/mock"
)

type mockRootExecutor struct {
	mock.Mock
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
	mock.Mock
}

func (m *mockServiceExecutor) ServiceByID(id string) (*service.Service, error) {
	args := m.Called()
	return args.Get(0).(*service.Service), args.Error(1)
}

func (m *mockServiceExecutor) ServiceDeleteAll() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockServiceExecutor) ServiceDelete(ids ...string) error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockServiceExecutor) ServiceDeploy(path string) (id string, valid bool, err error) {
	args := m.Called()
	return args.String(0), args.Bool(1), args.Error(2)
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

func (m *mockServiceExecutor) ServiceExecuteTask(id, taskKey, inputData string, tags []string) (listenResults chan coreapi.ResultData, err error) {
	args := m.Called()
	return args.Get(0).(chan coreapi.ResultData), args.Error(1)
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

func (m *mockServiceExecutor) ServiceList() ([]*service.Service, error) {
	args := m.Called()
	return args.Get(0).([]*service.Service), args.Error(1)
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
