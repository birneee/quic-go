package handover

import (
	"errors"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/wire"
	"math"
	"net"
	"strconv"
)

// StreamState
//TODO not used yet
type StreamState struct {
	InitiatedBy          string
	UnidirectionalStream bool
	StreamID             uint32
	Offset               uint32
}

type ActiveConnectionID struct {
	SequenceNumber uint64
	ConnectionID   protocol.ConnectionID
	// 16 bytes
	StatelessResetToken []byte
}

// State is used to handover QUIC connection
type State struct {
	ActiveClientConnectionIDs []ActiveConnectionID
	ActiveServerConnectionIDs []ActiveConnectionID
	Version                   protocol.VersionNumber
	KeyPhase                  protocol.KeyPhase
	// id of the used TLS 1.3 cipher suites.
	// see RFC 8446 Appendix B.4. Cipher Suites.
	SuiteId uint16
	// used for header protection.
	// see RFC 9001 Section 5.4 Header Protection.
	FirstServerSendTrafficSecret []byte
	// used for header protection.
	// see RFC 9001 Section 5.4 Header Protection.
	FirstClientSendTrafficSecret []byte
	ServerSendTrafficSecret      []byte
	ClientSendTrafficSecret      []byte
	ServerAddress                string
	ClientAddress                string
	ClientTransportParameters    wire.TransportParameters
	ServerTransportParameters    wire.TransportParameters
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
		return copyBytes(s.ClientSendTrafficSecret)
	} else {
		return copyBytes(s.ServerSendTrafficSecret)
	}
}

func (s *State) SetSendTrafficSecret(perspective protocol.Perspective, ts []byte) {
	if perspective == protocol.PerspectiveClient {
		s.ClientSendTrafficSecret = ts
	} else {
		s.ServerSendTrafficSecret = ts
	}
}

func (s *State) FirstSendTrafficSecret(perspective protocol.Perspective) []byte {
	if perspective == protocol.PerspectiveClient {
		return copyBytes(s.FirstClientSendTrafficSecret)
	} else {
		return copyBytes(s.FirstServerSendTrafficSecret)
	}
}

func (s *State) SetFirstSendTrafficSecret(perspective protocol.Perspective, ts []byte) {
	if perspective == protocol.PerspectiveClient {
		s.FirstClientSendTrafficSecret = ts
	} else {
		s.FirstServerSendTrafficSecret = ts
	}
}

func (s *State) ReceiveTrafficSecret(perspective protocol.Perspective) []byte {
	if perspective == protocol.PerspectiveClient {
		return copyBytes(s.ServerSendTrafficSecret)
	} else {
		return copyBytes(s.ClientSendTrafficSecret)
	}
}

func (s *State) SetReceiveTrafficSecret(perspective protocol.Perspective, ts []byte) {
	if perspective == protocol.PerspectiveClient {
		s.ServerSendTrafficSecret = ts
	} else {
		s.ClientSendTrafficSecret = ts
	}
}

func (s *State) FirstReceiveTrafficSecret(perspective protocol.Perspective) []byte {
	if perspective == protocol.PerspectiveClient {
		return copyBytes(s.FirstServerSendTrafficSecret)
	} else {
		return copyBytes(s.FirstClientSendTrafficSecret)
	}
}

func (s *State) SetFirstReceiveTrafficSecret(perspective protocol.Perspective, ts []byte) {
	if perspective == protocol.PerspectiveClient {
		s.FirstServerSendTrafficSecret = ts
	} else {
		s.FirstClientSendTrafficSecret = ts
	}
}

func (s *State) ActiveSrcConnectionIDs(perspective protocol.Perspective) []ActiveConnectionID {
	if perspective == protocol.PerspectiveClient {
		return s.ActiveClientConnectionIDs
	} else {
		return s.ActiveServerConnectionIDs
	}
}

func (s *State) MinActiveSrcConnectionID(perspective protocol.Perspective) protocol.ConnectionID {
	var minSN uint64 = math.MaxUint64
	var minID protocol.ConnectionID
	for _, activeConnID := range s.ActiveSrcConnectionIDs(perspective) {
		if activeConnID.SequenceNumber <= minSN {
			minSN = activeConnID.SequenceNumber
			minID = activeConnID.ConnectionID
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
			minID = activeConnID.ConnectionID
		}
	}
	return minSN, minID
}

func (s *State) SetActiveSrcConnectionIDs(perspective protocol.Perspective, connIDs []ActiveConnectionID) {
	if perspective == protocol.PerspectiveClient {
		s.ActiveClientConnectionIDs = connIDs
	} else {
		s.ActiveServerConnectionIDs = connIDs
	}
}

func (s *State) ActiveDestConnectionIDs(perspective protocol.Perspective) []ActiveConnectionID {
	if perspective == protocol.PerspectiveClient {
		return s.ActiveServerConnectionIDs
	} else {
		return s.ActiveClientConnectionIDs
	}
}

func (s *State) MinActiveDestConnectionID(perspective protocol.Perspective) protocol.ConnectionID {
	var minSN uint64 = math.MaxUint64
	var minID protocol.ConnectionID
	for _, activeConnID := range s.ActiveDestConnectionIDs(perspective) {
		if activeConnID.SequenceNumber <= minSN {
			minSN = activeConnID.SequenceNumber
			minID = activeConnID.ConnectionID
		}
	}
	return minID
}

func (s *State) SetActiveDestConnectionIDs(perspective protocol.Perspective, connIDs []ActiveConnectionID) {
	if perspective == protocol.PerspectiveClient {
		s.ActiveServerConnectionIDs = connIDs
	} else {
		s.ActiveClientConnectionIDs = connIDs
	}
}

func (s *State) SrcConnectionIDLength(perspective protocol.Perspective) int {
	for _, activeConnID := range s.ActiveSrcConnectionIDs(perspective) {
		return activeConnID.ConnectionID.Len()
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

// Clone
//TODO deep copy
func (s *State) Clone() *State {
	return &*s
}
