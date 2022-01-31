package xse

import (
	"context"
	"io"
	"time"

	"github.com/lucas-clemente/quic-go/internal/ackhandler"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/qerr"
	"github.com/lucas-clemente/quic-go/internal/wire"
)

// TODO add key update mechanism
type CryptoSetup interface {
	NewStream(base Stream) Stream
	NewReceiveStream(base ReceiveStream) ReceiveStream
	NewSendStream(base SendStream) SendStream
}

type Stream interface {
	ReceiveStream
	SendStream
	SetDeadline(t time.Time) error
	ReceiveStream() ReceiveStream
	SendStream() SendStream
}

type SendStream interface {
	StreamID() protocol.StreamID
	io.Writer
	io.Closer
	CancelWrite(qerr.StreamErrorCode)
	Context() context.Context
	SetWriteDeadline(t time.Time) error
	HasData() bool
	HandleStopSendingFrame(*wire.StopSendingFrame)
	PopStreamFrame(maxBytes protocol.ByteCount) (*ackhandler.Frame, bool)
	CloseForShutdown(error)
	UpdateSendWindow(protocol.ByteCount)
}

type ReceiveStream interface {
	StreamID() protocol.StreamID
	io.Reader
	CancelRead(qerr.StreamErrorCode)
	SetReadDeadline(t time.Time) error
	HandleStreamFrame(*wire.StreamFrame) error
	HandleResetStreamFrame(*wire.ResetStreamFrame) error
	CloseForShutdown(error)
	GetWindowUpdate() protocol.ByteCount
}
