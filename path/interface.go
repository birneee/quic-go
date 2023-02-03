package path

import (
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/wire"
	"net"
)

type PathManager interface {
	CurrentSendPath() Path
	OnReceiveNonProbingPacket(addr *net.UDPAddr)
	OnReceivePathResponseFrame(addr *net.UDPAddr, frame *wire.PathResponseFrame)
	// GetOrCreatePath does not initiate migration.
	// migration is initiated once a non-probing packet is received on that path
	GetOrCreatePath(addr net.Addr) Path
}

type Path interface {
	PeerCompletedAddressValidation() bool
	SetPeerCompletedAddressValidation()
	IncrementBytesReceived(n protocol.ByteCount)
	PeerAddressValidated() bool
	SetPeerAddressValidated()
	IncrementBytesSent(n protocol.ByteCount)
	BytesReceived() protocol.ByteCount
	BytesSent() protocol.ByteCount
	Addr() net.Addr
	ChallengeData() [8]byte
	BytesInFlight() protocol.ByteCount
	RemoveFromBytesInFlight(length protocol.ByteCount)
	AddToBytesInFlight(length protocol.ByteCount)
	ResetBytesInFlight()
}
