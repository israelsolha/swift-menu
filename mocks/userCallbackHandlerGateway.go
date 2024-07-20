// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	entities "swift-menu-session/internal/domain/entities"

	mock "github.com/stretchr/testify/mock"
)

// UserCallbackHandlerGateway is an autogenerated mock type for the userCallbackHandlerGateway type
type UserCallbackHandlerGateway struct {
	mock.Mock
}

// CreateUser provides a mock function with given fields: user
func (_m *UserCallbackHandlerGateway) CreateUser(user entities.User) (entities.User, error) {
	ret := _m.Called(user)

	if len(ret) == 0 {
		panic("no return value specified for CreateUser")
	}

	var r0 entities.User
	var r1 error
	if rf, ok := ret.Get(0).(func(entities.User) (entities.User, error)); ok {
		return rf(user)
	}
	if rf, ok := ret.Get(0).(func(entities.User) entities.User); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Get(0).(entities.User)
	}

	if rf, ok := ret.Get(1).(func(entities.User) error); ok {
		r1 = rf(user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByEmail provides a mock function with given fields: email
func (_m *UserCallbackHandlerGateway) GetUserByEmail(email string) (entities.User, error) {
	ret := _m.Called(email)

	if len(ret) == 0 {
		panic("no return value specified for GetUserByEmail")
	}

	var r0 entities.User
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (entities.User, error)); ok {
		return rf(email)
	}
	if rf, ok := ret.Get(0).(func(string) entities.User); ok {
		r0 = rf(email)
	} else {
		r0 = ret.Get(0).(entities.User)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUserCallbackHandlerGateway creates a new instance of UserCallbackHandlerGateway. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserCallbackHandlerGateway(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserCallbackHandlerGateway {
	mock := &UserCallbackHandlerGateway{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
