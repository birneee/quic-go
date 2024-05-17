//go:generate msgp
//go:generate go run github.com/mailru/easyjson/...@latest --all
package qstate

type Connection struct {
	State     ConnectionState `msg:"state" json:"state"`
	Transport Transport       `msg:"transport" json:"transport"`
	Crypto    Crypto          `msg:"crypto" json:"crypto"`
	Metrics   Metrics         `msg:"metrics" json:"metrics"`
}

// ChangeVantagePoint creates an estimated connection state of the peer
func (c *Connection) ChangeVantagePoint(DestinationIP string, DestinationPort uint16) Connection {
	if c.State != ConnectionStateConfirmed {
		panic("unexpected connection state")
	}
	return Connection{
		State:     ConnectionStateConfirmed,
		Transport: c.Transport.ChangeVantagePoint(DestinationIP, DestinationPort),
		Crypto:    c.Crypto.ChangeVantagePoint(),
		Metrics:   c.Metrics,
	}
}
