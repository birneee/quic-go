package quic

import (
	"github.com/quic-go/quic-go/internal/ackhandler"
	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/wire"
	"github.com/quic-go/quic-go/internal/xads"
)

// this compatibility wrapper is required to use a xads.Stream as a quic.streamI
type xadsStreamI struct {
	xads.Stream
}

var _ streamI = &xadsStreamI{}

func (x xadsStreamI) closeForShutdown(err error) {
	x.Stream.CloseForShutdown(err)
}

func (x xadsStreamI) handleStreamFrame(frame *wire.StreamFrame) error {
	return x.Stream.HandleStreamFrame(frame)
}

func (x xadsStreamI) handleResetStreamFrame(frame *wire.ResetStreamFrame) error {
	return x.Stream.HandleResetStreamFrame(frame)
}

func (x xadsStreamI) getWindowUpdate() protocol.ByteCount {
	return x.Stream.GetWindowUpdate()
}

func (x xadsStreamI) hasData() bool {
	return x.Stream.HasData()
}

func (x xadsStreamI) handleStopSendingFrame(frame *wire.StopSendingFrame) {
	x.Stream.HandleStopSendingFrame(frame)
}

func (x xadsStreamI) popStreamFrame(maxBytes protocol.ByteCount, v protocol.VersionNumber) (ackhandler.StreamFrame, bool, bool) {
	return x.Stream.PopStreamFrame(maxBytes, v)
}

func (x xadsStreamI) updateSendWindow(count protocol.ByteCount) {
	x.Stream.UpdateSendWindow(count)
}

// this compatibility wrapper is required to use a xads.SendStream as a quic.sendStreamI
type xadsSendStreamI struct {
	xads.SendStream
}

var _ sendStreamI = &xadsSendStreamI{}

func (x xadsSendStreamI) handleStopSendingFrame(frame *wire.StopSendingFrame) {
	x.SendStream.HandleStopSendingFrame(frame)
}

func (x xadsSendStreamI) hasData() bool {
	return x.SendStream.HasData()
}

func (x xadsSendStreamI) popStreamFrame(maxBytes protocol.ByteCount, v protocol.VersionNumber) (ackhandler.StreamFrame, bool, bool) {
	return x.SendStream.PopStreamFrame(maxBytes, v)
}

func (x xadsSendStreamI) closeForShutdown(err error) {
	x.SendStream.CloseForShutdown(err)
}

func (x xadsSendStreamI) updateSendWindow(count protocol.ByteCount) {
	x.SendStream.UpdateSendWindow(count)
}

// this compatibility wrapper is required to use a xads.ReceiveStream as a quic.xadsReceiveStreamI
type xadsReceiveStreamI struct {
	xads.ReceiveStream
}

var _ receiveStreamI = &xadsReceiveStreamI{}

func (x xadsReceiveStreamI) handleStreamFrame(frame *wire.StreamFrame) error {
	return x.ReceiveStream.HandleStreamFrame(frame)
}

func (x xadsReceiveStreamI) handleResetStreamFrame(frame *wire.ResetStreamFrame) error {
	return x.ReceiveStream.HandleResetStreamFrame(frame)
}

func (x xadsReceiveStreamI) closeForShutdown(err error) {
	x.ReceiveStream.CloseForShutdown(err)
}

func (x xadsReceiveStreamI) getWindowUpdate() protocol.ByteCount {
	return x.ReceiveStream.GetWindowUpdate()
}

var _ xads.SendStream = &sendStream{}

func (s *sendStream) HasData() bool {
	return s.hasData()
}

func (s *sendStream) HandleStopSendingFrame(frame *wire.StopSendingFrame) {
	s.handleStopSendingFrame(frame)
}

func (s *sendStream) PopStreamFrame(maxBytes protocol.ByteCount, v protocol.VersionNumber) (ackhandler.StreamFrame, bool, bool) {
	return s.popStreamFrame(maxBytes, v)
}

func (s *sendStream) CloseForShutdown(err error) {
	s.closeForShutdown(err)
}

func (s *sendStream) UpdateSendWindow(count protocol.ByteCount) {
	s.updateSendWindow(count)
}

var _ xads.ReceiveStream = &receiveStream{}

func (s *receiveStream) HandleStreamFrame(frame *wire.StreamFrame) error {
	return s.handleStreamFrame(frame)
}

func (s *receiveStream) HandleResetStreamFrame(frame *wire.ResetStreamFrame) error {
	return s.handleResetStreamFrame(frame)
}

func (s *receiveStream) CloseForShutdown(err error) {
	s.closeForShutdown(err)
}

func (s *receiveStream) GetWindowUpdate() protocol.ByteCount {
	return s.getWindowUpdate()
}

var _ streamI = &stream{}
var _ xads.Stream = &stream{}

func (s *stream) ReceiveStream() xads.ReceiveStream {
	return &s.receiveStream
}

func (s *stream) SendStream() xads.SendStream {
	return &s.sendStream
}

func (s *stream) CloseForShutdown(err error) {
	s.closeForShutdown(err)
}
