package path

import (
	"net"
	"sync"
)

type pathManager struct {
	mutex         sync.RWMutex
	ignoreReceive AddressMap[interface{}]
	ignoreMigrate AddressMap[interface{}]
	ignoreSend    AddressMap[interface{}]
}

var _ PathManager = &pathManager{}

func NewPathManager() PathManager {
	return &pathManager{
		ignoreReceive: newAddressMap[interface{}](),
		ignoreMigrate: newAddressMap[interface{}](),
		ignoreSend:    newAddressMap[interface{}](),
	}
}

func (p *pathManager) IgnoreReceiveFrom(addr net.Addr) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.ignoreReceive.Put(addr, nil)
}

func (p *pathManager) IgnoreMigrateTo(addr net.Addr) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.ignoreMigrate.Put(addr, nil)
}

func (p *pathManager) IgnoreSendTo(addr net.Addr) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.ignoreSend.Put(addr, nil)
}

func (p *pathManager) IsIgnoreReceiveFrom(addr net.Addr) bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	_, ok := p.ignoreReceive.Get(addr)
	return ok
}

func (p *pathManager) IsIgnoreMigrateTo(addr net.Addr) bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	_, ok := p.ignoreMigrate.Get(addr)
	return ok
}

func (p *pathManager) IsIgnoreSendTo(addr net.Addr) bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	_, ok := p.ignoreSend.Get(addr)
	return ok
}
