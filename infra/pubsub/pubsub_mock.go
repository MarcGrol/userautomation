// Code generated by MockGen. DO NOT EDIT.
// Source: api.go

// Package pubsub is a generated GoMock package.
package pubsub

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	mux "github.com/gorilla/mux"
)

// MockPubsub is a mock of Pubsub interface.
type MockPubsub struct {
	ctrl     *gomock.Controller
	recorder *MockPubsubMockRecorder
}

// MockPubsubMockRecorder is the mock recorder for MockPubsub.
type MockPubsubMockRecorder struct {
	mock *MockPubsub
}

// NewMockPubsub creates a new mock instance.
func NewMockPubsub(ctrl *gomock.Controller) *MockPubsub {
	mock := &MockPubsub{ctrl: ctrl}
	mock.recorder = &MockPubsubMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPubsub) EXPECT() *MockPubsubMockRecorder {
	return m.recorder
}

// Subscribe mocks base method.
func (m *MockPubsub) Subscribe(ctx context.Context, topic string, onEvent OnEventFunc) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subscribe", ctx, topic, onEvent)
	ret0, _ := ret[0].(error)
	return ret0
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockPubsubMockRecorder) Subscribe(ctx, topic, onEvent interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockPubsub)(nil).Subscribe), ctx, topic, onEvent)
}

// Publish mocks base method.
func (m *MockPubsub) Publish(ctx context.Context, topic string, event interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Publish", ctx, topic, event)
	ret0, _ := ret[0].(error)
	return ret0
}

// Publish indicates an expected call of Publish.
func (mr *MockPubsubMockRecorder) Publish(ctx, topic, event interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockPubsub)(nil).Publish), ctx, topic, event)
}

// MockSubscribingService is a mock of SubscribingService interface.
type MockSubscribingService struct {
	ctrl     *gomock.Controller
	recorder *MockSubscribingServiceMockRecorder
}

// MockSubscribingServiceMockRecorder is the mock recorder for MockSubscribingService.
type MockSubscribingServiceMockRecorder struct {
	mock *MockSubscribingService
}

// NewMockSubscribingService creates a new mock instance.
func NewMockSubscribingService(ctrl *gomock.Controller) *MockSubscribingService {
	mock := &MockSubscribingService{ctrl: ctrl}
	mock.recorder = &MockSubscribingServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSubscribingService) EXPECT() *MockSubscribingServiceMockRecorder {
	return m.recorder
}

// IamSubscribing mocks base method.
func (m *MockSubscribingService) IamSubscribing() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "IamSubscribing")
}

// IamSubscribing indicates an expected call of IamSubscribing.
func (mr *MockSubscribingServiceMockRecorder) IamSubscribing() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IamSubscribing", reflect.TypeOf((*MockSubscribingService)(nil).IamSubscribing))
}

// Subscribe mocks base method.
func (m *MockSubscribingService) Subscribe(c context.Context, router *mux.Router) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subscribe", c, router)
	ret0, _ := ret[0].(error)
	return ret0
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockSubscribingServiceMockRecorder) Subscribe(c, router interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockSubscribingService)(nil).Subscribe), c, router)
}

// OnEvent mocks base method.
func (m *MockSubscribingService) OnEvent(c context.Context, topic string, event interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OnEvent", c, topic, event)
	ret0, _ := ret[0].(error)
	return ret0
}

// OnEvent indicates an expected call of OnEvent.
func (mr *MockSubscribingServiceMockRecorder) OnEvent(c, topic, event interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnEvent", reflect.TypeOf((*MockSubscribingService)(nil).OnEvent), c, topic, event)
}
