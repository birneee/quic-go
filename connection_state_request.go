package quic

import (
	"github.com/quic-go/quic-go/handover"
	"github.com/quic-go/quic-go/qstate"
)

type HandoverStateRequest struct {
	Destroy bool
	Return  chan HandoverStateResponse
	Config  *handover.ConnectionStateStoreConf
}

type HandoverStateResponse struct {
	// is invalid if error
	State qstate.Connection
	Error error
	Early bool
}

type ProxySetupResponse struct {
	Error error
	Early bool
}
