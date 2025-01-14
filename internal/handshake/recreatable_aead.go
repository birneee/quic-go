package handshake

import (
	"crypto/cipher"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/qtls"
)

type RecreatableAEAD interface {
	cipher.AEAD
	Suite() *qtls.CipherSuiteTLS13
	TrafficSecret() []byte
	Version() protocol.VersionNumber
}

type recreatableAEAD struct {
	inner         cipher.AEAD
	suite         *qtls.CipherSuiteTLS13
	trafficSecret []byte
	version       protocol.VersionNumber
}

var _ RecreatableAEAD = &recreatableAEAD{}

func NewRecreatableAEAD(suite *qtls.CipherSuiteTLS13, trafficSecret []byte, version protocol.VersionNumber) RecreatableAEAD {
	return &recreatableAEAD{
		inner:         createAEAD(suite, trafficSecret, version),
		suite:         suite,
		trafficSecret: trafficSecret,
		version:       version,
	}
}

func (r *recreatableAEAD) NonceSize() int {
	return r.inner.NonceSize()
}

func (r *recreatableAEAD) Overhead() int {
	return r.inner.Overhead()
}

func (r *recreatableAEAD) Seal(dst, nonce, plaintext, additionalData []byte) []byte {
	return r.inner.Seal(dst, nonce, plaintext, additionalData)
}

func (r *recreatableAEAD) Open(dst, nonce, ciphertext, additionalData []byte) ([]byte, error) {
	return r.inner.Open(dst, nonce, ciphertext, additionalData)
}

func (r *recreatableAEAD) Suite() *qtls.CipherSuiteTLS13 {
	return r.suite
}

func (r *recreatableAEAD) TrafficSecret() []byte {
	return r.trafficSecret
}

func (r *recreatableAEAD) Version() protocol.VersionNumber {
	return r.version
}
