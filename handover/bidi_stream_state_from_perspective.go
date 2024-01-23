package handover

import (
	"github.com/quic-go/quic-go/internal/protocol"
)

type BidiStreamStateFromPerspective struct {
	state       *BidiStreamState
	perspective protocol.Perspective
}

var _ SendStreamStateFromPerspective = &BidiStreamStateFromPerspective{}

func (s *BidiStreamStateFromPerspective) IncomingOffset() protocol.ByteCount {
	if s.perspective == protocol.PerspectiveClient {
		return s.state.ClientDirectionOffset
	} else {
		return s.state.ServerDirectionOffset
	}
}

func (s *BidiStreamStateFromPerspective) OutgoingOffset() protocol.ByteCount {
	if s.perspective == protocol.PerspectiveClient {
		return s.state.ServerDirectionOffset
	} else {
		return s.state.ClientDirectionOffset
	}
}

func (s *BidiStreamStateFromPerspective) SetIncomingOffset(offset protocol.ByteCount) {
	if s.perspective == protocol.PerspectiveClient {
		s.state.ClientDirectionOffset = offset
	} else {
		s.state.ServerDirectionOffset = offset
	}
}

func (s *BidiStreamStateFromPerspective) SetOutgoingOffset(offset protocol.ByteCount) {
	if s.perspective == protocol.PerspectiveClient {
		s.state.ServerDirectionOffset = offset
	} else {
		s.state.ClientDirectionOffset = offset
	}
}

func (s *BidiStreamStateFromPerspective) OutgoingAcknowledgedOffset() protocol.ByteCount {
	if s.perspective == protocol.PerspectiveClient {
		return s.state.ServerDirectionAcknowledgedOffset
	} else {
		return s.state.ClientDirectionAcknowledgedOffset
	}
}

func (s *BidiStreamStateFromPerspective) SetOutgoingAcknowledgedOffset(offset protocol.ByteCount) {
	if s.perspective == protocol.PerspectiveClient {
		s.state.ServerDirectionAcknowledgedOffset = offset
	} else {
		s.state.ClientDirectionAcknowledgedOffset = offset
	}
}

func (s *BidiStreamStateFromPerspective) PendingIncomingFrames() map[protocol.ByteCount][]byte {
	if s.perspective == protocol.PerspectiveClient {
		return s.state.ClientDirectionPendingFrames
	} else {
		return s.state.ServerDirectionPendingFrames
	}
}

func (s *BidiStreamStateFromPerspective) PendingOutgoingFrames() map[protocol.ByteCount][]byte {
	if s.perspective == protocol.PerspectiveClient {
		return s.state.ServerDirectionPendingFrames
	} else {
		return s.state.ClientDirectionPendingFrames
	}
}

func (s *BidiStreamStateFromPerspective) PutBack(offset protocol.ByteCount, data []byte) {
	s.PendingIncomingFrames()[offset] = data
	s.SetIncomingOffset(min(s.IncomingOffset(), offset))
}

func (s *BidiStreamStateFromPerspective) SetOutgoingFinOffset(offset protocol.ByteCount) {
	if s.perspective == protocol.PerspectiveClient {
		s.state.ServerDirectionFinOffset = offset
	} else {
		s.state.ClientDirectionFinOffset = offset
	}
}

func (s *BidiStreamStateFromPerspective) SetPendingOutgoingFrames(frames map[protocol.ByteCount][]byte) {
	if s.perspective == protocol.PerspectiveClient {
		s.state.ServerDirectionPendingFrames = frames
	} else {
		s.state.ClientDirectionPendingFrames = frames
	}
}

func (s *BidiStreamStateFromPerspective) PendingSentData() map[protocol.ByteCount][]byte {
	if s.perspective == protocol.PerspectiveClient {
		return s.state.ServerDirectionPendingFrames
	} else {
		return s.state.ClientDirectionPendingFrames
	}
}

func (s *BidiStreamStateFromPerspective) WriteFinOffset() protocol.ByteCount {
	if s.perspective == protocol.PerspectiveClient {
		return s.state.ServerDirectionFinOffset
	} else {
		return s.state.ClientDirectionFinOffset
	}
}

func (s *BidiStreamStateFromPerspective) SetOutgoingMaxData(maxData protocol.ByteCount) {
	if s.perspective == protocol.PerspectiveClient {
		s.state.ServerDirectionMaxData = maxData
	} else {
		s.state.ClientDirectionMaxData = maxData
	}
}

func (s *BidiStreamStateFromPerspective) OutgoingMaxData() protocol.ByteCount {
	if s.perspective == protocol.PerspectiveClient {
		return s.state.ServerDirectionMaxData
	} else {
		return s.state.ClientDirectionMaxData
	}
}

func (s *BidiStreamStateFromPerspective) SetIncomingFinOffset(offset protocol.ByteCount) {
	if s.perspective == protocol.PerspectiveClient {
		s.state.ClientDirectionFinOffset = offset
	} else {
		s.state.ServerDirectionFinOffset = offset
	}
}

func (s *BidiStreamStateFromPerspective) SetPendingIncomingFrames(frames map[protocol.ByteCount][]byte) {
	if s.perspective == protocol.PerspectiveClient {
		s.state.ClientDirectionPendingFrames = frames
	} else {
		s.state.ServerDirectionPendingFrames = frames
	}
}

func (s *BidiStreamStateFromPerspective) IncomingFinOffset() protocol.ByteCount {
	if s.perspective == protocol.PerspectiveClient {
		return s.state.ClientDirectionFinOffset
	} else {
		return s.state.ServerDirectionFinOffset
	}
}

func (s *BidiStreamStateFromPerspective) IncomingMaxData() protocol.ByteCount {
	if s.perspective == protocol.PerspectiveClient {
		return s.state.ClientDirectionMaxData
	} else {
		return s.state.ServerDirectionMaxData
	}
}

func (s *BidiStreamStateFromPerspective) SetIncomingMaxData(maxData protocol.ByteCount) {
	if s.perspective == protocol.PerspectiveClient {
		s.state.ClientDirectionMaxData = maxData
	} else {
		s.state.ServerDirectionMaxData = maxData
	}
}

func (s *BidiStreamStateFromPerspective) SetIncomingAcknowledgedOffset(offset protocol.ByteCount) {
	if s.perspective == protocol.PerspectiveClient {
		s.state.ClientDirectionAcknowledgedOffset = offset
	} else {
		s.state.ServerDirectionAcknowledgedOffset = offset
	}
}
