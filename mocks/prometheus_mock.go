// Code generated by MockGen. DO NOT EDIT.
// Source: observability/prometheus.go
//
// Generated by this command:
//
//	mockgen -source=observability/prometheus.go -destination=./mocks/prometheus_mock.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	message "github.com/ThreeDotsLabs/watermill/message"
	gomock "go.uber.org/mock/gomock"
)

// MockMetrics is a mock of Metrics interface.
type MockMetrics struct {
	ctrl     *gomock.Controller
	recorder *MockMetricsMockRecorder
	isgomock struct{}
}

// MockMetricsMockRecorder is the mock recorder for MockMetrics.
type MockMetricsMockRecorder struct {
	mock *MockMetrics
}

// NewMockMetrics creates a new mock instance.
func NewMockMetrics(ctrl *gomock.Controller) *MockMetrics {
	mock := &MockMetrics{ctrl: ctrl}
	mock.recorder = &MockMetricsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetrics) EXPECT() *MockMetricsMockRecorder {
	return m.recorder
}

// AddPrometheusRouterMetrics mocks base method.
func (m *MockMetrics) AddPrometheusRouterMetrics(r *message.Router) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddPrometheusRouterMetrics", r)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddPrometheusRouterMetrics indicates an expected call of AddPrometheusRouterMetrics.
func (mr *MockMetricsMockRecorder) AddPrometheusRouterMetrics(r any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPrometheusRouterMetrics", reflect.TypeOf((*MockMetrics)(nil).AddPrometheusRouterMetrics), r)
}
