package quic

import (
	"net"
	"os"
	"reflect"
	"syscall"
	"time"
)

// check at compile time that interfaces satisfies interfaces
var _ net.PacketConn = &MigratableUDPConn{}
var _ net.Conn = &MigratableUDPConn{}
var _ interface {
	SyscallConn() (syscall.RawConn, error)
} = &MigratableUDPConn{}
var _ interface{ SetReadBuffer(int) error } = &MigratableUDPConn{}

// MigratableUDPConn
//
// Packet connection that supports migration of IP address and UDP port
type MigratableUDPConn struct {
	internal           *net.UDPConn
	maxRetries         int
	initialBackoffTime time.Duration
	// to restore after
	deadline      *time.Time
	readDeadline  *time.Time
	writeDeadline *time.Time
	readBuffer    *int
}

func ListenMigratableUDP(network string, laddr *net.UDPAddr) (*MigratableUDPConn, error) {
	conn, err := net.ListenUDP(network, laddr)
	if err != nil {
		return nil, err
	}
	return &MigratableUDPConn{
		internal:           conn,
		maxRetries:         5,
		initialBackoffTime: 1 * time.Millisecond,
	}, nil
}

func (m *MigratableUDPConn) Read(b []byte) (n int, err error) {
	//TODO handle errors caused by migration
	return m.internal.Read(b)
}

func (m *MigratableUDPConn) Write(b []byte) (n int, err error) {
	//TODO handle errors caused by migration
	return m.internal.Write(b)
}

func (m *MigratableUDPConn) RemoteAddr() net.Addr {
	return m.internal.RemoteAddr()
}

func (m *MigratableUDPConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	retryCount := 0
	backoffTime := m.initialBackoffTime
	for {
		n, addr, err = m.internal.ReadFrom(p)
		if err != nil {
			err = m.handleError(err, &retryCount, &backoffTime)
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
	backoffTime := m.initialBackoffTime
	for {
		n, err = m.internal.WriteTo(p, addr)
		if err != nil {
			err = m.handleError(err, &retryCount, &backoffTime)
			if err != nil {
				return n, err
			}
			continue // retry
		}
		return n, err
	}
}

// If error is returned it should no longer be retried
func (m *MigratableUDPConn) handleError(err error, retryCount *int, backoffTime *time.Duration) error {
	switch err := err.(type) {
	case *net.OpError:
		switch err := err.Err.(type) {
		case *os.SyscallError:
			switch err := err.Err.(type) {
			case syscall.Errno:
				if err.Error() == "network is unreachable" {
					// reopen and retry, because it could be caused of migration
					if *retryCount < m.maxRetries {
						time.Sleep(*backoffTime) // give socket migration some time
						*retryCount++
						*backoffTime *= 2
						_ = m.Reopen()
						return nil
					}
				}
			}
		default:
			// use reflect, because type is not public
			if reflect.TypeOf(err).String() == "poll.errNetClosing" {
				// retry, because it could be caused of migration
				if *retryCount < m.maxRetries {
					time.Sleep(*backoffTime) // give socket migration some time
					*retryCount++
					*backoffTime *= 2
					return nil
				}
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
	m.deadline = &t
	return m.internal.SetDeadline(t)
}

func (m *MigratableUDPConn) SetReadDeadline(t time.Time) error {
	m.readDeadline = &t
	return m.internal.SetReadDeadline(t)
}

func (m *MigratableUDPConn) SetWriteDeadline(t time.Time) error {
	m.writeDeadline = &t
	return m.internal.SetWriteDeadline(t)
}

func (m *MigratableUDPConn) SetReadBuffer(bytes int) error {
	m.readBuffer = &bytes
	return m.internal.SetReadBuffer(bytes)
}

func (m *MigratableUDPConn) SyscallConn() (syscall.RawConn, error) {
	return m.internal.SyscallConn()
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
	err = m.applyConfig()
	if err != nil {
		return err
	}
	return nil
}

// MigrateUDPSocket migrates to a new UDP socket.
// Returns new UDP address.
func (m *MigratableUDPConn) MigrateUDPSocket() (*net.UDPAddr, error) {
	oldSocket := m.internal
	defer oldSocket.Close()

	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	if err != nil {
		return nil, err
	}
	m.internal = conn
	err = m.applyConfig()
	if err != nil {
		return nil, err
	}
	return conn.LocalAddr().(*net.UDPAddr), nil
}

// apply config to the current internal UDP connection
func (m *MigratableUDPConn) applyConfig() error {
	if m.deadline != nil {
		err := m.internal.SetDeadline(*m.deadline)
		if err != nil {
			return err
		}
	}
	if m.readDeadline != nil {
		err := m.internal.SetReadDeadline(*m.readDeadline)
		if err != nil {
			return err
		}
	}
	if m.writeDeadline != nil {
		err := m.internal.SetWriteDeadline(*m.writeDeadline)
		if err != nil {
			return err
		}
	}
	if m.readBuffer != nil {
		err := m.internal.SetReadBuffer(*m.readBuffer)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MigratableUDPConn) ReadFromUDP(b []byte) (n int, addr *net.UDPAddr, err error) {
	//TODO handle errors caused by migration
	return m.internal.ReadFromUDP(b)
}

func (m *MigratableUDPConn) WriteToUDP(b []byte, addr *net.UDPAddr) (int, error) {
	//TODO handle errors caused by migration
	return m.internal.WriteToUDP(b, addr)
}
