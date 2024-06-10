//go:generate msgp
package qstate

type StreamRange struct {
	Offset int64  `msg:"offset" json:"offset" cbor:"1,keyasint"`
	Data   []byte `msg:"data" json:"data" cbor:"2,keyasint"`
}
