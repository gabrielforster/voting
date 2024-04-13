// Code generated by mockery v2.18.0. DO NOT EDIT.
// https://github.com/eminetto/api-o11y/blob/main/auth/user/mocks/Repository.go

package mocks

import (
	context "context"

	user "github.com/gabrielforster/voting/auth/user"
	mock "github.com/stretchr/testify/mock"
)

type Repository struct {
	mock.Mock
}

func (_m *Repository) Get(ctx context.Context, email string) (*user.User, error) {
	ret := _m.Called(ctx, email)

	var r0 *user.User
	if rf, ok := ret.Get(0).(func(context.Context, string) *user.User); ok {
		r0 = rf(ctx, email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*user.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewRepository interface {
	mock.TestingT
	Cleanup(func())
}

func NewRepository(t mockConstructorTestingTNewRepository) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
