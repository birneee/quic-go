package path

import (
	"net"
	"net/netip"
)

func IsAddrMigrated(a net.Addr, b net.Addr) bool {
	aUdp, ok := a.(*net.UDPAddr)
	if !ok {
		return false
	}
	bUdp, ok := b.(*net.UDPAddr)
	if !ok {
		return false
	}
	if aUdp.IP == nil {
		return false
	}
	if bUdp.IP == nil {
		return false
	}
	return !aUdp.IP.Equal(bUdp.IP) || aUdp.Port != bUdp.Port
}

func IsAddrMigrated2(a net.Addr, b net.Addr) bool {
	aUdp, ok := a.(*net.UDPAddr)
	if !ok {
		return false
	}
	bUdp, ok := b.(*net.UDPAddr)
	if !ok {
		return false
	}
	return aUdp.AddrPort() != bUdp.AddrPort()
}

func IsAddrMigrated3(a netip.AddrPort, b netip.AddrPort) bool {
	return a != b
}

func IsAddrMigrated4(a net.Addr, b net.Addr) bool {
	return a.(*net.UDPAddr).AddrPort() != b.(*net.UDPAddr).AddrPort()
}
