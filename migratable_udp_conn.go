package quic

import (
	"net"
	"os"
	"time"
)

// MigratableUDPConn
//
// Packet connection that supports migration of IP address and UDP port
type MigratableUDPConn struct {
	net.PacketConn
	internal   *net.UDPConn
	maxRetries int
}

func ListenMigratableUDP(network string, laddr *net.UDPAddr) (*MigratableUDPConn, error) {
	conn, err := net.ListenUDP(network, laddr)
	if err != nil {
		return nil, err
	}
	return &MigratableUDPConn{
		internal:   conn,
		maxRetries: 5,
	}, nil
}

func (m *MigratableUDPConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	retryCount := 0
	for {
		n, addr, err = m.internal.ReadFrom(p)
		if err != nil {
			err = m.handleError(err, &retryCount)
			if err != nil {
				return n, addr, err
			}
			continue // retry
		}
		return n, addr, err
	}
}

func (m *MigratableUDPConn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	retryCount := 0
	for {
		n, err = m.internal.WriteTo(p, addr)
		if err != nil {
			err = m.handleError(err, &retryCount)
			if err != nil {
				return n, err
			}
			continue // retry
		}
		return n, err
	}
}

// If error is returned it should no longer be retried
func (m *MigratableUDPConn) handleError(err error, retryCount *int) error {
	if opErr, ok := err.(*net.OpError); ok {
		if sysErr, ok := opErr.Err.(*os.SyscallError); ok {
			switch sysErr.Error() {
			case "sendto: network is unreachable":
				// reopen and retry, because it could be caused of migration
				if *retryCount < m.maxRetries {
					err = m.Reopen()
					if err != nil {
						return err
					}
					return nil
				}
			}
		} else if opErr.Err.Error() == "use of closed network connection" {
			// retry, because it could be caused of migration
			if *retryCount < m.maxRetries {
				// give socket migration some time
				time.Sleep(10 * time.Millisecond)
				*retryCount++
				return nil
			}
		}
	}
	return err
}

func (m *MigratableUDPConn) Close() error {
	return m.internal.Close()
}

func (m *MigratableUDPConn) LocalAddr() net.Addr {
	return m.internal.LocalAddr()
}

func (m *MigratableUDPConn) SetDeadline(t time.Time) error {
	return m.internal.SetDeadline(t)
}

func (m *MigratableUDPConn) SetReadDeadline(t time.Time) error {
	return m.internal.SetReadDeadline(t)
}

func (m *MigratableUDPConn) SetWriteDeadline(t time.Time) error {
	return m.internal.SetWriteDeadline(t)
}

// Reopen new UDP socket on same address
func (m *MigratableUDPConn) Reopen() error {
	err := m.internal.Close()
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", m.internal.LocalAddr().(*net.UDPAddr))
	if err != nil {
		return err
	}
	m.internal = conn
	return nil
}

// Migrate connection to new UDP socket.
// Returns new UDP address.
func (m *MigratableUDPConn) Migrate() (*net.UDPAddr, error) {
	err := m.internal.Close()
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	if err != nil {
		return nil, err
	}
	m.internal = conn

	return conn.LocalAddr().(*net.UDPAddr), nil
}
