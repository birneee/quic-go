// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/quic-go/quic-go/internal/ackhandler (interfaces: SentPacketHandler)
//
// Generated by this command:
//
//	mockgen -typed -build_flags=-tags=gomock -package mockackhandler -destination ackhandler/sent_packet_handler.go github.com/quic-go/quic-go/internal/ackhandler SentPacketHandler
//
// Package mockackhandler is a generated GoMock package.
package mockackhandler

import (
	reflect "reflect"
	time "time"

	handover "github.com/quic-go/quic-go/handover"
	ackhandler "github.com/quic-go/quic-go/internal/ackhandler"
	congestion "github.com/quic-go/quic-go/internal/congestion"
	protocol "github.com/quic-go/quic-go/internal/protocol"
	wire "github.com/quic-go/quic-go/internal/wire"
	qstate "github.com/quic-go/quic-go/qstate"
	gomock "go.uber.org/mock/gomock"
)

// MockSentPacketHandler is a mock of SentPacketHandler interface.
type MockSentPacketHandler struct {
	ctrl     *gomock.Controller
	recorder *MockSentPacketHandlerMockRecorder
}

// MockSentPacketHandlerMockRecorder is the mock recorder for MockSentPacketHandler.
type MockSentPacketHandlerMockRecorder struct {
	mock *MockSentPacketHandler
}

// NewMockSentPacketHandler creates a new mock instance.
func NewMockSentPacketHandler(ctrl *gomock.Controller) *MockSentPacketHandler {
	mock := &MockSentPacketHandler{ctrl: ctrl}
	mock.recorder = &MockSentPacketHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSentPacketHandler) EXPECT() *MockSentPacketHandlerMockRecorder {
	return m.recorder
}

// DropPackets mocks base method.
func (m *MockSentPacketHandler) DropPackets(arg0 protocol.EncryptionLevel) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DropPackets", arg0)
}

// DropPackets indicates an expected call of DropPackets.
func (mr *MockSentPacketHandlerMockRecorder) DropPackets(arg0 any) *SentPacketHandlerDropPacketsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DropPackets", reflect.TypeOf((*MockSentPacketHandler)(nil).DropPackets), arg0)
	return &SentPacketHandlerDropPacketsCall{Call: call}
}

// SentPacketHandlerDropPacketsCall wrap *gomock.Call
type SentPacketHandlerDropPacketsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *SentPacketHandlerDropPacketsCall) Return() *SentPacketHandlerDropPacketsCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *SentPacketHandlerDropPacketsCall) Do(f func(protocol.EncryptionLevel)) *SentPacketHandlerDropPacketsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *SentPacketHandlerDropPacketsCall) DoAndReturn(f func(protocol.EncryptionLevel)) *SentPacketHandlerDropPacketsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ECNMode mocks base method.
func (m *MockSentPacketHandler) ECNMode(arg0 bool) protocol.ECN {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ECNMode", arg0)
	ret0, _ := ret[0].(protocol.ECN)
	return ret0
}

// ECNMode indicates an expected call of ECNMode.
func (mr *MockSentPacketHandlerMockRecorder) ECNMode(arg0 any) *SentPacketHandlerECNModeCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ECNMode", reflect.TypeOf((*MockSentPacketHandler)(nil).ECNMode), arg0)
	return &SentPacketHandlerECNModeCall{Call: call}
}

// SentPacketHandlerECNModeCall wrap *gomock.Call
type SentPacketHandlerECNModeCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *SentPacketHandlerECNModeCall) Return(arg0 protocol.ECN) *SentPacketHandlerECNModeCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *SentPacketHandlerECNModeCall) Do(f func(bool) protocol.ECN) *SentPacketHandlerECNModeCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *SentPacketHandlerECNModeCall) DoAndReturn(f func(bool) protocol.ECN) *SentPacketHandlerECNModeCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetLossDetectionTimeout mocks base method.
func (m *MockSentPacketHandler) GetLossDetectionTimeout() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLossDetectionTimeout")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// GetLossDetectionTimeout indicates an expected call of GetLossDetectionTimeout.
func (mr *MockSentPacketHandlerMockRecorder) GetLossDetectionTimeout() *SentPacketHandlerGetLossDetectionTimeoutCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLossDetectionTimeout", reflect.TypeOf((*MockSentPacketHandler)(nil).GetLossDetectionTimeout))
	return &SentPacketHandlerGetLossDetectionTimeoutCall{Call: call}
}

