package quic

import (
	"context"
	"net"
)

type StateTransferServer interface {
	Accept(ctx context.Context) (StateTransferConnection, error)
	Addr() net.Addr
	Close() error
}

type transferServer struct {
	quicServer EarlyListener
}

var _ StateTransferServer = &transferServer{}

func ListenStateTransfer(addr net.Addr, config *StateTransferConfig) (StateTransferServer, error) {
	config = config.Populate()

	quicServer, err := ListenAddrEarly(addr.String(), config.TlsConfig, config.QuicConfig)
	if err != nil {
		return nil, err
	}
	ts := &transferServer{
		quicServer: quicServer,
	}
	return ts, nil
}

func (s *transferServer) Accept(ctx context.Context) (StateTransferConnection, error) {
	quicConn, err := s.quicServer.Accept(ctx)
	if err != nil {
		return nil, err
	}
	return NewStateTransferConnection(quicConn), nil
}

func (s *transferServer) Addr() net.Addr {
	return s.quicServer.Addr()
}

func (s *transferServer) Close() error {
	return s.quicServer.Close()
}
