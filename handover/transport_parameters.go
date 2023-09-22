//go:generate msgp
package handover

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
