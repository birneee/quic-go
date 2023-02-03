package path

import (
	"context"
	"crypto/rand"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"net"
	"time"
)

type path struct {
	addr                  *net.UDPAddr
	lastPacketReceiveTime time.Time
	challengeData         [8]byte
	// Do we know that the peer completed address validation yet?
	// Always true for the server.
	peerCompletedAddressValidation bool
	bytesReceived                  protocol.ByteCount
	bytesSent                      protocol.ByteCount
	// Have we validated the peer's address yet?
	// Always true for the client.
	validationCtx       context.Context
	validationCtxCancel context.CancelFunc
	bytesInFlight       protocol.ByteCount
}

var _ Path = &path{}

func NewPath(addr *net.UDPAddr, peerAddressValidated bool) *path {
	validationCtx, validationCtxCancel := context.WithCancel(context.Background())
	p := &path{
		addr:                           addr,
		peerCompletedAddressValidation: false,
		bytesReceived:                  0,
		bytesSent:                      0,
		validationCtx:                  validationCtx,
		validationCtxCancel:            validationCtxCancel,
	}
	_, err := rand.Read(p.challengeData[:])
	if err != nil {
		panic(err)
	}
	if peerAddressValidated {
		validationCtxCancel()
	}
	return p
}

func (p *path) PeerCompletedAddressValidation() bool {
	return p.peerCompletedAddressValidation
}

func (p *path) SetPeerCompletedAddressValidation() {
	p.peerCompletedAddressValidation = true
}

func (p *path) PeerAddressValidated() bool {
	select {
	case _, _ = <-p.validationCtx.Done():
		return true
	default:
		return false
	}
}

func (p *path) SetPeerAddressValidated() {
	p.validationCtxCancel()
}

func (p *path) IncrementBytesReceived(n protocol.ByteCount) {
	p.bytesReceived += n
}

func (p *path) IncrementBytesSent(n protocol.ByteCount) {
	p.bytesSent += n
}

func (p *path) BytesReceived() protocol.ByteCount {
	return p.bytesReceived
}

func (p *path) BytesSent() protocol.ByteCount {
	return p.bytesSent
}

func (p *path) Addr() net.Addr {
	return p.addr
}

func (p *path) ChallengeData() [8]byte {
	return p.challengeData
}

func (p *path) BytesInFlight() protocol.ByteCount {
	return p.bytesInFlight
}

func (p *path) RemoveFromBytesInFlight(length protocol.ByteCount) {
	if length > p.bytesInFlight {
		panic("negative bytes_in_flight")
	}
	p.bytesInFlight -= length
}

func (p *path) AddToBytesInFlight(length protocol.ByteCount) {
	p.bytesInFlight += length
}

func (p *path) ResetBytesInFlight() {
	p.bytesInFlight = 0
}
