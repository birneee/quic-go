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
