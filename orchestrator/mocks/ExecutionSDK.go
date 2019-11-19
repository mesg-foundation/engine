// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import api "github.com/mesg-foundation/engine/protobuf/api"
import execution "github.com/mesg-foundation/engine/execution"
import hash "github.com/mesg-foundation/engine/hash"
import mock "github.com/stretchr/testify/mock"

// ExecutionSDK is an autogenerated mock type for the ExecutionSDK type
type ExecutionSDK struct {
	mock.Mock
}

// Create provides a mock function with given fields: req, accountName, accountPassword
func (_m *ExecutionSDK) Create(req *api.CreateExecutionRequest, accountName string, accountPassword string) (*execution.Execution, error) {
	ret := _m.Called(req, accountName, accountPassword)

	var r0 *execution.Execution
	if rf, ok := ret.Get(0).(func(*api.CreateExecutionRequest, string, string) *execution.Execution); ok {
		r0 = rf(req, accountName, accountPassword)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*execution.Execution)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*api.CreateExecutionRequest, string, string) error); ok {
		r1 = rf(req, accountName, accountPassword)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: _a0
func (_m *ExecutionSDK) Get(_a0 hash.Hash) (*execution.Execution, error) {
	ret := _m.Called(_a0)

	var r0 *execution.Execution
	if rf, ok := ret.Get(0).(func(hash.Hash) *execution.Execution); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*execution.Execution)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(hash.Hash) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Stream provides a mock function with given fields: req
func (_m *ExecutionSDK) Stream(req *api.StreamExecutionRequest) (chan *execution.Execution, error) {
	ret := _m.Called(req)

	var r0 chan *execution.Execution
	if rf, ok := ret.Get(0).(func(*api.StreamExecutionRequest) chan *execution.Execution); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(chan *execution.Execution)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*api.StreamExecutionRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
