package xads

import (
	"net"
	"time"
)

type streamConnAdapter struct {
	stream Stream
}

var _ net.Conn = &streamConnAdapter{}

func (c *streamConnAdapter) Read(b []byte) (n int, err error) {
	return c.stream.Read(b)
}

func (c *streamConnAdapter) Write(b []byte) (n int, err error) {
	return c.stream.Write(b)
}

func (c *streamConnAdapter) Close() error {
	return c.stream.Close()
}

func (c *streamConnAdapter) LocalAddr() net.Addr {
	panic("unexpected call")
}

func (c *streamConnAdapter) RemoteAddr() net.Addr {
	panic("unexpected call")
}

func (c *streamConnAdapter) SetDeadline(t time.Time) error {
	return c.stream.SetDeadline(t)
}

func (c *streamConnAdapter) SetReadDeadline(t time.Time) error {
	return c.stream.SetReadDeadline(t)
}

func (c *streamConnAdapter) SetWriteDeadline(t time.Time) error {
	return c.stream.SetWriteDeadline(t)
}

type sendStreamConnAdapter struct {
	stream SendStream
}

var _ net.Conn = &sendStreamConnAdapter{}

func (s sendStreamConnAdapter) Read(b []byte) (n int, err error) {
	panic("unexpected call")
}

func (s sendStreamConnAdapter) Write(b []byte) (n int, err error) {
	return s.stream.Write(b)
}

func (s sendStreamConnAdapter) Close() error {
	return s.stream.Close()
}

func (s sendStreamConnAdapter) LocalAddr() net.Addr {
	panic("unexpected call")
}

func (s sendStreamConnAdapter) RemoteAddr() net.Addr {
	panic("unexpected call")
}

func (s sendStreamConnAdapter) SetDeadline(t time.Time) error {
	panic("unexpected call")
}

func (s sendStreamConnAdapter) SetReadDeadline(t time.Time) error {
	panic("unexpected call")
}

func (s sendStreamConnAdapter) SetWriteDeadline(t time.Time) error {
	return s.stream.SetWriteDeadline(t)
}

type receiveStreamConnAdapter struct {
	stream ReceiveStream
}

var _ net.Conn = &receiveStreamConnAdapter{}

func (r receiveStreamConnAdapter) Read(b []byte) (n int, err error) {
	return r.stream.Read(b)
}

func (r receiveStreamConnAdapter) Write(b []byte) (n int, err error) {
	panic("unexpected call")
}

func (r receiveStreamConnAdapter) Close() error {
	panic("unexpected call")
}

func (r receiveStreamConnAdapter) LocalAddr() net.Addr {
	panic("unexpected call")
}

func (r receiveStreamConnAdapter) RemoteAddr() net.Addr {
	panic("unexpected call")
}

func (r receiveStreamConnAdapter) SetDeadline(t time.Time) error {
	panic("unexpected call")
}

func (r receiveStreamConnAdapter) SetReadDeadline(t time.Time) error {
	return r.stream.SetReadDeadline(t)
}

func (r receiveStreamConnAdapter) SetWriteDeadline(t time.Time) error {
	panic("unexpected call")
}