// SentPacketHandlerGetLossDetectionTimeoutCall wrap *gomock.Call
type SentPacketHandlerGetLossDetectionTimeoutCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *SentPacketHandlerGetLossDetectionTimeoutCall) Return(arg0 time.Time) *SentPacketHandlerGetLossDetectionTimeoutCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *SentPacketHandlerGetLossDetectionTimeoutCall) Do(f func() time.Time) *SentPacketHandlerGetLossDetectionTimeoutCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *SentPacketHandlerGetLossDetectionTimeoutCall) DoAndReturn(f func() time.Time) *SentPacketHandlerGetLossDetectionTimeoutCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// OnLossDetectionTimeout mocks base method.
func (m *MockSentPacketHandler) OnLossDetectionTimeout() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OnLossDetectionTimeout")
	ret0, _ := ret[0].(error)
	return ret0
}

// OnLossDetectionTimeout indicates an expected call of OnLossDetectionTimeout.
func (mr *MockSentPacketHandlerMockRecorder) OnLossDetectionTimeout() *SentPacketHandlerOnLossDetectionTimeoutCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnLossDetectionTimeout", reflect.TypeOf((*MockSentPacketHandler)(nil).OnLossDetectionTimeout))
	return &SentPacketHandlerOnLossDetectionTimeoutCall{Call: call}
}

// SentPacketHandlerOnLossDetectionTimeoutCall wrap *gomock.Call
type SentPacketHandlerOnLossDetectionTimeoutCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *SentPacketHandlerOnLossDetectionTimeoutCall) Return(arg0 error) *SentPacketHandlerOnLossDetectionTimeoutCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *SentPacketHandlerOnLossDetectionTimeoutCall) Do(f func() error) *SentPacketHandlerOnLossDetectionTimeoutCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *SentPacketHandlerOnLossDetectionTimeoutCall) DoAndReturn(f func() error) *SentPacketHandlerOnLossDetectionTimeoutCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// PeekPacketNumber mocks base method.
func (m *MockSentPacketHandler) PeekPacketNumber(arg0 protocol.EncryptionLevel) (protocol.PacketNumber, protocol.PacketNumberLen) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PeekPacketNumber", arg0)
	ret0, _ := ret[0].(protocol.PacketNumber)
	ret1, _ := ret[1].(protocol.PacketNumberLen)
	return ret0, ret1
}

// PeekPacketNumber indicates an expected call of PeekPacketNumber.
func (mr *MockSentPacketHandlerMockRecorder) PeekPacketNumber(arg0 any) *SentPacketHandlerPeekPacketNumberCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PeekPacketNumber", reflect.TypeOf((*MockSentPacketHandler)(nil).PeekPacketNumber), arg0)
	return &SentPacketHandlerPeekPacketNumberCall{Call: call}
}

// SentPacketHandlerPeekPacketNumberCall wrap *gomock.Call
type SentPacketHandlerPeekPacketNumberCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *SentPacketHandlerPeekPacketNumberCall) Return(arg0 protocol.PacketNumber, arg1 protocol.PacketNumberLen) *SentPacketHandlerPeekPacketNumberCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *SentPacketHandlerPeekPacketNumberCall) Do(f func(protocol.EncryptionLevel) (protocol.PacketNumber, protocol.PacketNumberLen)) *SentPacketHandlerPeekPacketNumberCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *SentPacketHandlerPeekPacketNumberCall) DoAndReturn(f func(protocol.EncryptionLevel) (protocol.PacketNumber, protocol.PacketNumberLen)) *SentPacketHandlerPeekPacketNumberCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// PopPacketNumber mocks base method.
func (m *MockSentPacketHandler) PopPacketNumber(arg0 protocol.EncryptionLevel) protocol.PacketNumber {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PopPacketNumber", arg0)
	ret0, _ := ret[0].(protocol.PacketNumber)
	return ret0
}

