//go:generate msgp
package qstate

type ConnectionID struct {
	SequenceNumber      uint64               `msg:"sequence_number" json:"sequence_number"`
	ConnectionID        HexByteSlice         `msg:"connection_id" json:"connection_id"`
	StatelessResetToken *StatelessResetToken `msg:"stateless_reset_token" json:"stateless_reset_token"`
}
