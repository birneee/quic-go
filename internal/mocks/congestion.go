// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/quic-go/quic-go/internal/congestion (interfaces: SendAlgorithmWithDebugInfos)
//
// Generated by this command:
//
//	mockgen -typed -build_flags=-tags=gomock -package mocks -destination congestion.go github.com/quic-go/quic-go/internal/congestion SendAlgorithmWithDebugInfos
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	time "time"

	congestion "github.com/quic-go/quic-go/internal/congestion"
	protocol "github.com/quic-go/quic-go/internal/protocol"
	gomock "go.uber.org/mock/gomock"
)

// MockSendAlgorithmWithDebugInfos is a mock of SendAlgorithmWithDebugInfos interface.
type MockSendAlgorithmWithDebugInfos struct {
	ctrl     *gomock.Controller
	recorder *MockSendAlgorithmWithDebugInfosMockRecorder
}

// MockSendAlgorithmWithDebugInfosMockRecorder is the mock recorder for MockSendAlgorithmWithDebugInfos.
type MockSendAlgorithmWithDebugInfosMockRecorder struct {
	mock *MockSendAlgorithmWithDebugInfos
}

// NewMockSendAlgorithmWithDebugInfos creates a new mock instance.
func NewMockSendAlgorithmWithDebugInfos(ctrl *gomock.Controller) *MockSendAlgorithmWithDebugInfos {
	mock := &MockSendAlgorithmWithDebugInfos{ctrl: ctrl}
	mock.recorder = &MockSendAlgorithmWithDebugInfosMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSendAlgorithmWithDebugInfos) EXPECT() *MockSendAlgorithmWithDebugInfosMockRecorder {
	return m.recorder
}

// CanSend mocks base method.
func (m *MockSendAlgorithmWithDebugInfos) CanSend(arg0 protocol.ByteCount) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CanSend", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// CanSend indicates an expected call of CanSend.
func (mr *MockSendAlgorithmWithDebugInfosMockRecorder) CanSend(arg0 any) *MockSendAlgorithmWithDebugInfosCanSendCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CanSend", reflect.TypeOf((*MockSendAlgorithmWithDebugInfos)(nil).CanSend), arg0)
	return &MockSendAlgorithmWithDebugInfosCanSendCall{Call: call}
}

// MockSendAlgorithmWithDebugInfosCanSendCall wrap *gomock.Call
type MockSendAlgorithmWithDebugInfosCanSendCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockSendAlgorithmWithDebugInfosCanSendCall) Return(arg0 bool) *MockSendAlgorithmWithDebugInfosCanSendCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockSendAlgorithmWithDebugInfosCanSendCall) Do(f func(protocol.ByteCount) bool) *MockSendAlgorithmWithDebugInfosCanSendCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockSendAlgorithmWithDebugInfosCanSendCall) DoAndReturn(f func(protocol.ByteCount) bool) *MockSendAlgorithmWithDebugInfosCanSendCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetCongestionWindow mocks base method.
func (m *MockSendAlgorithmWithDebugInfos) GetCongestionWindow() protocol.ByteCount {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCongestionWindow")
	ret0, _ := ret[0].(protocol.ByteCount)
	return ret0
}

// GetCongestionWindow indicates an expected call of GetCongestionWindow.
func (mr *MockSendAlgorithmWithDebugInfosMockRecorder) GetCongestionWindow() *MockSendAlgorithmWithDebugInfosGetCongestionWindowCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCongestionWindow", reflect.TypeOf((*MockSendAlgorithmWithDebugInfos)(nil).GetCongestionWindow))
	return &MockSendAlgorithmWithDebugInfosGetCongestionWindowCall{Call: call}
}

