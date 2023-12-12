package ackhandler

import (
	"github.com/quic-go/quic-go/handover"
	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/utils"
	"github.com/quic-go/quic-go/logging"
)

// NewAckHandler creates a new SentPacketHandler and a new ReceivedPacketHandler.
// clientAddressValidated indicates whether the address was validated beforehand by an address validation token.
// clientAddressValidated has no effect for a client.
func NewAckHandler(
	initialPacketNumber protocol.PacketNumber,
	initialMaxDatagramSize protocol.ByteCount,
	rttStats *utils.RTTStats,
	clientAddressValidated bool,
	enableECN bool,
	pers protocol.Perspective,
	tracer *logging.ConnectionTracer,
	logger utils.Logger,
) (SentPacketHandler, ReceivedPacketHandler) {
	sph := newSentPacketHandler(initialPacketNumber, initialMaxDatagramSize, rttStats, clientAddressValidated, enableECN, pers, tracer, logger)
	return sph, newReceivedPacketHandler(sph, rttStats, logger)
}

var RestorePacketNumberSkip protocol.PacketNumber = 10000

func RestoreAckHandler(
	state *handover.StateFromPerspective,
	initialMaxDatagramSize protocol.ByteCount,
	rttStats *utils.RTTStats,
	enableECN bool,
	tracer *logging.ConnectionTracer,
	logger utils.Logger,
) (SentPacketHandler, ReceivedPacketHandler) {
	sph, rph := NewAckHandler(
		0,
		initialMaxDatagramSize,
		rttStats,
		true, // TODO path challenge
		enableECN,
		state.Perspective(),
		tracer,
		logger,
	)
	sph.DropPackets(protocol.EncryptionInitial)
	sph.DropPackets(protocol.EncryptionHandshake)

	// skip some packets for two reasons:
	//  - this number might be an estimate from the opposite perspective
	//  - some packets might be sent during the handshake
	sph.SetHighest1RTTPacketNumber(state.HighestSentPacketNumber() + RestorePacketNumberSkip)

	return sph, rph
}
