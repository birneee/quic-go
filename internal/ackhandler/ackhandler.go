package ackhandler

import (
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/utils"
	"github.com/lucas-clemente/quic-go/logging"
	"time"
)

// NewAckHandler creates a new SentPacketHandler and a new ReceivedPacketHandler.
// clientAddressValidated indicates whether the address was validated beforehand by an address validation token.
// clientAddressValidated has no effect for a client.
func NewAckHandler(
	initialPacketNumber protocol.PacketNumber,
	initialMaxDatagramSize protocol.ByteCount,
	initialCongestionWindow uint32, // number of packets
	minCongestionWindow uint32, // number of packets
	maxCongestionWindow uint32, // number of packets
	initialSlowStartThreshold protocol.ByteCount,
	minSlowStartThreshold protocol.ByteCount,
	maxSlowStartThreshold protocol.ByteCount,
	rttStats *utils.RTTStats,
	clientAddressValidated bool,
	pers protocol.Perspective,
	hyblaWestwood bool,
	fixedPTO time.Duration,
	tracer logging.ConnectionTracer,
	logger utils.Logger,
	version protocol.VersionNumber,
) (SentPacketHandler, ReceivedPacketHandler) {
	sph := newSentPacketHandler(initialPacketNumber, initialMaxDatagramSize, initialCongestionWindow, minCongestionWindow, maxCongestionWindow, initialSlowStartThreshold, minSlowStartThreshold, maxSlowStartThreshold, rttStats, clientAddressValidated, pers, hyblaWestwood, fixedPTO, tracer, logger)
	return sph, newReceivedPacketHandler(sph, rttStats, logger, version)
}
