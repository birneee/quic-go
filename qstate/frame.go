//go:generate msgp
package qstate

type Frame struct {
	Type           string    `msg:"frame_type" json:"frame_type"`
	StreamID       *int64    `msg:"stream_id,omitempty" json:"stream_id,omitempty"`
	Offset         *int64    `msg:"offset,omitempty" json:"offset,omitempty"`
	Length         *int64    `msg:"length,omitempty" json:"length,omitempty"`
	Token          ByteSlice `msg:"token,omitempty" json:"token,omitempty"`
	Data           ByteSlice `msg:"data,omitempty" json:"data,omitempty"`
	SequenceNumber *uint64   `msg:"sequence_number,omitempty" json:"sequence_number,omitempty"`
}
