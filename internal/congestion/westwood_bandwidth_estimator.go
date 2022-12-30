package congestion

import (
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/utils"
	"time"
)

type WestwoodBandwidthEstimator struct {
	lastAckTime         time.Time
	lastBandwidthSample Bandwidth
	bandwidthEstimate   Bandwidth
}

func NewWestwoodBandwidthEstimator() *WestwoodBandwidthEstimator {
	bwe := &WestwoodBandwidthEstimator{
		lastAckTime:         time.Now(),
		lastBandwidthSample: Bandwidth(0),
		bandwidthEstimate:   Bandwidth(0),
	}
	return bwe
}

func (w *WestwoodBandwidthEstimator) OnPacketAcked(ackedBytes protocol.ByteCount, ackTime time.Time) {
	delta := ackTime.Sub(w.lastAckTime)
	bandwidthSample := BandwidthFromDelta(ackedBytes, utils.Max(delta, time.Millisecond))

	w.bandwidthEstimate = lowPassFilter(w.bandwidthEstimate, bandwidthSample)

	w.lastAckTime = ackTime
	w.lastBandwidthSample = bandwidthSample
}

func (w *WestwoodBandwidthEstimator) Estimate() Bandwidth {
	return w.bandwidthEstimate
}

func lowPassFilter(bandwidth Bandwidth, sample Bandwidth) Bandwidth {
	return ((7 * bandwidth) + sample) / 8
}
