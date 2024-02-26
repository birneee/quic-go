//go:generate msgp
package handover

import (
	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/wire"
	"github.com/quic-go/quic-go/logging"
	"time"
)

// TransportParameters contains relevant parameters in SMAQ state
type TransportParameters struct {
	// nil if default
	InitialMaxStreamDataBidiLocal *int64
	// nil if default
	InitialMaxStreamDataBidiRemote *int64
	// nil if default
	InitialMaxStreamDataUni *int64
	// nil if default
	MaxAckDelay *int64
	// nil if default
	AckDelayExponent *uint8
	// nil if default
	DisableActiveMigration *bool
	// nil if default
	MaxUDPPayloadSize *int64
	// nil if default
	MaxIdleTimeout *int64
	// nil if client perspective
	OriginalDestinationConnectionID *[]byte
	ActiveConnectionIDLimit         uint64
	// nil if default
	MaxDatagramFrameSize *int64
}

func ToHandoverTransportParameters(p *wire.TransportParameters) TransportParameters {
	s := TransportParameters{
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
	if p.MaxIdleTimeout != 0 {
		ms := int64(p.MaxIdleTimeout / time.Millisecond)
		s.MaxIdleTimeout = &ms
	}
	if p.OriginalDestinationConnectionID.Len() > 0 {
		b := p.OriginalDestinationConnectionID.Bytes()
		s.OriginalDestinationConnectionID = &b
	}
	if p.MaxDatagramFrameSize > 0 {
		s.MaxDatagramFrameSize = (*int64)(&p.MaxDatagramFrameSize)
	}
	return s
}

func RestoreTransportParameters(s *TransportParameters) *wire.TransportParameters {
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
	if s.MaxIdleTimeout != nil {
		p.MaxIdleTimeout = time.Duration(*s.MaxIdleTimeout) * time.Millisecond
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
