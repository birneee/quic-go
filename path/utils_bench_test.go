package path

import (
	"github.com/stretchr/testify/require"
	"net"
	"testing"
)

func BenchmarkIsAddrMigrated(b *testing.B) {
	addr1, err := net.ResolveUDPAddr("udp", "127.0.0.1:1234")
	require.NoError(b, err)
	addr2, err := net.ResolveUDPAddr("udp", "127.0.0.1:1234")
	require.NoError(b, err)
	for i := 0; i < b.N; i++ {
		IsAddrMigrated(addr1, addr2)
	}
}

func BenchmarkIsAddrMigrated2(b *testing.B) {
	addr1, err := net.ResolveUDPAddr("udp", "127.0.0.1:1234")
	require.NoError(b, err)
	addr2, err := net.ResolveUDPAddr("udp", "127.0.0.1:1234")
	require.NoError(b, err)
	for i := 0; i < b.N; i++ {
		IsAddrMigrated2(addr1, addr2)
	}
}

func BenchmarkIsAddrMigrated3(b *testing.B) {
	addr1, err := net.ResolveUDPAddr("udp", "127.0.0.1:1234")
	require.NoError(b, err)
	addrPort1 := addr1.AddrPort()
	addr2, err := net.ResolveUDPAddr("udp", "127.0.0.1:1234")
	require.NoError(b, err)
	addrPort2 := addr2.AddrPort()
	for i := 0; i < b.N; i++ {
		IsAddrMigrated3(addrPort1, addrPort2)
	}
}

func BenchmarkIsAddrMigrated4(b *testing.B) {
	addr1, err := net.ResolveUDPAddr("udp", "127.0.0.1:1234")
	require.NoError(b, err)
	addr2, err := net.ResolveUDPAddr("udp", "127.0.0.1:1234")
	require.NoError(b, err)
	for i := 0; i < b.N; i++ {
		IsAddrMigrated4(addr1, addr2)
	}
}
