package quic

import (
	"context"
	"github.com/quic-go/quic-go/handover"
	"github.com/quic-go/quic-go/internal/utils"
	"github.com/quic-go/quic-go/logging"
	"net/netip"
	"time"
)

// for connection restore from H-QUIC state
type ConnectionRestoreConfig struct {
	Perspective                    logging.Perspective
	LocalAddr                      netip.Addr
	Tracer                         func(context.Context, logging.Perspective, ConnectionID) *logging.ConnectionTracer
	InitialCongestionWindow        uint32
	InitialConnectionReceiveWindow uint64
	InitialStreamReceiveWindow     uint64
	MaxConnectionReceiveWindow     uint64
	MaxIncomingStreams             int64
	MaxStreamReceiveWindow         uint64
	// DefaultRTT is set when state does not contain RTT
	DefaultRTT     *time.Duration
	MaxIdleTimeout time.Duration
}

func (restoreConf *ConnectionRestoreConfig) GenerateQuicConf(state *handover.State) (*Config, *handover.State) {
	if state.RTT == nil && restoreConf.DefaultRTT != nil {
		ms := restoreConf.DefaultRTT.Milliseconds()
		state.RTT = &ms
	}
	if restoreConf == nil {
		restoreConf = &ConnectionRestoreConfig{}
	}
	ownTransportParams := state.FromPerspective(restoreConf.Perspective).OwnTransportParameters()
	quicConf := &Config{
		InitialCongestionWindow:        restoreConf.InitialCongestionWindow,
		InitialConnectionReceiveWindow: restoreConf.InitialConnectionReceiveWindow,
		InitialStreamReceiveWindow:     restoreConf.InitialStreamReceiveWindow,
		MaxConnectionReceiveWindow:     restoreConf.MaxConnectionReceiveWindow,
		MaxIncomingStreams:             restoreConf.MaxIncomingStreams,
		MaxStreamReceiveWindow:         restoreConf.MaxStreamReceiveWindow,
		Tracer:                         restoreConf.Tracer,
		EnableDatagrams:                ownTransportParams.MaxDatagramFrameSize != nil && *ownTransportParams.MaxDatagramFrameSize != 0,
		MaxIdleTimeout:                 restoreConf.MaxIdleTimeout,
	}
	quicConf = populateConfig(quicConf)
	quicConf.MaxStreamReceiveWindow = utils.MaxV(
		quicConf.MaxStreamReceiveWindow,
		uint64(*ownTransportParams.InitialMaxStreamDataBidiLocal),
		uint64(*ownTransportParams.InitialMaxStreamDataBidiRemote),
		uint64(*ownTransportParams.InitialMaxStreamDataUni),
	)

	quicConf.MaxConnectionReceiveWindow = utils.MaxV(
		quicConf.MaxConnectionReceiveWindow,
		quicConf.MaxStreamReceiveWindow,
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
	return quicConf, state
}

func (c *ConnectionRestoreConfig) Validate() error {
	return nil
}

func (restoreConf *ConnectionRestoreConfig) Populate() *ConnectionRestoreConfig {
	if restoreConf == nil {
		restoreConf = &ConnectionRestoreConfig{}
	}
	return restoreConf
}

type RestoredStreams struct {
	BidiStreams    map[StreamID]Stream
	SendStreams    map[StreamID]SendStream
	ReceiveStreams map[StreamID]ReceiveStream
}