// PopPacketNumber indicates an expected call of PopPacketNumber.
func (mr *MockSentPacketHandlerMockRecorder) PopPacketNumber(arg0 any) *SentPacketHandlerPopPacketNumberCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PopPacketNumber", reflect.TypeOf((*MockSentPacketHandler)(nil).PopPacketNumber), arg0)
	return &SentPacketHandlerPopPacketNumberCall{Call: call}
}

// SentPacketHandlerPopPacketNumberCall wrap *gomock.Call
type SentPacketHandlerPopPacketNumberCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *SentPacketHandlerPopPacketNumberCall) Return(arg0 protocol.PacketNumber) *SentPacketHandlerPopPacketNumberCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *SentPacketHandlerPopPacketNumberCall) Do(f func(protocol.EncryptionLevel) protocol.PacketNumber) *SentPacketHandlerPopPacketNumberCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *SentPacketHandlerPopPacketNumberCall) DoAndReturn(f func(protocol.EncryptionLevel) protocol.PacketNumber) *SentPacketHandlerPopPacketNumberCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// QueueProbePacket mocks base method.
func (m *MockSentPacketHandler) QueueProbePacket(arg0 protocol.EncryptionLevel) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueueProbePacket", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// QueueProbePacket indicates an expected call of QueueProbePacket.
func (mr *MockSentPacketHandlerMockRecorder) QueueProbePacket(arg0 any) *SentPacketHandlerQueueProbePacketCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueueProbePacket", reflect.TypeOf((*MockSentPacketHandler)(nil).QueueProbePacket), arg0)
	return &SentPacketHandlerQueueProbePacketCall{Call: call}
}

// SentPacketHandlerQueueProbePacketCall wrap *gomock.Call
type SentPacketHandlerQueueProbePacketCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *SentPacketHandlerQueueProbePacketCall) Return(arg0 bool) *SentPacketHandlerQueueProbePacketCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *SentPacketHandlerQueueProbePacketCall) Do(f func(protocol.EncryptionLevel) bool) *SentPacketHandlerQueueProbePacketCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *SentPacketHandlerQueueProbePacketCall) DoAndReturn(f func(protocol.EncryptionLevel) bool) *SentPacketHandlerQueueProbePacketCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ReceivedAck mocks base method.
func (m *MockSentPacketHandler) ReceivedAck(arg0 *wire.AckFrame, arg1 protocol.EncryptionLevel, arg2 time.Time) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReceivedAck", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReceivedAck indicates an expected call of ReceivedAck.
func (mr *MockSentPacketHandlerMockRecorder) ReceivedAck(arg0, arg1, arg2 any) *SentPacketHandlerReceivedAckCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReceivedAck", reflect.TypeOf((*MockSentPacketHandler)(nil).ReceivedAck), arg0, arg1, arg2)
	return &SentPacketHandlerReceivedAckCall{Call: call}
}

// SentPacketHandlerReceivedAckCall wrap *gomock.Call
type SentPacketHandlerReceivedAckCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *SentPacketHandlerReceivedAckCall) Return(arg0 bool, arg1 error) *SentPacketHandlerReceivedAckCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *SentPacketHandlerReceivedAckCall) Do(f func(*wire.AckFrame, protocol.EncryptionLevel, time.Time) (bool, error)) *SentPacketHandlerReceivedAckCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *SentPacketHandlerReceivedAckCall) DoAndReturn(f func(*wire.AckFrame, protocol.EncryptionLevel, time.Time) (bool, error)) *SentPacketHandlerReceivedAckCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ReceivedBytes mocks base method.
func (m *MockSentPacketHandler) ReceivedBytes(arg0 protocol.ByteCount) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ReceivedBytes", arg0)
}

// ReceivedBytes indicates an expected call of ReceivedBytes.
func (mr *MockSentPacketHandlerMockRecorder) ReceivedBytes(arg0 any) *SentPacketHandlerReceivedBytesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReceivedBytes", reflect.TypeOf((*MockSentPacketHandler)(nil).ReceivedBytes), arg0)
	return &SentPacketHandlerReceivedBytesCall{Call: call}
}

