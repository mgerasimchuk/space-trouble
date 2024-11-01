// Code generated by MockGen. DO NOT EDIT.
// Source: landpad.go
//
// Generated by this command:
//
//	mockgen -source=landpad.go -destination=mock/landpad.go -package=mock
//
// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockLandpadRepository is a mock of LandpadRepository interface.
type MockLandpadRepository struct {
	ctrl     *gomock.Controller
	recorder *MockLandpadRepositoryMockRecorder
}

// MockLandpadRepositoryMockRecorder is the mock recorder for MockLandpadRepository.
type MockLandpadRepositoryMockRecorder struct {
	mock *MockLandpadRepository
}

// NewMockLandpadRepository creates a new mock instance.
func NewMockLandpadRepository(ctrl *gomock.Controller) *MockLandpadRepository {
	mock := &MockLandpadRepository{ctrl: ctrl}
	mock.recorder = &MockLandpadRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLandpadRepository) EXPECT() *MockLandpadRepositoryMockRecorder {
	return m.recorder
}

// IsActive mocks base method.
func (m *MockLandpadRepository) IsActive(id string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsActive", id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsActive indicates an expected call of IsActive.
func (mr *MockLandpadRepositoryMockRecorder) IsActive(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsActive", reflect.TypeOf((*MockLandpadRepository)(nil).IsActive), id)
}

// IsExists mocks base method.
func (m *MockLandpadRepository) IsExists(id string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsExists", id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsExists indicates an expected call of IsExists.
func (mr *MockLandpadRepositoryMockRecorder) IsExists(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsExists", reflect.TypeOf((*MockLandpadRepository)(nil).IsExists), id)
}
