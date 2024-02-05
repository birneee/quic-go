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

func StoreAckHandler(state handover.StateFromPerspective, config *handover.ConnectionStateStoreConf, sph SentPacketHandler, rph ReceivedPacketHandler) {
	state.SetHighestSentPacketNumber(sph.Highest1RTTPacketNumber())
	state.SetHighestReceivedPacketNumber(rph.Highest1RTTPacketNumber()) // higher packet numbers might be in flight
	if config.IncludeCongestionState {
		sph.StoreState(state)
	}
	rph.Store(state)
}

func RestoreAckHandler(
	state handover.StateFromPerspective,
	initialMaxDatagramSize protocol.ByteCount,
	rttStats *utils.RTTStats,
	enableECN bool,
	tracer *logging.ConnectionTracer,
	logger utils.Logger,
) (SentPacketHandler, ReceivedPacketHandler) {
	sph := restoreSendPacketHandler(state, initialMaxDatagramSize, rttStats, enableECN, tracer, logger)
	rph := newReceivedPacketHandler(sph, rttStats, logger)

	return sph, rph
}
