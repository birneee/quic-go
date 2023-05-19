package quic

import (
	"crypto/tls"
)

type StateTransferConfig struct {
	QuicConfig *Config
	TlsConfig  *tls.Config
}

func (c *StateTransferConfig) Populate() *StateTransferConfig {
	if c == nil {
		c = &StateTransferConfig{}
	}
	if c.QuicConfig == nil {
		c.QuicConfig = &Config{}
	}
	if c.TlsConfig == nil {
		c.TlsConfig = &tls.Config{}
	}
	if c.TlsConfig.NextProtos == nil {
		c.TlsConfig.NextProtos = []string{HQUICStateTransferALPN}
	}
	return c
}
