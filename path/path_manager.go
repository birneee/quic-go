package path

import (
	"net"
	"sync"
)

type PathKey = string

type pathManager struct {
	mutex         sync.RWMutex
	ignoreReceive map[PathKey]interface{}
	ignoreMigrate map[PathKey]interface{}
	ignoreSend    map[PathKey]interface{}
}

var _ PathManager = &pathManager{}

func NewPathManager() PathManager {
	return &pathManager{
		ignoreReceive: map[PathKey]interface{}{},
		ignoreMigrate: map[PathKey]interface{}{},
		ignoreSend:    map[PathKey]interface{}{},
	}
}

func getPathKey(addr net.Addr) PathKey {
	return addr.Network() + " " + addr.String()
}

func (p *pathManager) IgnoreReceiveFrom(addr net.Addr) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.ignoreReceive[getPathKey(addr)] = nil
}

func (p *pathManager) IgnoreMigrateTo(addr net.Addr) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.ignoreMigrate[getPathKey(addr)] = nil
}

func (p *pathManager) IgnoreSendTo(addr net.Addr) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.ignoreSend[getPathKey(addr)] = nil
}

func (p *pathManager) IsIgnoreReceiveFrom(addr net.Addr) bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	_, ok := p.ignoreReceive[getPathKey(addr)]
	return ok
}

func (p *pathManager) IsIgnoreMigrateTo(addr net.Addr) bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	_, ok := p.ignoreMigrate[getPathKey(addr)]
	return ok
}

func (p *pathManager) IsIgnoreSendTo(addr net.Addr) bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	_, ok := p.ignoreSend[getPathKey(addr)]
	return ok
}
