package xads

import (
	"context"
	"fmt"
	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/qtls"
	"github.com/quic-go/quic-go/logging"
)

const (
	// label to export XADS-QUIC master secret from TLS exporter_master_secret
	// (see RFC8446 Section 7.5)
	//TODO register non experimental exporter label (see RFC5705 Section 4)
	xadsMasterSecretLabel = "EXPERIMENTAL xads master"
)

// label to export XADS-QUIC stream traffic secret from xads_master_secret.
// (see RFC8446 Section 7.5)
// first char is 'c' if the sender is the client.
// first char is 's' if the sender is the server.
// followed by a space and the decimal stream id.
func trafficSecretLabel(perspective protocol.Perspective, streamID protocol.StreamID) string {
	if perspective == logging.PerspectiveClient {
		return fmt.Sprintf("c %d", streamID)
	} else {
		return fmt.Sprintf("s %d", streamID)

	}
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
		// TODO allow 0-RTT XADS-QUIC streams
		// await handshake completion
		err := quicTls.HandshakeContext(context.Background())
		if err != nil {
			panic(fmt.Errorf("failed to export xads_master_secret: %w", err))
		}
		cs := quicTls.ConnectionState()
		c.suite = qtls.CipherSuiteTLS13ByID(cs.CipherSuite)
		c.masterSecret, err = (&cs).ExportKeyingMaterial(xadsMasterSecretLabel, nil, c.suite.Hash.Size())
		if err != nil {
			panic(fmt.Errorf("failed to export xads_master_secret: %w", err))
		}
		c.secretReadyCtxCancel()
	}()
	return c
}

func (c *cryptoSetup) NewStream(baseStream Stream) Stream {
	<-c.secretReadyCtx.Done()

	rcvTrafficSecret := qtls.DeriveSecret(c.suite, c.masterSecret, trafficSecretLabel(c.perspective.Opposite(), baseStream.StreamID()), nil)
	sendTrafficSecret := qtls.DeriveSecret(c.suite, c.masterSecret, trafficSecretLabel(c.perspective, baseStream.StreamID()), nil)

	tlsConn := qtls.FromTrafficSecret(&streamConnAdapter{baseStream}, c.suite.ID, rcvTrafficSecret, sendTrafficSecret, &qtls.Config{}, c.extraConfigForStream(baseStream.StreamID()), c.perspective == protocol.PerspectiveClient)

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

	tlsConn := qtls.FromTrafficSecret(&receiveStreamConnAdapter{baseStream}, c.suite.ID, rcvTrafficSecret, nil, &qtls.Config{}, c.extraConfigForStream(baseStream.StreamID()), c.perspective == protocol.PerspectiveClient)

	return &receiveStream{
		ReceiveStream: baseStream,
		conn:          tlsConn,
	}
}

func (c *cryptoSetup) NewSendStream(baseStream SendStream) SendStream {
	<-c.secretReadyCtx.Done()

	sendTrafficSecret := qtls.DeriveSecret(c.suite, c.masterSecret, trafficSecretLabel(c.perspective, baseStream.StreamID()), nil)

	tlsConn := qtls.FromTrafficSecret(&sendStreamConnAdapter{baseStream}, c.suite.ID, nil, sendTrafficSecret, &qtls.Config{}, c.extraConfigForStream(baseStream.StreamID()), c.perspective == protocol.PerspectiveClient)

	return &sendStream{
		SendStream: baseStream,
		conn:       tlsConn,
	}
}

func (c *cryptoSetup) extraConfigForStream(streamID protocol.StreamID) *qtls.ExtraConfig {
	if c.tracer != nil {
		return &qtls.ExtraConfig{
			OnReceiveApplicationDataRecord: func(rawLength int, dataLength int) {
				c.tracer.XadsReceiveRecord(streamID, rawLength, dataLength)
			},
		}
	}
	return nil
}
