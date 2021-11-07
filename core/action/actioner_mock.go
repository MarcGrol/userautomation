// Code generated by MockGen. DO NOT EDIT.
// Source: action.go

// Package action is a generated GoMock package.
package action

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUserActioner is a mock of UserActioner interface.
type MockUserActioner struct {
	ctrl     *gomock.Controller
	recorder *MockUserActionerMockRecorder
}

// MockUserActionerMockRecorder is the mock recorder for MockUserActioner.
type MockUserActionerMockRecorder struct {
	mock *MockUserActioner
}

// NewMockUserActioner creates a new mock instance.
func NewMockUserActioner(ctrl *gomock.Controller) *MockUserActioner {
	mock := &MockUserActioner{ctrl: ctrl}
	mock.recorder = &MockUserActionerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserActioner) EXPECT() *MockUserActionerMockRecorder {
	return m.recorder
}

// Perform mocks base method.
func (m *MockUserActioner) Perform(ctx context.Context, action UserAction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Perform", ctx, action)
	ret0, _ := ret[0].(error)
	return ret0
}

// Perform indicates an expected call of Perform.
func (mr *MockUserActionerMockRecorder) Perform(ctx, action interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Perform", reflect.TypeOf((*MockUserActioner)(nil).Perform), ctx, action)
}