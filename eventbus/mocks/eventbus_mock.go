// Code generated by MockGen. DO NOT EDIT.
// Source: eventbus/eventbus.go
//
// Generated by this command:
//
//	mockgen -source=eventbus/eventbus.go -destination=eventbus/mocks/eventbus_mock.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	message "github.com/ThreeDotsLabs/watermill/message"
	gomock "go.uber.org/mock/gomock"
)

// MockEventBus is a mock of EventBus interface.
type MockEventBus struct {
	ctrl     *gomock.Controller
	recorder *MockEventBusMockRecorder
	isgomock struct{}
}

// MockEventBusMockRecorder is the mock recorder for MockEventBus.
type MockEventBusMockRecorder struct {
	mock *MockEventBus
}

// NewMockEventBus creates a new mock instance.
func NewMockEventBus(ctrl *gomock.Controller) *MockEventBus {
	mock := &MockEventBus{ctrl: ctrl}
	mock.recorder = &MockEventBusMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEventBus) EXPECT() *MockEventBusMockRecorder {
	return m.recorder
}

// CancelScheduledMessage mocks base method.
func (m *MockEventBus) CancelScheduledMessage(ctx context.Context, roundID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CancelScheduledMessage", ctx, roundID)
	ret0, _ := ret[0].(error)
	return ret0
}

// CancelScheduledMessage indicates an expected call of CancelScheduledMessage.
func (mr *MockEventBusMockRecorder) CancelScheduledMessage(ctx, roundID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CancelScheduledMessage", reflect.TypeOf((*MockEventBus)(nil).CancelScheduledMessage), ctx, roundID)
}

// Close mocks base method.
func (m *MockEventBus) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockEventBusMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockEventBus)(nil).Close))
}

// ProcessDelayedMessages mocks base method.
func (m *MockEventBus) ProcessDelayedMessages(ctx context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ProcessDelayedMessages", ctx)
}

// ProcessDelayedMessages indicates an expected call of ProcessDelayedMessages.
func (mr *MockEventBusMockRecorder) ProcessDelayedMessages(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessDelayedMessages", reflect.TypeOf((*MockEventBus)(nil).ProcessDelayedMessages), ctx)
}

// Publish mocks base method.
func (m *MockEventBus) Publish(topic string, messages ...*message.Message) error {
	m.ctrl.T.Helper()
	varargs := []any{topic}
	for _, a := range messages {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Publish", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Publish indicates an expected call of Publish.
func (mr *MockEventBusMockRecorder) Publish(topic any, messages ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{topic}, messages...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockEventBus)(nil).Publish), varargs...)
}

// Subscribe mocks base method.
func (m *MockEventBus) Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subscribe", ctx, topic)
	ret0, _ := ret[0].(<-chan *message.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockEventBusMockRecorder) Subscribe(ctx, topic any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockEventBus)(nil).Subscribe), ctx, topic)
}
