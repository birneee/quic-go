package handover

import (
	"errors"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/wire"
	"math"
	"net"
	"strconv"
)

type ActiveConnectionID struct {
	SequenceNumber uint64
	ConnectionID   []byte
	// 16 bytes
	StatelessResetToken []byte
}

// State is used to handover QUIC connection
type State struct {
	// used for connection identification e.g. for qlog
	LogConnectionID protocol.ConnectionID
	// active client connection IDs
	ClientConnectionIDs []ActiveConnectionID
	// active server connection IDs
	ServerConnectionIDs []ActiveConnectionID
	Version             protocol.VersionNumber
	KeyPhase            protocol.KeyPhase
	// id of the used TLS 1.3 cipher suites.
	// see RFC 8446 Appendix B.4. Cipher Suites.
	CipherSuiteId uint16
	// used for header protection.
	// see RFC 9001 Section 5.4 Header Protection.
	// TODO use header protection key instead
	// TODO security concern: a H-QUIC Proxy can derived all past traffic secrets from this
	InitialServerTrafficSecret []byte
	// used for header protection.
	// see RFC 9001 Section 5.4 Header Protection.
	// TODO use header protection key instead
	// TODO security concern: a H-QUIC Proxy can derived all past traffic secrets from this
	InitialClientTrafficSecret []byte
	ServerTrafficSecret        []byte
	ClientTrafficSecret        []byte
	ServerAddress              string
	ClientAddress              string
	// TODO only include non-default parameters
	ClientTransportParameters wire.TransportParameters
	// TODO only include non-default parameters
	ServerTransportParameters wire.TransportParameters
	// might be an estimate from the opposite perspective
	ClientHighestSentPacketNumber protocol.PacketNumber
	// might be an estimate from the opposite perspective
	ServerHighestSentPacketNumber protocol.PacketNumber
}

func parseAddress(stringAddr string) (*net.UDPAddr, error) {
	ipString, portString, err := net.SplitHostPort(stringAddr)
	if err != nil {
		return nil, err
	}
	ip := net.ParseIP(ipString)
	if ip == nil {
		return nil, errors.New("invalid remote ip")
	}
	port, err := strconv.Atoi(portString)
	if err != nil {
		return nil, err
	}
	return &net.UDPAddr{
		IP:   ip,
		Port: port,
	}, nil
}

// RemoteAddress can be nil
func (s *State) RemoteAddress(perspective protocol.Perspective) *net.UDPAddr {
	var addrString string
	if perspective == protocol.PerspectiveClient {
		addrString = s.ServerAddress
	} else {
		addrString = s.ClientAddress
	}
	if len(addrString) == 0 {
		return nil
	}
	addr, err := parseAddress(addrString)
	if err != nil {
		panic(err)
	}
	return addr
}

func (s *State) SetRemoteAddress(perspective protocol.Perspective, addr net.UDPAddr) {
	if perspective == protocol.PerspectiveClient {
		s.ServerAddress = addr.String()
	} else {
		s.ClientAddress = addr.String()
	}
}

func (s *State) SetLocalAddress(perspective protocol.Perspective, addr net.UDPAddr) {
	if perspective == protocol.PerspectiveClient {
		s.ClientAddress = addr.String()
	} else {
		s.ServerAddress = addr.String()
	}
}

func copyBytes(a []byte) []byte {
	b := make([]byte, len(a))
	copy(b, a)
	return b
}

func (s *State) SendTrafficSecret(perspective protocol.Perspective) []byte {
	if perspective == protocol.PerspectiveClient {
		return copyBytes(s.ClientTrafficSecret)
	} else {
		return copyBytes(s.ServerTrafficSecret)
	}
}

func (s *State) SetSendTrafficSecret(perspective protocol.Perspective, ts []byte) {
	if perspective == protocol.PerspectiveClient {
		s.ClientTrafficSecret = ts
	} else {
		s.ServerTrafficSecret = ts
	}
}

func (s *State) FirstSendTrafficSecret(perspective protocol.Perspective) []byte {
	if perspective == protocol.PerspectiveClient {
		return copyBytes(s.InitialClientTrafficSecret)
	} else {
		return copyBytes(s.InitialServerTrafficSecret)
	}
}

func (s *State) SetFirstSendTrafficSecret(perspective protocol.Perspective, ts []byte) {
	if perspective == protocol.PerspectiveClient {
		s.InitialClientTrafficSecret = ts
	} else {
		s.InitialServerTrafficSecret = ts
	}
}

func (s *State) ReceiveTrafficSecret(perspective protocol.Perspective) []byte {
	if perspective == protocol.PerspectiveClient {
		return copyBytes(s.ServerTrafficSecret)
	} else {
		return copyBytes(s.ClientTrafficSecret)
	}
}

func (s *State) SetReceiveTrafficSecret(perspective protocol.Perspective, ts []byte) {
	if perspective == protocol.PerspectiveClient {
		s.ServerTrafficSecret = ts
	} else {
		s.ClientTrafficSecret = ts
	}
}

