// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/quic-go/quic-go (interfaces: TokenStore)
//
// Generated by this command:
//
//	mockgen -typed -package quic -self_package github.com/quic-go/quic-go -self_package github.com/quic-go/quic-go -destination mock_token_store_test.go github.com/quic-go/quic-go TokenStore
//

// Package quic is a generated GoMock package.
package quic

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockTokenStore is a mock of TokenStore interface.
type MockTokenStore struct {
	ctrl     *gomock.Controller
	recorder *MockTokenStoreMockRecorder
}

// MockTokenStoreMockRecorder is the mock recorder for MockTokenStore.
type MockTokenStoreMockRecorder struct {
	mock *MockTokenStore
}

// NewMockTokenStore creates a new mock instance.
func NewMockTokenStore(ctrl *gomock.Controller) *MockTokenStore {
	mock := &MockTokenStore{ctrl: ctrl}
	mock.recorder = &MockTokenStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokenStore) EXPECT() *MockTokenStoreMockRecorder {
	return m.recorder
}

// Pop mocks base method.
func (m *MockTokenStore) Pop(arg0 string) *ClientToken {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Pop", arg0)
	ret0, _ := ret[0].(*ClientToken)
	return ret0
}

// Pop indicates an expected call of Pop.
func (mr *MockTokenStoreMockRecorder) Pop(arg0 any) *MockTokenStorePopCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Pop", reflect.TypeOf((*MockTokenStore)(nil).Pop), arg0)
	return &MockTokenStorePopCall{Call: call}
}

// MockTokenStorePopCall wrap *gomock.Call
type MockTokenStorePopCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockTokenStorePopCall) Return(arg0 *ClientToken) *MockTokenStorePopCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockTokenStorePopCall) Do(f func(string) *ClientToken) *MockTokenStorePopCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockTokenStorePopCall) DoAndReturn(f func(string) *ClientToken) *MockTokenStorePopCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Put mocks base method.
func (m *MockTokenStore) Put(arg0 string, arg1 *ClientToken) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Put", arg0, arg1)
}

// Put indicates an expected call of Put.
func (mr *MockTokenStoreMockRecorder) Put(arg0, arg1 any) *MockTokenStorePutCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockTokenStore)(nil).Put), arg0, arg1)
	return &MockTokenStorePutCall{Call: call}
}

// MockTokenStorePutCall wrap *gomock.Call
type MockTokenStorePutCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockTokenStorePutCall) Return() *MockTokenStorePutCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockTokenStorePutCall) Do(f func(string, *ClientToken)) *MockTokenStorePutCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockTokenStorePutCall) DoAndReturn(f func(string, *ClientToken)) *MockTokenStorePutCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
