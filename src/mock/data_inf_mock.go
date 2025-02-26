// Code generated by MockGen. DO NOT EDIT.
// Source: mock_test.go

// Package main is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockDataInf is a mock of DataInf interface.
type MockDataInf struct {
	ctrl     *gomock.Controller
	recorder *MockDataInfMockRecorder
}

// MockDataInfMockRecorder is the mock recorder for MockDataInf.
type MockDataInfMockRecorder struct {
	mock *MockDataInf
}

// NewMockDataInf creates a new mock instance.
func NewMockDataInf(ctrl *gomock.Controller) *MockDataInf {
	mock := &MockDataInf{ctrl: ctrl}
	mock.recorder = &MockDataInfMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDataInf) EXPECT() *MockDataInfMockRecorder {
	return m.recorder
}

// GetData mocks base method.
func (m *MockDataInf) GetData(key string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetData", key)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetData indicates an expected call of GetData.
func (mr *MockDataInfMockRecorder) GetData(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetData", reflect.TypeOf((*MockDataInf)(nil).GetData), key)
}
