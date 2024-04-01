package path

import (
	"github.com/stretchr/testify/require"
	"net"
	"testing"
)

func TestIPv4(t *testing.T) {
	pm := NewPathManager()
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:18080")
	require.NoError(t, err)
	require.False(t, pm.IsIgnoreMigrateTo(addr))
	pm.IgnoreMigrateTo(addr)
	require.True(t, pm.IsIgnoreMigrateTo(addr))
}

func TestIPv6(t *testing.T) {
	pm := NewPathManager()
	addr, err := net.ResolveUDPAddr("udp", "[::1]:18080")
	require.NoError(t, err)
	require.False(t, pm.IsIgnoreMigrateTo(addr))
	pm.IgnoreMigrateTo(addr)
	require.True(t, pm.IsIgnoreMigrateTo(addr))
}
