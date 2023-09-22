//go:generate msgp
package handover

import (
	"github.com/quic-go/quic-go/internal/protocol"
)

type UniStreamState struct {
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

func (u *UniStreamState) PendingIncomingFrames(perspective protocol.Perspective) map[protocol.ByteCount][]byte {
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

func (u *UniStreamState) PendingSentData(perspective protocol.Perspective) map[protocol.ByteCount][]byte {
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

func (u *UniStreamState) SendStreamFromPerspective(perspective protocol.Perspective) SendStreamStateFromPerspective {
	//TODO implement me
	panic("implement me")
}

type SendStreamState interface {
	SetOutgoingOffset(perspective protocol.Perspective, offset protocol.ByteCount)
	SetOutgoingFinOffset(perspective protocol.Perspective, offset protocol.ByteCount)
	SetPendingOutgoingFrames(perspective protocol.Perspective, frames map[protocol.ByteCount][]byte)
	OutgoingOffset(perspective protocol.Perspective) protocol.ByteCount
	PendingSentData(perspective protocol.Perspective) map[protocol.ByteCount][]byte
	WriteFinOffset(perspective protocol.Perspective) protocol.ByteCount
	SetOutgoingMaxData(perspective protocol.Perspective, window protocol.ByteCount)
	OutgoingMaxData(perspective protocol.Perspective) protocol.ByteCount
	SendStreamFromPerspective(perspective protocol.Perspective) SendStreamStateFromPerspective
}

type SendStreamStateFromPerspective interface {
	SetOutgoingOffset(offset protocol.ByteCount)
	SetOutgoingFinOffset(offset protocol.ByteCount)
	SetPendingOutgoingFrames(frames map[protocol.ByteCount][]byte)
	OutgoingOffset() protocol.ByteCount
	PendingSentData() map[protocol.ByteCount][]byte
	WriteFinOffset() protocol.ByteCount
	SetOutgoingMaxData(window protocol.ByteCount)
	OutgoingMaxData() protocol.ByteCount
	OutgoingAcknowledgedOffset() protocol.ByteCount
	SetOutgoingAcknowledgedOffset(offset protocol.ByteCount)
}

type ReceiveStreamState interface {
	SetIncomingOffset(perspective protocol.Perspective, offset protocol.ByteCount)
	SetIncomingFinOffset(perspective protocol.Perspective, offset protocol.ByteCount)
	SetPendingIncomingFrames(perspective protocol.Perspective, frames map[protocol.ByteCount][]byte)
	IncomingOffset(perspective protocol.Perspective) protocol.ByteCount
	PendingIncomingFrames(perspective protocol.Perspective) map[protocol.ByteCount][]byte
	IncomingFinOffset(perspective protocol.Perspective) protocol.ByteCount
	IncomingMaxData(perspective protocol.Perspective) protocol.ByteCount
	SetIncomingMaxData(perspective protocol.Perspective, window protocol.ByteCount)
}
