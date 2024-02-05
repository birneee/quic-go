//go:generate msgp
package handover

import (
	"bytes"
	"encoding/gob"
	"errors"
	jsoniter "github.com/json-iterator/go"
	"github.com/quic-go/quic-go/internal/indi_utils"
	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/tinylib/msgp/msgp"
	"math"
	"net"
	"strconv"
)

type ConnectionIDSequenceNumber uint64

func (n *ConnectionIDSequenceNumber) MsgpStrMapKey() string {
	return strconv.FormatInt(int64(*n), 10)
}

func (n *ConnectionIDSequenceNumber) MsgpFromStrMapKey(str string) msgp.NonStrMapKey {
	i, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		panic(err)
	}
	*n = ConnectionIDSequenceNumber(i)
	return n
}

func (n *ConnectionIDSequenceNumber) MsgpStrMapKeySize() int {
	return indi_utils.Base10Digits(*n)
}

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
	Version             protocol.Version
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
	ClientTransportParameters TransportParameters
	// TODO only include non-default parameters
	ServerTransportParameters TransportParameters
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
	// in byte
	ClientCongestionWindow *int64
	// in byte
	ServerCongestionWindow *int64
	// in ms
	RTT *int64
	// max stream id
	MaxClientUniStream int64
	// max stream id
	MaxServerUniStream int64
	// max stream id
	MaxClientBidiStream int64
	// max stream id
	MaxServerBidiStream  int64
	ALPN                 string
	ClientReceivedRanges [][2]int64
	ServerReceivedRanges [][2]int64
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

// Clone
// TODO deep copy
func (s *State) Clone() *State {
	return &*s
}

func (s *State) Serialize() ([]byte, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Marshal(s)
}

func (s *State) Parse(buf []byte) (*State, error) {
	if s == nil {
		s = &State{}
	}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(buf, s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *State) SerializeGob() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 10000))
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(s)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *State) ParseGob(buf []byte) (*State, error) {
	if s == nil {
		s = &State{}
	}
	reader := bytes.NewReader(buf)
	decoder := gob.NewDecoder(reader)
	err := decoder.Decode(s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *State) SerializeMsgp() ([]byte, error) {
	return s.MarshalMsg(nil)
}

func (s *State) ParseMsgp(buf []byte) (*State, error) {
	if s == nil {
		s = &State{}
	}
	_, err := s.UnmarshalMsg(buf)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *State) ConnIDLen(p protocol.Perspective) int {
	if p == protocol.PerspectiveClient {
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

func (s *State) FromPerspective(perspective protocol.Perspective) StateFromPerspective {
	return StateFromPerspective{
		state:       s,
		perspective: perspective,
	}
}
