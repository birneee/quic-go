package ackhandler

import (
	"fmt"
	"github.com/quic-go/quic-go/internal/utils"
	"github.com/quic-go/quic-go/internal/wire"
	"github.com/quic-go/quic-go/qstate"
	"reflect"
	"sync"
	"time"

	"github.com/quic-go/quic-go/internal/protocol"
)

// A Packet is a packet
type packet struct {
	SendTime        time.Time
	PacketNumber    protocol.PacketNumber
	StreamFrames    []StreamFrame
	Frames          []Frame
	LargestAcked    protocol.PacketNumber // InvalidPacketNumber if the packet doesn't contain an ACK
	Length          protocol.ByteCount
	EncryptionLevel protocol.EncryptionLevel

	IsPathMTUProbePacket bool // We don't report the loss of Path MTU probe packets to the congestion controller.

	includedInBytesInFlight bool
	declaredLost            bool
	skippedPacket           bool
}

func (p *packet) outstanding() bool {
	return !p.declaredLost && !p.skippedPacket && !p.IsPathMTUProbePacket
}

var packetPool = sync.Pool{New: func() any { return &packet{} }}

func getPacket() *packet {
	p := packetPool.Get().(*packet)
	p.PacketNumber = 0
	p.StreamFrames = nil
	p.Frames = nil
	p.LargestAcked = 0
	p.Length = 0
	p.EncryptionLevel = protocol.EncryptionLevel(0)
	p.SendTime = time.Time{}
	p.IsPathMTUProbePacket = false
	p.includedInBytesInFlight = false
	p.declaredLost = false
	p.skippedPacket = false
	return p
}

// We currently only return Packets back into the pool when they're acknowledged (not when they're lost).
// This simplifies the code, and gives the vast majority of the performance benefit we can gain from using the pool.
func putPacket(p *packet) {
	p.Frames = nil
	p.StreamFrames = nil
	packetPool.Put(p)
}

func (p *packet) Qstate() qstate.Packet {
	ps := qstate.Packet{
		PacketNumber: int64(p.PacketNumber),
	}
	for _, frame := range p.Frames {
		switch f := frame.Frame.(type) {
		case *wire.HandshakeDoneFrame:
			ps.Frames = append(ps.Frames, qstate.Frame{Type: "handshake_done"})
		case *wire.NewTokenFrame:
			ps.Frames = append(ps.Frames, qstate.Frame{Type: "new_token", Token: f.Token})
		case *wire.CryptoFrame:
			ps.Frames = append(ps.Frames, qstate.Frame{Type: "crypto", Offset: utils.New(int64(f.Offset)), Data: f.Data})
		case *wire.PingFrame:
			ps.Frames = append(ps.Frames, qstate.Frame{Type: "ping"})
		case *wire.RetireConnectionIDFrame:
			ps.Frames = append(ps.Frames, qstate.Frame{Type: "retire_connection_id", SequenceNumber: utils.New(f.SequenceNumber)})
		case *wire.NewConnectionIDFrame:
			ps.Frames = append(ps.Frames, qstate.Frame{Type: "new_connection_id", SequenceNumber: utils.New(f.SequenceNumber)}) // the connection_id, stateless_reset_token, retire_prior_to is already part of the transport state
		case *wire.MaxDataFrame:
			ps.Frames = append(ps.Frames, qstate.Frame{Type: "max_data"})
		case *wire.MaxStreamDataFrame:
			ps.Frames = append(ps.Frames, qstate.Frame{Type: "max_stream_data", StreamID: utils.New(int64(f.StreamID))})
		case *wire.StreamDataBlockedFrame:
			ps.Frames = append(ps.Frames, qstate.Frame{Type: "stream_data_blocked", StreamID: utils.New(int64(f.StreamID))}) // max_stream_data is already part of the stream state
		case *wire.MaxStreamsFrame:
			var streamType string
			switch f.Type {
			case protocol.StreamTypeBidi:
				streamType = "bidirectional"
			case protocol.StreamTypeUni:
				streamType = "unidirectional"
			}
			ps.Frames = append(ps.Frames, qstate.Frame{Type: "max_streams", StreamType: streamType}) // current value is already part of the transport state
		case *wire.DatagramFrame:
			// ignore
		default:
			panic(fmt.Sprintf("unexpected frame type: %s", reflect.ValueOf(frame.Frame).Type().String()))
		}
	}
	for _, f := range p.StreamFrames {
		ps.Frames = append(ps.Frames, qstate.Frame{Type: "stream", StreamID: utils.New(int64(f.Frame.StreamID)), Offset: utils.New(int64(f.Frame.Offset)), Length: utils.New(int64(f.Frame.DataLen()))})
	}
	return ps
}
