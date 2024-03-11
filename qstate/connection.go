//go:generate msgp
package qstate

type Connection struct {
	Transport Transport `msg:"transport" json:"transport"`
	Metrics   Metrics   `msg:"metrics" json:"metrics"`
}

// ChangeVantagePoint creates an estimated connection state of the peer
func (c *Connection) ChangeVantagePoint(DestinationIP string, DestinationPort uint16) Connection {
	return Connection{
		Transport: c.Transport.ChangeVantagePoint(DestinationIP, DestinationPort),
		Metrics:   c.Metrics,
	}
}
