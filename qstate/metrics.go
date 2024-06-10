//go:generate msgp
package qstate

type Metrics struct {
	// in byte
	CongestionWindow *int64 `msg:"congestion_window,omitempty" json:"congestion_window,omitempty" cbor:"1,keyasint,omitempty"`
	// in ms
	SmoothedRTT *int64 `msg:"smoothed_rtt,omitempty" json:"smoothed_rtt,omitempty" cbor:"2,keyasint,omitempty"`
}
