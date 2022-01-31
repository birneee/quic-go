package xse

import (
	"github.com/lucas-clemente/quic-go/internal/qtls"
)

type receiveStream struct {
	ReceiveStream // Base
	conn          *qtls.Conn
}

func (r *receiveStream) Read(p []byte) (n int, err error) {
	return r.conn.Read(p)
}
