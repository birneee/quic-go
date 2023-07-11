package handover

import (
	"encoding/json"
	"errors"
	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/wire"
	"github.com/quic-go/quic-go/logging"
	"math"
	"net"
	"strconv"
	"time"
)

type ConnectionIDSequenceNumber uint64

type ConnectionIDWithResetToken struct {
	ConnectionID []byte
	// 16 bytes
	StatelessResetToken []byte
}

// State is used to handover QUIC connection
type State struct {
	// active client connection IDs
	ClientConnectionIDs map[ConnectionIDSequenceNumber]*ConnectionIDWithResetToken
	// active server connection IDs
	ServerConnectionIDs map[ConnectionIDSequenceNumber]*ConnectionIDWithResetToken
	Version             protocol.VersionNumber
	KeyPhase            protocol.KeyPhase
	// id of the used TLS 1.3 cipher suites.
	// see RFC 8446 Appendix B.4. Cipher Suites.
	CipherSuiteId uint16
	// used for header protection.
	// see RFC 9001 Section 5.4 Header Protection.
	ServerHeaderProtectionKey []byte
	// used for header protection.
	// see RFC 9001 Section 5.4 Header Protection.
	ClientHeaderProtectionKey []byte
	ServerTrafficSecret       []byte
	ClientTrafficSecret       []byte
	ServerAddress             string
	ClientAddress             string
	// TODO only include non-default parameters
	ClientTransportParameters []byte
	// TODO only include non-default parameters
	ServerTransportParameters []byte
	// might be an estimate from the opposite perspective
	ClientHighestSentPacketNumber protocol.PacketNumber
	// might be an estimate from the opposite perspective
	ServerHighestSentPacketNumber protocol.PacketNumber
	UniStreams                    map[protocol.StreamID]*UniStreamState
	BidiStreams                   map[protocol.StreamID]*BidiStreamState
	ClientNextUniStream           protocol.StreamID
	ServerNextUniStream           protocol.StreamID
	ClientNextBidiStream          protocol.StreamID
	ServerNextBidiStream          protocol.StreamID
	ClientDirectionMaxData        protocol.ByteCount
	ServerDirectionMaxData        protocol.ByteCount
	ServerDirectionBytes          protocol.ByteCount
	ClientDirectionBytes          protocol.ByteCount
	ClientCongestionWindow        protocol.ByteCount
	ServerCongestionWindow        protocol.ByteCount
	RTT                           time.Duration
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

func (s *State) SendHeaderProtectionKey(perspective protocol.Perspective) []byte {
	if perspective == protocol.PerspectiveClient {
		return copyBytes(s.ClientHeaderProtectionKey)
	} else {
		return copyBytes(s.ServerHeaderProtectionKey)
	}
}

func (s *State) SetSendHeaderProtectionKey(perspective protocol.Perspective, ts []byte) {
	if perspective == protocol.PerspectiveClient {
		s.ClientHeaderProtectionKey = ts
	} else {
		s.ServerHeaderProtectionKey = ts
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

func (s *State) ReceiveHeaderProtectionKey(perspective protocol.Perspective) []byte {
	if perspective == protocol.PerspectiveClient {
		return copyBytes(s.ServerHeaderProtectionKey)
	} else {
		return copyBytes(s.ClientHeaderProtectionKey)
	}
}

func (s *State) SetReceiveHeaderProtectionKey(perspective protocol.Perspective, ts []byte) {
	if perspective == protocol.PerspectiveClient {
		s.ServerHeaderProtectionKey = ts
	} else {
		s.ClientHeaderProtectionKey = ts
	}
}

func (s *State) ActiveSrcConnectionIDs(perspective protocol.Perspective) map[ConnectionIDSequenceNumber]*ConnectionIDWithResetToken {
	if perspective == protocol.PerspectiveClient {
		return s.ClientConnectionIDs
	} else {
		return s.ServerConnectionIDs
	}
}

func (s *State) MinActiveSrcConnectionID(perspective protocol.Perspective) protocol.ConnectionID {
	minSN := ConnectionIDSequenceNumber(math.MaxUint64)
	var minID protocol.ConnectionID
	for sequenceNumber, activeConnID := range s.ActiveSrcConnectionIDs(perspective) {
		if sequenceNumber <= minSN {
			minSN = sequenceNumber
			minID = protocol.ParseConnectionID(activeConnID.ConnectionID)
		}
	}
	return minID
}

func (s *State) MaxActiveSrcConnectionID(perspective protocol.Perspective) (ConnectionIDSequenceNumber, protocol.ConnectionID) {
	minSN := ConnectionIDSequenceNumber(0)
	var minID protocol.ConnectionID
	for sequenceNumber, activeConnID := range s.ActiveSrcConnectionIDs(perspective) {
		if sequenceNumber >= minSN {
			minSN = sequenceNumber
			minID = protocol.ParseConnectionID(activeConnID.ConnectionID)
		}
	}
	return minSN, minID
}

func (s *State) SetActiveSrcConnectionIDs(perspective protocol.Perspective, connIDs map[ConnectionIDSequenceNumber]*ConnectionIDWithResetToken) {
	if perspective == protocol.PerspectiveClient {
		s.ClientConnectionIDs = connIDs
	} else {
		s.ServerConnectionIDs = connIDs
	}
}

func (s *State) ActiveDestConnectionIDs(perspective protocol.Perspective) map[ConnectionIDSequenceNumber]*ConnectionIDWithResetToken {
	if perspective == protocol.PerspectiveClient {
		return s.ServerConnectionIDs
	} else {
		return s.ClientConnectionIDs
	}
}

func (s *State) MinActiveDestConnectionID(perspective protocol.Perspective) *protocol.ConnectionID {
	minSN := ConnectionIDSequenceNumber(math.MaxUint64)
	var minID protocol.ConnectionID
	for sequenceNumber, activeConnID := range s.ActiveDestConnectionIDs(perspective) {
		if sequenceNumber <= minSN {
			minSN = sequenceNumber
			minID = protocol.ParseConnectionID(activeConnID.ConnectionID)
		}
	}
	return &minID
}

func (s *State) SetActiveDestConnectionIDs(perspective protocol.Perspective, connIDs map[ConnectionIDSequenceNumber]*ConnectionIDWithResetToken) {
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
	var bytes []byte
	if perspective == protocol.PerspectiveClient {
		bytes = s.ClientTransportParameters
	} else {
		bytes = s.ServerTransportParameters
	}

	tp := wire.TransportParameters{}
	err := tp.Unmarshal(bytes, perspective)
	if err != nil {
		panic(err)
	}
	return tp
}

func (s *State) SetOwnTransportParameters(perspective protocol.Perspective, tp wire.TransportParameters) {
	if perspective == protocol.PerspectiveClient {
		s.ClientTransportParameters = tp.MarshalForHandover(perspective)
	} else {
		s.ServerTransportParameters = tp.MarshalForHandover(perspective)
	}
}

func (s *State) PeerTransportParameters(perspective protocol.Perspective) wire.TransportParameters {
	return s.OwnTransportParameters(perspective.Opposite())
}

func (s *State) SetPeerTransportParameters(perspective protocol.Perspective, tp wire.TransportParameters) {
	s.SetOwnTransportParameters(perspective.Opposite(), tp)
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

func (s *State) IncomingMaxData(perspective protocol.Perspective) protocol.ByteCount {
	if perspective == protocol.PerspectiveClient {
		return s.ClientDirectionMaxData
	} else {
		return s.ServerDirectionMaxData
	}
}

func (s *State) OutgoingMaxData(perspective protocol.Perspective) protocol.ByteCount {
	return s.IncomingMaxData(perspective.Opposite())
}

func (s *State) SetIncomingMaxData(perspective protocol.Perspective, maxData protocol.ByteCount) {
	if perspective == protocol.PerspectiveClient {
		s.ClientDirectionMaxData = maxData
	} else {
		s.ServerDirectionMaxData = maxData
	}
}

func (s *State) SetOutgoingMaxData(perspective protocol.Perspective, maxData protocol.ByteCount) {
	s.SetIncomingMaxData(perspective.Opposite(), maxData)
}

// Clone
// TODO deep copy
func (s *State) Clone() *State {
	return &*s
}

func (s *State) BytesSent(perspective protocol.Perspective) protocol.ByteCount {
	if perspective == protocol.PerspectiveClient {
		return s.ServerDirectionBytes
	} else {
		return s.ClientDirectionBytes
	}
}

func (s *State) SetBytesSent(perspective protocol.Perspective, sent protocol.ByteCount) {
	if perspective == protocol.PerspectiveClient {
		s.ServerDirectionBytes = sent
	} else {
		s.ClientDirectionBytes = sent
	}
}

func (s *State) Serialize() ([]byte, error) {
	return json.Marshal(s)
}

func (s *State) CongestionWindow(perspective protocol.Perspective) protocol.ByteCount {
	if perspective == protocol.PerspectiveClient {
		return s.ClientCongestionWindow
	} else {
		return s.ServerCongestionWindow
	}
}

func (s *State) SetCongestionWindow(window protocol.ByteCount, perspective protocol.Perspective) {
	if perspective == protocol.PerspectiveClient {
		s.ClientCongestionWindow = window
	} else {
		s.ServerCongestionWindow = window
	}
}

func (s *State) ConnIDLen(p logging.Perspective) int {
	if p == logging.PerspectiveClient {
		for _, value := range s.ClientConnectionIDs {
			return len(value.ConnectionID)
		}
	} else {
		for _, value := range s.ServerConnectionIDs {
			return len(value.ConnectionID)
		}
	}
	panic("unexpected empty set")
}

func (s *State) FromPerspective(perspective protocol.Perspective) *StateFromPerspective {
	return &StateFromPerspective{
		state:       s,
		perspective: perspective,
	}
}
