package qtls

import (
	"context"
	"crypto/tls"
)

type QUICConn interface {
	SetTransportParameters(params []byte)
	Start(ctx context.Context) error
	NextEvent() tls.QUICEvent
	Close() error
	HandleData(level tls.QUICEncryptionLevel, data []byte) error
	SendSessionTicket(opts tls.QUICSessionTicketOptions) error
	ConnectionState() tls.ConnectionState
}
