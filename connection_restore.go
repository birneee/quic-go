package quic

import (
	"github.com/lucas-clemente/quic-go/handover"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/utils"
	"net"
)

// for connection restore from H-QUIC state
type ConnectionRestoreConfig struct {
	Perspective Perspective
	QuicConf    *Config
	PacketConn  net.PacketConn
	// if null a new one is created
	PacketHandlerManager packetHandlerManager
}

func (c *ConnectionRestoreConfig) Populate(state *handover.State) *ConnectionRestoreConfig {
	if c == nil {
		c = &ConnectionRestoreConfig{}
	}
	if c.QuicConf == nil {
		c.QuicConf = &Config{}
	}
	// must be called before populateClientConfig or populateServerConfig
	c.QuicConf.ConnectionIDLength = state.SrcConnectionIDLength(c.Perspective)
	if c.Perspective == protocol.PerspectiveClient {
		c.QuicConf = populateClientConfig(c.QuicConf, c.QuicConf.ConnectionIDLength == 0)
	} else {
		c.QuicConf = populateServerConfig(c.QuicConf)
	}
	ownTransportParams := state.OwnTransportParameters(c.Perspective)
	c.QuicConf.MaxStreamReceiveWindow = utils.MaxUint64V(
		c.QuicConf.MaxStreamReceiveWindow,
		uint64(ownTransportParams.InitialMaxStreamDataBidiLocal),
		uint64(ownTransportParams.InitialMaxStreamDataBidiRemote),
		uint64(ownTransportParams.InitialMaxStreamDataUni),
	)

	c.QuicConf.MaxConnectionReceiveWindow = utils.MaxUint64V(
		c.QuicConf.MaxConnectionReceiveWindow,
		c.QuicConf.MaxStreamReceiveWindow,
		uint64(ownTransportParams.InitialMaxData),
	)
	if c.PacketConn == nil {
		if c.PacketHandlerManager == nil {
			var err error
			c.PacketConn, err = ListenMigratableUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
			if err != nil {
				panic(err)
			}
		} else {
			c.PacketConn = c.PacketHandlerManager.PacketConn()
		}
	}
	if c.PacketHandlerManager == nil {
		var err error
		c.PacketHandlerManager, err = getMultiplexer().AddConn(c.PacketConn, state.SrcConnectionIDLength(c.Perspective), c.QuicConf.StatelessResetKey, c.QuicConf.Tracer)
		if err != nil {
			panic(err)
		}
	}
	return c
}

type RestoredStreams struct {
	BidiStreams    map[StreamID]Stream
	SendStreams    map[StreamID]SendStream
	ReceiveStreams map[StreamID]ReceiveStream
}
