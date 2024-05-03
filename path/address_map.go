package path

import (
	"net"
)

type AddressMap[T any] interface {
	Put(key net.Addr, value T)
	Get(key net.Addr) (T, bool)
	Delete(key net.Addr)
	Length() int
}

type addressMap[T any] struct {
	innerMap map[[16]byte]map[uint16]T
}

func (a addressMap[T]) Length() int {
	return len(a.innerMap)
}

func newAddressMap[T any]() AddressMap[T] {
	m := &addressMap[T]{
		innerMap: map[[16]byte]map[uint16]T{},
	}
	return m
}

func (a addressMap[T]) Put(key net.Addr, value T) {
	addr, ok := key.(*net.UDPAddr)
	if !ok {
		panic("implement me")
	}
	ipKey := *(*[16]byte)(addr.IP.To16())
	portMap, ok := a.innerMap[ipKey]
	if !ok {
		portMap = map[uint16]T{}
		a.innerMap[ipKey] = portMap
	}
	portMap[uint16(addr.Port)] = value
}

func (a addressMap[T]) Get(key net.Addr) (T, bool) {
	addr, ok := key.(*net.UDPAddr)
	if !ok {
		panic("implement me")
	}
	ipKey := IpTo16(addr.IP)
	portMap, ok := a.innerMap[ipKey]
	if !ok {
		return *new(T), false
	}
	value, ok := portMap[uint16(addr.Port)]
	return value, ok
}

func IpTo16(ip net.IP) [16]byte {
	switch len(ip) {
	case net.IPv4len:
		return *(*[16]byte)(net.IPv4(ip[0], ip[1], ip[2], ip[3]))
	case net.IPv6len:
		return *(*[16]byte)(ip)
	case 0:
		return [16]byte{}
	default:
		panic("unexpected length")
	}
}

func (a addressMap[T]) Delete(key net.Addr) {
	//TODO implement me
	panic("implement me")
}
