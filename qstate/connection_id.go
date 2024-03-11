//go:generate msgp
package qstate

type ConnectionID struct {
	SequenceNumber      uint64    `msg:"sequence_number" json:"sequence_number"`
	ConnectionID        []byte    `msg:"connection_id" json:"connection_id"`
	StatelessResetToken *[16]byte `msg:"stateless_reset_token" json:"stateless_reset_token"`
}
