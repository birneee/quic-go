package congestion

import (
	"fmt"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/utils"
	"github.com/lucas-clemente/quic-go/logging"
	"time"
)

const (
	// RTT0 is the reference RTT
	RTT0        = 25 * time.Millisecond
	BitsPerByte = 8
)

type hyblaWestwoodSender struct {
	pacer                      *pacer
	rttStats                   *utils.RTTStats
	westwoodBandwidthEstimator *WestwoodBandwidthEstimator
	clock                      Clock
	// Congestion window in bytes.
	congestionWindow        protocol.ByteCount
	initialCongestionWindow protocol.ByteCount
	minCongestionWindow     protocol.ByteCount
	maxCongestionWindow     protocol.ByteCount
	maxDatagramSize         protocol.ByteCount
	// Slow start congestion window in bytes, aka ssthresh.
	slowStartThreshold        protocol.ByteCount
	initialSlowStartThreshold protocol.ByteCount
	minSlowStartThreshold     protocol.ByteCount
	maxSlowStartThreshold     protocol.ByteCount
	// Track the largest packet that has been sent.
	largestSentPacketNumber protocol.PacketNumber
	// Track the largest packet that has been acked.
	largestAckedPacketNumber protocol.PacketNumber
	// Track the largest packet number outstanding when a CWND cutback occurs.
	largestSentAtLastCutback protocol.PacketNumber
	bytesInFlight            protocol.ByteCount
	// Whether the last loss event caused us to exit slowstart.
	// Used for stats collection of slowstartPacketsLost
	lastCutbackExitedSlowstart bool
	lastState                  logging.CongestionState
	tracer                     logging.ConnectionTracer
}

var (
	_ SendAlgorithm               = &hyblaWestwoodSender{}
	_ SendAlgorithmWithDebugInfos = &hyblaWestwoodSender{}
)

func NewHyblaWestwoodSender(
	clock Clock,
	rttStats *utils.RTTStats,
	initialMaxDatagramSize protocol.ByteCount,
	initialCongestionWindow protocol.ByteCount,
	minCongestionWindow protocol.ByteCount,
	maxCongestionWindow protocol.ByteCount,
	initialSlowStartThreshold protocol.ByteCount,
	minSlowStartThreshold protocol.ByteCount,
	maxSlowStartThreshold protocol.ByteCount,
	tracer logging.ConnectionTracer,
) *hyblaWestwoodSender {
	h := &hyblaWestwoodSender{
		clock:                      clock,
		rttStats:                   rttStats,
		westwoodBandwidthEstimator: NewWestwoodBandwidthEstimator(),
		congestionWindow:           initialCongestionWindow,
		initialCongestionWindow:    initialCongestionWindow,
		minCongestionWindow:        minCongestionWindow,
		maxCongestionWindow:        maxCongestionWindow,
		maxDatagramSize:            initialMaxDatagramSize,
		slowStartThreshold:         initialSlowStartThreshold,
		initialSlowStartThreshold:  initialSlowStartThreshold,
		minSlowStartThreshold:      minSlowStartThreshold,
		maxSlowStartThreshold:      maxSlowStartThreshold,
		largestSentPacketNumber:    protocol.InvalidPacketNumber,
		largestAckedPacketNumber:   protocol.InvalidPacketNumber,
		largestSentAtLastCutback:   protocol.InvalidPacketNumber,
		tracer:                     tracer,
	}
	h.pacer = newPacer(h.pacerBandwidthEstimate)
	if h.tracer != nil {
		h.lastState = logging.CongestionStateSlowStart
		h.tracer.UpdatedCongestionState(logging.CongestionStateSlowStart)
	}
	return h
}

func (h *hyblaWestwoodSender) TimeUntilSend(bytesInFlight protocol.ByteCount) time.Time {
	return h.pacer.TimeUntilSend()
}

func (h *hyblaWestwoodSender) HasPacingBudget() bool {
	return h.pacer.Budget(h.clock.Now()) >= h.maxDatagramSize
}

