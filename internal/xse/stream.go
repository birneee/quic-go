package xse

import (
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"time"
)

type stream struct {
	*receiveStream
	*sendStream
}

var _ Stream = &stream{}

// need to define StreamID() here, since both receiveStream and readStream have a StreamID()
func (s *stream) StreamID() protocol.StreamID {
	// the result is same for ReceiveStream and SendStream
	return s.sendStream.StreamID()
}

func (s *stream) Close() error {
	return s.sendStream.Close()
}

func (s *stream) SetDeadline(t time.Time) error {
	_ = s.receiveStream.SetReadDeadline(t) // SetReadDeadline never errors
	_ = s.sendStream.SetWriteDeadline(t)   // SetWriteDeadline never errors
	return nil
}

// CloseForShutdown closes a stream abruptly.
// It makes Read and Write unblock (and return the error) immediately.
// The peer will NOT be informed about this: the stream is closed without sending a FIN or RST.
func (s *stream) CloseForShutdown(err error) {
	s.sendStream.CloseForShutdown(err)
	s.receiveStream.CloseForShutdown(err)
}

func (s *stream) ReceiveStream() ReceiveStream {
	return s.receiveStream
}

func (s *stream) SendStream() SendStream {
	return s.sendStream
}
