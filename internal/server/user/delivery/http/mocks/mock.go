// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package mock_http is a generated GoMock package.
package mock_http

import (
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUserHandlerInterface is a mock of UserHandlerInterface interface.
type MockUserHandlerInterface struct {
	ctrl     *gomock.Controller
	recorder *MockUserHandlerInterfaceMockRecorder
}

// MockUserHandlerInterfaceMockRecorder is the mock recorder for MockUserHandlerInterface.
type MockUserHandlerInterfaceMockRecorder struct {
	mock *MockUserHandlerInterface
}

// NewMockUserHandlerInterface creates a new mock instance.
func NewMockUserHandlerInterface(ctrl *gomock.Controller) *MockUserHandlerInterface {
	mock := &MockUserHandlerInterface{ctrl: ctrl}
	mock.recorder = &MockUserHandlerInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserHandlerInterface) EXPECT() *MockUserHandlerInterfaceMockRecorder {
	return m.recorder
}

// UpdateUser mocks base method.
func (m *MockUserHandlerInterface) UpdateUser(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateUser", w, r)
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockUserHandlerInterfaceMockRecorder) UpdateUser(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockUserHandlerInterface)(nil).UpdateUser), w, r)
}
