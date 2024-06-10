//go:generate msgp
package qstate

type ConnectionID struct {
	SequenceNumber      uint64               `msg:"sequence_number" json:"sequence_number" cbor:"1,keyasint"`
	ConnectionID        HexByteSlice         `msg:"connection_id" json:"connection_id" cbor:"2,keyasint"`
	StatelessResetToken *StatelessResetToken `msg:"stateless_reset_token" json:"stateless_reset_token" cbor:"3,keyasint"`
}
