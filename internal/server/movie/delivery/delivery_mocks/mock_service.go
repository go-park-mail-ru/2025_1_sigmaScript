// Code generated by MockGen. DO NOT EDIT.
// Source: movie_handler.go

// Package delivery_mocks is a generated GoMock package.
package delivery_mocks

import (
	context "context"
	reflect "reflect"

	mocks "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	gomock "github.com/golang/mock/gomock"
)

// MockMovieServiceInterface is a mock of MovieServiceInterface interface.
type MockMovieServiceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockMovieServiceInterfaceMockRecorder
}

// MockMovieServiceInterfaceMockRecorder is the mock recorder for MockMovieServiceInterface.
type MockMovieServiceInterfaceMockRecorder struct {
	mock *MockMovieServiceInterface
}

// NewMockMovieServiceInterface creates a new mock instance.
func NewMockMovieServiceInterface(ctrl *gomock.Controller) *MockMovieServiceInterface {
	mock := &MockMovieServiceInterface{ctrl: ctrl}
	mock.recorder = &MockMovieServiceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMovieServiceInterface) EXPECT() *MockMovieServiceInterfaceMockRecorder {
	return m.recorder
}

// GetMovieByID mocks base method.
func (m *MockMovieServiceInterface) GetMovieByID(ctx context.Context, movieID int) (*mocks.MovieJSON, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMovieByID", ctx, movieID)
	ret0, _ := ret[0].(*mocks.MovieJSON)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMovieByID indicates an expected call of GetMovieByID.
func (mr *MockMovieServiceInterfaceMockRecorder) GetMovieByID(ctx, movieID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMovieByID", reflect.TypeOf((*MockMovieServiceInterface)(nil).GetMovieByID), ctx, movieID)
}
