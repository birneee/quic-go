//go:generate msgp
package qstate

import (
	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/wire"
	"github.com/quic-go/quic-go/logging"
	"time"
)

type Parameters struct {
	// nil if default
	InitialMaxStreamDataBidiLocal *int64 `msg:"initial_max_stream_data_bidi_local,omitempty" json:"initial_max_stream_data_bidi_local,omitempty" cbor:"1,keyasint,omitempty"`
	// nil if default
	InitialMaxStreamDataBidiRemote *int64 `msg:"initial_max_stream_data_bidi_remote,omitempty" json:"initial_max_stream_data_bidi_remote,omitempty" cbor:"2,keyasint,omitempty"`
	// nil if default
	InitialMaxStreamDataUni *int64 `msg:"initial_max_stream_data_uni,omitempty" json:"initial_max_stream_data_uni,omitempty" cbor:"3,keyasint,omitempty"`
	// nil if default
	MaxAckDelay *int64 `msg:"max_ack_delay,omitempty" json:"max_ack_delay,omitempty" cbor:"4,keyasint,omitempty"`
	// nil if default
	AckDelayExponent *uint8 `msg:"ack_delay_exponent,omitempty" json:"ack_delay_exponent,omitempty" cbor:"5,keyasint,omitempty"`
	// nil if default
	DisableActiveMigration *bool `msg:"disable_active_migration,omitempty" json:"disable_active_migration,omitempty" cbor:"6,keyasint,omitempty"`
	// nil if default
	MaxUDPPayloadSize *int64 `msg:"max_udp_payload_size,omitempty" json:"max_udp_payload_size,omitempty" cbor:"7,keyasint,omitempty"`
	// nil if client perspective
	OriginalDestinationConnectionID *HexByteSlice `msg:"original_destination_connection_id,omitempty" json:"original_destination_connection_id,omitempty" cbor:"8,keyasint,omitempty"`
	ActiveConnectionIDLimit         uint64        `msg:"active_connection_id_limit,omitempty" json:"active_connection_id_limit,omitempty" cbor:"9,keyasint,omitempty"`
	// nil if default
	MaxDatagramFrameSize *int64 `msg:"max_datagram_frame_size,omitempty" json:"max_datagram_frame_size,omitempty" cbor:"10,keyasint,omitempty"`
}

func ToQStateParameters(p *wire.TransportParameters) Parameters {
	s := Parameters{
		ActiveConnectionIDLimit: p.ActiveConnectionIDLimit,
	}
	if p.InitialMaxStreamDataBidiLocal != 0 {
		s.InitialMaxStreamDataBidiLocal = (*int64)(&p.InitialMaxStreamDataBidiLocal)
	}
	if p.InitialMaxStreamDataBidiRemote != 0 {
		s.InitialMaxStreamDataBidiRemote = (*int64)(&p.InitialMaxStreamDataBidiRemote)
	}
	if p.InitialMaxStreamDataUni != 0 {
		s.InitialMaxStreamDataUni = (*int64)(&(p.InitialMaxStreamDataUni))
	}
	if p.MaxAckDelay != protocol.DefaultMaxAckDelay {
		ms := int64(p.MaxAckDelay / time.Millisecond)
		s.MaxAckDelay = &ms
	}
	if p.AckDelayExponent != protocol.DefaultAckDelayExponent {
		s.AckDelayExponent = &p.AckDelayExponent
	}
	if p.DisableActiveMigration {
		s.DisableActiveMigration = &p.DisableActiveMigration
	}
	if p.MaxUDPPayloadSize < 65527 {
		s.MaxUDPPayloadSize = (*int64)(&p.MaxUDPPayloadSize)
	}
	if p.OriginalDestinationConnectionID.Len() > 0 {
		b := HexByteSlice(p.OriginalDestinationConnectionID.Bytes())
		s.OriginalDestinationConnectionID = &b
	}
	if p.MaxDatagramFrameSize > 0 {
		s.MaxDatagramFrameSize = (*int64)(&p.MaxDatagramFrameSize)
	}
	return s
}

func RestoreTransportParameters(s *Parameters) *wire.TransportParameters {
	p := &wire.TransportParameters{}
	p.AckDelayExponent = protocol.DefaultAckDelayExponent
	p.MaxAckDelay = protocol.DefaultMaxAckDelay
	p.MaxDatagramFrameSize = protocol.InvalidByteCount
	if s.InitialMaxStreamDataBidiLocal != nil {
		p.InitialMaxStreamDataBidiLocal = logging.ByteCount(*s.InitialMaxStreamDataBidiLocal)
	}
	if s.InitialMaxStreamDataBidiRemote != nil {
		p.InitialMaxStreamDataBidiRemote = protocol.ByteCount(*s.InitialMaxStreamDataBidiLocal)
	}
	if s.InitialMaxStreamDataUni != nil {
		p.InitialMaxStreamDataUni = protocol.ByteCount(*s.InitialMaxStreamDataUni)
	}
	if s.MaxAckDelay != nil {
		p.MaxAckDelay = time.Duration(*s.MaxAckDelay) * time.Millisecond
	}
	if s.AckDelayExponent != nil {
		p.AckDelayExponent = *s.AckDelayExponent
	}
	if s.DisableActiveMigration != nil {
		p.DisableActiveMigration = *s.DisableActiveMigration
	}
	if s.MaxUDPPayloadSize != nil {
		p.MaxUDPPayloadSize = protocol.ByteCount(*s.MaxUDPPayloadSize)
	}
	if s.OriginalDestinationConnectionID != nil {
		p.OriginalDestinationConnectionID = protocol.ParseConnectionID(*s.OriginalDestinationConnectionID)
	}
	p.ActiveConnectionIDLimit = s.ActiveConnectionIDLimit
	if s.MaxDatagramFrameSize != nil {
		p.MaxDatagramFrameSize = protocol.ByteCount(*s.MaxDatagramFrameSize)
	}
	return p
}
