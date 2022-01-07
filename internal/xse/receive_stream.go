package xse

import (
	"encoding/binary"
	"io"
)

type receiveStream struct {
	ReceiveStream
	xseSealer        CryptoSetup
	nextRecordNumber RecordNumber
	buf              []byte
	nextBytes        []byte
}

// NewReceiveStream creates a XSE-QUIC receive stream
func NewReceiveStream(baseStream ReceiveStream, xseSealer CryptoSetup) *receiveStream {
	return &receiveStream{
		ReceiveStream: baseStream,
		xseSealer:     xseSealer,
		buf:           make([]byte, 0, xseSealer.MaxEncryptedRecordPayloadLength()),
	}
}

//TODO read multiple Xse records if available
func (x *receiveStream) Read(p []byte) (n int, err error) {
	read := 0
	for read < len(p) {
		if len(x.nextBytes) == 0 {
			err := x.readNextRecord()
			if err != nil {
				return 0, err
			}
		}
		copied := copy(p[read:], x.nextBytes)
		read += copied
		x.nextBytes = x.nextBytes[copied:]
		break
	}
	return read, err
}

func (x *receiveStream) readRecordHeader() (RecordHeader, error) {
	x.buf = x.buf[:2]
	_, err := io.ReadFull(x.ReceiveStream, x.buf)
	if err != nil {
		return 0, err
	}
	return RecordHeader(binary.BigEndian.Uint16(x.buf)), nil
}

func (x *receiveStream) readNextRecord() error {
	hdr, err := x.readRecordHeader()
	if err != nil {
		return err
	}
	size := x.xseSealer.EncryptedRecordPayloadLength(hdr.DecryptedPayloadLength())
	x.buf = x.buf[:size]
	_, err = io.ReadFull(x.ReceiveStream, x.buf)
	if err != nil {
		return err
	}
	x.buf, err = x.xseSealer.Open(x.buf, x.StreamID(), x.nextRecordNumber)
	x.nextRecordNumber++
	if err != nil {
		return err
	}
	x.nextBytes = x.buf
	return nil
}
