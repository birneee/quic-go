// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/lucas-clemente/quic-go (interfaces: EarlyConnection)

// Package mockquic is a generated GoMock package.
package mockquic

import (
	context "context"
	net "net"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	quic "github.com/lucas-clemente/quic-go"
	protocol "github.com/lucas-clemente/quic-go/internal/protocol"
	qerr "github.com/lucas-clemente/quic-go/internal/qerr"
	logging "github.com/lucas-clemente/quic-go/logging"
)

// MockEarlyConnection is a mock of EarlyConnection interface.
type MockEarlyConnection struct {
	ctrl     *gomock.Controller
	recorder *MockEarlyConnectionMockRecorder
}

// MockEarlyConnectionMockRecorder is the mock recorder for MockEarlyConnection.
type MockEarlyConnectionMockRecorder struct {
	mock *MockEarlyConnection
}

// NewMockEarlyConnection creates a new mock instance.
func NewMockEarlyConnection(ctrl *gomock.Controller) *MockEarlyConnection {
	mock := &MockEarlyConnection{ctrl: ctrl}
	mock.recorder = &MockEarlyConnectionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEarlyConnection) EXPECT() *MockEarlyConnectionMockRecorder {
	return m.recorder
}

// AcceptStream mocks base method.
func (m *MockEarlyConnection) AcceptStream(arg0 context.Context) (quic.Stream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AcceptStream", arg0)
	ret0, _ := ret[0].(quic.Stream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AcceptStream indicates an expected call of AcceptStream.
func (mr *MockEarlyConnectionMockRecorder) AcceptStream(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AcceptStream", reflect.TypeOf((*MockEarlyConnection)(nil).AcceptStream), arg0)
}

// AcceptUniStream mocks base method.
func (m *MockEarlyConnection) AcceptUniStream(arg0 context.Context) (quic.ReceiveStream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AcceptUniStream", arg0)
	ret0, _ := ret[0].(quic.ReceiveStream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AcceptUniStream indicates an expected call of AcceptUniStream.
func (mr *MockEarlyConnectionMockRecorder) AcceptUniStream(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AcceptUniStream", reflect.TypeOf((*MockEarlyConnection)(nil).AcceptUniStream), arg0)
}

// AddProxy mocks base method.
func (m *MockEarlyConnection) AddProxy(arg0 *quic.ProxyConfig) quic.ProxySetupResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddProxy", arg0)
	ret0, _ := ret[0].(quic.ProxySetupResponse)
	return ret0
}

// AddProxy indicates an expected call of AddProxy.
func (mr *MockEarlyConnectionMockRecorder) AddProxy(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddProxy", reflect.TypeOf((*MockEarlyConnection)(nil).AddProxy), arg0)
}

// AwaitPathUpdate mocks base method.
func (m *MockEarlyConnection) AwaitPathUpdate() <-chan struct{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AwaitPathUpdate")
	ret0, _ := ret[0].(<-chan struct{})
	return ret0
}

// AwaitPathUpdate indicates an expected call of AwaitPathUpdate.
func (mr *MockEarlyConnectionMockRecorder) AwaitPathUpdate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AwaitPathUpdate", reflect.TypeOf((*MockEarlyConnection)(nil).AwaitPathUpdate))
}

// CloseWithError mocks base method.
func (m *MockEarlyConnection) CloseWithError(arg0 qerr.ApplicationErrorCode, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseWithError", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseWithError indicates an expected call of CloseWithError.
func (mr *MockEarlyConnectionMockRecorder) CloseWithError(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseWithError", reflect.TypeOf((*MockEarlyConnection)(nil).CloseWithError), arg0, arg1)
}

// ConnectionState mocks base method.
func (m *MockEarlyConnection) ConnectionState() quic.ConnectionState {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConnectionState")
	ret0, _ := ret[0].(quic.ConnectionState)
	return ret0
}

// ConnectionState indicates an expected call of ConnectionState.
func (mr *MockEarlyConnectionMockRecorder) ConnectionState() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConnectionState", reflect.TypeOf((*MockEarlyConnection)(nil).ConnectionState))
}

// Context mocks base method.
func (m *MockEarlyConnection) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockEarlyConnectionMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockEarlyConnection)(nil).Context))
}

