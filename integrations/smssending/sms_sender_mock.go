// Code generated by MockGen. DO NOT EDIT.
// Source: user_action.go

// Package smssending is a generated GoMock package.
package smssending

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockSmsSender is a mock of SmsSender interface.
type MockSmsSender struct {
	ctrl     *gomock.Controller
	recorder *MockSmsSenderMockRecorder
}

// MockSmsSenderMockRecorder is the mock recorder for MockSmsSender.
type MockSmsSenderMockRecorder struct {
	mock *MockSmsSender
}

// NewMockSmsSender creates a new mock instance.
func NewMockSmsSender(ctrl *gomock.Controller) *MockSmsSender {
	mock := &MockSmsSender{ctrl: ctrl}
	mock.recorder = &MockSmsSenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSmsSender) EXPECT() *MockSmsSenderMockRecorder {
	return m.recorder
}

// Send mocks base method.
func (m *MockSmsSender) Send(c context.Context, recipient, body string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", c, recipient, body)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockSmsSenderMockRecorder) Send(c, recipient, body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockSmsSender)(nil).Send), c, recipient, body)
}
