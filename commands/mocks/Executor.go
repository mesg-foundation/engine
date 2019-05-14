// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import container "github.com/mesg-foundation/core/container"
import coreapi "github.com/mesg-foundation/core/protobuf/coreapi"
import io "io"
import mock "github.com/stretchr/testify/mock"
import provider "github.com/mesg-foundation/core/commands/provider"
import servicetemplate "github.com/mesg-foundation/core/utils/servicetemplate"

// Executor is an autogenerated mock type for the Executor type
type Executor struct {
	mock.Mock
}

// Create provides a mock function with given fields: passphrase
func (_m *Executor) Create(passphrase string) (string, error) {
	ret := _m.Called(passphrase)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(passphrase)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(passphrase)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateServiceOffer provides a mock function with given fields: sid, price, duration, from
func (_m *Executor) CreateServiceOffer(sid string, price string, duration string, from string) (provider.Transaction, error) {
	ret := _m.Called(sid, price, duration, from)

	var r0 provider.Transaction
	if rf, ok := ret.Get(0).(func(string, string, string, string) provider.Transaction); ok {
		r0 = rf(sid, price, duration, from)
	} else {
		r0 = ret.Get(0).(provider.Transaction)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string, string) error); ok {
		r1 = rf(sid, price, duration, from)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: address, passphrase
func (_m *Executor) Delete(address string, passphrase string) (string, error) {
	ret := _m.Called(address, passphrase)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(address, passphrase)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(address, passphrase)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Export provides a mock function with given fields: address, passphrase
func (_m *Executor) Export(address string, passphrase string) (provider.WalletEncryptedKeyJSONV3, error) {
	ret := _m.Called(address, passphrase)

	var r0 provider.WalletEncryptedKeyJSONV3
	if rf, ok := ret.Get(0).(func(string, string) provider.WalletEncryptedKeyJSONV3); ok {
		r0 = rf(address, passphrase)
	} else {
		r0 = ret.Get(0).(provider.WalletEncryptedKeyJSONV3)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(address, passphrase)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetService provides a mock function with given fields: sid
func (_m *Executor) GetService(sid string) (provider.MarketplaceService, error) {
	ret := _m.Called(sid)

	var r0 provider.MarketplaceService
	if rf, ok := ret.Get(0).(func(string) provider.MarketplaceService); ok {
		r0 = rf(sid)
	} else {
		r0 = ret.Get(0).(provider.MarketplaceService)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(sid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Import provides a mock function with given fields: account, passphrase
func (_m *Executor) Import(account provider.WalletEncryptedKeyJSONV3, passphrase string) (string, error) {
	ret := _m.Called(account, passphrase)

	var r0 string
	if rf, ok := ret.Get(0).(func(provider.WalletEncryptedKeyJSONV3, string) string); ok {
		r0 = rf(account, passphrase)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(provider.WalletEncryptedKeyJSONV3, string) error); ok {
		r1 = rf(account, passphrase)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ImportFromPrivateKey provides a mock function with given fields: privateKey, passphrase
func (_m *Executor) ImportFromPrivateKey(privateKey string, passphrase string) (string, error) {
	ret := _m.Called(privateKey, passphrase)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(privateKey, passphrase)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(privateKey, passphrase)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields:
func (_m *Executor) List() ([]string, error) {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Logs provides a mock function with given fields:
func (_m *Executor) Logs() (io.ReadCloser, error) {
	ret := _m.Called()

	var r0 io.ReadCloser
	if rf, ok := ret.Get(0).(func() io.ReadCloser); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(io.ReadCloser)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PublishServiceVersion provides a mock function with given fields: service, from
func (_m *Executor) PublishServiceVersion(service provider.MarketplaceManifestServiceData, from string) (provider.Transaction, error) {
	ret := _m.Called(service, from)

	var r0 provider.Transaction
	if rf, ok := ret.Get(0).(func(provider.MarketplaceManifestServiceData, string) provider.Transaction); ok {
		r0 = rf(service, from)
	} else {
		r0 = ret.Get(0).(provider.Transaction)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(provider.MarketplaceManifestServiceData, string) error); ok {
		r1 = rf(service, from)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Purchase provides a mock function with given fields: sid, offerIndex, from
func (_m *Executor) Purchase(sid string, offerIndex string, from string) ([]provider.Transaction, error) {
	ret := _m.Called(sid, offerIndex, from)

	var r0 []provider.Transaction
	if rf, ok := ret.Get(0).(func(string, string, string) []provider.Transaction); ok {
		r0 = rf(sid, offerIndex, from)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]provider.Transaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string) error); ok {
		r1 = rf(sid, offerIndex, from)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SendSignedTransaction provides a mock function with given fields: signedTransaction
func (_m *Executor) SendSignedTransaction(signedTransaction string) (provider.TransactionReceipt, error) {
	ret := _m.Called(signedTransaction)

	var r0 provider.TransactionReceipt
	if rf, ok := ret.Get(0).(func(string) provider.TransactionReceipt); ok {
		r0 = rf(signedTransaction)
	} else {
		r0 = ret.Get(0).(provider.TransactionReceipt)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(signedTransaction)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ServiceByID provides a mock function with given fields: id
func (_m *Executor) ServiceByID(id string) (*coreapi.Service, error) {
	ret := _m.Called(id)

	var r0 *coreapi.Service
	if rf, ok := ret.Get(0).(func(string) *coreapi.Service); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*coreapi.Service)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ServiceDelete provides a mock function with given fields: deleteData, ids
func (_m *Executor) ServiceDelete(deleteData bool, ids ...string) error {
	_va := make([]interface{}, len(ids))
	for _i := range ids {
		_va[_i] = ids[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, deleteData)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(bool, ...string) error); ok {
		r0 = rf(deleteData, ids...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ServiceDeleteAll provides a mock function with given fields: deleteData
func (_m *Executor) ServiceDeleteAll(deleteData bool) error {
	ret := _m.Called(deleteData)

	var r0 error
	if rf, ok := ret.Get(0).(func(bool) error); ok {
		r0 = rf(deleteData)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ServiceDeploy provides a mock function with given fields: path, env, statuses
func (_m *Executor) ServiceDeploy(path string, env map[string]string, statuses chan provider.DeployStatus) (string, string, error, error) {
	ret := _m.Called(path, env, statuses)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, map[string]string, chan provider.DeployStatus) string); ok {
		r0 = rf(path, env, statuses)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(string, map[string]string, chan provider.DeployStatus) string); ok {
		r1 = rf(path, env, statuses)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string, map[string]string, chan provider.DeployStatus) error); ok {
		r2 = rf(path, env, statuses)
	} else {
		r2 = ret.Error(2)
	}

	var r3 error
	if rf, ok := ret.Get(3).(func(string, map[string]string, chan provider.DeployStatus) error); ok {
		r3 = rf(path, env, statuses)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

// ServiceExecuteTask provides a mock function with given fields: id, taskKey, inputData, tags
func (_m *Executor) ServiceExecuteTask(id string, taskKey string, inputData string, tags []string) (string, error) {
	ret := _m.Called(id, taskKey, inputData, tags)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string, string, []string) string); ok {
		r0 = rf(id, taskKey, inputData, tags)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string, []string) error); ok {
		r1 = rf(id, taskKey, inputData, tags)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ServiceGenerateDocs provides a mock function with given fields: path
func (_m *Executor) ServiceGenerateDocs(path string) error {
	ret := _m.Called(path)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ServiceInitDownloadTemplate provides a mock function with given fields: t, dst
func (_m *Executor) ServiceInitDownloadTemplate(t *servicetemplate.Template, dst string) error {
	ret := _m.Called(t, dst)

	var r0 error
	if rf, ok := ret.Get(0).(func(*servicetemplate.Template, string) error); ok {
		r0 = rf(t, dst)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ServiceInitTemplateList provides a mock function with given fields:
func (_m *Executor) ServiceInitTemplateList() ([]*servicetemplate.Template, error) {
	ret := _m.Called()

	var r0 []*servicetemplate.Template
	if rf, ok := ret.Get(0).(func() []*servicetemplate.Template); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*servicetemplate.Template)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ServiceList provides a mock function with given fields:
func (_m *Executor) ServiceList() ([]*coreapi.Service, error) {
	ret := _m.Called()

	var r0 []*coreapi.Service
	if rf, ok := ret.Get(0).(func() []*coreapi.Service); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*coreapi.Service)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ServiceListenEvents provides a mock function with given fields: id, eventFilter
func (_m *Executor) ServiceListenEvents(id string, eventFilter string) (chan *coreapi.EventData, chan error, error) {
	ret := _m.Called(id, eventFilter)

	var r0 chan *coreapi.EventData
	if rf, ok := ret.Get(0).(func(string, string) chan *coreapi.EventData); ok {
		r0 = rf(id, eventFilter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(chan *coreapi.EventData)
		}
	}

	var r1 chan error
	if rf, ok := ret.Get(1).(func(string, string) chan error); ok {
		r1 = rf(id, eventFilter)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(chan error)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string, string) error); ok {
		r2 = rf(id, eventFilter)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// ServiceListenResults provides a mock function with given fields: id, taskFilter, outputFilter, tagFilters
func (_m *Executor) ServiceListenResults(id string, taskFilter string, outputFilter string, tagFilters []string) (chan *coreapi.ResultData, chan error, error) {
	ret := _m.Called(id, taskFilter, outputFilter, tagFilters)

	var r0 chan *coreapi.ResultData
	if rf, ok := ret.Get(0).(func(string, string, string, []string) chan *coreapi.ResultData); ok {
		r0 = rf(id, taskFilter, outputFilter, tagFilters)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(chan *coreapi.ResultData)
		}
	}

	var r1 chan error
	if rf, ok := ret.Get(1).(func(string, string, string, []string) chan error); ok {
		r1 = rf(id, taskFilter, outputFilter, tagFilters)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(chan error)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string, string, string, []string) error); ok {
		r2 = rf(id, taskFilter, outputFilter, tagFilters)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// ServiceLogs provides a mock function with given fields: id, dependencies
func (_m *Executor) ServiceLogs(id string, dependencies ...string) ([]*provider.Log, func(), chan error, error) {
	_va := make([]interface{}, len(dependencies))
	for _i := range dependencies {
		_va[_i] = dependencies[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, id)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []*provider.Log
	if rf, ok := ret.Get(0).(func(string, ...string) []*provider.Log); ok {
		r0 = rf(id, dependencies...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*provider.Log)
		}
	}

	var r1 func()
	if rf, ok := ret.Get(1).(func(string, ...string) func()); ok {
		r1 = rf(id, dependencies...)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(func())
		}
	}

	var r2 chan error
	if rf, ok := ret.Get(2).(func(string, ...string) chan error); ok {
		r2 = rf(id, dependencies...)
	} else {
		if ret.Get(2) != nil {
			r2 = ret.Get(2).(chan error)
		}
	}

	var r3 error
	if rf, ok := ret.Get(3).(func(string, ...string) error); ok {
		r3 = rf(id, dependencies...)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

// ServiceStart provides a mock function with given fields: id
func (_m *Executor) ServiceStart(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ServiceStop provides a mock function with given fields: id
func (_m *Executor) ServiceStop(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ServiceValidate provides a mock function with given fields: path
func (_m *Executor) ServiceValidate(path string) (string, error) {
	ret := _m.Called(path)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Sign provides a mock function with given fields: address, passphrase, transaction
func (_m *Executor) Sign(address string, passphrase string, transaction provider.Transaction) (string, error) {
	ret := _m.Called(address, passphrase, transaction)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string, provider.Transaction) string); ok {
		r0 = rf(address, passphrase, transaction)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, provider.Transaction) error); ok {
		r1 = rf(address, passphrase, transaction)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Start provides a mock function with given fields:
func (_m *Executor) Start() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Status provides a mock function with given fields:
func (_m *Executor) Status() (container.StatusType, error) {
	ret := _m.Called()

	var r0 container.StatusType
	if rf, ok := ret.Get(0).(func() container.StatusType); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(container.StatusType)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Stop provides a mock function with given fields:
func (_m *Executor) Stop() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UploadSource provides a mock function with given fields: path
func (_m *Executor) UploadSource(path string) (provider.MarketplaceDeployedSource, error) {
	ret := _m.Called(path)

	var r0 provider.MarketplaceDeployedSource
	if rf, ok := ret.Get(0).(func(string) provider.MarketplaceDeployedSource); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Get(0).(provider.MarketplaceDeployedSource)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}