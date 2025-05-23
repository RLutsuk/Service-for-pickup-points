// Code generated by MockGen. DO NOT EDIT.
// Source: .\repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	reflect "reflect"

	models "github.com/RLutsuk/Service-for-pickup-points/app/models"
	gomock "github.com/golang/mock/gomock"
)

// MockRepositoryI is a mock of RepositoryI interface.
type MockRepositoryI struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryIMockRecorder
}

// MockRepositoryIMockRecorder is the mock recorder for MockRepositoryI.
type MockRepositoryIMockRecorder struct {
	mock *MockRepositoryI
}

// NewMockRepositoryI creates a new mock instance.
func NewMockRepositoryI(ctrl *gomock.Controller) *MockRepositoryI {
	mock := &MockRepositoryI{ctrl: ctrl}
	mock.recorder = &MockRepositoryIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepositoryI) EXPECT() *MockRepositoryIMockRecorder {
	return m.recorder
}

// CloseReception mocks base method.
func (m *MockRepositoryI) CloseReception(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseReception", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseReception indicates an expected call of CloseReception.
func (mr *MockRepositoryIMockRecorder) CloseReception(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseReception", reflect.TypeOf((*MockRepositoryI)(nil).CloseReception), id)
}

// CreateReception mocks base method.
func (m *MockRepositoryI) CreateReception(reception *models.Reception) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateReception", reception)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateReception indicates an expected call of CreateReception.
func (mr *MockRepositoryIMockRecorder) CreateReception(reception interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateReception", reflect.TypeOf((*MockRepositoryI)(nil).CreateReception), reception)
}

// GetOpenReceptionByPPID mocks base method.
func (m *MockRepositoryI) GetOpenReceptionByPPID(pickupPointID string) (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOpenReceptionByPPID", pickupPointID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetOpenReceptionByPPID indicates an expected call of GetOpenReceptionByPPID.
func (mr *MockRepositoryIMockRecorder) GetOpenReceptionByPPID(pickupPointID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOpenReceptionByPPID", reflect.TypeOf((*MockRepositoryI)(nil).GetOpenReceptionByPPID), pickupPointID)
}