// SentPacketHandlerReceivedBytesCall wrap *gomock.Call
type SentPacketHandlerReceivedBytesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *SentPacketHandlerReceivedBytesCall) Return() *SentPacketHandlerReceivedBytesCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *SentPacketHandlerReceivedBytesCall) Do(f func(protocol.ByteCount)) *SentPacketHandlerReceivedBytesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *SentPacketHandlerReceivedBytesCall) DoAndReturn(f func(protocol.ByteCount)) *SentPacketHandlerReceivedBytesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ResetForRetry mocks base method.
func (m *MockSentPacketHandler) ResetForRetry(arg0 time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetForRetry", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ResetForRetry indicates an expected call of ResetForRetry.
func (mr *MockSentPacketHandlerMockRecorder) ResetForRetry(arg0 any) *SentPacketHandlerResetForRetryCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetForRetry", reflect.TypeOf((*MockSentPacketHandler)(nil).ResetForRetry), arg0)
	return &SentPacketHandlerResetForRetryCall{Call: call}
}

// SentPacketHandlerResetForRetryCall wrap *gomock.Call
type SentPacketHandlerResetForRetryCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *SentPacketHandlerResetForRetryCall) Return(arg0 error) *SentPacketHandlerResetForRetryCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *SentPacketHandlerResetForRetryCall) Do(f func(time.Time) error) *SentPacketHandlerResetForRetryCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *SentPacketHandlerResetForRetryCall) DoAndReturn(f func(time.Time) error) *SentPacketHandlerResetForRetryCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SendMode mocks base method.
func (m *MockSentPacketHandler) SendMode(arg0 time.Time) ackhandler.SendMode {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMode", arg0)
	ret0, _ := ret[0].(ackhandler.SendMode)
	return ret0
}

// SendMode indicates an expected call of SendMode.
func (mr *MockSentPacketHandlerMockRecorder) SendMode(arg0 any) *SentPacketHandlerSendModeCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMode", reflect.TypeOf((*MockSentPacketHandler)(nil).SendMode), arg0)
	return &SentPacketHandlerSendModeCall{Call: call}
}

// SentPacketHandlerSendModeCall wrap *gomock.Call
type SentPacketHandlerSendModeCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *SentPacketHandlerSendModeCall) Return(arg0 ackhandler.SendMode) *SentPacketHandlerSendModeCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *SentPacketHandlerSendModeCall) Do(f func(time.Time) ackhandler.SendMode) *SentPacketHandlerSendModeCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *SentPacketHandlerSendModeCall) DoAndReturn(f func(time.Time) ackhandler.SendMode) *SentPacketHandlerSendModeCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SentPacket mocks base method.
func (m *MockSentPacketHandler) SentPacket(arg0 time.Time, arg1, arg2 protocol.PacketNumber, arg3 []ackhandler.StreamFrame, arg4 []ackhandler.Frame, arg5 protocol.EncryptionLevel, arg6 protocol.ECN, arg7 protocol.ByteCount, arg8 bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SentPacket", arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8)
}

// SentPacket indicates an expected call of SentPacket.
func (mr *MockSentPacketHandlerMockRecorder) SentPacket(arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8 any) *SentPacketHandlerSentPacketCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SentPacket", reflect.TypeOf((*MockSentPacketHandler)(nil).SentPacket), arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8)
	return &SentPacketHandlerSentPacketCall{Call: call}
}

// SentPacketHandlerSentPacketCall wrap *gomock.Call
type SentPacketHandlerSentPacketCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *SentPacketHandlerSentPacketCall) Return() *SentPacketHandlerSentPacketCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *SentPacketHandlerSentPacketCall) Do(f func(time.Time, protocol.PacketNumber, protocol.PacketNumber, []ackhandler.StreamFrame, []ackhandler.Frame, protocol.EncryptionLevel, protocol.ECN, protocol.ByteCount, bool)) *SentPacketHandlerSentPacketCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *SentPacketHandlerSentPacketCall) DoAndReturn(f func(time.Time, protocol.PacketNumber, protocol.PacketNumber, []ackhandler.StreamFrame, []ackhandler.Frame, protocol.EncryptionLevel, protocol.ECN, protocol.ByteCount, bool)) *SentPacketHandlerSentPacketCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SetHandshakeConfirmed mocks base method.
func (m *MockSentPacketHandler) SetHandshakeConfirmed() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetHandshakeConfirmed")
}