func (h *hyblaWestwoodSender) OnPacketSent(sentTime time.Time, bytesInFlight protocol.ByteCount, packetNumber protocol.PacketNumber, bytes protocol.ByteCount, isRetransmittable bool) {
	h.pacer.SentPacket(sentTime, bytes)
	if !isRetransmittable {
		return
	}
	h.largestSentPacketNumber = packetNumber
}

func (h *hyblaWestwoodSender) CanSend(bytesInFlight protocol.ByteCount) bool {
	return bytesInFlight < h.GetCongestionWindow()
}

func (h *hyblaWestwoodSender) MaybeExitSlowStart() {
	return
}

func (h *hyblaWestwoodSender) OnPacketAcked(ackedPacketNumber protocol.PacketNumber, ackedBytes protocol.ByteCount, priorInFlight protocol.ByteCount, eventTime time.Time) {
	h.bytesInFlight = priorInFlight
	h.westwoodBandwidthEstimator.OnPacketAcked(ackedBytes, eventTime)
	h.largestAckedPacketNumber = utils.Max(ackedPacketNumber, h.largestAckedPacketNumber)
	if h.InRecovery() {
		return
	}
	h.maybeIncreaseCwnd(ackedPacketNumber, ackedBytes, priorInFlight, eventTime)
}

func (h *hyblaWestwoodSender) OnPacketLost(packetNumber protocol.PacketNumber, lostBytes protocol.ByteCount, priorInFlight protocol.ByteCount) {
	h.bytesInFlight = priorInFlight
	// TCP NewReno (RFC6582) says that once a loss occurs, any losses in packets
	// already sent should be treated as a single loss event, since it's expected.
	if packetNumber <= h.largestSentAtLastCutback {
		return
	}
	h.lastCutbackExitedSlowstart = h.InSlowStart()
	h.maybeTraceStateChange(logging.CongestionStateRecovery)

	// from https://doi.org/10.1109/NUICONE.2012.6493205 version-3
	bwe := protocol.ByteCount(h.westwoodBandwidthEstimator.Estimate())
	rttMin := protocol.ByteCount(h.rttStats.MinRTT().Milliseconds())
	h.setCongestionWindow(utils.Max(
		bwe*rttMin/1000/BitsPerByte*3/2,
		priorInFlight*3/4,
	))
	h.setSlowStartThreshold(h.GetCongestionWindow())
	h.largestSentAtLastCutback = h.largestSentPacketNumber
}

func (h *hyblaWestwoodSender) OnRetransmissionTimeout(packetsRetransmitted bool) {
	h.largestSentAtLastCutback = protocol.InvalidPacketNumber
	if !packetsRetransmitted {
		return
	}

	// from https://doi.org/10.1109/NUICONE.2012.6493205 version-3
	bwe := protocol.ByteCount(h.westwoodBandwidthEstimator.Estimate())
	rttMin := protocol.ByteCount(h.rttStats.MinRTT().Milliseconds())
	h.setSlowStartThreshold(utils.Max(
		bwe*rttMin/1000/BitsPerByte*3/2,
		h.bytesInFlight*3/4,
	))
	h.setCongestionWindow(protocol.ByteCount(h.rho()) * h.maxDatagramSize)
}

func (h *hyblaWestwoodSender) SetMaxDatagramSize(s protocol.ByteCount) {
	if s < h.maxDatagramSize {
		panic(fmt.Sprintf("congestion BUG: decreased max datagram size from %d to %d", h.maxDatagramSize, s))
	}
	h.maxDatagramSize = s
	h.pacer.SetMaxDatagramSize(s)
}

func (h *hyblaWestwoodSender) OnConnectionMigration() {
	h.largestSentPacketNumber = protocol.InvalidPacketNumber
	h.largestAckedPacketNumber = protocol.InvalidPacketNumber
	h.largestSentAtLastCutback = protocol.InvalidPacketNumber
	h.lastCutbackExitedSlowstart = false
	h.setCongestionWindow(h.initialCongestionWindow)
	h.setSlowStartThreshold(h.initialSlowStartThreshold)
	h.bytesInFlight = 0
}

