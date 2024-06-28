//go:generate msgp
package qstate

type Crypto struct {
	KeyPhase uint64 `msg:"key_phase" json:"key_phase" cbor:"1,keyasint"`
	// id of the used TLS 1.3 cipher suites.
	// see RFC 8446 Appendix B.4. Cipher Suites.
	// e.g."AES_128_GCM_SHA256"
	TlsCipher string `msg:"tls_cipher" json:"tls_cipher" cbor:"2,keyasint"`
	// used for header protection sent by peer.
	// see RFC 9001 Section 5.4 Header Protection.
	RemoteHeaderProtectionKey HexByteSlice `msg:"remote_header_protection_key" json:"remote_header_protection_key" cbor:"3,keyasint"`
	// used for header protection sent to peer.
	// see RFC 9001 Section 5.4 Header Protection.
	HeaderProtectionKey HexByteSlice `msg:"header_protection_key" json:"header_protection_key" cbor:"4,keyasint"`
	// secret used on packets sent from peer.
	RemoteTrafficSecret HexByteSlice `msg:"remote_traffic_secret" json:"remote_traffic_secret" cbor:"5,keyasint"`
	// secret used on packets sent to peer.
	TrafficSecret HexByteSlice `msg:"traffic_secret" json:"traffic_secret" cbor:"6,keyasint"`
}

func (c *Crypto) ChangeVantagePoint() Crypto {
	newCrypto := Crypto{
		KeyPhase:                  c.KeyPhase,
		TlsCipher:                 c.TlsCipher,
		RemoteHeaderProtectionKey: c.HeaderProtectionKey,
		HeaderProtectionKey:       c.RemoteHeaderProtectionKey,
		RemoteTrafficSecret:       c.TrafficSecret,
		TrafficSecret:             c.RemoteTrafficSecret,
	}
	return newCrypto
}