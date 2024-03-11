//go:generate msgp
package qstate

type Packet struct {
	PacketNumber int64   `msg:"packet_number" json:"packet_number"`
	Frames       []Frame `msg:"frames" json:"frames"`
}
