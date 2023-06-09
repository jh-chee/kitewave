// Code generated by mockery v2.28.1. DO NOT EDIT.

package mocks

import (
	models "github.com/jh-chee/kitewave/rpc-server/models"
	mock "github.com/stretchr/testify/mock"
)

// MessageService is an autogenerated mock type for the MessageService type
type MessageService struct {
	mock.Mock
}

// Pull provides a mock function with given fields: req
func (_m *MessageService) Pull(req *models.Request) (*models.Response, error) {
	ret := _m.Called(req)

	var r0 *models.Response
	var r1 error
	if rf, ok := ret.Get(0).(func(*models.Request) (*models.Response, error)); ok {
		return rf(req)
	}
	if rf, ok := ret.Get(0).(func(*models.Request) *models.Response); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Response)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.Request) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Send provides a mock function with given fields: msg
func (_m *MessageService) Send(msg *models.Message) error {
	ret := _m.Called(msg)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.Message) error); ok {
		r0 = rf(msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewMessageService interface {
	mock.TestingT
	Cleanup(func())
}

// NewMessageService creates a new instance of MessageService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMessageService(t mockConstructorTestingTNewMessageService) *MessageService {
	mock := &MessageService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
