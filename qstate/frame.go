//go:generate msgp
package qstate

type Frame struct {
	Type     string `msg:"frame_type" json:"frame_type" cbor:"1,keyasint"`
	StreamID *int64 `msg:"stream_id,omitempty" json:"stream_id,omitempty" cbor:"2,keyasint,omitempty"`
	Offset   *int64 `msg:"offset,omitempty" json:"offset,omitempty" cbor:"3,keyasint,omitempty"`
	Length   *int64 `msg:"length,omitempty" json:"length,omitempty" cbor:"4,keyasint,omitempty"`
	//TODO msgp omitempty seems not to work for custom types
	Token          HexByteSlice `msg:"token,omitempty" json:"token,omitempty" cbor:"5,keyasint,omitempty"`
	Data           []byte       `msg:"data,omitempty" json:"data,omitempty" cbor:"6,keyasint,omitempty"`
	SequenceNumber *uint64      `msg:"sequence_number,omitempty" json:"sequence_number,omitempty" cbor:"7,keyasint,omitempty"`
	// bidirectional or unidirectional
	StreamType string `msg:"stream_type,omitempty" json:"stream_type,omitempty" cbor:"8,keyasint,omitempty"`
}
