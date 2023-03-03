package quic

import (
	"github.com/lucas-clemente/quic-go/handover"
	"github.com/lucas-clemente/quic-go/internal/ackhandler"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/wire"
	"github.com/lucas-clemente/quic-go/internal/xse"
)

// this compatibility wrapper is required to use a xse.Stream as a quic.streamI
type xseStreamI struct {
	xse.Stream
}

var _ streamI = &xseStreamI{}

func (x xseStreamI) closeForShutdown(err error) {
	x.Stream.CloseForShutdown(err)
}

func (x xseStreamI) handleStreamFrame(frame *wire.StreamFrame) error {
	return x.Stream.HandleStreamFrame(frame)
}

func (x xseStreamI) handleResetStreamFrame(frame *wire.ResetStreamFrame) error {
	return x.Stream.HandleResetStreamFrame(frame)
}

func (x xseStreamI) getWindowUpdate() protocol.ByteCount {
	return x.Stream.GetWindowUpdate()
}

func (x xseStreamI) hasData() bool {
	return x.Stream.HasData()
}

func (x xseStreamI) handleStopSendingFrame(frame *wire.StopSendingFrame) {
	x.Stream.HandleStopSendingFrame(frame)
}

func (x xseStreamI) popStreamFrame(maxBytes protocol.ByteCount) (*ackhandler.Frame, bool) {
	return x.Stream.PopStreamFrame(maxBytes)
}

func (x xseStreamI) updateSendWindow(count protocol.ByteCount) {
	x.Stream.UpdateSendWindow(count)
}

func (x xseStreamI) storeReceiveState(state *handover.BidiStreamState, perspective protocol.Perspective) {
	//TODO implement me
	panic("implement me")
}

func (x xseStreamI) restoreReceiveState(state *handover.BidiStreamState, perspective protocol.Perspective) {
	//TODO implement me
	panic("implement me")
}

func (x xseStreamI) storeSendState(state *handover.BidiStreamState, perspective protocol.Perspective) {
	//TODO implement me
	panic("implement me")
}

func (x xseStreamI) restoreSendState(state *handover.BidiStreamState, perspective protocol.Perspective) {
	//TODO implement me
	panic("implement me")
}

// this compatibility wrapper is required to use a xse.SendStream as a quic.sendStreamI
type xseSendStreamI struct {
	xse.SendStream
}

var _ sendStreamI = &xseSendStreamI{}

func (x xseSendStreamI) handleStopSendingFrame(frame *wire.StopSendingFrame) {
	x.SendStream.HandleStopSendingFrame(frame)
}

func (x xseSendStreamI) hasData() bool {
	return x.SendStream.HasData()
}

func (x xseSendStreamI) popStreamFrame(maxBytes protocol.ByteCount) (*ackhandler.Frame, bool) {
	return x.SendStream.PopStreamFrame(maxBytes)
}

func (x xseSendStreamI) closeForShutdown(err error) {
	x.SendStream.CloseForShutdown(err)
}

func (x xseSendStreamI) updateSendWindow(count protocol.ByteCount) {
	x.SendStream.UpdateSendWindow(count)
}

func (x xseSendStreamI) storeSendState(state *handover.BidiStreamState, perspective protocol.Perspective) {
	//TODO implement me
	panic("implement me")
}

func (x xseSendStreamI) restoreSendState(state *handover.BidiStreamState, perspective protocol.Perspective) {
	//TODO implement me
	panic("implement me")
}

// this compatibility wrapper is required to use a xse.ReceiveStream as a quic.xseReceiveStreamI
type xseReceiveStreamI struct {
	xse.ReceiveStream
}

var _ receiveStreamI = &xseReceiveStreamI{}

func (x xseReceiveStreamI) handleStreamFrame(frame *wire.StreamFrame) error {
	return x.ReceiveStream.HandleStreamFrame(frame)
}

func (x xseReceiveStreamI) handleResetStreamFrame(frame *wire.ResetStreamFrame) error {
	return x.ReceiveStream.HandleResetStreamFrame(frame)
}

func (x xseReceiveStreamI) closeForShutdown(err error) {
	x.ReceiveStream.CloseForShutdown(err)
}

func (x xseReceiveStreamI) getWindowUpdate() protocol.ByteCount {
	return x.ReceiveStream.GetWindowUpdate()
}

func (x xseReceiveStreamI) storeReceiveState(state *handover.BidiStreamState, perspective protocol.Perspective) {
	//TODO implement me
	panic("implement me")
}

func (x xseReceiveStreamI) restoreReceiveState(state *handover.BidiStreamState, perspective protocol.Perspective) {
	//TODO implement me
	panic("implement me")
}
