package ackhandler

import (
	"github.com/quic-go/quic-go/handover"
	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/utils"
	"github.com/quic-go/quic-go/logging"
	"github.com/quic-go/quic-go/qstate"
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
	return sph, newReceivedPacketHandler(sph, logger)
}

func StoreAckHandler(state *qstate.Connection, config *handover.ConnectionStateStoreConf, sph SentPacketHandler, rph ReceivedPacketHandler) {
	sph.StoreState(state, config)
	rph.Store(state)
}

func RestoreAckHandler(
	state *qstate.Connection,
	initialMaxDatagramSize protocol.ByteCount,
	rttStats *utils.RTTStats,
	enableECN bool,
	tracer *logging.ConnectionTracer,
	logger utils.Logger,
) (SentPacketHandler, ReceivedPacketHandler) {
	sph := restoreSendPacketHandler(state, initialMaxDatagramSize, rttStats, enableECN, tracer, logger)
	rph := restoreReceivedPacketHandler(state, sph, logger)

	return sph, rph
}
