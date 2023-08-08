//go:generate msgp
package handover

import (
	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/logging"
)

type BidiStreamState struct {
	// highest sent or received offset
	ClientDirectionOffset protocol.ByteCount
	// highest sent or received offset
	ServerDirectionOffset protocol.ByteCount
	// offset until stream data is acknowledged or read by application layer
	ServerDirectionAcknowledgedOffset protocol.ByteCount
	// offset until stream data is acknowledged or read by application layer
	ClientDirectionAcknowledgedOffset protocol.ByteCount
	// MaxByteCount if not known yet
	ClientDirectionFinOffset protocol.ByteCount
	// MaxByteCount if not known yet
	ServerDirectionFinOffset     protocol.ByteCount
	ClientDirectionPendingFrames map[protocol.ByteCount][]byte
	ServerDirectionPendingFrames map[protocol.ByteCount][]byte
	// also required for sending RESET_STREAM frames
	ClientDirectionMaxData protocol.ByteCount
	// also required for sending RESET_STREAM frames
	ServerDirectionMaxData protocol.ByteCount
}

var _ SendStreamState = &BidiStreamState{}
var _ ReceiveStreamState = &BidiStreamState{}

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

func (s *BidiStreamState) IncomingOffset(perspective protocol.Perspective) protocol.ByteCount {
	if perspective == logging.PerspectiveClient {
		return s.ClientDirectionOffset
	} else {
		return s.ServerDirectionOffset
	}
}

func (s *BidiStreamState) OutgoingOffset(perspective protocol.Perspective) protocol.ByteCount {
	return s.IncomingOffset(perspective.Opposite())
}

func (s *BidiStreamState) IncomingFinOffset(perspective protocol.Perspective) protocol.ByteCount {
	if perspective == logging.PerspectiveClient {
		return s.ClientDirectionFinOffset
	} else {
		return s.ServerDirectionFinOffset
	}
}

func (s *BidiStreamState) WriteFinOffset(perspective protocol.Perspective) protocol.ByteCount {
	return s.IncomingFinOffset(perspective.Opposite())
}

func (s *BidiStreamState) PendingIncomingFrames(perspective protocol.Perspective) map[logging.ByteCount][]byte {
	if perspective == logging.PerspectiveClient {
		return s.ClientDirectionPendingFrames
	} else {
		return s.ServerDirectionPendingFrames
	}
}

func (s *BidiStreamState) PendingSentData(perspective protocol.Perspective) map[logging.ByteCount][]byte {
	return s.PendingIncomingFrames(perspective.Opposite())
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

func (s *BidiStreamState) SetIncomingMaxData(perspective protocol.Perspective, maxData protocol.ByteCount) {
	if perspective == protocol.PerspectiveClient {
		s.ClientDirectionMaxData = maxData
	} else {
		s.ServerDirectionMaxData = maxData
	}
}

func (s *BidiStreamState) SetOutgoingMaxData(perspective protocol.Perspective, maxData protocol.ByteCount) {
	s.SetIncomingMaxData(perspective.Opposite(), maxData)
}

func (s *BidiStreamState) IncomingMaxData(perspective protocol.Perspective) protocol.ByteCount {
	if perspective == protocol.PerspectiveClient {
		return s.ClientDirectionMaxData
	} else {
		return s.ServerDirectionMaxData
	}
}

func (s *BidiStreamState) OutgoingMaxData(perspective protocol.Perspective) protocol.ByteCount {
	return s.IncomingMaxData(perspective.Opposite())
}

func (s *BidiStreamState) FromPerspective(perspective protocol.Perspective) *BidiStreamStateFromPerspective {
	return &BidiStreamStateFromPerspective{
		state:       s,
		perspective: perspective,
	}
}

func (s *BidiStreamState) SendStreamFromPerspective(perspective protocol.Perspective) SendStreamStateFromPerspective {
	return s.FromPerspective(perspective)
}
