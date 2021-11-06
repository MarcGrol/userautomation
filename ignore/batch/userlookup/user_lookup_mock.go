// Code generated by MockGen. DO NOT EDIT.
// Source: user_lookup.go

// Package userlookup is a generated GoMock package.
package userlookup

import (
	context "context"
	reflect "reflect"

	batchcore "github.com/MarcGrol/userautomation/ignore/batch/batchcore"
	gomock "github.com/golang/mock/gomock"
)

// MockUserLookuper is a mock of UserLookuper interface.
type MockUserLookuper struct {
	ctrl     *gomock.Controller
	recorder *MockUserLookuperMockRecorder
}

// MockUserLookuperMockRecorder is the mock recorder for MockUserLookuper.
type MockUserLookuperMockRecorder struct {
	mock *MockUserLookuper
}

// NewMockUserLookuper creates a new mock instance.
func NewMockUserLookuper(ctrl *gomock.Controller) *MockUserLookuper {
	mock := &MockUserLookuper{ctrl: ctrl}
	mock.recorder = &MockUserLookuperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserLookuper) EXPECT() *MockUserLookuperMockRecorder {
	return m.recorder
}

// GetUserOnUid mocks base method.
func (m *MockUserLookuper) GetUserOnUid(c context.Context, uid string) (batchcore.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserOnUid", c, uid)
	ret0, _ := ret[0].(batchcore.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserOnUid indicates an expected call of GetUserOnUid.
func (mr *MockUserLookuperMockRecorder) GetUserOnUid(c, uid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserOnUid", reflect.TypeOf((*MockUserLookuper)(nil).GetUserOnUid), c, uid)
}

// GetUserOnQuery mocks base method.
func (m *MockUserLookuper) GetUserOnQuery(c context.Context, whereClause string) ([]batchcore.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserOnQuery", c, whereClause)
	ret0, _ := ret[0].([]batchcore.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserOnQuery indicates an expected call of GetUserOnQuery.
func (mr *MockUserLookuperMockRecorder) GetUserOnQuery(c, whereClause interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserOnQuery", reflect.TypeOf((*MockUserLookuper)(nil).GetUserOnQuery), c, whereClause)
}
