package handshake

import (
	"crypto/cipher"
	"github.com/quic-go/quic-go/internal/protocol"
)

type RecreatableAEAD interface {
	cipher.AEAD
	Suite() *cipherSuite
	TrafficSecret() []byte
	Version() protocol.Version
}

type recreatableAEAD struct {
	inner         cipher.AEAD
	suite         *cipherSuite
	trafficSecret []byte
	version       protocol.Version
}

var _ RecreatableAEAD = &recreatableAEAD{}

func NewRecreatableAEAD(suite *cipherSuite, trafficSecret []byte, version protocol.Version) RecreatableAEAD {
	return &recreatableAEAD{
		inner:         createAEAD(suite, trafficSecret, version),
		suite:         suite,
		trafficSecret: trafficSecret,
		version:       version,
	}
}

func NewRecreatableAEADNoAlloc(suite *cipherSuite, trafficSecret []byte, version protocol.Version, key []byte, iv []byte, tmp []byte) RecreatableAEAD {
	return &recreatableAEAD{
		inner:         createAEADNoAlloc(suite, trafficSecret, version, key, iv, tmp),
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

func (r *recreatableAEAD) Suite() *cipherSuite {
	return r.suite
}

func (r *recreatableAEAD) TrafficSecret() []byte {
	return r.trafficSecret
}

func (r *recreatableAEAD) Version() protocol.Version {
	return r.version
}
