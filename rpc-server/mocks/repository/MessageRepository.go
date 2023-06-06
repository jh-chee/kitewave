// Code generated by mockery v2.28.1. DO NOT EDIT.

package mocks

import (
	models "github.com/jh-chee/kitewave/rpc-server/models"
	mock "github.com/stretchr/testify/mock"
)

// MessageRepository is an autogenerated mock type for the MessageRepository type
type MessageRepository struct {
	mock.Mock
}

// CheckCursorExistence provides a mock function with given fields: cursor
func (_m *MessageRepository) CheckCursorExistence(cursor int64) error {
	ret := _m.Called(cursor)

	var r0 error
	if rf, ok := ret.Get(0).(func(int64) error); ok {
		r0 = rf(cursor)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Pull provides a mock function with given fields: req
func (_m *MessageRepository) Pull(req *models.Request) ([]*models.Message, int64, error) {
	ret := _m.Called(req)

	var r0 []*models.Message
	var r1 int64
	var r2 error
	if rf, ok := ret.Get(0).(func(*models.Request) ([]*models.Message, int64, error)); ok {
		return rf(req)
	}
	if rf, ok := ret.Get(0).(func(*models.Request) []*models.Message); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Message)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.Request) int64); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Get(1).(int64)
	}

	if rf, ok := ret.Get(2).(func(*models.Request) error); ok {
		r2 = rf(req)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Save provides a mock function with given fields: message
func (_m *MessageRepository) Save(message *models.Message) error {
	ret := _m.Called(message)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.Message) error); ok {
		r0 = rf(message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewMessageRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewMessageRepository creates a new instance of MessageRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMessageRepository(t mockConstructorTestingTNewMessageRepository) *MessageRepository {
	mock := &MessageRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
