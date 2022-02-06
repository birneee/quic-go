package xse

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/qtls"
	"github.com/lucas-clemente/quic-go/logging"
	qtls2 "github.com/marten-seemann/qtls-go1-17"
)

const (
	// label to export XSE-QUIC master secret from TLS exporter_master_secret
	// (see RFC8446 Section 7.5)
	xseMasterSecretLabel = "xse master"
)

// label to export XSE-QUIC stream traffic secret from xse_master_secret.
// (see RFC8446 Section 7.5)
// return a label of length 9.
// first char is 'c' if the sender is the client.
// first char is 's' if the sender is the server.
// following bytes are the binary big endian stream id.
func trafficSecretLabel(perspective protocol.Perspective, streamID protocol.StreamID) string {
	buf := make([]byte, 9)
	if perspective == logging.PerspectiveClient {
		buf[0] = byte('c')
	} else {
		buf[0] = byte('s')
	}
	binary.BigEndian.PutUint64(buf[1:], uint64(streamID))
	return string(buf)
}

type cryptoSetup struct {
	perspective          protocol.Perspective
	secretReadyCtx       context.Context
	secretReadyCtxCancel context.CancelFunc
	// must wait for secretReadyCtx before access
	suite *qtls.CipherSuiteTLS13
	// must wait for secretReadyCtx before access
	masterSecret []byte
	tracer       logging.ConnectionTracer
}

func NewCryptoSetup(suite *qtls.CipherSuiteTLS13, masterSecret []byte, perspective protocol.Perspective, tracer logging.ConnectionTracer) CryptoSetup {
	c := &cryptoSetup{
		suite:        suite,
		masterSecret: masterSecret,
		perspective:  perspective,
		tracer:       tracer,
	}
	c.secretReadyCtx, c.secretReadyCtxCancel = context.WithCancel(context.Background())
	c.secretReadyCtxCancel()
	return c
}

func NewCryptoSetupFromConn(quicTls *qtls.Conn, perspective protocol.Perspective, tracer logging.ConnectionTracer) CryptoSetup {
	c := &cryptoSetup{
		perspective: perspective,
		tracer:      tracer,
	}
	c.secretReadyCtx, c.secretReadyCtxCancel = context.WithCancel(context.Background())
	go func() {
		// TODO allow 0-RTT XSE-QUIC streams
		// ConnectionState() blocks until handshake is done
		cs := quicTls.ConnectionState()
		c.suite = qtls.CipherSuiteTLS13ByID(cs.CipherSuite)
		var err error
		c.masterSecret, err = (&cs).ExportKeyingMaterial(xseMasterSecretLabel, nil, c.suite.Hash.Size())
		if err != nil {
			panic(fmt.Errorf("failed to export xse_master_secret: %w", err))
		}
		c.secretReadyCtxCancel()
	}()
	return c
}

func (c *cryptoSetup) NewStream(baseStream Stream) Stream {
	<-c.secretReadyCtx.Done()

	rcvTrafficSecret := qtls.DeriveSecret(c.suite, c.masterSecret, trafficSecretLabel(c.perspective.Opposite(), baseStream.StreamID()), nil)
	sendTrafficSecret := qtls.DeriveSecret(c.suite, c.masterSecret, trafficSecretLabel(c.perspective, baseStream.StreamID()), nil)

	tlsConn := qtls2.FromTrafficSecret(&streamConnAdapter{baseStream}, c.suite.ID, rcvTrafficSecret, sendTrafficSecret, &qtls.Config{}, c.extraConfigForStream(baseStream.StreamID()), c.perspective == protocol.PerspectiveClient)

	receiveStream := &receiveStream{
		ReceiveStream: baseStream.ReceiveStream(),
		conn:          tlsConn,
	}

	sendStream := &sendStream{
		SendStream: baseStream.SendStream(),
		conn:       tlsConn,
	}

	return &stream{
		sendStream:    sendStream,
		receiveStream: receiveStream,
	}
}

func (c *cryptoSetup) NewReceiveStream(baseStream ReceiveStream) ReceiveStream {
	<-c.secretReadyCtx.Done()

	rcvTrafficSecret := qtls.DeriveSecret(c.suite, c.masterSecret, trafficSecretLabel(c.perspective.Opposite(), baseStream.StreamID()), nil)

	tlsConn := qtls2.FromTrafficSecret(&receiveStreamConnAdapter{baseStream}, c.suite.ID, rcvTrafficSecret, nil, &qtls.Config{}, c.extraConfigForStream(baseStream.StreamID()), c.perspective == protocol.PerspectiveClient)

	return &receiveStream{
		ReceiveStream: baseStream,
		conn:          tlsConn,
	}
}

func (c *cryptoSetup) NewSendStream(baseStream SendStream) SendStream {
	<-c.secretReadyCtx.Done()

	sendTrafficSecret := qtls.DeriveSecret(c.suite, c.masterSecret, trafficSecretLabel(c.perspective, baseStream.StreamID()), nil)

	tlsConn := qtls2.FromTrafficSecret(&sendStreamConnAdapter{baseStream}, c.suite.ID, nil, sendTrafficSecret, &qtls.Config{}, c.extraConfigForStream(baseStream.StreamID()), c.perspective == protocol.PerspectiveClient)

	return &sendStream{
		SendStream: baseStream,
		conn:       tlsConn,
	}
}

func (c *cryptoSetup) extraConfigForStream(streamID protocol.StreamID) *qtls2.ExtraConfig {
	if c.tracer != nil {
		return &qtls2.ExtraConfig{
			OnReceiveApplicationDataRecord: func(rawLength int, dataLength int) {
				c.tracer.XseReceiveRecord(streamID, rawLength, dataLength)
			},
		}
	}
	return nil
}
