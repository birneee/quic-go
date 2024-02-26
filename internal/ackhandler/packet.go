package ackhandler

import (
	"bytes"
	"github.com/quic-go/quic-go/handover"
	"github.com/quic-go/quic-go/internal/wire"
	"github.com/quic-go/quic-go/quicvarint"
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

func (p *packet) PacketState() handover.PacketState {
	frameTypes := make([]uint64, 0)
	for _, frame := range p.Frames {
		frameTypes = append(frameTypes, TypeOfFrame(frame.Frame))
	}
	//TODO add stream frames in more compact form
	//for _ = range p.StreamFrames {
	//	frameTypes = append(frameTypes, 8)
	//}
	return handover.PacketState{
		Number: p.PacketNumber,
		Frames: frameTypes,
	}
}

func TypeOfFrame(f wire.Frame) uint64 {
	buf := make([]byte, 0, 2000)
	buf2, err := f.Append(buf, protocol.Version2) //TODO use actual version
	if err != nil {
		panic(err)
	}
	typ, err := quicvarint.Read(bytes.NewReader(buf2))
	if err != nil {
		panic(err)
	}
	return typ
}
