package quic

import (
	"crypto/rand"
	"github.com/lucas-clemente/quic-go/internal/utils"
	"github.com/lucas-clemente/quic-go/internal/wire"
	"github.com/lucas-clemente/quic-go/logging"
	"net"
	"time"
)

type PathManager interface {
	OnReceiveNonProbingPacket(addr *net.UDPAddr)
	OnReceivePathResponseFrame(addr *net.UDPAddr, frame *wire.PathResponseFrame)
}

type pathManager struct {
	currentSendPath         *path
	paths                   map[string]*path
	updateSendPathObservers []func(addr *net.UDPAddr)
	logger                  utils.Logger
	tracer                  logging.ConnectionTracer
}

var _ PathManager = &pathManager{}

func NewPathManager(initialAddr *net.UDPAddr, updateSendPathObserver func(addr *net.UDPAddr), logger utils.Logger, tracer logging.ConnectionTracer) PathManager {
	currentSendPath := NewPath(initialAddr)
	return &pathManager{
		currentSendPath:         currentSendPath,
		paths:                   map[string]*path{initialAddr.String(): currentSendPath},
		updateSendPathObservers: []func(addr *net.UDPAddr){updateSendPathObserver},
		logger:                  logger,
		tracer:                  tracer,
	}
}

func (m pathManager) OnReceiveNonProbingPacket(addr *net.UDPAddr) {
	if addr.String() == m.currentSendPath.addr.String() {
		m.currentSendPath.lastPacketReceiveTime = time.Now()
		return
	}
	path := m.getOrCreatePathByAddr(addr)
	path.lastPacketReceiveTime = time.Now()
	m.setSendPath(path)
	if !path.validated {
		m.probe(path)
	}
}

func (m pathManager) OnReceivePathResponseFrame(addr *net.UDPAddr, frame *wire.PathResponseFrame) {
	path := m.getOrCreatePathByAddr(addr)
	if path.challengeData == frame.Data {
		path.validated = true
	}
}

func (m pathManager) setSendPath(p *path) {
	m.logger.Debugf("Migrated from %s to %s", m.currentSendPath.addr, p.addr)
	//TODO change message when standardized https://datatracker.ietf.org/doc/html/draft-marx-qlog-event-definitions-quic-h3#section-5.1.8
	if m.tracer != nil {
		//m.tracer.UpdatedPath(remoteAddr)
	}
	for _, observer := range m.updateSendPathObservers {
		observer(p.addr)
	}
}

func (m pathManager) getOrCreatePathByAddr(addr *net.UDPAddr) *path {
	key := addr.String()
	p, ok := m.paths[key]
	if ok {
		return p
	} else {
		p := NewPath(addr)
		m.paths[key] = p
		return p
	}
}

func (m pathManager) probe(p *path) {
	panic("TODO")
}

type path struct {
	addr                  *net.UDPAddr
	validated             bool
	lastPacketReceiveTime time.Time
	challengeData         [8]byte
}

func NewPath(addr *net.UDPAddr) *path {
	p := &path{
		addr:      addr,
		validated: false,
	}
	_, err := rand.Read(p.challengeData[:])
	if err != nil {
		panic(err)
	}
	return p
}
