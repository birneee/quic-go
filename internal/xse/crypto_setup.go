package xse

import (
	"crypto/cipher"
	"encoding/binary"
	"fmt"
	"github.com/lucas-clemente/quic-go/internal/handshake"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/qtls"
	"sync"
)

// label to export XSE-QUIC server key from TLS exporter_master_secret
// (see RFC8446 Section 7.5)
const serverXseLabel = "s xse"

// label to export XSE-QUIC client key from TLS exporter_master_secret
// (see RFC8446 Section 7.5)
const clientXseLabel = "c xse"

type baseCryptoSetup struct {
	nonceBuf []byte
	rcvAead  cipher.AEAD
	sendAead cipher.AEAD
}

// NewCryptoSetup creates a XSE-QUIC crypto setup
func NewCryptoSetup(rcvAead cipher.AEAD, sendAead cipher.AEAD) *baseCryptoSetup {
	if rcvAead.NonceSize() != sendAead.NonceSize() {
		panic("AEAD nonce sizes are different")
	}
	return &baseCryptoSetup{
		nonceBuf: make([]byte, rcvAead.NonceSize()),
		rcvAead:  rcvAead,
		sendAead: sendAead,
	}
}

var _ CryptoSetup = &baseCryptoSetup{}

func (c *baseCryptoSetup) Seal(dst []byte, src []byte, sid protocol.StreamID, rn RecordNumber) []byte {
	binary.BigEndian.PutUint64(c.nonceBuf[len(c.nonceBuf)-8:], uint64(sid)^uint64(rn))
	buf := c.sendAead.Seal(dst, c.nonceBuf, src, nil)
	return buf
}

func (c *baseCryptoSetup) Open(encryptedPayload RecordEncryptedPayload, sid protocol.StreamID, rn RecordNumber) ([]byte, error) {
	binary.BigEndian.PutUint64(c.nonceBuf[len(c.nonceBuf)-8:], uint64(sid)^uint64(rn))
	return c.rcvAead.Open(encryptedPayload[:0], c.nonceBuf, encryptedPayload, nil)
}

func (c *baseCryptoSetup) EncryptedRecordPayloadLength(payload DecryptedPayloadLength) uint32 {
	return uint32(payload) + uint32(c.sendAead.Overhead())
}

func (c *baseCryptoSetup) MaxEncryptedRecordPayloadLength() uint32 {
	return c.EncryptedRecordPayloadLength(MaxDecryptedPayloadLength)
}

type cryptoSetup struct {
	conn         *qtls.Conn
	perspective  protocol.Perspective
	initAeadOnce sync.Once
	*baseCryptoSetup
}

var _ CryptoSetup = &cryptoSetup{}

// NewCryptoSetupFromConn creates a XSE-QUIC crypto setup
func NewCryptoSetupFromConn(conn *qtls.Conn, perspective protocol.Perspective) CryptoSetup {
	return &cryptoSetup{
		conn:        conn,
		perspective: perspective,
	}
}

//TODO improve error handling
func (c *cryptoSetup) initAead() {
	state := c.conn.ConnectionState()
	if !state.HandshakeComplete {
		panic("handshake must be completed")
	}
	suite := qtls.CipherSuiteTLS13ByID(state.CipherSuite)
	serverKey, err := (&state).ExportKeyingMaterial(serverXseLabel, nil, suite.KeyLen)
	if err != nil {
		panic(fmt.Errorf("failed to export key: %w", err))
	}
	serverAead := handshake.CreateAEAD(suite, serverKey)
	clientKey, err := (&state).ExportKeyingMaterial(clientXseLabel, nil, suite.KeyLen)
	if err != nil {
		panic(fmt.Errorf("failed to export key: %w", err))
	}
	clientAead := handshake.CreateAEAD(suite, clientKey)
	if c.perspective == protocol.PerspectiveClient {
		c.baseCryptoSetup = NewCryptoSetup(serverAead, clientAead)
	} else {
		c.baseCryptoSetup = NewCryptoSetup(clientAead, serverAead)

	}
}

func (c *cryptoSetup) Seal(dst []byte, src []byte, sid protocol.StreamID, rn RecordNumber) []byte {
	c.initAeadOnce.Do(c.initAead)
	return c.baseCryptoSetup.Seal(dst, src, sid, rn)
}

func (c *cryptoSetup) Open(payload RecordEncryptedPayload, id protocol.StreamID, number RecordNumber) ([]byte, error) {
	c.initAeadOnce.Do(c.initAead)
	return c.baseCryptoSetup.Open(payload, id, number)

}

func (c *cryptoSetup) EncryptedRecordPayloadLength(length DecryptedPayloadLength) uint32 {
	c.initAeadOnce.Do(c.initAead)
	return c.baseCryptoSetup.EncryptedRecordPayloadLength(length)
}

func (c *cryptoSetup) MaxEncryptedRecordPayloadLength() uint32 {
	c.initAeadOnce.Do(c.initAead)
	return c.EncryptedRecordPayloadLength(MaxDecryptedPayloadLength)
}
