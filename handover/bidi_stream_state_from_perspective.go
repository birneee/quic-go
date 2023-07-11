package handover

import (
	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/utils"
	"github.com/quic-go/quic-go/logging"
)

type BidiStreamStateFromPerspective struct {
	state       *BidiStreamState
	perspective protocol.Perspective
}

func (s *BidiStreamStateFromPerspective) IncomingOffset() protocol.ByteCount {
	if s.perspective == logging.PerspectiveClient {
		return s.state.ClientDirectionOffset
	} else {
		return s.state.ServerDirectionOffset
	}
}

func (s *BidiStreamStateFromPerspective) OutgoingOffset() protocol.ByteCount {
	if s.perspective == logging.PerspectiveClient {
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

func (s *BidiStreamStateFromPerspective) PendingIncomingFrames() map[logging.ByteCount][]byte {
	if s.perspective == logging.PerspectiveClient {
		return s.state.ClientDirectionPendingFrames
	} else {
		return s.state.ServerDirectionPendingFrames
	}
}

func (s *BidiStreamStateFromPerspective) PendingOutgoingFrames() map[logging.ByteCount][]byte {
	if s.perspective == logging.PerspectiveClient {
		return s.state.ServerDirectionPendingFrames
	} else {
		return s.state.ClientDirectionPendingFrames
	}
}

func (s *BidiStreamStateFromPerspective) PutBack(offset protocol.ByteCount, data []byte) {
	s.PendingIncomingFrames()[offset] = data
	s.SetIncomingOffset(utils.Min(s.IncomingOffset(), offset))
}
