package ackhandler

import (
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/utils"
	"github.com/lucas-clemente/quic-go/logging"
)

// NewAckHandler creates a new SentPacketHandler and a new ReceivedPacketHandler.
// peerAddressValidated defines weather the address was validated beforehand by an address validation token.
func NewAckHandler(
	initialPacketNumber protocol.PacketNumber,
	initialMaxDatagramSize protocol.ByteCount,
	initialCongestionWindow uint32,
	minCongestionWindow uint32,
	maxCongestionWindow uint32,
	rttStats *utils.RTTStats,
	pers protocol.Perspective,
	tracer logging.ConnectionTracer,
	logger utils.Logger,
	version protocol.VersionNumber,
	peerAddressValidated bool,
) (SentPacketHandler, ReceivedPacketHandler) {
	sph := newSentPacketHandler(initialPacketNumber, initialMaxDatagramSize, initialCongestionWindow, minCongestionWindow, maxCongestionWindow, rttStats, pers, tracer, logger, peerAddressValidated)
	return sph, newReceivedPacketHandler(sph, rttStats, logger, version)
}
