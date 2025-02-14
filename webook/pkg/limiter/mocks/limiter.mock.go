// Code generated by MockGen. DO NOT EDIT.
// Source: ./webook/pkg/limiter/types.go
//
// Generated by this command:
//
//	mockgen -source=./webook/pkg/limiter/types.go -package=limitermocks -destination=./webook/pkg/limiter/mocks/limiter.mock.go
//

// Package limitermocks is a generated GoMock package.
package limitermocks

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockLimiter is a mock of Limiter interface.
type MockLimiter struct {
	ctrl     *gomock.Controller
	recorder *MockLimiterMockRecorder
	isgomock struct{}
}

// MockLimiterMockRecorder is the mock recorder for MockLimiter.
type MockLimiterMockRecorder struct {
	mock *MockLimiter
}

// NewMockLimiter creates a new mock instance.
func NewMockLimiter(ctrl *gomock.Controller) *MockLimiter {
	mock := &MockLimiter{ctrl: ctrl}
	mock.recorder = &MockLimiterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLimiter) EXPECT() *MockLimiterMockRecorder {
	return m.recorder
}

// Limit mocks base method.
func (m *MockLimiter) Limit(ctx context.Context, key string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Limit", ctx, key)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Limit indicates an expected call of Limit.
func (mr *MockLimiterMockRecorder) Limit(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Limit", reflect.TypeOf((*MockLimiter)(nil).Limit), ctx, key)
}
