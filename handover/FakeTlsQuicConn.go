package handover

import (
	"context"
	"crypto/tls"
	"github.com/quic-go/quic-go/internal/qtls"
)

type fakeTlsQuicConn struct {
	connectionState tls.ConnectionState
}

func (f fakeTlsQuicConn) SetTransportParameters(params []byte) {
	//TODO implement me
	panic("implement me")
}

func (f fakeTlsQuicConn) Start(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (f fakeTlsQuicConn) NextEvent() tls.QUICEvent {
	// do nothing
	return tls.QUICEvent{Kind: tls.QUICNoEvent}
}

func (f fakeTlsQuicConn) Close() error {
	// do nothing
	return nil
}

func (f fakeTlsQuicConn) HandleData(level tls.QUICEncryptionLevel, data []byte) error {
	// do nothing
	return nil
}

func (f fakeTlsQuicConn) SendSessionTicket(opts tls.QUICSessionTicketOptions) error {
	//TODO implement me
	panic("implement me")
}

func (f fakeTlsQuicConn) ConnectionState() tls.ConnectionState {
	return f.connectionState
}

func NewFakeTlsQuicConn(alpn string) qtls.QUICConn {
	return &fakeTlsQuicConn{
		connectionState: tls.ConnectionState{
			HandshakeComplete:  true,
			NegotiatedProtocol: alpn,
		},
	}
}
