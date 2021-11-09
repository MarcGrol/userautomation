// Code generated by MockGen. DO NOT EDIT.
// Source: api.go

// Package taskqueue is a generated GoMock package.
package taskqueue

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockTaskQueue is a mock of TaskQueue interface.
type MockTaskQueue struct {
	ctrl     *gomock.Controller
	recorder *MockTaskQueueMockRecorder
}

// MockTaskQueueMockRecorder is the mock recorder for MockTaskQueue.
type MockTaskQueueMockRecorder struct {
	mock *MockTaskQueue
}

// NewMockTaskQueue creates a new mock instance.
func NewMockTaskQueue(ctrl *gomock.Controller) *MockTaskQueue {
	mock := &MockTaskQueue{ctrl: ctrl}
	mock.recorder = &MockTaskQueueMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTaskQueue) EXPECT() *MockTaskQueueMockRecorder {
	return m.recorder
}

// Enqueue mocks base method.
func (m *MockTaskQueue) Enqueue(ctx context.Context, task Task) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Enqueue", ctx, task)
	ret0, _ := ret[0].(error)
	return ret0
}

// Enqueue indicates an expected call of Enqueue.
func (mr *MockTaskQueueMockRecorder) Enqueue(ctx, task interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Enqueue", reflect.TypeOf((*MockTaskQueue)(nil).Enqueue), ctx, task)
}

// MockTaskQueueReceiver is a mock of TaskQueueReceiver interface.
type MockTaskQueueReceiver struct {
	ctrl     *gomock.Controller
	recorder *MockTaskQueueReceiverMockRecorder
}

// MockTaskQueueReceiverMockRecorder is the mock recorder for MockTaskQueueReceiver.
type MockTaskQueueReceiverMockRecorder struct {
	mock *MockTaskQueueReceiver
}

// NewMockTaskQueueReceiver creates a new mock instance.
func NewMockTaskQueueReceiver(ctrl *gomock.Controller) *MockTaskQueueReceiver {
	mock := &MockTaskQueueReceiver{ctrl: ctrl}
	mock.recorder = &MockTaskQueueReceiverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTaskQueueReceiver) EXPECT() *MockTaskQueueReceiverMockRecorder {
	return m.recorder
}

// IamReceivingTasks mocks base method.
func (m *MockTaskQueueReceiver) IamReceivingTasks() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "IamReceivingTasks")
}

// IamReceivingTasks indicates an expected call of IamReceivingTasks.
func (mr *MockTaskQueueReceiverMockRecorder) IamReceivingTasks() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IamReceivingTasks", reflect.TypeOf((*MockTaskQueueReceiver)(nil).IamReceivingTasks))
}

// OnTaskReceived mocks base method.
func (m *MockTaskQueueReceiver) OnTaskReceived(c context.Context, task Task) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OnTaskReceived", c, task)
	ret0, _ := ret[0].(error)
	return ret0
}

// OnTaskReceived indicates an expected call of OnTaskReceived.
func (mr *MockTaskQueueReceiverMockRecorder) OnTaskReceived(c, task interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnTaskReceived", reflect.TypeOf((*MockTaskQueueReceiver)(nil).OnTaskReceived), c, task)
}