func (s *State) FirstReceiveTrafficSecret(perspective protocol.Perspective) []byte {
	if perspective == protocol.PerspectiveClient {
		return copyBytes(s.InitialServerTrafficSecret)
	} else {
		return copyBytes(s.InitialClientTrafficSecret)
	}
}

func (s *State) SetFirstReceiveTrafficSecret(perspective protocol.Perspective, ts []byte) {
	if perspective == protocol.PerspectiveClient {
		s.InitialServerTrafficSecret = ts
	} else {
		s.InitialClientTrafficSecret = ts
	}
}

func (s *State) ActiveSrcConnectionIDs(perspective protocol.Perspective) []ActiveConnectionID {
	if perspective == protocol.PerspectiveClient {
		return s.ClientConnectionIDs
	} else {
		return s.ServerConnectionIDs
	}
}

func (s *State) MinActiveSrcConnectionID(perspective protocol.Perspective) protocol.ConnectionID {
	var minSN uint64 = math.MaxUint64
	var minID protocol.ConnectionID
	for _, activeConnID := range s.ActiveSrcConnectionIDs(perspective) {
		if activeConnID.SequenceNumber <= minSN {
			minSN = activeConnID.SequenceNumber
			minID = protocol.ParseConnectionID(activeConnID.ConnectionID)
		}
	}
	return minID
}

func (s *State) MaxActiveSrcConnectionID(perspective protocol.Perspective) (uint64, protocol.ConnectionID) {
	var minSN uint64 = 0
	var minID protocol.ConnectionID
	for _, activeConnID := range s.ActiveSrcConnectionIDs(perspective) {
		if activeConnID.SequenceNumber >= minSN {
			minSN = activeConnID.SequenceNumber
			minID = protocol.ParseConnectionID(activeConnID.ConnectionID)
		}
	}
	return minSN, minID
}

func (s *State) SetActiveSrcConnectionIDs(perspective protocol.Perspective, connIDs []ActiveConnectionID) {
	if perspective == protocol.PerspectiveClient {
		s.ClientConnectionIDs = connIDs
	} else {
		s.ServerConnectionIDs = connIDs
	}
}

func (s *State) ActiveDestConnectionIDs(perspective protocol.Perspective) []ActiveConnectionID {
	if perspective == protocol.PerspectiveClient {
		return s.ServerConnectionIDs
	} else {
		return s.ClientConnectionIDs
	}
}

func (s *State) MinActiveDestConnectionID(perspective protocol.Perspective) *protocol.ConnectionID {
	var minSN uint64 = math.MaxUint64
	var minID protocol.ConnectionID
	for _, activeConnID := range s.ActiveDestConnectionIDs(perspective) {
		if activeConnID.SequenceNumber <= minSN {
			minSN = activeConnID.SequenceNumber
			minID = protocol.ParseConnectionID(activeConnID.ConnectionID)
		}
	}
	return &minID
}

func (s *State) SetActiveDestConnectionIDs(perspective protocol.Perspective, connIDs []ActiveConnectionID) {
	if perspective == protocol.PerspectiveClient {
		s.ServerConnectionIDs = connIDs
	} else {
		s.ClientConnectionIDs = connIDs
	}
}

func (s *State) SrcConnectionIDLength(perspective protocol.Perspective) int {
	for _, activeConnID := range s.ActiveSrcConnectionIDs(perspective) {
		return len(activeConnID.ConnectionID)
	}
	panic("no connection ids")
}

func (s *State) OwnTransportParameters(perspective protocol.Perspective) wire.TransportParameters {
	if perspective == protocol.PerspectiveClient {
		return s.ClientTransportParameters
	} else {
		return s.ServerTransportParameters
	}
}

func (s *State) SetOwnTransportParameters(perspective protocol.Perspective, tp wire.TransportParameters) {
	if perspective == protocol.PerspectiveClient {
		s.ClientTransportParameters = tp
	} else {
		s.ServerTransportParameters = tp
	}
}

func (s *State) PeerTransportParameters(perspective protocol.Perspective) wire.TransportParameters {
	if perspective == protocol.PerspectiveClient {
		return s.ServerTransportParameters
	} else {
		return s.ClientTransportParameters
	}
}

func (s *State) SetPeerTransportParameters(perspective protocol.Perspective, tp wire.TransportParameters) {
	if perspective == protocol.PerspectiveClient {
		s.ServerTransportParameters = tp
	} else {
		s.ClientTransportParameters = tp
	}
}

func (s *State) HighestSentPacketNumber(perspective protocol.Perspective) protocol.PacketNumber {
	if perspective == protocol.PerspectiveClient {
		return s.ServerHighestSentPacketNumber
	} else {
		return s.ClientHighestSentPacketNumber
	}
}

func (s *State) SetHighestSentPacketNumber(perspective protocol.Perspective, pn protocol.PacketNumber) {
	if perspective == protocol.PerspectiveClient {
		s.ServerHighestSentPacketNumber = pn
	} else {
		s.ClientHighestSentPacketNumber = pn
	}
}

// Clone
// TODO deep copy
func (s *State) Clone() *State {
	return &*s
}
