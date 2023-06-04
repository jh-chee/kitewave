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

// Save provides a mock function with given fields: message
func (_m *MessageRepository) Save(message *models.Message) (*models.Message, error) {
	ret := _m.Called(message)

	var r0 *models.Message
	var r1 error
	if rf, ok := ret.Get(0).(func(*models.Message) (*models.Message, error)); ok {
		return rf(message)
	}
	if rf, ok := ret.Get(0).(func(*models.Message) *models.Message); ok {
		r0 = rf(message)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Message)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.Message) error); ok {
		r1 = rf(message)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
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