// MockSendAlgorithmWithDebugInfosGetCongestionWindowCall wrap *gomock.Call
type MockSendAlgorithmWithDebugInfosGetCongestionWindowCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockSendAlgorithmWithDebugInfosGetCongestionWindowCall) Return(arg0 protocol.ByteCount) *MockSendAlgorithmWithDebugInfosGetCongestionWindowCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockSendAlgorithmWithDebugInfosGetCongestionWindowCall) Do(f func() protocol.ByteCount) *MockSendAlgorithmWithDebugInfosGetCongestionWindowCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockSendAlgorithmWithDebugInfosGetCongestionWindowCall) DoAndReturn(f func() protocol.ByteCount) *MockSendAlgorithmWithDebugInfosGetCongestionWindowCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// HasPacingBudget mocks base method.
func (m *MockSendAlgorithmWithDebugInfos) HasPacingBudget(arg0 time.Time) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasPacingBudget", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// HasPacingBudget indicates an expected call of HasPacingBudget.
func (mr *MockSendAlgorithmWithDebugInfosMockRecorder) HasPacingBudget(arg0 any) *MockSendAlgorithmWithDebugInfosHasPacingBudgetCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasPacingBudget", reflect.TypeOf((*MockSendAlgorithmWithDebugInfos)(nil).HasPacingBudget), arg0)
	return &MockSendAlgorithmWithDebugInfosHasPacingBudgetCall{Call: call}
}

// MockSendAlgorithmWithDebugInfosHasPacingBudgetCall wrap *gomock.Call
type MockSendAlgorithmWithDebugInfosHasPacingBudgetCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockSendAlgorithmWithDebugInfosHasPacingBudgetCall) Return(arg0 bool) *MockSendAlgorithmWithDebugInfosHasPacingBudgetCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockSendAlgorithmWithDebugInfosHasPacingBudgetCall) Do(f func(time.Time) bool) *MockSendAlgorithmWithDebugInfosHasPacingBudgetCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockSendAlgorithmWithDebugInfosHasPacingBudgetCall) DoAndReturn(f func(time.Time) bool) *MockSendAlgorithmWithDebugInfosHasPacingBudgetCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// InRecovery mocks base method.
func (m *MockSendAlgorithmWithDebugInfos) InRecovery() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InRecovery")
	ret0, _ := ret[0].(bool)
	return ret0
}

// InRecovery indicates an expected call of InRecovery.
func (mr *MockSendAlgorithmWithDebugInfosMockRecorder) InRecovery() *MockSendAlgorithmWithDebugInfosInRecoveryCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InRecovery", reflect.TypeOf((*MockSendAlgorithmWithDebugInfos)(nil).InRecovery))
	return &MockSendAlgorithmWithDebugInfosInRecoveryCall{Call: call}
}

// MockSendAlgorithmWithDebugInfosInRecoveryCall wrap *gomock.Call
type MockSendAlgorithmWithDebugInfosInRecoveryCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockSendAlgorithmWithDebugInfosInRecoveryCall) Return(arg0 bool) *MockSendAlgorithmWithDebugInfosInRecoveryCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockSendAlgorithmWithDebugInfosInRecoveryCall) Do(f func() bool) *MockSendAlgorithmWithDebugInfosInRecoveryCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockSendAlgorithmWithDebugInfosInRecoveryCall) DoAndReturn(f func() bool) *MockSendAlgorithmWithDebugInfosInRecoveryCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// InSlowStart mocks base method.
func (m *MockSendAlgorithmWithDebugInfos) InSlowStart() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InSlowStart")
	ret0, _ := ret[0].(bool)
	return ret0
}

// InSlowStart indicates an expected call of InSlowStart.
func (mr *MockSendAlgorithmWithDebugInfosMockRecorder) InSlowStart() *MockSendAlgorithmWithDebugInfosInSlowStartCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InSlowStart", reflect.TypeOf((*MockSendAlgorithmWithDebugInfos)(nil).InSlowStart))
	return &MockSendAlgorithmWithDebugInfosInSlowStartCall{Call: call}
}

