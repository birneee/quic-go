// Code generated by MockGen. DO NOT EDIT.
// Source: connection.go

// Package quic is a generated GoMock package.
package quic

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	handover "github.com/lucas-clemente/quic-go/handover"
	protocol "github.com/lucas-clemente/quic-go/internal/protocol"
	wire "github.com/lucas-clemente/quic-go/internal/wire"
	xse "github.com/lucas-clemente/quic-go/internal/xse"
)

// MockStreamManager is a mock of StreamManager interface.
type MockStreamManager struct {
	ctrl     *gomock.Controller
	recorder *MockStreamManagerMockRecorder
}

// MockStreamManagerMockRecorder is the mock recorder for MockStreamManager.
type MockStreamManagerMockRecorder struct {
	mock *MockStreamManager
}

// NewMockStreamManager creates a new mock instance.
func NewMockStreamManager(ctrl *gomock.Controller) *MockStreamManager {
	mock := &MockStreamManager{ctrl: ctrl}
	mock.recorder = &MockStreamManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStreamManager) EXPECT() *MockStreamManagerMockRecorder {
	return m.recorder
}

// AcceptStream mocks base method.
func (m *MockStreamManager) AcceptStream(arg0 context.Context) (Stream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AcceptStream", arg0)
	ret0, _ := ret[0].(Stream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AcceptStream indicates an expected call of AcceptStream.
func (mr *MockStreamManagerMockRecorder) AcceptStream(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AcceptStream", reflect.TypeOf((*MockStreamManager)(nil).AcceptStream), arg0)
}

// AcceptUniStream mocks base method.
func (m *MockStreamManager) AcceptUniStream(arg0 context.Context) (ReceiveStream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AcceptUniStream", arg0)
	ret0, _ := ret[0].(ReceiveStream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AcceptUniStream indicates an expected call of AcceptUniStream.
func (mr *MockStreamManagerMockRecorder) AcceptUniStream(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AcceptUniStream", reflect.TypeOf((*MockStreamManager)(nil).AcceptUniStream), arg0)
}

// BidiStreamStates mocks base method.
func (m *MockStreamManager) BidiStreamStates(config *ConnectionStateStoreConf) map[StreamID]handover.BidiStreamState {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BidiStreamStates", config)
	ret0, _ := ret[0].(map[StreamID]handover.BidiStreamState)
	return ret0
}

// BidiStreamStates indicates an expected call of BidiStreamStates.
func (mr *MockStreamManagerMockRecorder) BidiStreamStates(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BidiStreamStates", reflect.TypeOf((*MockStreamManager)(nil).BidiStreamStates), config)
}

// CloseWithError mocks base method.
func (m *MockStreamManager) CloseWithError(arg0 error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CloseWithError", arg0)
}

// CloseWithError indicates an expected call of CloseWithError.
func (mr *MockStreamManagerMockRecorder) CloseWithError(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseWithError", reflect.TypeOf((*MockStreamManager)(nil).CloseWithError), arg0)
}

// DeleteStream mocks base method.
func (m *MockStreamManager) DeleteStream(arg0 protocol.StreamID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteStream", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteStream indicates an expected call of DeleteStream.
func (mr *MockStreamManagerMockRecorder) DeleteStream(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteStream", reflect.TypeOf((*MockStreamManager)(nil).DeleteStream), arg0)
}

// GetOrOpenReceiveStream mocks base method.
func (m *MockStreamManager) GetOrOpenReceiveStream(arg0 protocol.StreamID) (receiveStreamI, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrOpenReceiveStream", arg0)
	ret0, _ := ret[0].(receiveStreamI)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrOpenReceiveStream indicates an expected call of GetOrOpenReceiveStream.
func (mr *MockStreamManagerMockRecorder) GetOrOpenReceiveStream(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrOpenReceiveStream", reflect.TypeOf((*MockStreamManager)(nil).GetOrOpenReceiveStream), arg0)
}

// GetOrOpenSendStream mocks base method.
func (m *MockStreamManager) GetOrOpenSendStream(arg0 protocol.StreamID) (sendStreamI, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrOpenSendStream", arg0)
	ret0, _ := ret[0].(sendStreamI)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrOpenSendStream indicates an expected call of GetOrOpenSendStream.
func (mr *MockStreamManagerMockRecorder) GetOrOpenSendStream(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrOpenSendStream", reflect.TypeOf((*MockStreamManager)(nil).GetOrOpenSendStream), arg0)
}

// HandleMaxStreamsFrame mocks base method.
func (m *MockStreamManager) HandleMaxStreamsFrame(arg0 *wire.MaxStreamsFrame) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "HandleMaxStreamsFrame", arg0)
}

// HandleMaxStreamsFrame indicates an expected call of HandleMaxStreamsFrame.
func (mr *MockStreamManagerMockRecorder) HandleMaxStreamsFrame(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleMaxStreamsFrame", reflect.TypeOf((*MockStreamManager)(nil).HandleMaxStreamsFrame), arg0)
}

// OpenStream mocks base method.
func (m *MockStreamManager) OpenStream() (Stream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenStream")
	ret0, _ := ret[0].(Stream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenStream indicates an expected call of OpenStream.
func (mr *MockStreamManagerMockRecorder) OpenStream() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenStream", reflect.TypeOf((*MockStreamManager)(nil).OpenStream))
}

// OpenStreamSync mocks base method.
func (m *MockStreamManager) OpenStreamSync(arg0 context.Context) (Stream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenStreamSync", arg0)
	ret0, _ := ret[0].(Stream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenStreamSync indicates an expected call of OpenStreamSync.
func (mr *MockStreamManagerMockRecorder) OpenStreamSync(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenStreamSync", reflect.TypeOf((*MockStreamManager)(nil).OpenStreamSync), arg0)
}

// OpenUniStream mocks base method.
func (m *MockStreamManager) OpenUniStream() (SendStream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenUniStream")
	ret0, _ := ret[0].(SendStream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenUniStream indicates an expected call of OpenUniStream.
func (mr *MockStreamManagerMockRecorder) OpenUniStream() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenUniStream", reflect.TypeOf((*MockStreamManager)(nil).OpenUniStream))
}

// OpenUniStreamSync mocks base method.
func (m *MockStreamManager) OpenUniStreamSync(arg0 context.Context) (SendStream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenUniStreamSync", arg0)
	ret0, _ := ret[0].(SendStream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenUniStreamSync indicates an expected call of OpenUniStreamSync.
func (mr *MockStreamManagerMockRecorder) OpenUniStreamSync(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenUniStreamSync", reflect.TypeOf((*MockStreamManager)(nil).OpenUniStreamSync), arg0)
}

// OpenedBidiStream mocks base method.
func (m *MockStreamManager) OpenedBidiStream(id StreamID) (Stream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenedBidiStream", id)
	ret0, _ := ret[0].(Stream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenedBidiStream indicates an expected call of OpenedBidiStream.
func (mr *MockStreamManagerMockRecorder) OpenedBidiStream(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenedBidiStream", reflect.TypeOf((*MockStreamManager)(nil).OpenedBidiStream), id)
}

// ResetFor0RTT mocks base method.
func (m *MockStreamManager) ResetFor0RTT() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ResetFor0RTT")
}

// ResetFor0RTT indicates an expected call of ResetFor0RTT.
func (mr *MockStreamManagerMockRecorder) ResetFor0RTT() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetFor0RTT", reflect.TypeOf((*MockStreamManager)(nil).ResetFor0RTT))
}

// RestoreBidiStream mocks base method.
func (m *MockStreamManager) RestoreBidiStream(state *handover.BidiStreamState) (Stream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RestoreBidiStream", state)
	ret0, _ := ret[0].(Stream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RestoreBidiStream indicates an expected call of RestoreBidiStream.
func (mr *MockStreamManagerMockRecorder) RestoreBidiStream(state interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RestoreBidiStream", reflect.TypeOf((*MockStreamManager)(nil).RestoreBidiStream), state)
}

// RestoreReceiveStream mocks base method.
func (m *MockStreamManager) RestoreReceiveStream(state *handover.UniStreamState) (ReceiveStream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RestoreReceiveStream", state)
	ret0, _ := ret[0].(ReceiveStream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RestoreReceiveStream indicates an expected call of RestoreReceiveStream.
func (mr *MockStreamManagerMockRecorder) RestoreReceiveStream(state interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RestoreReceiveStream", reflect.TypeOf((*MockStreamManager)(nil).RestoreReceiveStream), state)
}

// RestoreSendStream mocks base method.
func (m *MockStreamManager) RestoreSendStream(state *handover.UniStreamState) (SendStream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RestoreSendStream", state)
	ret0, _ := ret[0].(SendStream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RestoreSendStream indicates an expected call of RestoreSendStream.
func (mr *MockStreamManagerMockRecorder) RestoreSendStream(state interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RestoreSendStream", reflect.TypeOf((*MockStreamManager)(nil).RestoreSendStream), state)
}

// SetXseCryptoSetup mocks base method.
func (m *MockStreamManager) SetXseCryptoSetup(arg0 xse.CryptoSetup) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetXseCryptoSetup", arg0)
}

// SetXseCryptoSetup indicates an expected call of SetXseCryptoSetup.
func (mr *MockStreamManagerMockRecorder) SetXseCryptoSetup(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetXseCryptoSetup", reflect.TypeOf((*MockStreamManager)(nil).SetXseCryptoSetup), arg0)
}

// UniStreamStates mocks base method.
func (m *MockStreamManager) UniStreamStates(config *ConnectionStateStoreConf) map[protocol.StreamID]handover.UniStreamState {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UniStreamStates", config)
	ret0, _ := ret[0].(map[protocol.StreamID]handover.UniStreamState)
	return ret0
}

// UniStreamStates indicates an expected call of UniStreamStates.
func (mr *MockStreamManagerMockRecorder) UniStreamStates(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UniStreamStates", reflect.TypeOf((*MockStreamManager)(nil).UniStreamStates), config)
}

// UpdateLimits mocks base method.
func (m *MockStreamManager) UpdateLimits(arg0 *wire.TransportParameters) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateLimits", arg0)
}

// UpdateLimits indicates an expected call of UpdateLimits.
func (mr *MockStreamManagerMockRecorder) UpdateLimits(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateLimits", reflect.TypeOf((*MockStreamManager)(nil).UpdateLimits), arg0)
}

// UseResetMaps mocks base method.
func (m *MockStreamManager) UseResetMaps() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UseResetMaps")
}

// UseResetMaps indicates an expected call of UseResetMaps.
func (mr *MockStreamManagerMockRecorder) UseResetMaps() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UseResetMaps", reflect.TypeOf((*MockStreamManager)(nil).UseResetMaps))
}
