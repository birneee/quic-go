//go:generate msgp
package qstate

type ConnectionState string

const (
	ConnectionStateComplete  ConnectionState = "handshake_complete"
	ConnectionStateConfirmed ConnectionState = "handshake_confirmed"
)
