package handover

import (
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/logging"
)

type UniStreamState struct {
	ID logging.StreamID
	// offset until stream data is acknowledged
	Offset protocol.ByteCount
	// -1 if not known yet
	FinOffset     protocol.ByteCount
	PendingFrames map[protocol.ByteCount][]byte
}

type BidiStreamState struct {
	ID logging.StreamID
	// offset until stream data is acknowledged or read by application layer
	ClientDirectionOffset protocol.ByteCount
	// offset until stream data is acknowledged or read by application layer
	ServerDirectionOffset protocol.ByteCount
	// MaxByteCount if not known yet
	ClientDirectionFinOffset protocol.ByteCount
	// MaxByteCount if not known yet
	ServerDirectionFinOffset     protocol.ByteCount
	ClientDirectionPendingFrames map[protocol.ByteCount][]byte
	ServerDirectionPendingFrames map[protocol.ByteCount][]byte
}

func NewBidiStreamStateFromPerspective(perspective protocol.Perspective, id logging.StreamID, incomingOffset logging.ByteCount, outgoingOffset logging.ByteCount, incomingFinOffset protocol.ByteCount, outgoingFinOffset protocol.ByteCount, pendingIncomingFrames map[protocol.ByteCount][]byte, pendingOutgoingFrames map[protocol.ByteCount][]byte) BidiStreamState {
	ss := BidiStreamState{
		ID: id,
	}
	ss.SetIncomingOffset(perspective, incomingOffset)
	ss.SetOutgoingOffset(perspective, outgoingOffset)
	ss.SetIncomingFinOffset(perspective, incomingFinOffset)
	ss.SetOutgoingFinOffset(perspective, outgoingFinOffset)
	ss.SetPendingIncomingFrames(perspective, pendingIncomingFrames)
	ss.SetPendingOutgoingFrames(perspective, pendingOutgoingFrames)
	return ss
}

// getter and setter for client and server perspective

func (s *BidiStreamState) SetIncomingOffset(perspective protocol.Perspective, offset protocol.ByteCount) {
	if perspective == protocol.PerspectiveClient {
		s.ClientDirectionOffset = offset
	} else {
		s.ServerDirectionOffset = offset
	}
}

func (s *BidiStreamState) SetOutgoingOffset(perspective protocol.Perspective, offset protocol.ByteCount) {
	s.SetIncomingOffset(perspective.Opposite(), offset)
}

func (s *BidiStreamState) SetPendingIncomingFrames(perspective protocol.Perspective, data map[protocol.ByteCount][]byte) {
	if perspective == protocol.PerspectiveClient {
		s.ClientDirectionPendingFrames = data
	} else {
		s.ServerDirectionPendingFrames = data
	}
}

func (s *BidiStreamState) SetPendingOutgoingFrames(perspective protocol.Perspective, data map[protocol.ByteCount][]byte) {
	s.SetPendingIncomingFrames(perspective.Opposite(), data)
}

func (s *BidiStreamState) ReadOffset(perspective protocol.Perspective) protocol.ByteCount {
	if perspective == logging.PerspectiveClient {
		return s.ClientDirectionOffset
	} else {
		return s.ServerDirectionOffset
	}
}

func (s *BidiStreamState) WriteOffset(perspective protocol.Perspective) protocol.ByteCount {
	return s.ReadOffset(perspective.Opposite())
}

func (s *BidiStreamState) ReadFinOffset(perspective protocol.Perspective) protocol.ByteCount {
	if perspective == logging.PerspectiveClient {
		return s.ClientDirectionFinOffset
	} else {
		return s.ServerDirectionFinOffset
	}
}

func (s *BidiStreamState) WriteFinOffset(perspective protocol.Perspective) protocol.ByteCount {
	return s.ReadFinOffset(perspective.Opposite())
}

func (s *BidiStreamState) PendingReceivedData(perspective protocol.Perspective) map[logging.ByteCount][]byte {
	if perspective == logging.PerspectiveClient {
		return s.ClientDirectionPendingFrames
	} else {
		return s.ServerDirectionPendingFrames
	}
}

func (s *BidiStreamState) PendingSentData(perspective protocol.Perspective) map[logging.ByteCount][]byte {
	return s.PendingReceivedData(perspective.Opposite())
}

func (s *BidiStreamState) SetIncomingFinOffset(perspective protocol.Perspective, offset protocol.ByteCount) {
	if perspective == protocol.PerspectiveClient {
		s.ClientDirectionFinOffset = offset
	} else {
		s.ServerDirectionFinOffset = offset
	}
}

func (s *BidiStreamState) SetOutgoingFinOffset(perspective protocol.Perspective, offset protocol.ByteCount) {
	s.SetIncomingFinOffset(perspective.Opposite(), offset)
}
