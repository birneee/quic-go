package qtls

import (
	"context"
	"crypto/tls"
)

type QUICConn interface {
	SetTransportParameters(params []byte)
	Start(ctx context.Context) error
	NextEvent() QUICEvent
	Close() error
	HandleData(level QUICEncryptionLevel, data []byte) error
	SendSessionTicket(earlyData bool) error
	ConnectionState() tls.ConnectionState
}
