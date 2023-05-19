package handover

import (
	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/logging"
)

type UniStreamState struct {
	ID logging.StreamID
	// offset until stream data is acknowledged or read by application layer
	Offset protocol.ByteCount
	// MaxByteCount if not known yet
	FinOffset     protocol.ByteCount
	PendingFrames map[protocol.ByteCount][]byte
	MaxData       protocol.ByteCount
}

func (u *UniStreamState) SetIncomingOffset(perspective protocol.Perspective, offset protocol.ByteCount) {
	u.Offset = offset
}

func (u *UniStreamState) SetIncomingFinOffset(perspective protocol.Perspective, offset protocol.ByteCount) {
	u.FinOffset = offset
}

func (u *UniStreamState) SetPendingIncomingFrames(perspective protocol.Perspective, frames map[protocol.ByteCount][]byte) {
	u.PendingFrames = frames
}

func (u *UniStreamState) IncomingOffset(perspective protocol.Perspective) protocol.ByteCount {
	return u.Offset
}

func (u *UniStreamState) PendingIncomingFrames(perspective protocol.Perspective) map[logging.ByteCount][]byte {
	return u.PendingFrames
}

func (u *UniStreamState) IncomingFinOffset(perspective protocol.Perspective) protocol.ByteCount {
	return u.FinOffset
}

func (u *UniStreamState) IncomingMaxData(perspective protocol.Perspective) protocol.ByteCount {
	return u.MaxData
}

func (u *UniStreamState) SetIncomingMaxData(perspective protocol.Perspective, window protocol.ByteCount) {
	u.MaxData = window
}

func (u *UniStreamState) SetOutgoingOffset(perspective protocol.Perspective, offset protocol.ByteCount) {
	u.Offset = offset
}

func (u *UniStreamState) SetOutgoingFinOffset(perspective protocol.Perspective, offset protocol.ByteCount) {
	u.FinOffset = offset
}

func (u *UniStreamState) SetPendingOutgoingFrames(perspective protocol.Perspective, frames map[protocol.ByteCount][]byte) {
	u.PendingFrames = frames
}

func (u *UniStreamState) OutgoingOffset(perspective protocol.Perspective) protocol.ByteCount {
	return u.Offset
}

func (u *UniStreamState) PendingSentData(perspective protocol.Perspective) map[logging.ByteCount][]byte {
	return u.PendingFrames
}

func (u *UniStreamState) WriteFinOffset(perspective protocol.Perspective) protocol.ByteCount {
	return u.FinOffset
}

func (u *UniStreamState) SetOutgoingMaxData(perspective protocol.Perspective, window protocol.ByteCount) {
	u.MaxData = window
}

func (u *UniStreamState) OutgoingMaxData(perspective protocol.Perspective) protocol.ByteCount {
	return u.MaxData
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
	ClientDirectionMaxData       protocol.ByteCount
	ServerDirectionMaxData       protocol.ByteCount
}

var _ SendStreamState = &BidiStreamState{}

type SendStreamState interface {
	SetOutgoingOffset(perspective protocol.Perspective, offset protocol.ByteCount)
	SetOutgoingFinOffset(perspective protocol.Perspective, offset protocol.ByteCount)
	SetPendingOutgoingFrames(perspective protocol.Perspective, frames map[protocol.ByteCount][]byte)
	OutgoingOffset(perspective protocol.Perspective) protocol.ByteCount
	PendingSentData(perspective protocol.Perspective) map[logging.ByteCount][]byte
	WriteFinOffset(perspective protocol.Perspective) protocol.ByteCount
	SetOutgoingMaxData(perspective protocol.Perspective, window protocol.ByteCount)
	OutgoingMaxData(perspective protocol.Perspective) protocol.ByteCount
}

type ReceiveStreamState interface {
	SetIncomingOffset(perspective protocol.Perspective, offset protocol.ByteCount)
	SetIncomingFinOffset(perspective protocol.Perspective, offset protocol.ByteCount)
	SetPendingIncomingFrames(perspective protocol.Perspective, frames map[protocol.ByteCount][]byte)
	IncomingOffset(perspective protocol.Perspective) protocol.ByteCount
	PendingIncomingFrames(perspective protocol.Perspective) map[logging.ByteCount][]byte
	IncomingFinOffset(perspective protocol.Perspective) protocol.ByteCount
	IncomingMaxData(perspective protocol.Perspective) protocol.ByteCount
	SetIncomingMaxData(perspective protocol.Perspective, window protocol.ByteCount)
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
