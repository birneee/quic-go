package xse

import (
	"crypto/cipher"
	"encoding/binary"
	"fmt"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/qtls"
	"github.com/lucas-clemente/quic-go/logging"
	"math/bits"
	"sync"
)

const (
	// label to export XSE-QUIC master secret from TLS exporter_master_secret
	// (see RFC8446 Section 7.5)
	xseMasterSecretLabel          = "xse master"
	clientApplicationTrafficLabel = "c ap traffic"
	serverApplicationTrafficLabel = "s ap traffic"
	quicKeyUpdateLabel            = "quic ku"
)

func applicationTrafficLabel(perspective protocol.Perspective) string {
	if perspective == logging.PerspectiveClient {
		return clientApplicationTrafficLabel
	} else {
		return serverApplicationTrafficLabel
	}
}

type cryptoSetup struct {
	// protects all members below
	mutex sync.Mutex
	suite *qtls.CipherSuiteTLS13
	// reused buffer for storing the nonce
	nonceBuf                        []byte
	rcvXseApplicationTrafficSecret  []byte
	sendXseApplicationTrafficSecret []byte
	rcvAead                         cipher.AEAD
	sendAead                        cipher.AEAD
}

// NewCryptoSetup creates a XSE-QUIC crypto setup
func NewCryptoSetup(rcvXseApplicationTrafficSecret []byte, sendXseApplicationTrafficSecret []byte, suite *qtls.CipherSuiteTLS13) *cryptoSetup {
	c := &cryptoSetup{}

	c.suite = suite

	c.rcvXseApplicationTrafficSecret = rcvXseApplicationTrafficSecret
	c.sendXseApplicationTrafficSecret = sendXseApplicationTrafficSecret

	c.rcvAead = qtls.CreateAEAD(suite, c.rcvXseApplicationTrafficSecret)
	c.sendAead = qtls.CreateAEAD(suite, c.sendXseApplicationTrafficSecret)

	c.nonceBuf = make([]byte, c.rcvAead.NonceSize())
	return c
}

// NewCryptoSetupFromConn creates a XSE-QUIC crypto setup
// TODO 0-RTT is not supported, block until 1-RTT key is available
// TODO improve error handling
func NewCryptoSetupFromConn(conn *qtls.Conn, perspective protocol.Perspective) *cryptoSetup {
	c := &cryptoSetup{}
	c.mutex.Lock()
	go func() {
		// TODO allow 0-RTT XSE-QUIC streams
		// ConnectionState() blocks until handshake is done
		cs := conn.ConnectionState()
		c.suite = qtls.CipherSuiteTLS13ByID(cs.CipherSuite)
		xseMasterSecret, err := (&cs).ExportKeyingMaterial(xseMasterSecretLabel, nil, c.suite.Hash.Size())
		if err != nil {
			panic(fmt.Errorf("failed to export xse_master_secret: %w", err))
		}

		c.rcvXseApplicationTrafficSecret = qtls.DeriveSecret(c.suite, xseMasterSecret, applicationTrafficLabel(perspective.Opposite()), nil)
		c.sendXseApplicationTrafficSecret = qtls.DeriveSecret(c.suite, xseMasterSecret, applicationTrafficLabel(perspective), nil)

		c.rcvAead = qtls.CreateAEAD(c.suite, c.rcvXseApplicationTrafficSecret)
		c.sendAead = qtls.CreateAEAD(c.suite, c.sendXseApplicationTrafficSecret)

		c.nonceBuf = make([]byte, c.rcvAead.NonceSize())
		c.mutex.Unlock()
	}()

	return c
}

var _ CryptoSetup = &cryptoSetup{}

// Update trafficSecrets and AEADs for next key phase
// TODO remember old keys, until all streams use the new keys
func (c *cryptoSetup) Update() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.rcvXseApplicationTrafficSecret = qtls.DeriveSecret(c.suite, c.rcvXseApplicationTrafficSecret, quicKeyUpdateLabel, nil)
	c.sendXseApplicationTrafficSecret = qtls.DeriveSecret(c.suite, c.sendXseApplicationTrafficSecret, quicKeyUpdateLabel, nil)

	c.rcvAead = qtls.CreateAEAD(c.suite, c.rcvXseApplicationTrafficSecret)
	c.sendAead = qtls.CreateAEAD(c.suite, c.sendXseApplicationTrafficSecret)
}

func (c *cryptoSetup) Seal(dst []byte, plaintext []byte, sid protocol.StreamID, rn RecordNumber) []byte {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	binary.BigEndian.PutUint64(c.nonceBuf[len(c.nonceBuf)-8:], bits.Reverse64(uint64(sid))^uint64(rn))
	buf := c.sendAead.Seal(dst, c.nonceBuf, plaintext, nil)
	return buf
}

func (c *cryptoSetup) Open(dst []byte, ciphertext RecordEncryptedPayload, sid protocol.StreamID, rn RecordNumber) ([]byte, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	binary.BigEndian.PutUint64(c.nonceBuf[len(c.nonceBuf)-8:], bits.Reverse64(uint64(sid))^uint64(rn))
	return c.rcvAead.Open(dst, c.nonceBuf, ciphertext, nil)
}

func (c *cryptoSetup) EncryptedRecordPayloadLength(payload DecryptedPayloadLength) uint32 {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return uint32(payload) + uint32(c.sendAead.Overhead())
}

func (c *cryptoSetup) MaxEncryptedRecordPayloadLength() uint32 {
	// is locked by EncryptedRecordPayloadLength
	return c.EncryptedRecordPayloadLength(MaxDecryptedPayloadLength)
}
