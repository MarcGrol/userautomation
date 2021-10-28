// Code generated by MockGen. DO NOT EDIT.
// Source: group_api.go

// Package actions is a generated GoMock package.
package actions

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockGroupApi is a mock of GroupApi interface.
type MockGroupApi struct {
	ctrl     *gomock.Controller
	recorder *MockGroupApiMockRecorder
}

// MockGroupApiMockRecorder is the mock recorder for MockGroupApi.
type MockGroupApiMockRecorder struct {
	mock *MockGroupApi
}

// NewMockGroupApi creates a new mock instance.
func NewMockGroupApi(ctrl *gomock.Controller) *MockGroupApi {
	mock := &MockGroupApi{ctrl: ctrl}
	mock.recorder = &MockGroupApiMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGroupApi) EXPECT() *MockGroupApiMockRecorder {
	return m.recorder
}

// GroupExists mocks base method.
func (m *MockGroupApi) GroupExists(groupUid string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GroupExists", groupUid)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GroupExists indicates an expected call of GroupExists.
func (mr *MockGroupApiMockRecorder) GroupExists(groupUid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GroupExists", reflect.TypeOf((*MockGroupApi)(nil).GroupExists), groupUid)
}

// AddUserToGroup mocks base method.
func (m *MockGroupApi) AddUserToGroup(groupName, userUid string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUserToGroup", groupName, userUid)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddUserToGroup indicates an expected call of AddUserToGroup.
func (mr *MockGroupApiMockRecorder) AddUserToGroup(groupName, userUid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUserToGroup", reflect.TypeOf((*MockGroupApi)(nil).AddUserToGroup), groupName, userUid)
}
