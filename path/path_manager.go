package path

import (
	"github.com/lucas-clemente/quic-go/internal/utils"
	"github.com/lucas-clemente/quic-go/internal/wire"
	"github.com/lucas-clemente/quic-go/logging"
	"net"
	"time"
)

type pathManager struct {
	currentSendPath         *path
	paths                   map[string]*path
	updateSendPathObservers []func(path Path)
	logger                  utils.Logger
	tracer                  logging.ConnectionTracer
	queuePathChallengeFrame func(path Path)
}

var _ PathManager = &pathManager{}

func NewPathManager(initialAddr *net.UDPAddr, pathValidated bool, updateSendPathObserver func(path Path), queuePathChallengeFrame func(path Path), logger utils.Logger, tracer logging.ConnectionTracer) PathManager {
	currentSendPath := NewPath(initialAddr, pathValidated)
	return &pathManager{
		currentSendPath:         currentSendPath,
		paths:                   map[string]*path{initialAddr.String(): currentSendPath},
		updateSendPathObservers: []func(path Path){updateSendPathObserver},
		logger:                  logger,
		tracer:                  tracer,
		queuePathChallengeFrame: queuePathChallengeFrame,
	}
}

func (m *pathManager) CurrentSendPath() Path {
	return m.currentSendPath
}

func (m *pathManager) OnReceiveNonProbingPacket(addr *net.UDPAddr) {
	var p *path
	if addr.String() == m.currentSendPath.addr.String() {
		p = m.currentSendPath
	} else {
		p = m.getOrCreatePath(addr)
		m.setSendPath(p)
		if !p.PeerAddressValidated() {
			m.queuePathChallengeFrame(p)
		}
	}
	p.lastPacketReceiveTime = time.Now()
}

func (m *pathManager) OnReceivePathResponseFrame(addr *net.UDPAddr, frame *wire.PathResponseFrame) {
	path := m.getOrCreatePath(addr)
	if path.challengeData == frame.Data {
		path.SetPeerAddressValidated()
	}
}

func (m *pathManager) setSendPath(p *path) {
	m.currentSendPath = p
	m.logger.Debugf("Migrated from %s to %s", m.currentSendPath.addr, p.addr)
	//TODO change message when standardized https://datatracker.ietf.org/doc/html/draft-marx-qlog-event-definitions-quic-h3#section-5.1.8
	if m.tracer != nil {
		m.tracer.UpdatedPath(p.addr)
	}
	for _, observer := range m.updateSendPathObservers {
		observer(p)
	}
}

func (m *pathManager) getOrCreatePath(addr *net.UDPAddr) *path {
	key := addr.String()
	p, ok := m.paths[key]
	if ok {
		return p
	} else {
		p := NewPath(addr, false)
		m.paths[key] = p
		return p
	}
}

func (m *pathManager) GetOrCreatePath(addr net.Addr) Path {
	return m.getOrCreatePath(addr.(*net.UDPAddr))
}
