package xse

import (
	"context"
	"github.com/lucas-clemente/quic-go/internal/ackhandler"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/qerr"
	"github.com/lucas-clemente/quic-go/internal/wire"
	"io"
	"time"
)

type RecordNumber uint64

type RecordHeader DecryptedPayloadLength

func (r RecordHeader) DecryptedPayloadLength() DecryptedPayloadLength {
	return DecryptedPayloadLength(r)
}

type DecryptedPayloadLength uint16

type RecordEncryptedPayload []byte

const MaxDecryptedPayloadLength = ^DecryptedPayloadLength(0)

type CryptoSetup interface {
	Seal(dst []byte, src []byte, sid protocol.StreamID, rn RecordNumber) []byte
	Open(RecordEncryptedPayload, protocol.StreamID, RecordNumber) ([]byte, error)
	EncryptedRecordPayloadLength(DecryptedPayloadLength) uint32
	MaxEncryptedRecordPayloadLength() uint32
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
