package path

import (
	"net"
)

type AddressMap[T any] interface {
	Put(key net.Addr, value T)
	Get(key net.Addr) (T, bool)
	Delete(key net.Addr)
}

type addressMap[T any] struct {
	ipV4Map map[[4]byte]map[int]T
	ipV6Map map[[16]byte]map[int]T
}

func newAddressMap[T any]() AddressMap[T] {
	m := &addressMap[T]{
		ipV4Map: map[[4]byte]map[int]T{},
		ipV6Map: map[[16]byte]map[int]T{},
	}
	return m
}

func (a addressMap[T]) Put(key net.Addr, value T) {
	addr, ok := key.(*net.UDPAddr)
	if !ok {
		panic("implement me")
	}
	if len(addr.IP) == 4 {
		ipKey := *(*[4]byte)(addr.IP)
		portMap, ok := a.ipV4Map[ipKey]
		if !ok {
			portMap = map[int]T{}
			a.ipV4Map[ipKey] = portMap
		}
		portMap[addr.Port] = value
		return
	}
	if len(addr.IP) == 16 {
		ipKey := *(*[16]byte)(addr.IP)
		portMap, ok := a.ipV6Map[ipKey]
		if !ok {
			portMap = map[int]T{}
			a.ipV6Map[ipKey] = portMap
		}
		portMap[addr.Port] = value
		return
	}
	panic("unexpected state")
}

func (a addressMap[T]) Get(key net.Addr) (T, bool) {
	addr, ok := key.(*net.UDPAddr)
	if !ok {
		return *new(T), false
	}
	if len(addr.IP) == 4 {
		ipKey := *(*[4]byte)(addr.IP)
		portMap, ok := a.ipV4Map[ipKey]
		if !ok {
			return *new(T), false
		}
		value, ok := portMap[addr.Port]
		return value, ok
	}
	if len(addr.IP) == 16 {
		ipKey := *(*[16]byte)(addr.IP)
		portMap, ok := a.ipV6Map[ipKey]
		if !ok {
			return *new(T), false
		}
		value, ok := portMap[addr.Port]
		return value, ok
	}
	return *new(T), false
}

func (a addressMap[T]) Delete(key net.Addr) {
	//TODO implement me
	panic("implement me")
}
