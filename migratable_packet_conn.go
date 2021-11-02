package quic

import (
	"net"
	"time"
)

type MigratablePacketConn struct {
	net.PacketConn
	internal net.PacketConn
}

func ListenUDP(network string, laddr *net.UDPAddr) (*MigratablePacketConn, error) {
	conn, err := net.ListenUDP(network, laddr)
	if err != nil {
		return nil, err
	}
	return &MigratablePacketConn{
		internal: conn,
	}, nil
}

func (m *MigratablePacketConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	return m.internal.ReadFrom(p)
}

func (m *MigratablePacketConn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	return m.internal.WriteTo(p, addr)
}

func (m *MigratablePacketConn) Close() error {
	return m.internal.Close()
}

func (m *MigratablePacketConn) LocalAddr() net.Addr {
	return m.internal.LocalAddr()
}

func (m *MigratablePacketConn) SetDeadline(t time.Time) error {
	return m.internal.SetDeadline(t)
}

func (m *MigratablePacketConn) SetReadDeadline(t time.Time) error {
	return m.internal.SetReadDeadline(t)
}

func (m *MigratablePacketConn) SetWriteDeadline(t time.Time) error {
	return m.internal.SetWriteDeadline(t)
}

func (m *MigratablePacketConn) Migrate() error {
	oldConn := m.internal

	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	if err != nil {
		return err
	}
	m.internal = conn

	go func() {
		time.Sleep(time.Second)
		err = oldConn.Close()
		if err != nil {
			panic(err)
		}
	}()

	//err = oldConn.Close()
	//if err != nil {
	//	return err
	//}

	return nil
}
