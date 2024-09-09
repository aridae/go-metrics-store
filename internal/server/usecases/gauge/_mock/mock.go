// Code generated by MockGen. DO NOT EDIT.
// Source: ./../handler.go
//
// Generated by this command:
//
//	mockgen -package=_mock -destination=./mock.go -source=./../handler.go
//

// Package _mock is a generated GoMock package.
package _mock

import (
	context "context"
	reflect "reflect"

	models "github.com/aridae/go-metrics-store/internal/server/models"
	gomock "go.uber.org/mock/gomock"
)

// MockmetricsRepo is a mock of metricsRepo interface.
type MockmetricsRepo struct {
	ctrl     *gomock.Controller
	recorder *MockmetricsRepoMockRecorder
}

// MockmetricsRepoMockRecorder is the mock recorder for MockmetricsRepo.
type MockmetricsRepoMockRecorder struct {
	mock *MockmetricsRepo
}

// NewMockmetricsRepo creates a new mock instance.
func NewMockmetricsRepo(ctrl *gomock.Controller) *MockmetricsRepo {
	mock := &MockmetricsRepo{ctrl: ctrl}
	mock.recorder = &MockmetricsRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockmetricsRepo) EXPECT() *MockmetricsRepoMockRecorder {
	return m.recorder
}

// Save mocks base method.
func (m *MockmetricsRepo) Save(ctx context.Context, metric models.ScalarMetric) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, metric)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockmetricsRepoMockRecorder) Save(ctx, metric any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockmetricsRepo)(nil).Save), ctx, metric)
}
