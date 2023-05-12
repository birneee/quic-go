package xads

import "github.com/quic-go/quic-go/internal/qtls"

type sendStream struct {
	SendStream // base
	conn       *qtls.Conn
}

func (s *sendStream) Write(p []byte) (n int, err error) {
	return s.conn.Write(p)
}
