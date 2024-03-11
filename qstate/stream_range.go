//go:generate msgp
package qstate

type StreamRange struct {
	Offset int64  `msg:"offset" json:"offset"`
	Data   []byte `msg:"data" json:"data"`
}
