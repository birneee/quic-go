package quic

import (
	"sync"
)

type EncryptedPacketQueues struct {
	mutex           sync.Mutex
	pendingRestores map[string] /* connection id */ []*receivedPacket
}

// TODO add maximum
// TODO add timeout
func NewEncryptedPacketQueues() *EncryptedPacketQueues {
	return &EncryptedPacketQueues{
		pendingRestores: map[string][]*receivedPacket{},
	}
}

// true if new entry for connection id was created
func (p *EncryptedPacketQueues) Enqueue(connID ConnectionID, packet *receivedPacket) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	queue, ok := p.pendingRestores[connID.String()]
	if !ok {
		queue = []*receivedPacket{}
		p.pendingRestores[connID.String()] = queue
	}
	p.pendingRestores[connID.String()] = append(queue, packet)
	return !ok
}

func (p *EncryptedPacketQueues) Dequeue(connID ConnectionID, conn Connection) []*receivedPacket {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	conn2 := conn.(*connection)
	queue := p.pendingRestores[connID.String()]
	for _, packet := range queue {
		conn2.handlePacket(packet)
	}
	delete(p.pendingRestores, connID.String())
	return queue
}