// SetHandshakeConfirmed indicates an expected call of SetHandshakeConfirmed.
func (mr *MockSentPacketHandlerMockRecorder) SetHandshakeConfirmed() *SentPacketHandlerSetHandshakeConfirmedCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHandshakeConfirmed", reflect.TypeOf((*MockSentPacketHandler)(nil).SetHandshakeConfirmed))
	return &SentPacketHandlerSetHandshakeConfirmedCall{Call: call}
}

// SentPacketHandlerSetHandshakeConfirmedCall wrap *gomock.Call
type SentPacketHandlerSetHandshakeConfirmedCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *SentPacketHandlerSetHandshakeConfirmedCall) Return() *SentPacketHandlerSetHandshakeConfirmedCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *SentPacketHandlerSetHandshakeConfirmedCall) Do(f func()) *SentPacketHandlerSetHandshakeConfirmedCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *SentPacketHandlerSetHandshakeConfirmedCall) DoAndReturn(f func()) *SentPacketHandlerSetHandshakeConfirmedCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SetInitialCongestionWindow mocks base method.
func (m *MockSentPacketHandler) SetInitialCongestionWindow(arg0 uint32) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetInitialCongestionWindow", arg0)
}

// SetInitialCongestionWindow indicates an expected call of SetInitialCongestionWindow.
func (mr *MockSentPacketHandlerMockRecorder) SetInitialCongestionWindow(arg0 any) *SentPacketHandlerSetInitialCongestionWindowCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetInitialCongestionWindow", reflect.TypeOf((*MockSentPacketHandler)(nil).SetInitialCongestionWindow), arg0)
	return &SentPacketHandlerSetInitialCongestionWindowCall{Call: call}
}

// SentPacketHandlerSetInitialCongestionWindowCall wrap *gomock.Call
type SentPacketHandlerSetInitialCongestionWindowCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *SentPacketHandlerSetInitialCongestionWindowCall) Return() *SentPacketHandlerSetInitialCongestionWindowCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *SentPacketHandlerSetInitialCongestionWindowCall) Do(f func(uint32)) *SentPacketHandlerSetInitialCongestionWindowCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *SentPacketHandlerSetInitialCongestionWindowCall) DoAndReturn(f func(uint32)) *SentPacketHandlerSetInitialCongestionWindowCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SetMaxBandwidth mocks base method.
func (m *MockSentPacketHandler) SetMaxBandwidth(arg0 congestion.Bandwidth) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetMaxBandwidth", arg0)
}

// SetMaxBandwidth indicates an expected call of SetMaxBandwidth.
func (mr *MockSentPacketHandlerMockRecorder) SetMaxBandwidth(arg0 any) *SentPacketHandlerSetMaxBandwidthCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMaxBandwidth", reflect.TypeOf((*MockSentPacketHandler)(nil).SetMaxBandwidth), arg0)
	return &SentPacketHandlerSetMaxBandwidthCall{Call: call}
}

// SentPacketHandlerSetMaxBandwidthCall wrap *gomock.Call
type SentPacketHandlerSetMaxBandwidthCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *SentPacketHandlerSetMaxBandwidthCall) Return() *SentPacketHandlerSetMaxBandwidthCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *SentPacketHandlerSetMaxBandwidthCall) Do(f func(congestion.Bandwidth)) *SentPacketHandlerSetMaxBandwidthCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *SentPacketHandlerSetMaxBandwidthCall) DoAndReturn(f func(congestion.Bandwidth)) *SentPacketHandlerSetMaxBandwidthCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SetMaxDatagramSize mocks base method.
func (m *MockSentPacketHandler) SetMaxDatagramSize(arg0 protocol.ByteCount) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetMaxDatagramSize", arg0)
}

// SetMaxDatagramSize indicates an expected call of SetMaxDatagramSize.
func (mr *MockSentPacketHandlerMockRecorder) SetMaxDatagramSize(arg0 any) *SentPacketHandlerSetMaxDatagramSizeCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMaxDatagramSize", reflect.TypeOf((*MockSentPacketHandler)(nil).SetMaxDatagramSize), arg0)
	return &SentPacketHandlerSetMaxDatagramSizeCall{Call: call}
}

