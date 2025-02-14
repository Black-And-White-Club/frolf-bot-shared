// Code generated by MockGen. DO NOT EDIT.
// Source: observability/loki.go
//
// Generated by this command:
//
//	mockgen -source=observability/loki.go -destination=./mocks/loki_mock.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockLogger is a mock of Logger interface.
type MockLogger struct {
	ctrl     *gomock.Controller
	recorder *MockLoggerMockRecorder
	isgomock struct{}
}

// MockLoggerMockRecorder is the mock recorder for MockLogger.
type MockLoggerMockRecorder struct {
	mock *MockLogger
}

// NewMockLogger creates a new mock instance.
func NewMockLogger(ctrl *gomock.Controller) *MockLogger {
	mock := &MockLogger{ctrl: ctrl}
	mock.recorder = &MockLoggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogger) EXPECT() *MockLoggerMockRecorder {
	return m.recorder
}

// Debug mocks base method.
func (m *MockLogger) Debug(ctx context.Context, msg string, fields map[string]any) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Debug", ctx, msg, fields)
}

// Debug indicates an expected call of Debug.
func (mr *MockLoggerMockRecorder) Debug(ctx, msg, fields any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Debug", reflect.TypeOf((*MockLogger)(nil).Debug), ctx, msg, fields)
}

// Error mocks base method.
func (m *MockLogger) Error(ctx context.Context, msg string, fields map[string]any) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Error", ctx, msg, fields)
}

// Error indicates an expected call of Error.
func (mr *MockLoggerMockRecorder) Error(ctx, msg, fields any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockLogger)(nil).Error), ctx, msg, fields)
}

// Info mocks base method.
func (m *MockLogger) Info(ctx context.Context, msg string, fields map[string]any) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Info", ctx, msg, fields)
}

// Info indicates an expected call of Info.
func (mr *MockLoggerMockRecorder) Info(ctx, msg, fields any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockLogger)(nil).Info), ctx, msg, fields)
}

// Shutdown mocks base method.
func (m *MockLogger) Shutdown() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Shutdown")
}

// Shutdown indicates an expected call of Shutdown.
func (mr *MockLoggerMockRecorder) Shutdown() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Shutdown", reflect.TypeOf((*MockLogger)(nil).Shutdown))
}

// Warn mocks base method.
func (m *MockLogger) Warn(ctx context.Context, msg string, fields map[string]any) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Warn", ctx, msg, fields)
}

// Warn indicates an expected call of Warn.
func (mr *MockLoggerMockRecorder) Warn(ctx, msg, fields any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Warn", reflect.TypeOf((*MockLogger)(nil).Warn), ctx, msg, fields)
}
