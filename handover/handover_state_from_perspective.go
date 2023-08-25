package handover

import (
	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/wire"
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

func (s *StateFromPerspective) OwnTransportParameters() *wire.TransportParameters {
	var bytes []byte
	if s.perspective == protocol.PerspectiveClient {
		bytes = s.state.ClientTransportParameters
	} else {
		bytes = s.state.ServerTransportParameters
	}

	tp := wire.TransportParameters{}
	err := tp.Unmarshal(bytes, s.perspective)
	if err != nil {
		panic(err)
	}
	return &tp
}

func (s *StateFromPerspective) PeerTransportParameters() *wire.TransportParameters {
	return s.Opposite().OwnTransportParameters()
}

func (s *StateFromPerspective) Opposite() *StateFromPerspective {
	return &StateFromPerspective{
		state:       s.state,
		perspective: s.perspective.Opposite(),
	}
}
