// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	auth "github.com/gidyon/antibug/internal/pkg/auth"

	mock "github.com/stretchr/testify/mock"
)

// AuthAPI is an autogenerated mock type for the AuthAPI type
type AuthAPI struct {
	mock.Mock
}

// AuthenticateRequest provides a mock function with given fields: _a0
func (_m *AuthAPI) AuthenticateRequest(_a0 context.Context) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AuthorizeActor provides a mock function with given fields: ctx, actorID
func (_m *AuthAPI) AuthorizeActor(ctx context.Context, actorID string) (*auth.Payload, error) {
	ret := _m.Called(ctx, actorID)

	var r0 *auth.Payload
	if rf, ok := ret.Get(0).(func(context.Context, string) *auth.Payload); ok {
		r0 = rf(ctx, actorID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*auth.Payload)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, actorID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AuthorizeGroup provides a mock function with given fields: ctx, allowedGroups
func (_m *AuthAPI) AuthorizeGroup(ctx context.Context, allowedGroups ...string) (*auth.Payload, error) {
	_va := make([]interface{}, len(allowedGroups))
	for _i := range allowedGroups {
		_va[_i] = allowedGroups[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *auth.Payload
	if rf, ok := ret.Get(0).(func(context.Context, ...string) *auth.Payload); ok {
		r0 = rf(ctx, allowedGroups...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*auth.Payload)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, ...string) error); ok {
		r1 = rf(ctx, allowedGroups...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AuthorizeStrict provides a mock function with given fields: ctx, actorID, allowedGroups
func (_m *AuthAPI) AuthorizeStrict(ctx context.Context, actorID string, allowedGroups ...string) (*auth.Payload, error) {
	_va := make([]interface{}, len(allowedGroups))
	for _i := range allowedGroups {
		_va[_i] = allowedGroups[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, actorID)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *auth.Payload
	if rf, ok := ret.Get(0).(func(context.Context, string, ...string) *auth.Payload); ok {
		r0 = rf(ctx, actorID, allowedGroups...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*auth.Payload)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, ...string) error); ok {
		r1 = rf(ctx, actorID, allowedGroups...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GenToken provides a mock function with given fields: _a0, _a1, _a2
func (_m *AuthAPI) GenToken(_a0 context.Context, _a1 *auth.Payload, _a2 int64) (string, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, *auth.Payload, int64) string); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *auth.Payload, int64) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
