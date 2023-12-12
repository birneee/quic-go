package quic

import (
	"fmt"
	"github.com/quic-go/quic-go/internal/congestion"
	"net"
	"time"

	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/quicvarint"
)

// Clone clones a Config
func (c *Config) Clone() *Config {
	copy := *c
	return &copy
}

func (c *Config) handshakeTimeout() time.Duration {
	return 2 * c.HandshakeIdleTimeout
}

func (c *Config) maxRetryTokenAge() time.Duration {
	return c.handshakeTimeout()
}

func validateConfig(config *Config) error {
	if config == nil {
		return nil
	}
	const maxStreams = 1 << 60
	if config.MaxIncomingStreams > maxStreams {
		config.MaxIncomingStreams = maxStreams
	}
	if config.MaxIncomingUniStreams > maxStreams {
		config.MaxIncomingUniStreams = maxStreams
	}
	if config.MaxStreamReceiveWindow > quicvarint.Max {
		config.MaxStreamReceiveWindow = quicvarint.Max
	}
	if config.MaxConnectionReceiveWindow > quicvarint.Max {
		config.MaxConnectionReceiveWindow = quicvarint.Max
	}
	// check that all QUIC versions are actually supported
	for _, v := range config.Versions {
		if !protocol.IsValidVersion(v) {
			return fmt.Errorf("invalid QUIC version: %s", v)
		}
	}
	return nil
}

// populateServerConfig populates fields in the quic.Config with their default values, if none are set
// it may be called with nil
func populateServerConfig(config *Config) *Config {
	config = populateConfig(config)
	if config.RequireAddressValidation == nil {
		config.RequireAddressValidation = func(net.Addr) bool { return false }
	}
	if config.QlogLabel == "client" {
		config.QlogLabel = "server"
	}
	return config
}

// populateConfig populates fields in the quic.Config with their default values, if none are set
// it may be called with nil
func populateConfig(config *Config) *Config {
	if config == nil {
		config = &Config{}
	}
	versions := config.Versions
	if len(versions) == 0 {
		versions = protocol.SupportedVersions
	}
	handshakeIdleTimeout := protocol.DefaultHandshakeIdleTimeout
	if config.HandshakeIdleTimeout != 0 {
		handshakeIdleTimeout = config.HandshakeIdleTimeout
	}
	idleTimeout := protocol.DefaultIdleTimeout
	if config.MaxIdleTimeout != 0 {
		idleTimeout = config.MaxIdleTimeout
	}
	initialStreamReceiveWindow := config.InitialStreamReceiveWindow
	if initialStreamReceiveWindow == 0 {
		initialStreamReceiveWindow = protocol.DefaultInitialMaxStreamData
	}
	maxStreamReceiveWindow := config.MaxStreamReceiveWindow
	if maxStreamReceiveWindow == 0 {
		maxStreamReceiveWindow = protocol.DefaultMaxReceiveStreamFlowControlWindow
	}
	initialConnectionReceiveWindow := config.InitialConnectionReceiveWindow
	if initialConnectionReceiveWindow == 0 {
		initialConnectionReceiveWindow = protocol.DefaultInitialMaxData
	}
	maxConnectionReceiveWindow := config.MaxConnectionReceiveWindow
	if maxConnectionReceiveWindow == 0 {
		maxConnectionReceiveWindow = protocol.DefaultMaxReceiveConnectionFlowControlWindow
	}
	maxIncomingStreams := config.MaxIncomingStreams
	if maxIncomingStreams == 0 {
		maxIncomingStreams = protocol.DefaultMaxIncomingStreams
	} else if maxIncomingStreams < 0 {
		maxIncomingStreams = 0
	}
	maxIncomingUniStreams := config.MaxIncomingUniStreams
	if maxIncomingUniStreams == 0 {
		maxIncomingUniStreams = protocol.DefaultMaxIncomingUniStreams
	} else if maxIncomingUniStreams < 0 {
		maxIncomingUniStreams = 0
	}
	initialCongestionWindow := config.InitialCongestionWindow
	if initialCongestionWindow == 0 {
		initialCongestionWindow = protocol.DefaultInitialCongestionWindow
	}
	qlogLabel := config.QlogLabel
	if qlogLabel == "" {
		qlogLabel = "client"
	}
	maxBandwidth := config.MaxBandwidth
	if maxBandwidth == 0 {
		maxBandwidth = congestion.InfBandwidth
	}

	return &Config{
		GetConfigForClient:             config.GetConfigForClient,
		Versions:                       versions,
		HandshakeIdleTimeout:           handshakeIdleTimeout,
		MaxIdleTimeout:                 idleTimeout,
		RequireAddressValidation:       config.RequireAddressValidation,
		KeepAlivePeriod:                config.KeepAlivePeriod,
		InitialStreamReceiveWindow:     initialStreamReceiveWindow,
		MaxStreamReceiveWindow:         maxStreamReceiveWindow,
		InitialConnectionReceiveWindow: initialConnectionReceiveWindow,
		MaxConnectionReceiveWindow:     maxConnectionReceiveWindow,
		AllowConnectionWindowIncrease:  config.AllowConnectionWindowIncrease,
		MaxIncomingStreams:             maxIncomingStreams,
		MaxIncomingUniStreams:          maxIncomingUniStreams,
		TokenStore:                     config.TokenStore,
		EnableDatagrams:                config.EnableDatagrams,
		DisablePathMTUDiscovery:        config.DisablePathMTUDiscovery,
		Allow0RTT:                      config.Allow0RTT,
		Tracer:                         config.Tracer,
		ProxyConf:                      config.ProxyConf,
		AllowEarlyHandover:             config.AllowEarlyHandover,
		QlogLabel:                      qlogLabel,
		DisableQlog:                    config.DisableQlog,
		InitialCongestionWindow:        initialCongestionWindow,
		MaxBandwidth:                   maxBandwidth,
	}
}
