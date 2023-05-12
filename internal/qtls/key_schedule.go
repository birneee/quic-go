package qtls

import (
	"crypto"
	"crypto/cipher"
	"encoding/binary"
	"golang.org/x/crypto/hkdf"
	"hash"
)

func CreateAEAD(suite *CipherSuiteTLS13, trafficSecret []byte) cipher.AEAD {
	key := HkdfExpandLabel(suite.Hash, trafficSecret, []byte{}, "quic key", suite.KeyLen)
	iv := HkdfExpandLabel(suite.Hash, trafficSecret, []byte{}, "quic iv", suite.IVLen())
	return suite.AEAD(key, iv)
}

// DeriveSecret implements Derive-Secret from RFC 8446, Section 7.1.
func DeriveSecret(suite *CipherSuiteTLS13, secret []byte, label string, transcript hash.Hash) []byte {
	if transcript == nil {
		transcript = suite.Hash.New()
	}
	return HkdfExpandLabel(suite.Hash, secret, transcript.Sum(nil), label, suite.Hash.Size())
}

// HkdfExpandLabel HKDF expands a label.
// Since this implementation avoids using a cryptobyte.Builder, it is about 15% faster than the
// hkdfExpandLabel in the standard library.
func HkdfExpandLabel(hash crypto.Hash, secret, context []byte, label string, length int) []byte {
	b := make([]byte, 3, 3+6+len(label)+1+len(context))
	binary.BigEndian.PutUint16(b, uint16(length))
	b[2] = uint8(6 + len(label))
	b = append(b, []byte("tls13 ")...)
	b = append(b, []byte(label)...)
	b = b[:3+6+len(label)+1]
	b[3+6+len(label)] = uint8(len(context))
	b = append(b, context...)

	out := make([]byte, length)
	n, err := hkdf.Expand(hash.New, secret, b).Read(out)
	if err != nil || n != length {
		panic("quic: HKDF-Expand-Label invocation failed unexpectedly")
	}
	return out
}