// MockSendAlgorithmWithDebugInfosInSlowStartCall wrap *gomock.Call
type MockSendAlgorithmWithDebugInfosInSlowStartCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockSendAlgorithmWithDebugInfosInSlowStartCall) Return(arg0 bool) *MockSendAlgorithmWithDebugInfosInSlowStartCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockSendAlgorithmWithDebugInfosInSlowStartCall) Do(f func() bool) *MockSendAlgorithmWithDebugInfosInSlowStartCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockSendAlgorithmWithDebugInfosInSlowStartCall) DoAndReturn(f func() bool) *MockSendAlgorithmWithDebugInfosInSlowStartCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MaybeExitSlowStart mocks base method.
func (m *MockSendAlgorithmWithDebugInfos) MaybeExitSlowStart() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "MaybeExitSlowStart")
}

// MaybeExitSlowStart indicates an expected call of MaybeExitSlowStart.
func (mr *MockSendAlgorithmWithDebugInfosMockRecorder) MaybeExitSlowStart() *MockSendAlgorithmWithDebugInfosMaybeExitSlowStartCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MaybeExitSlowStart", reflect.TypeOf((*MockSendAlgorithmWithDebugInfos)(nil).MaybeExitSlowStart))
	return &MockSendAlgorithmWithDebugInfosMaybeExitSlowStartCall{Call: call}
}

// MockSendAlgorithmWithDebugInfosMaybeExitSlowStartCall wrap *gomock.Call
type MockSendAlgorithmWithDebugInfosMaybeExitSlowStartCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockSendAlgorithmWithDebugInfosMaybeExitSlowStartCall) Return() *MockSendAlgorithmWithDebugInfosMaybeExitSlowStartCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockSendAlgorithmWithDebugInfosMaybeExitSlowStartCall) Do(f func()) *MockSendAlgorithmWithDebugInfosMaybeExitSlowStartCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockSendAlgorithmWithDebugInfosMaybeExitSlowStartCall) DoAndReturn(f func()) *MockSendAlgorithmWithDebugInfosMaybeExitSlowStartCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// OnCongestionEvent mocks base method.
func (m *MockSendAlgorithmWithDebugInfos) OnCongestionEvent(arg0 protocol.PacketNumber, arg1, arg2 protocol.ByteCount) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnCongestionEvent", arg0, arg1, arg2)
}

// OnCongestionEvent indicates an expected call of OnCongestionEvent.
func (mr *MockSendAlgorithmWithDebugInfosMockRecorder) OnCongestionEvent(arg0, arg1, arg2 any) *MockSendAlgorithmWithDebugInfosOnCongestionEventCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnCongestionEvent", reflect.TypeOf((*MockSendAlgorithmWithDebugInfos)(nil).OnCongestionEvent), arg0, arg1, arg2)
	return &MockSendAlgorithmWithDebugInfosOnCongestionEventCall{Call: call}
}

// MockSendAlgorithmWithDebugInfosOnCongestionEventCall wrap *gomock.Call
type MockSendAlgorithmWithDebugInfosOnCongestionEventCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockSendAlgorithmWithDebugInfosOnCongestionEventCall) Return() *MockSendAlgorithmWithDebugInfosOnCongestionEventCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockSendAlgorithmWithDebugInfosOnCongestionEventCall) Do(f func(protocol.PacketNumber, protocol.ByteCount, protocol.ByteCount)) *MockSendAlgorithmWithDebugInfosOnCongestionEventCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockSendAlgorithmWithDebugInfosOnCongestionEventCall) DoAndReturn(f func(protocol.PacketNumber, protocol.ByteCount, protocol.ByteCount)) *MockSendAlgorithmWithDebugInfosOnCongestionEventCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// OnPacketAcked mocks base method.
func (m *MockSendAlgorithmWithDebugInfos) OnPacketAcked(arg0 protocol.PacketNumber, arg1, arg2 protocol.ByteCount, arg3 time.Time) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnPacketAcked", arg0, arg1, arg2, arg3)
}

