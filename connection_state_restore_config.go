package quic

import (
	"github.com/quic-go/quic-go/handover"
	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/utils"
	"github.com/quic-go/quic-go/logging"
)

// for connection restore from H-QUIC state
type ConnectionRestoreConfig struct {
	Perspective logging.Perspective
	QuicConf    *Config
}

func (c *ConnectionRestoreConfig) Populate(state *handover.State) *ConnectionRestoreConfig {
	if c == nil {
		c = &ConnectionRestoreConfig{}
	}
	if c.QuicConf == nil {
		c.QuicConf = &Config{}
	}

	if c.Perspective == protocol.PerspectiveClient {
		c.QuicConf = populateConfig(c.QuicConf)
	} else {
		c.QuicConf = populateServerConfig(c.QuicConf)
	}
	ownTransportParams := state.FromPerspective(c.Perspective).OwnTransportParameters()
	c.QuicConf.MaxStreamReceiveWindow = utils.MaxV(
		c.QuicConf.MaxStreamReceiveWindow,
		uint64(*ownTransportParams.InitialMaxStreamDataBidiLocal),
		uint64(*ownTransportParams.InitialMaxStreamDataBidiRemote),
		uint64(*ownTransportParams.InitialMaxStreamDataUni),
	)

	c.QuicConf.MaxConnectionReceiveWindow = utils.MaxV(
		c.QuicConf.MaxConnectionReceiveWindow,
		c.QuicConf.MaxStreamReceiveWindow,
	)
	//if c.PacketConn == nil {
	//	if c.PacketHandlerManager == nil {
	//		var err error
	//		c.PacketConn, err = net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	//		if err != nil {
	//			panic(err)
	//		}
	//	} else {
	//		//c.PacketConn = c.PacketHandlerManager.PacketConn()
	//		panic("implement me")
	//	}
	//}
	//if c.PacketHandlerManager == nil {
	//	tr := &Transport{Conn: c.PacketConn, isSingleUse: true, restoredHQUIC: true, ConnectionIDLength: state.ConnIDLen(c.Perspective)}
	//	if err := tr.init(false); err != nil {
	//		panic(err)
	//	}
	//	c.PacketHandlerManager = tr.handlerMap
	//}
	c.QuicConf.EnableDatagrams = ownTransportParams.MaxDatagramFrameSize != nil && *ownTransportParams.MaxDatagramFrameSize != 0
	return c
}

func (c *ConnectionRestoreConfig) Validate() error {
	return nil
}

type RestoredStreams struct {
	BidiStreams    map[StreamID]Stream
	SendStreams    map[StreamID]SendStream
	ReceiveStreams map[StreamID]ReceiveStream
}
