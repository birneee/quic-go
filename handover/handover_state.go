package handover

import "github.com/lucas-clemente/quic-go/internal/protocol"

type AeadState struct {
	KeyPhase protocol.KeyPhase
	HighestRcvdPN protocol.PacketNumber
	SuiteId uint16
	FirstRcvTrafficSecret []byte
	FirstSendTrafficSecret []byte
	RcvTrafficSecret []byte
	SendTrafficSecret []byte
}

// State is used to handover QUIC connection
type State struct {
	AeadState AeadState
}