func (h *hyblaWestwoodSender) InSlowStart() bool {
	return h.GetCongestionWindow() < h.slowStartThreshold
}

func (h *hyblaWestwoodSender) InRecovery() bool {
	return h.largestAckedPacketNumber != protocol.InvalidPacketNumber && h.largestAckedPacketNumber <= h.largestSentAtLastCutback
}

func (h *hyblaWestwoodSender) GetCongestionWindow() protocol.ByteCount {
	return h.congestionWindow
}

func (h *hyblaWestwoodSender) setCongestionWindow(cw protocol.ByteCount) {
	h.congestionWindow = cw
	if h.congestionWindow < h.minCongestionWindow {
		h.congestionWindow = h.minCongestionWindow
	}
	if h.congestionWindow > h.maxCongestionWindow {
		h.congestionWindow = h.maxCongestionWindow
	}
}

// BandwidthEstimate returns the current bandwidth estimate
func (h *hyblaWestwoodSender) pacerBandwidthEstimate() Bandwidth {
	srtt := h.rttStats.SmoothedRTT()
	if srtt == 0 {
		// If we haven't measured an rtt, the bandwidth estimate is unknown.
		return infBandwidth
	}
	return BandwidthFromDelta(h.GetCongestionWindow(), srtt)
}

// rho = RTT/RTT0
// defined by https://doi.org/10.1002/sat.799
func (h *hyblaWestwoodSender) rho() float64 {
	return h.rttStats.SmoothedRTT().Seconds() / RTT0.Seconds()
}

func (h *hyblaWestwoodSender) setSlowStartThreshold(sst protocol.ByteCount) {
	h.slowStartThreshold = sst
	if h.slowStartThreshold < h.minSlowStartThreshold {
		h.slowStartThreshold = h.minSlowStartThreshold
	}
	if h.slowStartThreshold > h.maxSlowStartThreshold {
		h.slowStartThreshold = h.maxSlowStartThreshold
	}
}

func (h *hyblaWestwoodSender) maybeTraceStateChange(new logging.CongestionState) {
	if h.tracer == nil || new == h.lastState {
		return
	}
	h.tracer.UpdatedCongestionState(new)
	h.lastState = new
}

// Called when we receive an ack. Normal TCP tracks how many packets one ack
// represents, but quic has a separate ack for each packet.
func (h *hyblaWestwoodSender) maybeIncreaseCwnd(
	_ protocol.PacketNumber,
	ackedBytes protocol.ByteCount,
	priorInFlight protocol.ByteCount,
	eventTime time.Time,
) {
	// Do not increase the congestion window unless the sender is close to using
	// the current window.
	if !h.isCwndLimited(priorInFlight) {
		h.maybeTraceStateChange(logging.CongestionStateApplicationLimited)
		return
	}
	if h.congestionWindow >= h.maxCongestionWindow {
		return
	}
	rho := protocol.ByteCount(h.rho())

	if h.InSlowStart() {
		h.maybeTraceStateChange(logging.CongestionStateSlowStart)
		// inspired by https://doi.org/10.1002/sat.799
		// reduce aggressiveness by using rho^2 instead of 2^rho-1
		h.setCongestionWindow(h.GetCongestionWindow() + rho*rho*ackedBytes)
		return
	}
	// Congestion avoidance
	h.maybeTraceStateChange(logging.CongestionStateCongestionAvoidance)
	// inspired by https://doi.org/10.1002/sat.799
	// reduce aggressiveness by using rho instead of rho^2
	h.setCongestionWindow(h.GetCongestionWindow() + rho*ackedBytes*ackedBytes/h.GetCongestionWindow())
}

func (h *hyblaWestwoodSender) isCwndLimited(bytesInFlight protocol.ByteCount) bool {
	congestionWindow := h.GetCongestionWindow()
	if bytesInFlight >= congestionWindow {
		return true
	}
	availableBytes := congestionWindow - bytesInFlight
	slowStartLimited := h.InSlowStart() && bytesInFlight > congestionWindow/2
	return slowStartLimited || availableBytes <= maxBurstPackets*h.maxDatagramSize
}
