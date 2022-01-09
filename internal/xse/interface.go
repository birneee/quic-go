package xse

import (
	"context"
	"encoding/binary"
	"github.com/lucas-clemente/quic-go/internal/ackhandler"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/qerr"
	"github.com/lucas-clemente/quic-go/internal/wire"
	"io"
	"time"
)

type RecordNumber uint64

// RecordHeader contains the decrypted payload length.
// Must be 2 bytes long.
type RecordHeader []byte

func (r RecordHeader) DecryptedPayloadLength() DecryptedPayloadLength {
	return DecryptedPayloadLength(binary.BigEndian.Uint16(r))
}

// SetDecryptedPayloadLength sets length of plaintext.
// Must not be 0.
// TODO error handling
func (r RecordHeader) SetDecryptedPayloadLength(length DecryptedPayloadLength) {
	if length == 0 {
		panic("XSE-QUIC protocol violation: record plaintext has length 0")
	}
	binary.BigEndian.PutUint16(r, uint16(length))
}

type DecryptedPayloadLength uint16

// RecordEncryptedPayload does not include the header
type RecordEncryptedPayload []byte

const MaxDecryptedPayloadLength = ^DecryptedPayloadLength(0)

// TODO add key update mechanism
type CryptoSetup interface {
	// Seal encrypts and authenticates plaintext, authenticates the
	// additional data and appends the result to dst, returning the updated
	// slice.
	//
	// To reuse plaintext's storage for the encrypted output, use plaintext[:0]
	// as dst. Otherwise, the remaining capacity of dst must not overlap plaintext.
	Seal(dst []byte, plaintext []byte, sid protocol.StreamID, rn RecordNumber) []byte
	// Open decrypts and authenticates ciphertext, authenticates the
	// additional data and, if successful, appends the resulting plaintext
	// to dst, returning the updated slice.
	// The StreamID and the RecordNumber must match the values passed to Seal.
	//
	// To reuse ciphertext's storage for the decrypted output, use ciphertext[:0]
	// as dst. Otherwise, the remaining capacity of dst must not overlap plaintext.
	//
	// Even if the function fails, the contents of dst, up to its capacity,
	// may be overwritten.
	Open(dst []byte, ciphertext RecordEncryptedPayload, sid protocol.StreamID, rn RecordNumber) ([]byte, error)
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
