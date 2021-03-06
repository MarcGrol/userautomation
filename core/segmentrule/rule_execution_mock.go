// Code generated by MockGen. DO NOT EDIT.
// Source: segment_rule.go

// Package segmentrule is a generated GoMock package.
package segmentrule

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockService is a mock of Management interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// Put mocks base method.
func (m *MockService) Put(ctx context.Context, rule Spec) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", ctx, rule)
	ret0, _ := ret[0].(error)
	return ret0
}

// Put indicates an expected call of Put.
func (mr *MockServiceMockRecorder) Put(ctx, rule interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockService)(nil).Put), ctx, rule)
}

// Get mocks base method.
func (m *MockService) Get(ctx context.Context, ruleUID string) (Spec, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, ruleUID)
	ret0, _ := ret[0].(Spec)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockServiceMockRecorder) Get(ctx, ruleUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockService)(nil).Get), ctx, ruleUID)
}

// Remove mocks base method.
func (m *MockService) Remove(ctx context.Context, ruleUID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", ctx, ruleUID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove.
func (mr *MockServiceMockRecorder) Remove(ctx, ruleUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockService)(nil).Remove), ctx, ruleUID)
}

// List mocks base method.
func (m *MockService) List(ctx context.Context) ([]Spec, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx)
	ret0, _ := ret[0].([]Spec)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockServiceMockRecorder) List(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockService)(nil).List), ctx)
}

// MockTriggerRuleExecution is a mock of TriggerRuleExecution interface.
type MockTriggerRuleExecution struct {
	ctrl     *gomock.Controller
	recorder *MockTriggerRuleExecutionMockRecorder
}

// MockTriggerRuleExecutionMockRecorder is the mock recorder for MockTriggerRuleExecution.
type MockTriggerRuleExecutionMockRecorder struct {
	mock *MockTriggerRuleExecution
}

// NewMockTriggerRuleExecution creates a new mock instance.
func NewMockTriggerRuleExecution(ctrl *gomock.Controller) *MockTriggerRuleExecution {
	mock := &MockTriggerRuleExecution{ctrl: ctrl}
	mock.recorder = &MockTriggerRuleExecutionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTriggerRuleExecution) EXPECT() *MockTriggerRuleExecutionMockRecorder {
	return m.recorder
}

// Trigger mocks base method.
func (m *MockTriggerRuleExecution) Trigger(ctx context.Context, rule Spec) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trigger", ctx, rule)
	ret0, _ := ret[0].(error)
	return ret0
}

// Trigger indicates an expected call of Trigger.
func (mr *MockTriggerRuleExecutionMockRecorder) Trigger(ctx, rule interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trigger", reflect.TypeOf((*MockTriggerRuleExecution)(nil).Trigger), ctx, rule)
}