// OnPacketAcked indicates an expected call of OnPacketAcked.
func (mr *MockSendAlgorithmWithDebugInfosMockRecorder) OnPacketAcked(arg0, arg1, arg2, arg3 any) *MockSendAlgorithmWithDebugInfosOnPacketAckedCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnPacketAcked", reflect.TypeOf((*MockSendAlgorithmWithDebugInfos)(nil).OnPacketAcked), arg0, arg1, arg2, arg3)
	return &MockSendAlgorithmWithDebugInfosOnPacketAckedCall{Call: call}
}

// MockSendAlgorithmWithDebugInfosOnPacketAckedCall wrap *gomock.Call
type MockSendAlgorithmWithDebugInfosOnPacketAckedCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockSendAlgorithmWithDebugInfosOnPacketAckedCall) Return() *MockSendAlgorithmWithDebugInfosOnPacketAckedCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockSendAlgorithmWithDebugInfosOnPacketAckedCall) Do(f func(protocol.PacketNumber, protocol.ByteCount, protocol.ByteCount, time.Time)) *MockSendAlgorithmWithDebugInfosOnPacketAckedCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockSendAlgorithmWithDebugInfosOnPacketAckedCall) DoAndReturn(f func(protocol.PacketNumber, protocol.ByteCount, protocol.ByteCount, time.Time)) *MockSendAlgorithmWithDebugInfosOnPacketAckedCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// OnPacketSent mocks base method.
func (m *MockSendAlgorithmWithDebugInfos) OnPacketSent(arg0 time.Time, arg1 protocol.ByteCount, arg2 protocol.PacketNumber, arg3 protocol.ByteCount, arg4 bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnPacketSent", arg0, arg1, arg2, arg3, arg4)
}

// OnPacketSent indicates an expected call of OnPacketSent.
func (mr *MockSendAlgorithmWithDebugInfosMockRecorder) OnPacketSent(arg0, arg1, arg2, arg3, arg4 any) *MockSendAlgorithmWithDebugInfosOnPacketSentCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnPacketSent", reflect.TypeOf((*MockSendAlgorithmWithDebugInfos)(nil).OnPacketSent), arg0, arg1, arg2, arg3, arg4)
	return &MockSendAlgorithmWithDebugInfosOnPacketSentCall{Call: call}
}

// MockSendAlgorithmWithDebugInfosOnPacketSentCall wrap *gomock.Call
type MockSendAlgorithmWithDebugInfosOnPacketSentCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockSendAlgorithmWithDebugInfosOnPacketSentCall) Return() *MockSendAlgorithmWithDebugInfosOnPacketSentCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockSendAlgorithmWithDebugInfosOnPacketSentCall) Do(f func(time.Time, protocol.ByteCount, protocol.PacketNumber, protocol.ByteCount, bool)) *MockSendAlgorithmWithDebugInfosOnPacketSentCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockSendAlgorithmWithDebugInfosOnPacketSentCall) DoAndReturn(f func(time.Time, protocol.ByteCount, protocol.PacketNumber, protocol.ByteCount, bool)) *MockSendAlgorithmWithDebugInfosOnPacketSentCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// OnRetransmissionTimeout mocks base method.
func (m *MockSendAlgorithmWithDebugInfos) OnRetransmissionTimeout(arg0 bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnRetransmissionTimeout", arg0)
}

// OnRetransmissionTimeout indicates an expected call of OnRetransmissionTimeout.
func (mr *MockSendAlgorithmWithDebugInfosMockRecorder) OnRetransmissionTimeout(arg0 any) *MockSendAlgorithmWithDebugInfosOnRetransmissionTimeoutCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnRetransmissionTimeout", reflect.TypeOf((*MockSendAlgorithmWithDebugInfos)(nil).OnRetransmissionTimeout), arg0)
	return &MockSendAlgorithmWithDebugInfosOnRetransmissionTimeoutCall{Call: call}
}

