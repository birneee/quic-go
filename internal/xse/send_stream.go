package xse

import (
	"encoding/binary"
)

type sendStream struct {
	SendStream
	cryptoSetup      CryptoSetup
	nextRecordNumber RecordNumber
	sendBuf          []byte
}

// NewSendStream creates a XSE-QUIC send stream
func NewSendStream(baseStream SendStream, cryptoSetup CryptoSetup) *sendStream {
	return &sendStream{
		SendStream:       baseStream,
		cryptoSetup:      cryptoSetup,
		nextRecordNumber: 0,
		sendBuf:          make([]byte, cryptoSetup.MaxEncryptedRecordPayloadLength()),
	}
}

//TODO do not send small records except on timeout, collect some data before writing
func (s *sendStream) Write(p []byte) (n int, err error) {
	writeNext := p
	for len(writeNext) > 0 {
		var decryptedPayloadLength DecryptedPayloadLength
		if len(writeNext) > int(MaxDecryptedPayloadLength) {
			decryptedPayloadLength = MaxDecryptedPayloadLength
		} else {
			decryptedPayloadLength = DecryptedPayloadLength(len(writeNext))
		}
		binary.BigEndian.PutUint16(s.sendBuf[:2], uint16(decryptedPayloadLength))
		_, err := s.SendStream.Write(s.sendBuf[:2])
		if err != nil {
			return len(p) - len(writeNext), err
		}
		s.sendBuf = s.cryptoSetup.Seal(s.sendBuf[:0], writeNext[:decryptedPayloadLength], s.StreamID(), s.nextRecordNumber)
		s.nextRecordNumber++
		_, err = s.SendStream.Write(s.sendBuf)
		if err != nil {
			return len(p) - len(writeNext), err
		}
		writeNext = writeNext[decryptedPayloadLength:]
	}
	return len(p) - len(writeNext), nil
}
