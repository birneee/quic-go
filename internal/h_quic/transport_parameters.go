package h_quic

import "github.com/quic-go/quic-go/internal/wire"

// FilterTransportParameters remove some parameters that must not be part of the handover state
func FilterTransportParameters(tp wire.TransportParameters) wire.TransportParameters {
	//TODO filter e.g. extra_application_data_security
	return tp
}
