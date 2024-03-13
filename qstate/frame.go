//go:generate msgp
package qstate

type Frame struct {
	Type     string `msg:"frame_type" json:"frame_type"`
	StreamID *int64 `msg:"stream_id,omitempty" json:"stream_id,omitempty"`
	Offset   *int64 `msg:"offset,omitempty" json:"offset,omitempty"`
	Length   *int64 `msg:"length,omitempty" json:"length,omitempty"`
	//TODO msgp omitempty seems not to work for custom types
	Token          HexByteSlice `msg:"token,omitempty" json:"token,omitempty"`
	Data           []byte       `msg:"data,omitempty" json:"data,omitempty"`
	SequenceNumber *uint64      `msg:"sequence_number,omitempty" json:"sequence_number,omitempty"`
	// bidirectional or unidirectional
	StreamType string `msg:"stream_type,omitempty" json:"stream_type,omitempty"`
}
