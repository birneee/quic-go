//go:generate msgp
package handover

import "github.com/quic-go/quic-go/internal/protocol"

type Frame struct {
	Type     string             `msg:"frame_type" json:"frame_type"`
	StreamID protocol.StreamID  `msg:"stream_id,omitzero" json:"stream_id,omitzero"`
	Offset   protocol.ByteCount `msg:"offset,omitzero" json:"offset,omitzero"`
	Length   protocol.ByteCount `msg:"length,omitzero" json:"length,omitzero"`
	Token    []byte             `msg:"token,omitempty" json:"token,omitempty"`
	Data     []byte             `msg:"data,omitempty" json:"data,omitempty"`
}

type PacketState struct {
	PacketNumber int64   `msg:"packet_number" json:"packet_number"`
	Frames       []Frame `msg:"frames" json:"frames"`
}
