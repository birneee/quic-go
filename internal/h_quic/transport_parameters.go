package h_quic

import "github.com/lucas-clemente/quic-go/internal/wire"

// FilterTransportParameters remove some parameters that must not be part of the handover state
func FilterTransportParameters(tp wire.TransportParameters) wire.TransportParameters {
	// extra_stream_encryption
	tp.ExtraStreamEncryption = false
	return tp
}