// ExtraStreamEncrypted mocks base method.
func (m *MockEarlyConnection) ExtraStreamEncrypted() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExtraStreamEncrypted")
	ret0, _ := ret[0].(bool)
	return ret0
}

// ExtraStreamEncrypted indicates an expected call of ExtraStreamEncrypted.
func (mr *MockEarlyConnectionMockRecorder) ExtraStreamEncrypted() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExtraStreamEncrypted", reflect.TypeOf((*MockEarlyConnection)(nil).ExtraStreamEncrypted))
}

// Handover mocks base method.
func (m *MockEarlyConnection) Handover(arg0 bool, arg1 *quic.ConnectionStateStoreConf) quic.HandoverStateResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handover", arg0, arg1)
	ret0, _ := ret[0].(quic.HandoverStateResponse)
	return ret0
}

// Handover indicates an expected call of Handover.
func (mr *MockEarlyConnectionMockRecorder) Handover(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handover", reflect.TypeOf((*MockEarlyConnection)(nil).Handover), arg0, arg1)
}

// HandshakeComplete mocks base method.
func (m *MockEarlyConnection) HandshakeComplete() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandshakeComplete")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// HandshakeComplete indicates an expected call of HandshakeComplete.
func (mr *MockEarlyConnectionMockRecorder) HandshakeComplete() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandshakeComplete", reflect.TypeOf((*MockEarlyConnection)(nil).HandshakeComplete))
}

// LocalAddr mocks base method.
func (m *MockEarlyConnection) LocalAddr() net.Addr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LocalAddr")
	ret0, _ := ret[0].(net.Addr)
	return ret0
}

// LocalAddr indicates an expected call of LocalAddr.
func (mr *MockEarlyConnectionMockRecorder) LocalAddr() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LocalAddr", reflect.TypeOf((*MockEarlyConnection)(nil).LocalAddr))
}

// MigrateUDPSocket mocks base method.
func (m *MockEarlyConnection) MigrateUDPSocket() (*net.UDPAddr, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MigrateUDPSocket")
	ret0, _ := ret[0].(*net.UDPAddr)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MigrateUDPSocket indicates an expected call of MigrateUDPSocket.
func (mr *MockEarlyConnectionMockRecorder) MigrateUDPSocket() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MigrateUDPSocket", reflect.TypeOf((*MockEarlyConnection)(nil).MigrateUDPSocket))
}

// NextConnection mocks base method.
func (m *MockEarlyConnection) NextConnection() quic.Connection {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NextConnection")
	ret0, _ := ret[0].(quic.Connection)
	return ret0
}

// NextConnection indicates an expected call of NextConnection.
func (mr *MockEarlyConnectionMockRecorder) NextConnection() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NextConnection", reflect.TypeOf((*MockEarlyConnection)(nil).NextConnection))
}

