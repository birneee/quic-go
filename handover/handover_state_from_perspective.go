package handover

import (
	"github.com/quic-go/quic-go/internal/protocol"
)

type StateFromPerspective struct {
	state       *State
	perspective protocol.Perspective
}

func (s *StateFromPerspective) NextIncomingBidiStream() protocol.StreamID {
	if s.perspective == protocol.PerspectiveClient {
		return s.state.ServerNextBidiStream
	} else {
		return s.state.ClientNextBidiStream
	}
}

func (s *StateFromPerspective) NextOutgoingBidiStream() protocol.StreamID {
	if s.perspective == protocol.PerspectiveClient {
		return s.state.ClientNextBidiStream
	} else {
		return s.state.ServerNextBidiStream
	}
}

func (s *StateFromPerspective) NextIncomingUniStream() protocol.StreamID {
	if s.perspective == protocol.PerspectiveClient {
		return s.state.ServerNextUniStream
	} else {
		return s.state.ClientNextUniStream
	}
}

func (s *StateFromPerspective) NextOutgoingUniStream() protocol.StreamID {
	if s.perspective == protocol.PerspectiveClient {
		return s.state.ClientNextUniStream
	} else {
		return s.state.ServerNextUniStream
	}
}

func (s *StateFromPerspective) SetNextIncomingBidiStream(value protocol.StreamID) {
	if s.perspective == protocol.PerspectiveClient {
		s.state.ServerNextBidiStream = value
	} else {
		s.state.ClientNextBidiStream = value
	}
}

func (s *StateFromPerspective) SetNextOutgoingBidiStream(value protocol.StreamID) {
	if s.perspective == protocol.PerspectiveClient {
		s.state.ClientNextBidiStream = value
	} else {
		s.state.ServerNextBidiStream = value
	}
}

func (s *StateFromPerspective) SetNextIncomingUniStream(value protocol.StreamID) {
	if s.perspective == protocol.PerspectiveClient {
		s.state.ServerNextUniStream = value
	} else {
		s.state.ClientNextUniStream = value
	}
}

func (s *StateFromPerspective) SetNextOutgoingUniStream(value protocol.StreamID) {
	if s.perspective == protocol.PerspectiveClient {
		s.state.ClientNextUniStream = value
	} else {
		s.state.ServerNextUniStream = value
	}
}

func (s *StateFromPerspective) PutBack(streamID protocol.StreamID, offset protocol.ByteCount, data []byte) {
	if streamID.Type() == protocol.StreamTypeBidi {
		stream := s.state.BidiStreams[streamID].FromPerspective(s.perspective)
		stream.PutBack(offset, data)
	} else {
		panic("implement me")
	}
}

func (s *StateFromPerspective) Version() protocol.VersionNumber {
	return s.state.Version
}

func (s *StateFromPerspective) SetVersion(version protocol.VersionNumber) {
	s.state.Version = version
}

func (s *StateFromPerspective) Perspective() protocol.Perspective {
	return s.perspective
}

func (s *StateFromPerspective) HighestSentPacketNumber() protocol.PacketNumber {
	if s.perspective == protocol.PerspectiveClient {
		return s.state.ServerHighestSentPacketNumber
	} else {
		return s.state.ClientHighestSentPacketNumber
	}
}

func (s *StateFromPerspective) Opposite() *StateFromPerspective {
	return &StateFromPerspective{
		state:       s.state,
		perspective: s.perspective.Opposite(),
	}
}

func (s *StateFromPerspective) OwnTransportParameters() *TransportParameters {
	if s.perspective == protocol.PerspectiveClient {
		return &s.state.ClientTransportParameters
	} else {
		return &s.state.ServerTransportParameters
	}
}

func (s *StateFromPerspective) PeerTransportParameters() *TransportParameters {
	return s.Opposite().OwnTransportParameters()
}

func (s *StateFromPerspective) SetOwnTransportParameters(params TransportParameters) {
	if s.perspective == protocol.PerspectiveClient {
		s.state.ClientTransportParameters = params
	} else {
		s.state.ServerTransportParameters = params
	}
}

func (s *StateFromPerspective) SetPeerTransportParameters(params TransportParameters) {
	s.Opposite().SetOwnTransportParameters(params)
}

func (s *StateFromPerspective) MaxOutgoingUniStream() int64 {
	if s.perspective == protocol.PerspectiveClient {
		return s.state.MaxClientUniStream
	} else {
		return s.state.MaxServerUniStream
	}
}

func (s *StateFromPerspective) MaxOutgoingBidiStream() int64 {
	if s.perspective == protocol.PerspectiveClient {
		return s.state.MaxClientBidiStream
	} else {
		return s.state.MaxServerBidiStream
	}
}

func (s *StateFromPerspective) MaxIncomingUniStream() int64 {
	return s.Opposite().MaxOutgoingUniStream()
}

func (s *StateFromPerspective) MaxIncomingBidiStream() int64 {
	return s.Opposite().MaxOutgoingBidiStream()
}

func (s *StateFromPerspective) SetMaxOutgoingUniStream(i int64) {
	if s.perspective == protocol.PerspectiveClient {
		s.state.MaxClientUniStream = i
	} else {
		s.state.MaxServerUniStream = i
	}
}

func (s *StateFromPerspective) SetMaxOutgoingBidiStream(i int64) {
	if s.perspective == protocol.PerspectiveClient {
		s.state.MaxClientBidiStream = i
	} else {
		s.state.MaxServerBidiStream = i
	}
}

func (s *StateFromPerspective) SetMaxIncomingUniStream(i int64) {
	s.Opposite().SetMaxOutgoingUniStream(i)
}

func (s *StateFromPerspective) SetMaxIncomingBidiStream(i int64) {
	s.Opposite().SetMaxOutgoingBidiStream(i)
}

func (s *StateFromPerspective) SetRTT(rtt *int64) {
	s.state.RTT = rtt
}

func (s *StateFromPerspective) SetCongestionWindow(cw *int64) {
	if s.perspective == protocol.PerspectiveClient {
		s.state.ClientCongestionWindow = cw
	} else {
		s.state.ServerCongestionWindow = cw
	}
}

func (s *StateFromPerspective) RTT() *int64 {
	return s.state.RTT
}

func (s *StateFromPerspective) CongestionWindow() *int64 {
	if s.perspective == protocol.PerspectiveClient {
		return s.state.ClientCongestionWindow
	} else {
		return s.state.ServerCongestionWindow
	}
}
