package xse

import (
	"github.com/lucas-clemente/quic-go/internal/qtls"
)

type sendStream struct {
	SendStream // base
	conn       *qtls.Conn
}

func (s *sendStream) Write(p []byte) (n int, err error) {
	return s.conn.Write(p)
}