// SentPacketHandlerSetMaxDatagramSizeCall wrap *gomock.Call
type SentPacketHandlerSetMaxDatagramSizeCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *SentPacketHandlerSetMaxDatagramSizeCall) Return() *SentPacketHandlerSetMaxDatagramSizeCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *SentPacketHandlerSetMaxDatagramSizeCall) Do(f func(protocol.ByteCount)) *SentPacketHandlerSetMaxDatagramSizeCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *SentPacketHandlerSetMaxDatagramSizeCall) DoAndReturn(f func(protocol.ByteCount)) *SentPacketHandlerSetMaxDatagramSizeCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// StoreState mocks base method.
func (m *MockSentPacketHandler) StoreState(arg0 *qstate.Connection, arg1 *handover.ConnectionStateStoreConf) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StoreState", arg0, arg1)
}

// StoreState indicates an expected call of StoreState.
func (mr *MockSentPacketHandlerMockRecorder) StoreState(arg0, arg1 any) *SentPacketHandlerStoreStateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreState", reflect.TypeOf((*MockSentPacketHandler)(nil).StoreState), arg0, arg1)
	return &SentPacketHandlerStoreStateCall{Call: call}
}

// SentPacketHandlerStoreStateCall wrap *gomock.Call
type SentPacketHandlerStoreStateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *SentPacketHandlerStoreStateCall) Return() *SentPacketHandlerStoreStateCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *SentPacketHandlerStoreStateCall) Do(f func(*qstate.Connection, *handover.ConnectionStateStoreConf)) *SentPacketHandlerStoreStateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *SentPacketHandlerStoreStateCall) DoAndReturn(f func(*qstate.Connection, *handover.ConnectionStateStoreConf)) *SentPacketHandlerStoreStateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// StreamFramesInFlight mocks base method.
func (m *MockSentPacketHandler) StreamFramesInFlight(arg0 protocol.StreamID, arg1 protocol.EncryptionLevel) []*wire.StreamFrame {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StreamFramesInFlight", arg0, arg1)
	ret0, _ := ret[0].([]*wire.StreamFrame)
	return ret0
}

// StreamFramesInFlight indicates an expected call of StreamFramesInFlight.
func (mr *MockSentPacketHandlerMockRecorder) StreamFramesInFlight(arg0, arg1 any) *SentPacketHandlerStreamFramesInFlightCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StreamFramesInFlight", reflect.TypeOf((*MockSentPacketHandler)(nil).StreamFramesInFlight), arg0, arg1)
	return &SentPacketHandlerStreamFramesInFlightCall{Call: call}
}

// SentPacketHandlerStreamFramesInFlightCall wrap *gomock.Call
type SentPacketHandlerStreamFramesInFlightCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *SentPacketHandlerStreamFramesInFlightCall) Return(arg0 []*wire.StreamFrame) *SentPacketHandlerStreamFramesInFlightCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *SentPacketHandlerStreamFramesInFlightCall) Do(f func(protocol.StreamID, protocol.EncryptionLevel) []*wire.StreamFrame) *SentPacketHandlerStreamFramesInFlightCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *SentPacketHandlerStreamFramesInFlightCall) DoAndReturn(f func(protocol.StreamID, protocol.EncryptionLevel) []*wire.StreamFrame) *SentPacketHandlerStreamFramesInFlightCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// TimeUntilSend mocks base method.
func (m *MockSentPacketHandler) TimeUntilSend() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TimeUntilSend")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// TimeUntilSend indicates an expected call of TimeUntilSend.
func (mr *MockSentPacketHandlerMockRecorder) TimeUntilSend() *SentPacketHandlerTimeUntilSendCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TimeUntilSend", reflect.TypeOf((*MockSentPacketHandler)(nil).TimeUntilSend))
	return &SentPacketHandlerTimeUntilSendCall{Call: call}
}

// SentPacketHandlerTimeUntilSendCall wrap *gomock.Call
type SentPacketHandlerTimeUntilSendCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *SentPacketHandlerTimeUntilSendCall) Return(arg0 time.Time) *SentPacketHandlerTimeUntilSendCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *SentPacketHandlerTimeUntilSendCall) Do(f func() time.Time) *SentPacketHandlerTimeUntilSendCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *SentPacketHandlerTimeUntilSendCall) DoAndReturn(f func() time.Time) *SentPacketHandlerTimeUntilSendCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
