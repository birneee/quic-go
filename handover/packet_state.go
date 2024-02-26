//go:generate msgp
package handover

import "github.com/quic-go/quic-go/internal/protocol"

type PacketState struct {
	Frames []uint64
	Number protocol.PacketNumber
}
