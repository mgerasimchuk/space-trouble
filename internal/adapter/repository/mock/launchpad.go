// Code generated by MockGen. DO NOT EDIT.
// Source: internal/domain/repository/launchpad.go

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	time "time"
)

// MockLaunchpadRepository is a mock of LaunchpadRepository interface
type MockLaunchpadRepository struct {
	ctrl     *gomock.Controller
	recorder *MockLaunchpadRepositoryMockRecorder
}

// MockLaunchpadRepositoryMockRecorder is the mock recorder for MockLaunchpadRepository
type MockLaunchpadRepositoryMockRecorder struct {
	mock *MockLaunchpadRepository
}

// NewMockLaunchpadRepository creates a new mock instance
func NewMockLaunchpadRepository(ctrl *gomock.Controller) *MockLaunchpadRepository {
	mock := &MockLaunchpadRepository{ctrl: ctrl}
	mock.recorder = &MockLaunchpadRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLaunchpadRepository) EXPECT() *MockLaunchpadRepositoryMockRecorder {
	return m.recorder
}

// IsExists mocks base method
func (m *MockLaunchpadRepository) IsExists(id string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsExists", id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsExists indicates an expected call of IsExists
func (mr *MockLaunchpadRepositoryMockRecorder) IsExists(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsExists", reflect.TypeOf((*MockLaunchpadRepository)(nil).IsExists), id)
}

// IsActive mocks base method
func (m *MockLaunchpadRepository) IsActive(id string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsActive", id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsActive indicates an expected call of IsActive
func (mr *MockLaunchpadRepositoryMockRecorder) IsActive(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsActive", reflect.TypeOf((*MockLaunchpadRepository)(nil).IsActive), id)
}

// IsDateAvailableForLaunch mocks base method
func (m *MockLaunchpadRepository) IsDateAvailableForLaunch(launchpadID string, launchDate time.Time) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsDateAvailableForLaunch", launchpadID, launchDate)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsDateAvailableForLaunch indicates an expected call of IsDateAvailableForLaunch
func (mr *MockLaunchpadRepositoryMockRecorder) IsDateAvailableForLaunch(launchpadID, launchDate interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsDateAvailableForLaunch", reflect.TypeOf((*MockLaunchpadRepository)(nil).IsDateAvailableForLaunch), launchpadID, launchDate)
}
