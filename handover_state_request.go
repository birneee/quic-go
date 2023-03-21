package quic

import "github.com/lucas-clemente/quic-go/handover"

type HandoverStateRequest struct {
	Destroy bool
	Return  chan HandoverStateResponse
	Config  *ConnectionStateStoreConf
}

type HandoverStateResponse struct {
	// is invalid if error
	State handover.State
	Error error
	Early bool
}

type ProxySetupResponse struct {
	Error error
	Early bool
}
