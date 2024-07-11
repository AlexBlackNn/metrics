// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/AlexBlackNn/metrics/internal/services/metricsservice (interfaces: MetricsStorage,HealthChecker)

// Package mockstorage is a generated GoMock package.
package mockstorage

import (
	context "context"
	reflect "reflect"

	models "github.com/AlexBlackNn/metrics/internal/domain/models"
	gomock "github.com/golang/mock/gomock"
)

// MockMetricsStorage is a mock of MetricsStorage interface.
type MockMetricsStorage struct {
	ctrl     *gomock.Controller
	recorder *MockMetricsStorageMockRecorder
}

// MockMetricsStorageMockRecorder is the mock recorder for MockMetricsStorage.
type MockMetricsStorageMockRecorder struct {
	mock *MockMetricsStorage
}

// NewMockMetricsStorage creates a new mock instance.
func NewMockMetricsStorage(ctrl *gomock.Controller) *MockMetricsStorage {
	mock := &MockMetricsStorage{ctrl: ctrl}
	mock.recorder = &MockMetricsStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricsStorage) EXPECT() *MockMetricsStorageMockRecorder {
	return m.recorder
}

// GetAllMetrics mocks base method.
func (m *MockMetricsStorage) GetAllMetrics(arg0 context.Context) ([]models.MetricGetter, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllMetrics", arg0)
	ret0, _ := ret[0].([]models.MetricGetter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllMetrics indicates an expected call of GetAllMetrics.
func (mr *MockMetricsStorageMockRecorder) GetAllMetrics(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllMetrics", reflect.TypeOf((*MockMetricsStorage)(nil).GetAllMetrics), arg0)
}

// GetMetric mocks base method.
func (m *MockMetricsStorage) GetMetric(arg0 context.Context, arg1 models.MetricGetter) (models.MetricGetter, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetric", arg0, arg1)
	ret0, _ := ret[0].(models.MetricGetter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMetric indicates an expected call of GetMetric.
func (mr *MockMetricsStorageMockRecorder) GetMetric(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetric", reflect.TypeOf((*MockMetricsStorage)(nil).GetMetric), arg0, arg1)
}

// UpdateMetric mocks base method.
func (m *MockMetricsStorage) UpdateMetric(arg0 context.Context, arg1 models.MetricGetter) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMetric", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMetric indicates an expected call of UpdateMetric.
func (mr *MockMetricsStorageMockRecorder) UpdateMetric(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetric", reflect.TypeOf((*MockMetricsStorage)(nil).UpdateMetric), arg0, arg1)
}

// UpdateSeveralMetrics mocks base method.
func (m *MockMetricsStorage) UpdateSeveralMetrics(arg0 context.Context, arg1 []models.MetricGetter) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSeveralMetrics", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateSeveralMetrics indicates an expected call of UpdateSeveralMetrics.
func (mr *MockMetricsStorageMockRecorder) UpdateSeveralMetrics(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSeveralMetrics", reflect.TypeOf((*MockMetricsStorage)(nil).UpdateSeveralMetrics), arg0, arg1)
}

// MockHealthChecker is a mock of HealthChecker interface.
type MockHealthChecker struct {
	ctrl     *gomock.Controller
	recorder *MockHealthCheckerMockRecorder
}

// MockHealthCheckerMockRecorder is the mock recorder for MockHealthChecker.
type MockHealthCheckerMockRecorder struct {
	mock *MockHealthChecker
}

// NewMockHealthChecker creates a new mock instance.
func NewMockHealthChecker(ctrl *gomock.Controller) *MockHealthChecker {
	mock := &MockHealthChecker{ctrl: ctrl}
	mock.recorder = &MockHealthCheckerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHealthChecker) EXPECT() *MockHealthCheckerMockRecorder {
	return m.recorder
}

// HealthCheck mocks base method.
func (m *MockHealthChecker) HealthCheck(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HealthCheck", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// HealthCheck indicates an expected call of HealthCheck.
func (mr *MockHealthCheckerMockRecorder) HealthCheck(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HealthCheck", reflect.TypeOf((*MockHealthChecker)(nil).HealthCheck), arg0)
}