// MockSendAlgorithmWithDebugInfosOnRetransmissionTimeoutCall wrap *gomock.Call
type MockSendAlgorithmWithDebugInfosOnRetransmissionTimeoutCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockSendAlgorithmWithDebugInfosOnRetransmissionTimeoutCall) Return() *MockSendAlgorithmWithDebugInfosOnRetransmissionTimeoutCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockSendAlgorithmWithDebugInfosOnRetransmissionTimeoutCall) Do(f func(bool)) *MockSendAlgorithmWithDebugInfosOnRetransmissionTimeoutCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockSendAlgorithmWithDebugInfosOnRetransmissionTimeoutCall) DoAndReturn(f func(bool)) *MockSendAlgorithmWithDebugInfosOnRetransmissionTimeoutCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SetCongestionWindow mocks base method.
func (m *MockSendAlgorithmWithDebugInfos) SetCongestionWindow(arg0 protocol.ByteCount) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetCongestionWindow", arg0)
}

// SetCongestionWindow indicates an expected call of SetCongestionWindow.
func (mr *MockSendAlgorithmWithDebugInfosMockRecorder) SetCongestionWindow(arg0 any) *SendAlgorithmWithDebugInfosSetCongestionWindowCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetCongestionWindow", reflect.TypeOf((*MockSendAlgorithmWithDebugInfos)(nil).SetCongestionWindow), arg0)
	return &SendAlgorithmWithDebugInfosSetCongestionWindowCall{Call: call}
}

// SendAlgorithmWithDebugInfosSetCongestionWindowCall wrap *gomock.Call
type SendAlgorithmWithDebugInfosSetCongestionWindowCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *SendAlgorithmWithDebugInfosSetCongestionWindowCall) Return() *SendAlgorithmWithDebugInfosSetCongestionWindowCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *SendAlgorithmWithDebugInfosSetCongestionWindowCall) Do(f func(protocol.ByteCount)) *SendAlgorithmWithDebugInfosSetCongestionWindowCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *SendAlgorithmWithDebugInfosSetCongestionWindowCall) DoAndReturn(f func(protocol.ByteCount)) *SendAlgorithmWithDebugInfosSetCongestionWindowCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SetInitialCongestionWindow mocks base method.
func (m *MockSendAlgorithmWithDebugInfos) SetInitialCongestionWindow(arg0 uint32) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetInitialCongestionWindow", arg0)
}

// SetInitialCongestionWindow indicates an expected call of SetInitialCongestionWindow.
func (mr *MockSendAlgorithmWithDebugInfosMockRecorder) SetInitialCongestionWindow(arg0 any) *SendAlgorithmWithDebugInfosSetInitialCongestionWindowCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetInitialCongestionWindow", reflect.TypeOf((*MockSendAlgorithmWithDebugInfos)(nil).SetInitialCongestionWindow), arg0)
	return &SendAlgorithmWithDebugInfosSetInitialCongestionWindowCall{Call: call}
}

// SendAlgorithmWithDebugInfosSetInitialCongestionWindowCall wrap *gomock.Call
type SendAlgorithmWithDebugInfosSetInitialCongestionWindowCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *SendAlgorithmWithDebugInfosSetInitialCongestionWindowCall) Return() *SendAlgorithmWithDebugInfosSetInitialCongestionWindowCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *SendAlgorithmWithDebugInfosSetInitialCongestionWindowCall) Do(f func(uint32)) *SendAlgorithmWithDebugInfosSetInitialCongestionWindowCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *SendAlgorithmWithDebugInfosSetInitialCongestionWindowCall) DoAndReturn(f func(uint32)) *SendAlgorithmWithDebugInfosSetInitialCongestionWindowCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SetMaxBandwidth mocks base method.
func (m *MockSendAlgorithmWithDebugInfos) SetMaxBandwidth(arg0 congestion.Bandwidth) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetMaxBandwidth", arg0)
}

// SetMaxBandwidth indicates an expected call of SetMaxBandwidth.
func (mr *MockSendAlgorithmWithDebugInfosMockRecorder) SetMaxBandwidth(arg0 any) *SendAlgorithmWithDebugInfosSetMaxBandwidthCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMaxBandwidth", reflect.TypeOf((*MockSendAlgorithmWithDebugInfos)(nil).SetMaxBandwidth), arg0)
	return &SendAlgorithmWithDebugInfosSetMaxBandwidthCall{Call: call}
}