// OpenStream mocks base method.
func (m *MockEarlyConnection) OpenStream() (quic.Stream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenStream")
	ret0, _ := ret[0].(quic.Stream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenStream indicates an expected call of OpenStream.
func (mr *MockEarlyConnectionMockRecorder) OpenStream() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenStream", reflect.TypeOf((*MockEarlyConnection)(nil).OpenStream))
}

// OpenStreamSync mocks base method.
func (m *MockEarlyConnection) OpenStreamSync(arg0 context.Context) (quic.Stream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenStreamSync", arg0)
	ret0, _ := ret[0].(quic.Stream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenStreamSync indicates an expected call of OpenStreamSync.
func (mr *MockEarlyConnectionMockRecorder) OpenStreamSync(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenStreamSync", reflect.TypeOf((*MockEarlyConnection)(nil).OpenStreamSync), arg0)
}

// OpenUniStream mocks base method.
func (m *MockEarlyConnection) OpenUniStream() (quic.SendStream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenUniStream")
	ret0, _ := ret[0].(quic.SendStream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenUniStream indicates an expected call of OpenUniStream.
func (mr *MockEarlyConnectionMockRecorder) OpenUniStream() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenUniStream", reflect.TypeOf((*MockEarlyConnection)(nil).OpenUniStream))
}

// OpenUniStreamSync mocks base method.
func (m *MockEarlyConnection) OpenUniStreamSync(arg0 context.Context) (quic.SendStream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenUniStreamSync", arg0)
	ret0, _ := ret[0].(quic.SendStream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenUniStreamSync indicates an expected call of OpenUniStreamSync.
func (mr *MockEarlyConnectionMockRecorder) OpenUniStreamSync(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenUniStreamSync", reflect.TypeOf((*MockEarlyConnection)(nil).OpenUniStreamSync), arg0)
}

// OpenedBidiStream mocks base method.
func (m *MockEarlyConnection) OpenedBidiStream(arg0 protocol.StreamID) (quic.Stream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenedBidiStream", arg0)
	ret0, _ := ret[0].(quic.Stream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenedBidiStream indicates an expected call of OpenedBidiStream.
func (mr *MockEarlyConnectionMockRecorder) OpenedBidiStream(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenedBidiStream", reflect.TypeOf((*MockEarlyConnection)(nil).OpenedBidiStream), arg0)
}

// OriginalDestinationConnectionID mocks base method.
func (m *MockEarlyConnection) OriginalDestinationConnectionID() protocol.ConnectionID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OriginalDestinationConnectionID")
	ret0, _ := ret[0].(protocol.ConnectionID)
	return ret0
}

// OriginalDestinationConnectionID indicates an expected call of OriginalDestinationConnectionID.
func (mr *MockEarlyConnectionMockRecorder) OriginalDestinationConnectionID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OriginalDestinationConnectionID", reflect.TypeOf((*MockEarlyConnection)(nil).OriginalDestinationConnectionID))
}

// QlogWriter mocks base method.
func (m *MockEarlyConnection) QlogWriter() logging.QlogWriter {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QlogWriter")
	ret0, _ := ret[0].(logging.QlogWriter)
	return ret0
}

// QlogWriter indicates an expected call of QlogWriter.
func (mr *MockEarlyConnectionMockRecorder) QlogWriter() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QlogWriter", reflect.TypeOf((*MockEarlyConnection)(nil).QlogWriter))
}

// QueueHandshakeDoneFrame mocks base method.
func (m *MockEarlyConnection) QueueHandshakeDoneFrame() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueueHandshakeDoneFrame")
	ret0, _ := ret[0].(error)
	return ret0
}

// QueueHandshakeDoneFrame indicates an expected call of QueueHandshakeDoneFrame.
func (mr *MockEarlyConnectionMockRecorder) QueueHandshakeDoneFrame() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueueHandshakeDoneFrame", reflect.TypeOf((*MockEarlyConnection)(nil).QueueHandshakeDoneFrame))
}

// ReceiveMessage mocks base method.
func (m *MockEarlyConnection) ReceiveMessage() ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReceiveMessage")
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReceiveMessage indicates an expected call of ReceiveMessage.
func (mr *MockEarlyConnectionMockRecorder) ReceiveMessage() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReceiveMessage", reflect.TypeOf((*MockEarlyConnection)(nil).ReceiveMessage))
}

// RemoteAddr mocks base method.
func (m *MockEarlyConnection) RemoteAddr() net.Addr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoteAddr")
	ret0, _ := ret[0].(net.Addr)
	return ret0
}

// RemoteAddr indicates an expected call of RemoteAddr.
func (mr *MockEarlyConnectionMockRecorder) RemoteAddr() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoteAddr", reflect.TypeOf((*MockEarlyConnection)(nil).RemoteAddr))
}

// SendMessage mocks base method.
func (m *MockEarlyConnection) SendMessage(arg0 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessage", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMessage indicates an expected call of SendMessage.
func (mr *MockEarlyConnectionMockRecorder) SendMessage(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MockEarlyConnection)(nil).SendMessage), arg0)
}

// UpdateRemoteAddr mocks base method.
func (m *MockEarlyConnection) UpdateRemoteAddr(arg0 net.UDPAddr, arg1, arg2 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateRemoteAddr", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateRemoteAddr indicates an expected call of UpdateRemoteAddr.
func (mr *MockEarlyConnectionMockRecorder) UpdateRemoteAddr(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRemoteAddr", reflect.TypeOf((*MockEarlyConnection)(nil).UpdateRemoteAddr), arg0, arg1, arg2)
}
