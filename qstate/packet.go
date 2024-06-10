//go:generate msgp
package qstate

type Packet struct {
	PacketNumber int64   `msg:"packet_number" json:"packet_number" cbor:"1,keyasint"`
	Frames       []Frame `msg:"frames" json:"frames" cbor:"2,keyasint"`
}