// SendAlgorithmWithDebugInfosSetMaxBandwidthCall wrap *gomock.Call
type SendAlgorithmWithDebugInfosSetMaxBandwidthCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *SendAlgorithmWithDebugInfosSetMaxBandwidthCall) Return() *SendAlgorithmWithDebugInfosSetMaxBandwidthCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *SendAlgorithmWithDebugInfosSetMaxBandwidthCall) Do(f func(congestion.Bandwidth)) *SendAlgorithmWithDebugInfosSetMaxBandwidthCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *SendAlgorithmWithDebugInfosSetMaxBandwidthCall) DoAndReturn(f func(congestion.Bandwidth)) *SendAlgorithmWithDebugInfosSetMaxBandwidthCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SetMaxDatagramSize mocks base method.
func (m *MockSendAlgorithmWithDebugInfos) SetMaxDatagramSize(arg0 protocol.ByteCount) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetMaxDatagramSize", arg0)
}

// SetMaxDatagramSize indicates an expected call of SetMaxDatagramSize.
func (mr *MockSendAlgorithmWithDebugInfosMockRecorder) SetMaxDatagramSize(arg0 any) *MockSendAlgorithmWithDebugInfosSetMaxDatagramSizeCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMaxDatagramSize", reflect.TypeOf((*MockSendAlgorithmWithDebugInfos)(nil).SetMaxDatagramSize), arg0)
	return &MockSendAlgorithmWithDebugInfosSetMaxDatagramSizeCall{Call: call}
}

// MockSendAlgorithmWithDebugInfosSetMaxDatagramSizeCall wrap *gomock.Call
type MockSendAlgorithmWithDebugInfosSetMaxDatagramSizeCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockSendAlgorithmWithDebugInfosSetMaxDatagramSizeCall) Return() *MockSendAlgorithmWithDebugInfosSetMaxDatagramSizeCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockSendAlgorithmWithDebugInfosSetMaxDatagramSizeCall) Do(f func(protocol.ByteCount)) *MockSendAlgorithmWithDebugInfosSetMaxDatagramSizeCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockSendAlgorithmWithDebugInfosSetMaxDatagramSizeCall) DoAndReturn(f func(protocol.ByteCount)) *MockSendAlgorithmWithDebugInfosSetMaxDatagramSizeCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// TimeUntilSend mocks base method.
func (m *MockSendAlgorithmWithDebugInfos) TimeUntilSend(arg0 protocol.ByteCount) time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TimeUntilSend", arg0)
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// TimeUntilSend indicates an expected call of TimeUntilSend.
func (mr *MockSendAlgorithmWithDebugInfosMockRecorder) TimeUntilSend(arg0 any) *MockSendAlgorithmWithDebugInfosTimeUntilSendCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TimeUntilSend", reflect.TypeOf((*MockSendAlgorithmWithDebugInfos)(nil).TimeUntilSend), arg0)
	return &MockSendAlgorithmWithDebugInfosTimeUntilSendCall{Call: call}
}

// MockSendAlgorithmWithDebugInfosTimeUntilSendCall wrap *gomock.Call
type MockSendAlgorithmWithDebugInfosTimeUntilSendCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockSendAlgorithmWithDebugInfosTimeUntilSendCall) Return(arg0 time.Time) *MockSendAlgorithmWithDebugInfosTimeUntilSendCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockSendAlgorithmWithDebugInfosTimeUntilSendCall) Do(f func(protocol.ByteCount) time.Time) *MockSendAlgorithmWithDebugInfosTimeUntilSendCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockSendAlgorithmWithDebugInfosTimeUntilSendCall) DoAndReturn(f func(protocol.ByteCount) time.Time) *MockSendAlgorithmWithDebugInfosTimeUntilSendCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
