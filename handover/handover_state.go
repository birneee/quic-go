package handover

import (
	"errors"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/wire"
	"net"
	"strconv"
)

// State is used to handover QUIC connection
type State struct {
	SrcConnectionID             protocol.ConnectionID
	DestConnectionID            protocol.ConnectionID
	Version                     protocol.VersionNumber
	KeyPhase                    protocol.KeyPhase
	SuiteId                     uint16
	FirstRcvTrafficSecret       []byte
	FirstSendTrafficSecret      []byte
	RcvTrafficSecret            []byte
	SendTrafficSecret           []byte
	RemoteAddress               string
	HighestSentPacketNumber     protocol.PacketNumber
	HighestReceivedPacketNumber protocol.PacketNumber
	SrcTransportParameters      wire.TransportParameters
	DestTransportParameters     wire.TransportParameters
}

func (s *State) GetParsedRemoteAddress() (*net.UDPAddr, error) {
	ipString, portString, err := net.SplitHostPort(s.RemoteAddress)
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
