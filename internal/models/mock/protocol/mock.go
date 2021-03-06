// Code generated by MockGen. DO NOT EDIT.
// Source: protocol/repository.go

// Package protocol is a generated GoMock package.
package protocol

import (
	protocolModel "github.com/baking-bad/bcdhub/internal/models/protocol"
	types "github.com/baking-bad/bcdhub/internal/models/types"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockRepository is a mock of Repository interface
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockRepository) Get(network types.Network, hash string, level int64) (protocolModel.Protocol, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", network, hash, level)
	ret0, _ := ret[0].(protocolModel.Protocol)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockRepositoryMockRecorder) Get(network, hash, level interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRepository)(nil).Get), network, hash, level)
}

// GetAll mocks base method
func (m *MockRepository) GetAll() ([]protocolModel.Protocol, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].([]protocolModel.Protocol)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll
func (mr *MockRepositoryMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockRepository)(nil).GetAll))
}

// GetByNetworkWithSort mocks base method
func (m *MockRepository) GetByNetworkWithSort(network types.Network, sortField, order string) ([]protocolModel.Protocol, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByNetworkWithSort", network, sortField, order)
	ret0, _ := ret[0].([]protocolModel.Protocol)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByNetworkWithSort indicates an expected call of GetByNetworkWithSort
func (mr *MockRepositoryMockRecorder) GetByNetworkWithSort(network, sortField, order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByNetworkWithSort", reflect.TypeOf((*MockRepository)(nil).GetByNetworkWithSort), network, sortField, order)
}
