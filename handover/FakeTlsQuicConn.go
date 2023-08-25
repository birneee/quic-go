package handover

import (
	"context"
	"crypto/tls"
	"github.com/quic-go/quic-go/internal/qtls"
)

type fakeTlsQuicConn struct {
}

func (f fakeTlsQuicConn) SetTransportParameters(params []byte) {
	//TODO implement me
	panic("implement me")
}

func (f fakeTlsQuicConn) Start(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (f fakeTlsQuicConn) NextEvent() qtls.QUICEvent {
	// do nothing
	return qtls.QUICEvent{Kind: qtls.QUICNoEvent}
}

func (f fakeTlsQuicConn) Close() error {
	// do nothing
	return nil
}

func (f fakeTlsQuicConn) HandleData(level qtls.QUICEncryptionLevel, data []byte) error {
	// do nothing
	return nil
}

func (f fakeTlsQuicConn) SendSessionTicket(earlyData bool) error {
	//TODO implement me
	panic("implement me")
}

func (f fakeTlsQuicConn) ConnectionState() tls.ConnectionState {
	return tls.ConnectionState{
		HandshakeComplete: true,
	}
}

func NewFakeTlsQuicConn() qtls.QUICConn {
	return &fakeTlsQuicConn{}
}
