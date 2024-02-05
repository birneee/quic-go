package wire

import (
	"github.com/quic-go/quic-go/internal/protocol"
)

// AckRange is an ACK range
type AckRange struct {
	Smallest protocol.PacketNumber
	Largest  protocol.PacketNumber
}

// Len returns the number of packets contained in this ACK range
func (r AckRange) Len() protocol.PacketNumber {
	return r.Largest - r.Smallest + 1
}

func AckRangesTo2DList(ackRanges []AckRange) [][2]int64 {
	list := make([][2]int64, len(ackRanges))
	for i, v := range ackRanges {
		list[i][0] = int64(v.Smallest)
		list[i][1] = int64(v.Largest)
	}
	return list
}
