package quic

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/lucas-clemente/quic-go/handover"
	"io"
)

// TODO IANA registration
const HQUICProxyALPN = "qproxy"

// TODO IANA registration
const DefaultHQUICProxyControlPort = 18081

type ProxyControlSession struct {
	session Connection
}

func (p *ProxyControlSession) Close() error {
	return p.session.CloseWithError(ApplicationErrorCode(0), "cancel")
}

func (p *ProxyControlSession) SendHandover(state *handover.State) error {
	//TODO compare Sync vs non-Sync version
	stream, err := p.session.OpenStream()
	if err != nil {
		return fmt.Errorf("failed to open stream: %w", err)
	}
	marshalledState, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal handover state: %w", err)
	}
	_, err = io.Copy(stream, bytes.NewReader(marshalledState))
	if err != nil {
		return fmt.Errorf("failed to send handover state: %w", err)
	}
	_ = stream.Close()
	return nil
}

// DialProxyAddr establish a new HQUIC-Proxy connection
func DialProxyAddr(addr string, tlsConf *tls.Config, config *Config) (*ProxyControlSession, error) {
	if tlsConf == nil {
		tlsConf = &tls.Config{}
	} else {
		tlsConf = tlsConf.Clone()
	}
	tlsConf.NextProtos = []string{HQUICProxyALPN}
	session, err := DialAddrEarly(addr, tlsConf, config)
	if err != nil {
		return nil, fmt.Errorf("proxy connection failed: %w", err)
	}
	return &ProxyControlSession{
		session: session,
	}, nil
}