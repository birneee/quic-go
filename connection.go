package quic

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lucas-clemente/quic-go/handover"
	"github.com/lucas-clemente/quic-go/internal/qtls"
	"github.com/lucas-clemente/quic-go/internal/xse"
	"github.com/lucas-clemente/quic-go/path"
	"io"
	"net"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/lucas-clemente/quic-go/internal/ackhandler"
	"github.com/lucas-clemente/quic-go/internal/flowcontrol"
	"github.com/lucas-clemente/quic-go/internal/handshake"
	"github.com/lucas-clemente/quic-go/internal/logutils"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/qerr"
	"github.com/lucas-clemente/quic-go/internal/utils"
	"github.com/lucas-clemente/quic-go/internal/wire"
	"github.com/lucas-clemente/quic-go/logging"
)

type unpacker interface {
	UnpackLongHeader(hdr *wire.Header, rcvTime time.Time, data []byte) (*unpackedPacket, error)
	UnpackShortHeader(rcvTime time.Time, data []byte) (protocol.PacketNumber, protocol.PacketNumberLen, protocol.KeyPhaseBit, []byte, error)
}

type streamGetter interface {
	GetOrOpenReceiveStream(protocol.StreamID) (receiveStreamI, error)
	GetOrOpenSendStream(protocol.StreamID) (sendStreamI, error)
}

type streamManager interface {
	GetOrOpenSendStream(protocol.StreamID) (sendStreamI, error)
	GetOrOpenReceiveStream(protocol.StreamID) (receiveStreamI, error)
	OpenStream() (Stream, error)
	OpenUniStream() (SendStream, error)
	OpenStreamSync(context.Context) (Stream, error)
	OpenUniStreamSync(context.Context) (SendStream, error)
	AcceptStream(context.Context) (Stream, error)
	AcceptUniStream(context.Context) (ReceiveStream, error)
	DeleteStream(protocol.StreamID) error
	UpdateLimits(*wire.TransportParameters)
	HandleMaxStreamsFrame(*wire.MaxStreamsFrame)
	CloseWithError(error)
	ResetFor0RTT()
	UseResetMaps()
	// TODO set this in constructor
	SetXseCryptoSetup(xse.CryptoSetup)
}

type cryptoStreamHandler interface {
	RunHandshake()
	ChangeConnectionID(protocol.ConnectionID)
	SetLargest1RTTAcked(protocol.PacketNumber) error
	SetHandshakeConfirmed()
	GetSessionTicket() ([]byte, error)
	io.Closer
	ConnectionState() handshake.ConnectionState
	StoreHandoverState(t *handover.State, p protocol.Perspective)
	TlsConn() *qtls.Conn
}

type packetInfo struct {
	addr    net.IP
	ifIndex uint32
}

type receivedPacket struct {
	buffer *packetBuffer

	remoteAddr net.Addr
	rcvTime    time.Time
	data       []byte

	ecn protocol.ECN

	info *packetInfo
}

func (p *receivedPacket) Size() protocol.ByteCount { return protocol.ByteCount(len(p.data)) }

func (p *receivedPacket) Clone() *receivedPacket {
	return &receivedPacket{
		remoteAddr: p.remoteAddr,
		rcvTime:    p.rcvTime,
		data:       p.data,
		buffer:     p.buffer,
		ecn:        p.ecn,
		info:       p.info,
	}
}

type connRunner interface {
	Add(protocol.ConnectionID, packetHandler) bool
	GetStatelessResetToken(protocol.ConnectionID) protocol.StatelessResetToken
	Retire(protocol.ConnectionID)
	Remove(protocol.ConnectionID)
	ReplaceWithClosed([]protocol.ConnectionID, protocol.Perspective, []byte)
	AddResetToken(protocol.StatelessResetToken, packetHandler)
	RemoveResetToken(protocol.StatelessResetToken)
}

type handshakeRunner struct {
	onReceivedParams    func(*wire.TransportParameters)
	onError             func(error)
	dropKeys            func(protocol.EncryptionLevel)
	onHandshakeComplete func()
}

func (r *handshakeRunner) OnReceivedParams(tp *wire.TransportParameters) { r.onReceivedParams(tp) }
func (r *handshakeRunner) OnError(e error)                               { r.onError(e) }
func (r *handshakeRunner) DropKeys(el protocol.EncryptionLevel)          { r.dropKeys(el) }
func (r *handshakeRunner) OnHandshakeComplete()                          { r.onHandshakeComplete() }

type closeError struct {
	err       error
	remote    bool
	immediate bool
}

type errCloseForRecreating struct {
	nextPacketNumber protocol.PacketNumber
	nextVersion      protocol.VersionNumber
}

func (e *errCloseForRecreating) Error() string {
	return "closing connection in order to recreate it"
}

var connTracingID uint64        // to be accessed atomically
func nextConnTracingID() uint64 { return atomic.AddUint64(&connTracingID, 1) }

// A Connection is a QUIC connection
type connection struct {
	// Destination connection ID used during the handshake.
	// Used to check source connection ID on incoming packets.
	handshakeDestConnID protocol.ConnectionID
	// Set for the client. Destination connection ID used on the first Initial sent.
	origDestConnID protocol.ConnectionID
	retrySrcConnID *protocol.ConnectionID // only set for the client (and if a Retry was performed)

	srcConnIDLen int

	perspective protocol.Perspective
	version     protocol.VersionNumber
	config      *Config

	conn      sendConn
	sendQueue sender

	streamsMap      streamManager
	connIDManager   *connIDManager
	connIDGenerator *connIDGenerator

	rttStats *utils.RTTStats

	cryptoStreamManager   *cryptoStreamManager
	sentPacketHandler     ackhandler.SentPacketHandler
	receivedPacketHandler ackhandler.ReceivedPacketHandler
	retransmissionQueue   *retransmissionQueue
	framer                framer
	windowUpdateQueue     *windowUpdateQueue
	connFlowController    flowcontrol.ConnectionFlowController
	tokenStoreKey         string                    // only set for the client
	tokenGenerator        *handshake.TokenGenerator // only set for the server

	unpacker      unpacker
	frameParser   wire.FrameParser
	packer        packer
	mtuDiscoverer mtuDiscoverer // initialized when the handshake completes

	oneRTTStream        cryptoStream // only set for the server
	cryptoStreamHandler cryptoStreamHandler

	receivedPackets  chan *receivedPacket
	sendingScheduled chan struct{}

	closeOnce sync.Once
	// closeChan is used to notify the run loop that it should terminate
	closeChan chan closeError

	ctx                context.Context
	ctxCancel          context.CancelFunc
	handshakeCtx       context.Context
	handshakeCtxCancel context.CancelFunc

	undecryptablePackets          []*receivedPacket // undecryptable packets, waiting for a change in encryption level
	undecryptablePacketsToProcess []*receivedPacket

	clientHelloWritten     <-chan *wire.TransportParameters
	earlyConnReadyChan     chan struct{}
	handshakeCompleteChan  chan struct{} // is closed when the handshake completes
	handshakeConfirmedChan chan struct{} // is closed when the handshake is confirmed
	sentFirstPacket        bool
	handshakeComplete      bool
	handshakeConfirmed     bool

	receivedRetry       bool
	versionNegotiated   bool
	receivedFirstPacket bool

	idleTimeout  time.Duration
	creationTime time.Time
	// The idle timeout is set based on the max of the time we received the last packet...
	lastPacketReceivedTime time.Time
	// ... and the time we sent a new ack-eliciting packet after receiving a packet.
	firstAckElicitingPacketAfterIdleSentTime time.Time
	// pacingDeadline is the time when the next packet should be sent
	pacingDeadline time.Time

	peerParams *wire.TransportParameters

	timer connectionTimer
	// keepAlivePingSent stores whether a keep alive PING is in flight.
	// It is reset as soon as we receive a packet from the peer.
	keepAlivePingSent bool
	keepAliveInterval time.Duration

	pathManager path.PathManager

	datagramQueue *datagramQueue

	logID  string
	tracer logging.ConnectionTracer
	logger utils.Logger

	ignoredRemoteAddresses []net.UDPAddr
	runner                 connRunner
	//TODO this should not be necessary
	cloned bool
	// closed when H-QUIC handover migration is done.
	// ignore received 1RTT packets and opening streams.
	handoverMigrationCtx       context.Context
	handoverMigrationCtxCancel context.CancelFunc
}

var (
	_ Connection      = &connection{}
	_ EarlyConnection = &connection{}
	_ streamSender    = &connection{}
)

var newConnection = func(
	conn sendConn,
	runner connRunner,
	origDestConnID protocol.ConnectionID,
	retrySrcConnID *protocol.ConnectionID,
	clientDestConnID protocol.ConnectionID,
	destConnID protocol.ConnectionID,
	srcConnID protocol.ConnectionID,
	statelessResetToken protocol.StatelessResetToken,
	conf *Config,
	tlsConf *tls.Config,
	tokenGenerator *handshake.TokenGenerator,
	enable0RTT bool,
	clientAddressValidated bool,
	tracer logging.ConnectionTracer,
	tracingID uint64,
	logger utils.Logger,
	v protocol.VersionNumber,
) quicConn {
	s := &connection{
		conn:                   conn,
		config:                 conf,
		handshakeDestConnID:    destConnID,
		srcConnIDLen:           srcConnID.Len(),
		tokenGenerator:         tokenGenerator,
		oneRTTStream:           newCryptoStream(),
		perspective:            protocol.PerspectiveServer,
		handshakeCompleteChan:  make(chan struct{}),
		handshakeConfirmedChan: make(chan struct{}),
		tracer:                 tracer,
		logger:                 logger,
		version:                v,
		runner:                 runner,
	}
	if origDestConnID.Len() > 0 {
		s.logID = origDestConnID.String()
	} else {
		s.logID = destConnID.String()
	}
	s.connIDManager = newConnIDManager(
		destConnID,
		func(token protocol.StatelessResetToken) { runner.AddResetToken(token, s) },
		runner.RemoveResetToken,
		s.queueControlFrame,
	)
	s.connIDGenerator = newConnIDGenerator(
		srcConnID,
		&clientDestConnID,
		func(connID protocol.ConnectionID) { runner.Add(connID, s) },
		runner.GetStatelessResetToken,
		runner.Remove,
		runner.Retire,
		runner.ReplaceWithClosed,
		s.queueControlFrame,
		s.config.ConnectionIDGenerator,
		s.version,
	)
	s.ctx, s.ctxCancel = context.WithCancel(context.WithValue(context.Background(), ConnectionTracingKey, tracingID))
	s.preSetup()
	s.sentPacketHandler, s.receivedPacketHandler = ackhandler.NewAckHandler(
		s.pathManager.CurrentSendPath(),
		0,
		getMaxPacketSize(s.conn.RemoteAddr()),
		s.config.InitialCongestionWindow,
		s.config.MinCongestionWindow,
		s.config.MaxCongestionWindow,
		s.config.InitialSlowStartThreshold,
		s.config.MinSlowStartThreshold,
		s.config.MaxSlowStartThreshold,
		s.rttStats,
		clientAddressValidated,
		s.perspective,
		s.config.HyblaWestwoodCongestionControl,
		s.config.FixedPTO,
		s.tracer,
		s.logger,
		s.version,
	)
	initialStream := newCryptoStream()
	handshakeStream := newCryptoStream()
	params := &wire.TransportParameters{
		InitialMaxStreamDataBidiLocal:   protocol.ByteCount(s.config.InitialStreamReceiveWindow),
		InitialMaxStreamDataBidiRemote:  protocol.ByteCount(s.config.InitialStreamReceiveWindow),
		InitialMaxStreamDataUni:         protocol.ByteCount(s.config.InitialStreamReceiveWindow),
		InitialMaxData:                  protocol.ByteCount(s.config.InitialConnectionReceiveWindow),
		MaxIdleTimeout:                  s.config.MaxIdleTimeout,
		MaxBidiStreamNum:                protocol.StreamNum(s.config.MaxIncomingStreams),
		MaxUniStreamNum:                 protocol.StreamNum(s.config.MaxIncomingUniStreams),
		MaxAckDelay:                     protocol.MaxAckDelayInclGranularity,
		AckDelayExponent:                protocol.AckDelayExponent,
		DisableActiveMigration:          !s.config.EnableActiveMigration,
		StatelessResetToken:             &statelessResetToken,
		OriginalDestinationConnectionID: origDestConnID,
		ActiveConnectionIDLimit:         protocol.MaxActiveConnectionIDs,
		InitialSourceConnectionID:       srcConnID,
		RetrySourceConnectionID:         retrySrcConnID,
		ExtraStreamEncryption:           s.config.ExtraStreamEncryption.enabled(),
	}
	if s.config.EnableDatagrams {
		params.MaxDatagramFrameSize = protocol.MaxDatagramFrameSize
	} else {
		params.MaxDatagramFrameSize = protocol.InvalidByteCount
	}
	if s.tracer != nil {
		s.tracer.SentTransportParameters(params)
	}
	cs := handshake.NewCryptoSetupServer(
		initialStream,
		handshakeStream,
		clientDestConnID,
		conn.LocalAddr(),
		conn.RemoteAddr(),
		params,
		&handshakeRunner{
			onReceivedParams: s.handleTransportParameters,
			onError:          s.closeLocal,
			dropKeys:         s.dropEncryptionLevel,
			onHandshakeComplete: func() {
				runner.Retire(clientDestConnID)
				close(s.handshakeCompleteChan)
			},
		},
		tlsConf,
		enable0RTT,
		s.rttStats,
		tracer,
		logger,
		s.version,
	)
	s.cryptoStreamHandler = cs
	s.packer = newPacketPacker(
		srcConnID,
		s.connIDManager.Get,
		initialStream,
		handshakeStream,
		s.sentPacketHandler,
		s.retransmissionQueue,
		s.RemoteAddr(),
		cs,
		s.framer,
		s.receivedPacketHandler,
		s.datagramQueue,
		s.perspective,
		s.version,
	)
	s.unpacker = newPacketUnpacker(cs, s.srcConnIDLen, s.version)
	s.cryptoStreamManager = newCryptoStreamManager(cs, initialStream, handshakeStream, s.oneRTTStream)
	return s
}

// declare this as a variable, such that we can it mock it in the tests
var newClientConnection = func(
	conn sendConn,
	runner connRunner,
	destConnID protocol.ConnectionID,
	srcConnID protocol.ConnectionID,
	conf *Config,
	tlsConf *tls.Config,
	initialPacketNumber protocol.PacketNumber,
	enable0RTT bool,
	hasNegotiatedVersion bool,
	tracer logging.ConnectionTracer,
	tracingID uint64,
	logger utils.Logger,
	v protocol.VersionNumber,
) quicConn {
	s := &connection{
		conn:                   conn,
		config:                 conf,
		origDestConnID:         destConnID,
		handshakeDestConnID:    destConnID,
		srcConnIDLen:           srcConnID.Len(),
		perspective:            protocol.PerspectiveClient,
		handshakeCompleteChan:  make(chan struct{}),
		handshakeConfirmedChan: make(chan struct{}),
		logID:                  destConnID.String(),
		logger:                 logger,
		tracer:                 tracer,
		versionNegotiated:      hasNegotiatedVersion,
		version:                v,
		runner:                 runner,
	}
	s.connIDManager = newConnIDManager(
		destConnID,
		func(token protocol.StatelessResetToken) { runner.AddResetToken(token, s) },
		runner.RemoveResetToken,
		s.queueControlFrame,
	)
	s.connIDGenerator = newConnIDGenerator(
		srcConnID,
		nil,
		func(connID protocol.ConnectionID) { runner.Add(connID, s) },
		runner.GetStatelessResetToken,
		runner.Remove,
		runner.Retire,
		runner.ReplaceWithClosed,
		s.queueControlFrame,
		s.config.ConnectionIDGenerator,
		s.version,
	)
	s.ctx, s.ctxCancel = context.WithCancel(context.WithValue(context.Background(), ConnectionTracingKey, tracingID))
	s.preSetup()
	s.sentPacketHandler, s.receivedPacketHandler = ackhandler.NewAckHandler(
		s.pathManager.CurrentSendPath(),
		initialPacketNumber,
		getMaxPacketSize(s.conn.RemoteAddr()),
		s.config.InitialCongestionWindow,
		s.config.MinCongestionWindow,
		s.config.MaxCongestionWindow,
		s.config.InitialSlowStartThreshold,
		s.config.MinSlowStartThreshold,
		s.config.MaxSlowStartThreshold,
		s.rttStats,
		false, /* has no effect */
		s.perspective,
		s.config.HyblaWestwoodCongestionControl,
		s.config.FixedPTO,
		s.tracer,
		s.logger,
		s.version,
	)
	initialStream := newCryptoStream()
	handshakeStream := newCryptoStream()
	params := &wire.TransportParameters{
		InitialMaxStreamDataBidiRemote: protocol.ByteCount(s.config.InitialStreamReceiveWindow),
		InitialMaxStreamDataBidiLocal:  protocol.ByteCount(s.config.InitialStreamReceiveWindow),
		InitialMaxStreamDataUni:        protocol.ByteCount(s.config.InitialStreamReceiveWindow),
		InitialMaxData:                 protocol.ByteCount(s.config.InitialConnectionReceiveWindow),
		MaxIdleTimeout:                 s.config.MaxIdleTimeout,
		MaxBidiStreamNum:               protocol.StreamNum(s.config.MaxIncomingStreams),
		MaxUniStreamNum:                protocol.StreamNum(s.config.MaxIncomingUniStreams),
		MaxAckDelay:                    protocol.MaxAckDelayInclGranularity,
		AckDelayExponent:               protocol.AckDelayExponent,
		DisableActiveMigration:         !s.config.EnableActiveMigration,
		ActiveConnectionIDLimit:        protocol.MaxActiveConnectionIDs,
		InitialSourceConnectionID:      srcConnID,
		ExtraStreamEncryption:          s.config.ExtraStreamEncryption.enabled(),
	}
	if s.config.EnableDatagrams {
		params.MaxDatagramFrameSize = protocol.MaxDatagramFrameSize
	} else {
		params.MaxDatagramFrameSize = protocol.InvalidByteCount
	}
	if s.tracer != nil {
		s.tracer.SentTransportParameters(params)
	}
	cs, clientHelloWritten := handshake.NewCryptoSetupClient(
		initialStream,
		handshakeStream,
		destConnID,
		conn.LocalAddr(),
		conn.RemoteAddr(),
		params,
		&handshakeRunner{
			onReceivedParams:    s.handleTransportParameters,
			onError:             s.closeLocal,
			dropKeys:            s.dropEncryptionLevel,
			onHandshakeComplete: func() { close(s.handshakeCompleteChan) },
		},
		tlsConf,
		enable0RTT,
		s.rttStats,
		tracer,
		logger,
		s.version,
	)
	s.clientHelloWritten = clientHelloWritten
	s.cryptoStreamHandler = cs
	s.cryptoStreamManager = newCryptoStreamManager(cs, initialStream, handshakeStream, newCryptoStream())
	s.unpacker = newPacketUnpacker(cs, s.srcConnIDLen, s.version)
	s.packer = newPacketPacker(
		srcConnID,
		s.connIDManager.Get,
		initialStream,
		handshakeStream,
		s.sentPacketHandler,
		s.retransmissionQueue,
		s.RemoteAddr(),
		cs,
		s.framer,
		s.receivedPacketHandler,
		s.datagramQueue,
		s.perspective,
		s.version,
	)
	if len(tlsConf.ServerName) > 0 {
		s.tokenStoreKey = tlsConf.ServerName
	} else {
		s.tokenStoreKey = conn.RemoteAddr().String()
	}
	if s.config.TokenStore != nil {
		if token := s.config.TokenStore.Pop(s.tokenStoreKey); token != nil {
			s.packer.SetToken(token.data)
		}
	}
	return s
}

func (s *connection) preSetup() {
	s.sendQueue = newSendQueue(s.conn)
	s.retransmissionQueue = newRetransmissionQueue(s.version)
	s.frameParser = wire.NewFrameParser(s.config.EnableDatagrams, s.version)
	s.rttStats = &utils.RTTStats{}
	s.connFlowController = flowcontrol.NewConnectionFlowController(
		protocol.ByteCount(s.config.InitialConnectionReceiveWindow),
		protocol.ByteCount(s.config.MaxConnectionReceiveWindow),
		s.onHasConnectionWindowUpdate,
		func(size protocol.ByteCount) bool {
			if s.config.AllowConnectionWindowIncrease == nil {
				return true
			}
			return s.config.AllowConnectionWindowIncrease(s, uint64(size))
		},
		s.rttStats,
		s.logger,
	)
	s.earlyConnReadyChan = make(chan struct{})
	s.streamsMap = newStreamsMap(
		s,
		s.newFlowController,
		uint64(s.config.MaxIncomingStreams),
		uint64(s.config.MaxIncomingUniStreams),
		s.perspective,
		s.version,
	)
	s.framer = newFramer(s.streamsMap, s.version)
	s.receivedPackets = make(chan *receivedPacket, protocol.MaxConnUnprocessedPackets)
	s.closeChan = make(chan closeError, 1)
	s.sendingScheduled = make(chan struct{}, 1)
	s.handshakeCtx, s.handshakeCtxCancel = context.WithCancel(context.Background())
	s.handoverMigrationCtx, s.handoverMigrationCtxCancel = context.WithCancel(context.Background())
	if s.config.ProxyConf == nil && !s.config.IgnoreReceived1RTTPacketsUntilFirstPathMigration {
		s.handoverMigrationCtxCancel()
	}

	now := time.Now()
	s.lastPacketReceivedTime = now
	s.creationTime = now

	s.windowUpdateQueue = newWindowUpdateQueue(s.streamsMap, s.connFlowController, s.framer.QueueControlFrame)
	s.datagramQueue = newDatagramQueue(s.scheduleSending, s.logger)

	s.pathManager = path.NewPathManager(s.conn.RemoteAddr().(*net.UDPAddr), false, s.updateSendPath, s.queuePathChallengeFrame, s.logger, s.tracer)
}

// run the connection main loop
func (s *connection) run() error {
	defer s.ctxCancel()

	s.timer = *newTimer()

	handshaking := make(chan struct{})
	go func() {
		defer close(handshaking)
		if !s.handshakeComplete || !s.cloned {
			go s.cryptoStreamHandler.RunHandshake()
		}
	}()
	go func() {
		if err := s.sendQueue.Run(); err != nil {
			s.destroyImpl(err)
		}
	}()

	if (!s.handshakeComplete || !s.cloned) && s.perspective == protocol.PerspectiveClient {
		select {
		case zeroRTTParams := <-s.clientHelloWritten:
			s.scheduleSending()
			if zeroRTTParams != nil {
				s.restoreTransportParameters(zeroRTTParams)
				close(s.earlyConnReadyChan)
			}
		case closeErr := <-s.closeChan:
			// put the close error back into the channel, so that the run loop can receive it
			s.closeChan <- closeErr
		}
	}

	var (
		closeErr           closeError
		sendQueueAvailable <-chan struct{}
	)

runLoop:
	for {
		// Close immediately if requested
		select {
		case closeErr = <-s.closeChan:
			break runLoop
		case <-s.handshakeCompleteChan:
			s.handleHandshakeComplete()
		default:
		}

		s.maybeResetTimer()

		var processedUndecryptablePacket bool
		if len(s.undecryptablePacketsToProcess) > 0 {
			queue := s.undecryptablePacketsToProcess
			s.undecryptablePacketsToProcess = nil
			for _, p := range queue {
				if processed := s.handlePacketImpl(p); processed {
					processedUndecryptablePacket = true
				}
				// Don't set timers and send packets if the packet made us close the connection.
				select {
				case closeErr = <-s.closeChan:
					break runLoop
				default:
				}
			}
		}
		// If we processed any undecryptable packets, jump to the resetting of the timers directly.
		if !processedUndecryptablePacket {
			select {
			case closeErr = <-s.closeChan:
				break runLoop
			case <-s.timer.Chan():
				s.timer.SetRead()
				// We do all the interesting stuff after the switch statement, so
				// nothing to see here.
			case <-s.sendingScheduled:
				// We do all the interesting stuff after the switch statement, so
				// nothing to see here.
			case <-sendQueueAvailable:
			case firstPacket := <-s.receivedPackets:
				wasProcessed := s.handlePacketImpl(firstPacket)
				// Don't set timers and send packets if the packet made us close the connection.
				select {
				case closeErr = <-s.closeChan:
					break runLoop
				default:
				}
				if s.handshakeComplete {
					// Now process all packets in the receivedPackets channel.
					// Limit the number of packets to the length of the receivedPackets channel,
					// so we eventually get a chance to send out an ACK when receiving a lot of packets.
					numPackets := len(s.receivedPackets)
				receiveLoop:
					for i := 0; i < numPackets; i++ {
						select {
						case p := <-s.receivedPackets:
							if processed := s.handlePacketImpl(p); processed {
								wasProcessed = true
							}
							select {
							case closeErr = <-s.closeChan:
								break runLoop
							default:
							}
						default:
							break receiveLoop
						}
					}
				}
				// Only reset the timers if this packet was actually processed.
				// This avoids modifying any state when handling undecryptable packets,
				// which could be injected by an attacker.
				if !wasProcessed {
					continue
				}
			case <-s.handshakeCompleteChan:
				s.handleHandshakeComplete()
			}
		}

		now := time.Now()
		if timeout := s.sentPacketHandler.GetLossDetectionTimeout(); !timeout.IsZero() && timeout.Before(now) {
			// This could cause packets to be retransmitted.
			// Check it before trying to send packets.
			if err := s.sentPacketHandler.OnLossDetectionTimeout(); err != nil {
				s.closeLocal(err)
			}
		}

		if keepAliveTime := s.nextKeepAliveTime(); !keepAliveTime.IsZero() && !now.Before(keepAliveTime) {
			// send a PING frame since there is no activity in the connection
			s.logger.Debugf("Sending a keep-alive PING to keep the connection alive.")
			s.framer.QueueControlFrame(&wire.PingFrame{})
			s.keepAlivePingSent = true
		} else if !s.handshakeComplete && now.Sub(s.creationTime) >= s.config.handshakeTimeout() {
			s.destroyImpl(qerr.ErrHandshakeTimeout)
			continue
		} else {
			idleTimeoutStartTime := s.idleTimeoutStartTime()
			if (!s.handshakeComplete && now.Sub(idleTimeoutStartTime) >= s.config.HandshakeIdleTimeout) ||
				(s.handshakeComplete && now.Sub(idleTimeoutStartTime) >= s.idleTimeout) {
				s.destroyImpl(qerr.ErrIdleTimeout)
				continue
			}
		}

		if s.sendQueue.WouldBlock() {
			// The send queue is still busy sending out packets.
			// Wait until there's space to enqueue new packets.
			sendQueueAvailable = s.sendQueue.Available()
			continue
		}
		if err := s.sendPackets(); err != nil {
			s.closeLocal(err)
		}
		if s.sendQueue.WouldBlock() {
			sendQueueAvailable = s.sendQueue.Available()
		} else {
			sendQueueAvailable = nil
		}
	}

	s.cryptoStreamHandler.Close()
	<-handshaking
	s.handleCloseError(&closeErr)
	if e := (&errCloseForRecreating{}); !errors.As(closeErr.err, &e) && s.tracer != nil {
		s.tracer.Close()
	}
	s.logger.Infof("Connection %s closed.", s.logID)
	s.sendQueue.Close()
	s.timer.Stop()
	return closeErr.err
}

// blocks until the early connection can be used
func (s *connection) earlyConnReady() <-chan struct{} {
	return s.earlyConnReadyChan
}

func (s *connection) HandshakeComplete() context.Context {
	return s.handshakeCtx
}

func (s *connection) Context() context.Context {
	return s.ctx
}

func (s *connection) supportsDatagrams() bool {
	return s.peerParams.MaxDatagramFrameSize > 0
}

func (s *connection) ConnectionState() ConnectionState {
	return ConnectionState{
		TLS:               s.cryptoStreamHandler.ConnectionState(),
		SupportsDatagrams: s.supportsDatagrams(),
		Version:           s.version,
	}
}

// Time when the next keep-alive packet should be sent.
// It returns a zero time if no keep-alive should be sent.
func (s *connection) nextKeepAliveTime() time.Time {
	if s.config.KeepAlivePeriod == 0 || s.keepAlivePingSent || !s.firstAckElicitingPacketAfterIdleSentTime.IsZero() {
		return time.Time{}
	}
	return s.lastPacketReceivedTime.Add(s.keepAliveInterval)
}

func (s *connection) maybeResetTimer() {
	var deadline time.Time
	if !s.handshakeComplete {
		deadline = utils.MinTime(
			s.creationTime.Add(s.config.handshakeTimeout()),
			s.idleTimeoutStartTime().Add(s.config.HandshakeIdleTimeout),
		)
	} else {
		if keepAliveTime := s.nextKeepAliveTime(); !keepAliveTime.IsZero() {
			deadline = keepAliveTime
		} else {
			deadline = s.idleTimeoutStartTime().Add(s.idleTimeout)
		}
	}

	s.timer.SetTimer(
		deadline,
		s.receivedPacketHandler.GetAlarmTimeout(),
		s.sentPacketHandler.GetLossDetectionTimeout(),
		s.pacingDeadline,
	)
}

func (s *connection) idleTimeoutStartTime() time.Time {
	return utils.MaxTime(s.lastPacketReceivedTime, s.firstAckElicitingPacketAfterIdleSentTime)
}

func (s *connection) handleHandshakeComplete() {
	s.handshakeComplete = true
	s.handshakeCompleteChan = nil // prevent this case from ever being selected again
	defer s.handshakeCtxCancel()
	// Once the handshake completes, we have derived 1-RTT keys.
	// There's no point in queueing undecryptable packets for later decryption any more.
	s.undecryptablePackets = nil

	s.connIDManager.SetHandshakeComplete()
	s.connIDGenerator.SetHandshakeComplete()

	if s.config.ProxyConf != nil && s.config.AllowEarlyHandover {
		err := s.useProxy()
		if err != nil {
			s.closeLocal(err)
		}
	}

	if s.perspective == protocol.PerspectiveClient {
		s.applyTransportParameters()
		return
	}

	s.handleHandshakeConfirmed()

	ticket, err := s.cryptoStreamHandler.GetSessionTicket()
	if err != nil {
		s.closeLocal(err)
	}
	if ticket != nil {
		s.oneRTTStream.Write(ticket)
		for s.oneRTTStream.HasData() {
			s.queueControlFrame(s.oneRTTStream.PopCryptoFrame(protocol.MaxPostHandshakeCryptoFrameSize))
		}
	}
	token, err := s.tokenGenerator.NewToken(s.conn.RemoteAddr())
	if err != nil {
		s.closeLocal(err)
	}
	s.queueControlFrame(&wire.NewTokenFrame{Token: token})
	s.queueControlFrame(&wire.HandshakeDoneFrame{})
}

func (s *connection) handleHandshakeConfirmed() {
	s.handshakeConfirmed = true
	close(s.handshakeConfirmedChan)
	s.sentPacketHandler.SetHandshakeConfirmed()
	s.cryptoStreamHandler.SetHandshakeConfirmed()

	if !s.config.DisablePathMTUDiscovery {
		maxPacketSize := s.peerParams.MaxUDPPayloadSize
		if maxPacketSize == 0 {
			maxPacketSize = protocol.MaxByteCount
		}
		maxPacketSize = utils.Min(maxPacketSize, protocol.MaxPacketBufferSize)
		s.mtuDiscoverer = newMTUDiscoverer(
			s.rttStats,
			getMaxPacketSize(s.conn.RemoteAddr()),
			maxPacketSize,
			func(size protocol.ByteCount) {
				s.sentPacketHandler.SetMaxDatagramSize(size)
				s.packer.SetMaxPacketSize(size)
			},
		)
	}

	if s.config.ProxyConf != nil && !s.config.AllowEarlyHandover {
		err := s.useProxy()
		if err != nil {
			s.closeLocal(err)
		}
	}
}

func (s *connection) handlePacketImpl(rp *receivedPacket) bool {

	// ignore packet if from a ignored remote
	if s.isRemoteAddressIgnored(rp.remoteAddr) {
		s.handleIgnoredAddressPacket(rp)
		return false
	}

	s.sentPacketHandler.ReceivedBytes(rp.Size(), s.pathManager.GetOrCreatePath(rp.remoteAddr))

	if wire.IsVersionNegotiationPacket(rp.data) {
		s.handleVersionNegotiationPacket(rp)
		return false
	}

	var counter uint8
	var lastConnID protocol.ConnectionID
	var processed bool
	data := rp.data
	p := rp
	for len(data) > 0 {
		var destConnID protocol.ConnectionID
		if counter > 0 {
			p = p.Clone()
			p.data = data

			var err error
			destConnID, err = wire.ParseConnectionID(p.data, s.srcConnIDLen)
			if err != nil {
				if s.tracer != nil {
					s.tracer.DroppedPacket(logging.PacketTypeNotDetermined, protocol.ByteCount(len(data)), logging.PacketDropHeaderParseError)
				}
				s.logger.Debugf("error parsing packet, couldn't parse connection ID: %s", err)
				break
			}
			if destConnID != lastConnID {
				if s.tracer != nil {
					s.tracer.DroppedPacket(logging.PacketTypeNotDetermined, protocol.ByteCount(len(data)), logging.PacketDropUnknownConnectionID)
				}
				s.logger.Debugf("coalesced packet has different destination connection ID: %s, expected %s", destConnID, lastConnID)
				break
			}
		}

		if wire.IsLongHeaderPacket(p.data[0]) {
			hdr, packetData, rest, err := wire.ParsePacket(p.data, s.srcConnIDLen)
			if err != nil {
				if s.tracer != nil {
					dropReason := logging.PacketDropHeaderParseError
					if err == wire.ErrUnsupportedVersion {
						dropReason = logging.PacketDropUnsupportedVersion
					}
					s.tracer.DroppedPacket(logging.PacketTypeNotDetermined, protocol.ByteCount(len(data)), dropReason)
				}
				s.logger.Debugf("error parsing packet: %s", err)
				break
			}
			lastConnID = hdr.DestConnectionID

			if hdr.Version != s.version {
				if s.tracer != nil {
					s.tracer.DroppedPacket(logging.PacketTypeFromHeader(hdr), protocol.ByteCount(len(data)), logging.PacketDropUnexpectedVersion)
				}
				s.logger.Debugf("Dropping packet with version %x. Expected %x.", hdr.Version, s.version)
				break
			}

			if counter > 0 {
				p.buffer.Split()
			}
			counter++

			// only log if this actually a coalesced packet
			if s.logger.Debug() && (counter > 1 || len(rest) > 0) {
				s.logger.Debugf("Parsed a coalesced packet. Part %d: %d bytes. Remaining: %d bytes.", counter, len(packetData), len(rest))
			}

			p.data = packetData

			if wasProcessed := s.handleLongHeaderPacket(p, hdr); wasProcessed {
				processed = true
			}
			data = rest
		} else {
			if counter > 0 {
				p.buffer.Split()
			}
			processed = s.handleShortHeaderPacket(p, destConnID)
			break
		}
	}

	p.buffer.MaybeRelease()
	return processed
}

func (s *connection) handleIgnoredAddressPacket(rp *receivedPacket) {
	if s.tracer != nil {
		s.tracer.DroppedPacket(logging.PacketTypeNotDetermined, rp.Size(), logging.PacketDropWaitForMigration)
	}
	if s.perspective == PerspectiveServer {
		return
	}
	if !s.config.AllowEarlyHandover {
		return
	}
	if s.handshakeConfirmed {
		return
	}
	// on early handover, the HANDSHAKE_DONE frame must be received from the original server
	if wire.IsLongHeaderPacket(rp.data[0]) {
		return
	}
	_, _, _, data, err := s.unpacker.UnpackShortHeader(rp.rcvTime, rp.data)
	if err != nil {
		panic("failed to unpack")
	}
	for len(data) > 0 {
		l, frame, err := s.frameParser.ParseNext(data, protocol.Encryption1RTT)
		if err != nil {
			panic("failed to parse frame")
		}
		data = data[l:]
		if frame == nil {
			break
		}
		switch frame.(type) {
		case *wire.HandshakeDoneFrame:
			err = s.handleHandshakeDoneFrame()
			if err != nil {
				panic("failed to handle handshake done frame")
			}
			if s.tracer != nil {
				s.tracer.Debug("accept_handshake_done_from_ignored_packet", "")
			}
			s.logger.Debugf("received HANDSHAKE_DONE frame from original server")
		}
	}
}

func (s *connection) handleShortHeaderPacket(p *receivedPacket, destConnID protocol.ConnectionID) bool {
	var wasQueued bool

	defer func() {
		// Put back the packet buffer if the packet wasn't queued for later decryption.
		if !wasQueued {
			p.buffer.Decrement()
		}
	}()

	pn, pnLen, keyPhase, data, err := s.unpacker.UnpackShortHeader(p.rcvTime, p.data)
	if err != nil {
		wasQueued = s.handleUnpackError(err, p, logging.PacketType1RTT)
		return false
	}

	if s.logger.Debug() {
		s.logger.Debugf("<- Reading packet %d (%d bytes) for connection %s, 1-RTT", pn, p.Size(), destConnID)
		wire.LogShortHeader(s.logger, destConnID, pn, pnLen, keyPhase)
	}

	// handle changed remote address (migration)
	// TODO only react to non-probing packets
	// TODO improve migration: make secure, send path challenge
	// TODO validate path
	// TODO use new connection ID from the peer
	// TODO reset congestion controller and RTT estimate
	// TODO re-validate ECN capability
	if err != handshake.ErrDecryptionFailed {
		s.pathManager.OnReceiveNonProbingPacket(p.remoteAddr.(*net.UDPAddr))
	}

	if s.receivedPacketHandler.IsPotentiallyDuplicate(pn, protocol.Encryption1RTT) {
		s.logger.Debugf("Dropping (potentially) duplicate packet.")
		if s.tracer != nil {
			s.tracer.DroppedPacket(logging.PacketType1RTT, p.Size(), logging.PacketDropDuplicate)
		}
		return false
	}

	var log func([]logging.Frame)
	if s.tracer != nil {
		log = func(frames []logging.Frame) {
			s.tracer.ReceivedShortHeaderPacket(
				&logging.ShortHeader{
					DestConnectionID: destConnID,
					PacketNumber:     pn,
					PacketNumberLen:  pnLen,
					KeyPhase:         keyPhase,
				},
				p.Size(),
				frames,
			)
		}
	}
	if err := s.handleUnpackedShortHeaderPacket(p.remoteAddr, destConnID, pn, data, p.ecn, p.rcvTime, log); err != nil {
		s.closeLocal(err)
		return false
	}
	return true
}

func (s *connection) handleLongHeaderPacket(p *receivedPacket, hdr *wire.Header) bool /* was the packet successfully processed */ {
	var wasQueued bool

	defer func() {
		// Put back the packet buffer if the packet wasn't queued for later decryption.
		if !wasQueued {
			p.buffer.Decrement()
		}
	}()

	if hdr.Type == protocol.PacketTypeRetry {
		return s.handleRetryPacket(hdr, p.data)
	}

	// The server can change the source connection ID with the first Handshake packet.
	// After this, all packets with a different source connection have to be ignored.
	if s.receivedFirstPacket && hdr.Type == protocol.PacketTypeInitial && hdr.SrcConnectionID != s.handshakeDestConnID {
		if s.tracer != nil {
			s.tracer.DroppedPacket(logging.PacketTypeInitial, p.Size(), logging.PacketDropUnknownConnectionID)
		}
		s.logger.Debugf("Dropping Initial packet (%d bytes) with unexpected source connection ID: %s (expected %s)", p.Size(), hdr.SrcConnectionID, s.handshakeDestConnID)
		return false
	}
	// drop 0-RTT packets, if we are a client
	if s.perspective == protocol.PerspectiveClient && hdr.Type == protocol.PacketType0RTT {
		if s.tracer != nil {
			s.tracer.DroppedPacket(logging.PacketType0RTT, p.Size(), logging.PacketDropKeyUnavailable)
		}
		return false
	}

	packet, err := s.unpacker.UnpackLongHeader(hdr, p.rcvTime, p.data)

	if err != nil {
		wasQueued = s.handleUnpackError(err, p, logging.PacketTypeFromHeader(hdr))
		return false
	}

	if s.logger.Debug() {
		s.logger.Debugf("<- Reading packet %d (%d bytes) for connection %s, %s", packet.hdr.PacketNumber, p.Size(), hdr.DestConnectionID, packet.encryptionLevel)
		packet.hdr.Log(s.logger)
	}

	if s.receivedPacketHandler.IsPotentiallyDuplicate(packet.hdr.PacketNumber, packet.encryptionLevel) {
		s.logger.Debugf("Dropping (potentially) duplicate packet.")
		if s.tracer != nil {
			s.tracer.DroppedPacket(logging.PacketTypeFromHeader(hdr), p.Size(), logging.PacketDropDuplicate)
		}
		return false
	}

	if err := s.handleUnpackedPacket(p.remoteAddr, packet, p.ecn, p.rcvTime, p.Size()); err != nil {
		s.closeLocal(err)
		return false
	}
	return true
}

func (s *connection) handleUnpackError(err error, p *receivedPacket, pt logging.PacketType) (wasQueued bool) {
	switch err {
	case handshake.ErrKeysDropped:
		if s.tracer != nil {
			s.tracer.DroppedPacket(pt, p.Size(), logging.PacketDropKeyUnavailable)
		}
		s.logger.Debugf("Dropping %s packet (%d bytes) because we already dropped the keys.", pt, p.Size())
	case handshake.ErrKeysNotYetAvailable:
		// Sealer for this encryption level not yet available.
		// Try again later.
		s.tryQueueingUndecryptablePacket(p, pt)
		return true
	case wire.ErrInvalidReservedBits:
		s.closeLocal(&qerr.TransportError{
			ErrorCode:    qerr.ProtocolViolation,
			ErrorMessage: err.Error(),
		})
	case handshake.ErrDecryptionFailed:
		// This might be a packet injected by an attacker. Drop it.
		if s.tracer != nil {
			s.tracer.DroppedPacket(pt, p.Size(), logging.PacketDropPayloadDecryptError)
		}
		s.logger.Debugf("Dropping %s packet (%d bytes) that could not be unpacked. Error: %s", pt, p.Size(), err)
	default:
		var headerErr *headerParseError
		if errors.As(err, &headerErr) {
			// This might be a packet injected by an attacker. Drop it.
			if s.tracer != nil {
				s.tracer.DroppedPacket(pt, p.Size(), logging.PacketDropHeaderParseError)
			}
			s.logger.Debugf("Dropping %s packet (%d bytes) for which we couldn't unpack the header. Error: %s", pt, p.Size(), err)
		} else {
			// This is an error returned by the AEAD (other than ErrDecryptionFailed).
			// For example, a PROTOCOL_VIOLATION due to key updates.
			s.closeLocal(err)
		}
	}
	return false
}

func (s *connection) handleRetryPacket(hdr *wire.Header, data []byte) bool /* was this a valid Retry */ {
	if s.perspective == protocol.PerspectiveServer {
		if s.tracer != nil {
			s.tracer.DroppedPacket(logging.PacketTypeRetry, protocol.ByteCount(len(data)), logging.PacketDropUnexpectedPacket)
		}
		s.logger.Debugf("Ignoring Retry.")
		return false
	}
	if s.receivedFirstPacket {
		if s.tracer != nil {
			s.tracer.DroppedPacket(logging.PacketTypeRetry, protocol.ByteCount(len(data)), logging.PacketDropUnexpectedPacket)
		}
		s.logger.Debugf("Ignoring Retry, since we already received a packet.")
		return false
	}
	destConnID := s.connIDManager.Get()
	if hdr.SrcConnectionID == destConnID {
		if s.tracer != nil {
			s.tracer.DroppedPacket(logging.PacketTypeRetry, protocol.ByteCount(len(data)), logging.PacketDropUnexpectedPacket)
		}
		s.logger.Debugf("Ignoring Retry, since the server didn't change the Source Connection ID.")
		return false
	}
	// If a token is already set, this means that we already received a Retry from the server.
	// Ignore this Retry packet.
	if s.receivedRetry {
		s.logger.Debugf("Ignoring Retry, since a Retry was already received.")
		return false
	}

	tag := handshake.GetRetryIntegrityTag(data[:len(data)-16], destConnID, hdr.Version)
	if !bytes.Equal(data[len(data)-16:], tag[:]) {
		if s.tracer != nil {
			s.tracer.DroppedPacket(logging.PacketTypeRetry, protocol.ByteCount(len(data)), logging.PacketDropPayloadDecryptError)
		}
		s.logger.Debugf("Ignoring spoofed Retry. Integrity Tag doesn't match.")
		return false
	}

	if s.logger.Debug() {
		s.logger.Debugf("<- Received Retry:")
		(&wire.ExtendedHeader{Header: *hdr}).Log(s.logger)
		s.logger.Debugf("Switching destination connection ID to: %s", hdr.SrcConnectionID)
	}
	if s.tracer != nil {
		s.tracer.ReceivedRetry(hdr)
	}
	newDestConnID := hdr.SrcConnectionID
	s.receivedRetry = true
	if err := s.sentPacketHandler.ResetForRetry(); err != nil {
		s.closeLocal(err)
		return false
	}
	s.handshakeDestConnID = newDestConnID
	s.retrySrcConnID = &newDestConnID
	s.cryptoStreamHandler.ChangeConnectionID(newDestConnID)
	s.packer.SetToken(hdr.Token)
	s.connIDManager.ChangeInitialConnID(newDestConnID)
	s.scheduleSending()
	return true
}

func (s *connection) handleVersionNegotiationPacket(p *receivedPacket) {
	if s.perspective == protocol.PerspectiveServer || // servers never receive version negotiation packets
		s.receivedFirstPacket || s.versionNegotiated { // ignore delayed / duplicated version negotiation packets
		if s.tracer != nil {
			s.tracer.DroppedPacket(logging.PacketTypeVersionNegotiation, p.Size(), logging.PacketDropUnexpectedPacket)
		}
		return
	}

	src, dest, supportedVersions, err := wire.ParseVersionNegotiationPacket(p.data)
	if err != nil {
		if s.tracer != nil {
			s.tracer.DroppedPacket(logging.PacketTypeVersionNegotiation, p.Size(), logging.PacketDropHeaderParseError)
		}
		s.logger.Debugf("Error parsing Version Negotiation packet: %s", err)
		return
	}

	for _, v := range supportedVersions {
		if v == s.version {
			if s.tracer != nil {
				s.tracer.DroppedPacket(logging.PacketTypeVersionNegotiation, p.Size(), logging.PacketDropUnexpectedVersion)
			}
			// The Version Negotiation packet contains the version that we offered.
			// This might be a packet sent by an attacker, or it was corrupted.
			return
		}
	}

	s.logger.Infof("Received a Version Negotiation packet. Supported Versions: %s", supportedVersions)
	if s.tracer != nil {
		s.tracer.ReceivedVersionNegotiationPacket(dest, src, supportedVersions)
	}
	newVersion, ok := protocol.ChooseSupportedVersion(s.config.Versions, supportedVersions)
	if !ok {
		s.destroyImpl(&VersionNegotiationError{
			Ours:   s.config.Versions,
			Theirs: supportedVersions,
		})
		s.logger.Infof("No compatible QUIC version found.")
		return
	}
	if s.tracer != nil {
		s.tracer.NegotiatedVersion(newVersion, s.config.Versions, supportedVersions)
	}

	s.logger.Infof("Switching to QUIC version %s.", newVersion)
	nextPN, _ := s.sentPacketHandler.PeekPacketNumber(protocol.EncryptionInitial)
	s.destroyImpl(&errCloseForRecreating{
		nextPacketNumber: nextPN,
		nextVersion:      newVersion,
	})
}

func (s *connection) handleUnpackedPacket(
	addr net.Addr,
	packet *unpackedPacket,
	ecn protocol.ECN,
	rcvTime time.Time,
	packetSize protocol.ByteCount, // only for logging
) error {
	if !s.receivedFirstPacket {
		s.receivedFirstPacket = true
		if !s.versionNegotiated && s.tracer != nil {
			var clientVersions, serverVersions []protocol.VersionNumber
			switch s.perspective {
			case protocol.PerspectiveClient:
				clientVersions = s.config.Versions
			case protocol.PerspectiveServer:
				serverVersions = s.config.Versions
			}
			s.tracer.NegotiatedVersion(s.version, clientVersions, serverVersions)
		}
		// The server can change the source connection ID with the first Handshake packet.
		if s.perspective == protocol.PerspectiveClient && packet.hdr.IsLongHeader && packet.hdr.SrcConnectionID != s.handshakeDestConnID {
			cid := packet.hdr.SrcConnectionID
			s.logger.Debugf("Received first packet. Switching destination connection ID to: %s", cid)
			s.handshakeDestConnID = cid
			s.connIDManager.ChangeInitialConnID(cid)
		}
		// We create the connection as soon as we receive the first packet from the client.
		// We do that before authenticating the packet.
		// That means that if the source connection ID was corrupted,
		// we might have created a connection with an incorrect source connection ID.
		// Once we authenticate the first packet, we need to update it.
		if s.perspective == protocol.PerspectiveServer {
			if packet.hdr.SrcConnectionID != s.handshakeDestConnID {
				s.handshakeDestConnID = packet.hdr.SrcConnectionID
				s.connIDManager.ChangeInitialConnID(packet.hdr.SrcConnectionID)
			}
			if s.tracer != nil {
				s.tracer.StartedConnection(
					s.conn.LocalAddr(),
					s.conn.RemoteAddr(),
					packet.hdr.SrcConnectionID,
					packet.hdr.DestConnectionID,
				)
			}
		}
	}

	s.lastPacketReceivedTime = rcvTime
	s.firstAckElicitingPacketAfterIdleSentTime = time.Time{}
	s.keepAlivePingSent = false

	var log func([]logging.Frame)
	if s.tracer != nil {
		log = func(frames []logging.Frame) {
			s.tracer.ReceivedLongHeaderPacket(packet.hdr, packetSize, frames)
		}
	}
	isAckEliciting, err := s.handleFrames(addr, packet.data, packet.hdr.DestConnectionID, packet.encryptionLevel, log)
	if err != nil {
		return err
	}
	return s.receivedPacketHandler.ReceivedPacket(packet.hdr.PacketNumber, ecn, packet.encryptionLevel, rcvTime, isAckEliciting)
}

func (s *connection) handleUnpackedShortHeaderPacket(
	addr net.Addr,
	destConnID protocol.ConnectionID,
	pn protocol.PacketNumber,
	data []byte,
	ecn protocol.ECN,
	rcvTime time.Time,
	log func([]logging.Frame),
) error {
	s.lastPacketReceivedTime = rcvTime
	s.firstAckElicitingPacketAfterIdleSentTime = time.Time{}
	s.keepAlivePingSent = false

	isAckEliciting, err := s.handleFrames(addr, data, destConnID, protocol.Encryption1RTT, log)
	if err != nil {
		return err
	}
	return s.receivedPacketHandler.ReceivedPacket(pn, ecn, protocol.Encryption1RTT, rcvTime, isAckEliciting)
}

func (s *connection) handleFrames(
	addr net.Addr,
	data []byte,
	destConnID protocol.ConnectionID,
	encLevel protocol.EncryptionLevel,
	log func([]logging.Frame),
) (isAckEliciting bool, _ error) {
	// Only used for tracing.
	// If we're not tracing, this slice will always remain empty.
	var frames []wire.Frame
	for len(data) > 0 {
		l, frame, err := s.frameParser.ParseNext(data, encLevel)
		if err != nil {
			return false, err
		}
		data = data[l:]
		if frame == nil {
			break
		}
		if ackhandler.IsFrameAckEliciting(frame) {
			isAckEliciting = true
		}
		// Only process frames now if we're not logging.
		// If we're logging, we need to make sure that the packet_received event is logged first.
		if log == nil {
			if err := s.handleFrame(addr, frame, encLevel, destConnID); err != nil {
				return false, err
			}
		} else {
			frames = append(frames, frame)
		}
	}

	if log != nil {
		fs := make([]logging.Frame, len(frames))
		for i, frame := range frames {
			fs[i] = logutils.ConvertFrame(frame)
		}
		log(fs)
		for _, frame := range frames {
			if err := s.handleFrame(addr, frame, encLevel, destConnID); err != nil {
				return false, err
			}
		}
	}
	return
}

func (s *connection) handleFrame(addr net.Addr, f wire.Frame, encLevel protocol.EncryptionLevel, destConnID protocol.ConnectionID) error {
	var err error
	wire.LogFrame(s.logger, f, false)
	switch frame := f.(type) {
	case *wire.CryptoFrame:
		err = s.handleCryptoFrame(frame, encLevel)
	case *wire.StreamFrame:
		err = s.handleStreamFrame(frame)
	case *wire.AckFrame:
		err = s.handleAckFrame(frame, encLevel)
		wire.PutAckFrame(frame)
	case *wire.ConnectionCloseFrame:
		s.handleConnectionCloseFrame(frame)
	case *wire.ResetStreamFrame:
		err = s.handleResetStreamFrame(frame)
	case *wire.MaxDataFrame:
		s.handleMaxDataFrame(frame)
	case *wire.MaxStreamDataFrame:
		err = s.handleMaxStreamDataFrame(frame)
	case *wire.MaxStreamsFrame:
		s.handleMaxStreamsFrame(frame)
	case *wire.DataBlockedFrame:
	case *wire.StreamDataBlockedFrame:
	case *wire.StreamsBlockedFrame:
	case *wire.StopSendingFrame:
		err = s.handleStopSendingFrame(frame)
	case *wire.PingFrame:
	case *wire.PathChallengeFrame:
		s.handlePathChallengeFrame(frame)
	case *wire.PathResponseFrame:
		s.handlePathResponseFrame(addr, frame)
	case *wire.NewTokenFrame:
		err = s.handleNewTokenFrame(frame)
	case *wire.NewConnectionIDFrame:
		err = s.handleNewConnectionIDFrame(frame)
	case *wire.RetireConnectionIDFrame:
		err = s.handleRetireConnectionIDFrame(frame, destConnID)
	case *wire.HandshakeDoneFrame:
		err = s.handleHandshakeDoneFrame()
	case *wire.DatagramFrame:
		err = s.handleDatagramFrame(frame)
	default:
		err = fmt.Errorf("unexpected frame type: %s", reflect.ValueOf(&frame).Elem().Type().Name())
	}
	return err
}

// handlePacket is called by the server with a new packet
func (s *connection) handlePacket(p *receivedPacket) {
	// Discard packets once the amount of queued packets is larger than
	// the channel size, protocol.MaxConnUnprocessedPackets
	select {
	case s.receivedPackets <- p:
	default:
		if s.tracer != nil {
			s.tracer.DroppedPacket(logging.PacketTypeNotDetermined, p.Size(), logging.PacketDropDOSPrevention)
		}
	}
}

func (s *connection) handleConnectionCloseFrame(frame *wire.ConnectionCloseFrame) {
	if frame.IsApplicationError {
		s.closeRemote(&qerr.ApplicationError{
			Remote:       true,
			ErrorCode:    qerr.ApplicationErrorCode(frame.ErrorCode),
			ErrorMessage: frame.ReasonPhrase,
		})
		return
	}
	s.closeRemote(&qerr.TransportError{
		Remote:       true,
		ErrorCode:    qerr.TransportErrorCode(frame.ErrorCode),
		FrameType:    frame.FrameType,
		ErrorMessage: frame.ReasonPhrase,
	})
}

func (s *connection) handleCryptoFrame(frame *wire.CryptoFrame, encLevel protocol.EncryptionLevel) error {
	encLevelChanged, err := s.cryptoStreamManager.HandleCryptoFrame(frame, encLevel)
	if err != nil {
		return err
	}
	if encLevelChanged {
		// Queue all packets for decryption that have been undecryptable so far.
		s.undecryptablePacketsToProcess = s.undecryptablePackets
		s.undecryptablePackets = nil
	}
	return nil
}

func (s *connection) handleStreamFrame(frame *wire.StreamFrame) error {
	if s.isWaitingForHandoverMigration() {
		// Cannot be handled before handover
		// ignore this StreamFrame
		return nil
	}
	str, err := s.streamsMap.GetOrOpenReceiveStream(frame.StreamID)
	if err != nil {
		return err
	}
	if str == nil {
		// Stream is closed and already garbage collected
		// ignore this StreamFrame
		return nil
	}
	return str.handleStreamFrame(frame)
}

func (s *connection) handleMaxDataFrame(frame *wire.MaxDataFrame) {
	s.connFlowController.UpdateSendWindow(frame.MaximumData)
}

func (s *connection) handleMaxStreamDataFrame(frame *wire.MaxStreamDataFrame) error {
	str, err := s.streamsMap.GetOrOpenSendStream(frame.StreamID)
	if err != nil {
		return err
	}
	if str == nil {
		// stream is closed and already garbage collected
		return nil
	}
	str.updateSendWindow(frame.MaximumStreamData)
	return nil
}

func (s *connection) handleMaxStreamsFrame(frame *wire.MaxStreamsFrame) {
	s.streamsMap.HandleMaxStreamsFrame(frame)
}

func (s *connection) handleResetStreamFrame(frame *wire.ResetStreamFrame) error {
	str, err := s.streamsMap.GetOrOpenReceiveStream(frame.StreamID)
	if err != nil {
		return err
	}
	if str == nil {
		// stream is closed and already garbage collected
		return nil
	}
	return str.handleResetStreamFrame(frame)
}

func (s *connection) handleStopSendingFrame(frame *wire.StopSendingFrame) error {
	str, err := s.streamsMap.GetOrOpenSendStream(frame.StreamID)
	if err != nil {
		return err
	}
	if str == nil {
		// stream is closed and already garbage collected
		return nil
	}
	str.handleStopSendingFrame(frame)
	return nil
}

func (s *connection) handlePathChallengeFrame(frame *wire.PathChallengeFrame) {
	s.queueControlFrame(&wire.PathResponseFrame{Data: frame.Data})
}

func (s *connection) handlePathResponseFrame(addr net.Addr, frame *wire.PathResponseFrame) {
	s.pathManager.OnReceivePathResponseFrame(addr.(*net.UDPAddr), frame)
}

func (s *connection) handleNewTokenFrame(frame *wire.NewTokenFrame) error {
	if s.perspective == protocol.PerspectiveServer {
		return &qerr.TransportError{
			ErrorCode:    qerr.ProtocolViolation,
			ErrorMessage: "received NEW_TOKEN frame from the client",
		}
	}
	if s.config.TokenStore != nil {
		s.config.TokenStore.Put(s.tokenStoreKey, &ClientToken{data: frame.Token})
	}
	return nil
}

func (s *connection) handleNewConnectionIDFrame(f *wire.NewConnectionIDFrame) error {
	return s.connIDManager.Add(f)
}

func (s *connection) handleRetireConnectionIDFrame(f *wire.RetireConnectionIDFrame, destConnID protocol.ConnectionID) error {
	return s.connIDGenerator.Retire(f.SequenceNumber, destConnID)
}

func (s *connection) handleHandshakeDoneFrame() error {
	if s.perspective == protocol.PerspectiveServer {
		return &qerr.TransportError{
			ErrorCode:    qerr.ProtocolViolation,
			ErrorMessage: "received a HANDSHAKE_DONE frame",
		}
	}
	if !s.handshakeConfirmed {
		s.handleHandshakeConfirmed()
	}
	return nil
}

func (s *connection) handleAckFrame(frame *wire.AckFrame, encLevel protocol.EncryptionLevel) error {
	acked1RTTPacket, err := s.sentPacketHandler.ReceivedAck(frame, encLevel, s.lastPacketReceivedTime)
	if err != nil {
		return err
	}
	if !acked1RTTPacket {
		return nil
	}
	if s.perspective == protocol.PerspectiveClient && !s.handshakeConfirmed {
		s.handleHandshakeConfirmed()
	}
	return s.cryptoStreamHandler.SetLargest1RTTAcked(frame.LargestAcked())
}

func (s *connection) handleDatagramFrame(f *wire.DatagramFrame) error {
	if f.Length(s.version) > protocol.MaxDatagramFrameSize {
		return &qerr.TransportError{
			ErrorCode:    qerr.ProtocolViolation,
			ErrorMessage: "DATAGRAM frame too large",
		}
	}
	s.datagramQueue.HandleDatagramFrame(f)
	return nil
}

// closeLocal closes the connection and send a CONNECTION_CLOSE containing the error
func (s *connection) closeLocal(e error) {
	s.closeOnce.Do(func() {
		if e == nil {
			s.logger.Infof("Closing connection.")
		} else {
			s.logger.Errorf("Closing connection with error: %s", e)
		}
		s.closeChan <- closeError{err: e, immediate: false, remote: false}
	})
}

// destroy closes the connection without sending the error on the wire
func (s *connection) destroy(e error) {
	s.destroyImpl(e)
	<-s.ctx.Done()
}

func (s *connection) destroyImpl(e error) {
	s.closeOnce.Do(func() {
		if nerr, ok := e.(net.Error); ok && nerr.Timeout() {
			s.logger.Errorf("Destroying connection: %s", e)
		} else {
			s.logger.Errorf("Destroying connection with error: %s", e)
		}
		s.closeChan <- closeError{err: e, immediate: true, remote: false}
	})
}

func (s *connection) closeRemote(e error) {
	s.closeOnce.Do(func() {
		s.logger.Errorf("Peer closed connection with error: %s", e)
		s.closeChan <- closeError{err: e, immediate: true, remote: true}
	})
}

// Close the connection. It sends a NO_ERROR application error.
// It waits until the run loop has stopped before returning
func (s *connection) shutdown() {
	s.closeLocal(nil)
	<-s.ctx.Done()
}

func (s *connection) CloseWithError(code ApplicationErrorCode, desc string) error {
	s.closeLocal(&qerr.ApplicationError{
		ErrorCode:    code,
		ErrorMessage: desc,
	})
	<-s.ctx.Done()
	return nil
}

func (s *connection) handleCloseError(closeErr *closeError) {
	e := closeErr.err
	if e == nil {
		e = &qerr.ApplicationError{}
	} else {
		defer func() {
			closeErr.err = e
		}()
	}

	var (
		statelessResetErr     *StatelessResetError
		versionNegotiationErr *VersionNegotiationError
		recreateErr           *errCloseForRecreating
		applicationErr        *ApplicationError
		transportErr          *TransportError
	)
	switch {
	case errors.Is(e, qerr.ErrIdleTimeout),
		errors.Is(e, qerr.ErrHandshakeTimeout),
		errors.As(e, &statelessResetErr),
		errors.As(e, &versionNegotiationErr),
		errors.As(e, &recreateErr),
		errors.As(e, &applicationErr),
		errors.As(e, &transportErr):
	default:
		e = &qerr.TransportError{
			ErrorCode:    qerr.InternalError,
			ErrorMessage: e.Error(),
		}
	}

	s.streamsMap.CloseWithError(e)
	s.connIDManager.Close()
	if s.datagramQueue != nil {
		s.datagramQueue.CloseWithError(e)
	}

	if s.tracer != nil && !errors.As(e, &recreateErr) {
		s.tracer.ClosedConnection(e)
	}

	// If this is a remote close we're done here
	if closeErr.remote {
		s.connIDGenerator.ReplaceWithClosed(s.perspective, nil)
		return
	}
	if closeErr.immediate {
		s.connIDGenerator.RemoveAll()
		return
	}
	// Don't send out any CONNECTION_CLOSE if this is an error that occurred
	// before we even sent out the first packet.
	if s.perspective == protocol.PerspectiveClient && !s.sentFirstPacket {
		s.connIDGenerator.RemoveAll()
		return
	}
	connClosePacket, err := s.sendConnectionClose(e)
	if err != nil {
		s.logger.Debugf("Error sending CONNECTION_CLOSE: %s", err)
	}
	s.connIDGenerator.ReplaceWithClosed(s.perspective, connClosePacket)
}

func (s *connection) dropEncryptionLevel(encLevel protocol.EncryptionLevel) {
	s.sentPacketHandler.DropPackets(encLevel)
	s.receivedPacketHandler.DropPackets(encLevel)
	if s.tracer != nil {
		s.tracer.DroppedEncryptionLevel(encLevel)
	}
	if encLevel == protocol.Encryption0RTT {
		s.streamsMap.ResetFor0RTT()
		if err := s.connFlowController.Reset(); err != nil {
			s.closeLocal(err)
		}
		if err := s.framer.Handle0RTTRejection(); err != nil {
			s.closeLocal(err)
		}
	}
}

// is called for the client, when restoring transport parameters saved for 0-RTT
func (s *connection) restoreTransportParameters(params *wire.TransportParameters) {
	if s.logger.Debug() {
		s.logger.Debugf("Restoring Transport Parameters: %s", params)
	}

	s.peerParams = params
	s.connIDGenerator.SetMaxActiveConnIDs(params.ActiveConnectionIDLimit)
	s.connFlowController.UpdateSendWindow(params.InitialMaxData)
	s.streamsMap.UpdateLimits(params)

	if s.config.ExtraStreamEncryption.enabled() && params.ExtraStreamEncryption {
		s.streamsMap.SetXseCryptoSetup(xse.NewCryptoSetupFromConn(s.cryptoStreamHandler.TlsConn(), s.perspective, s.tracer))
	}
}

func (s *connection) handleTransportParameters(params *wire.TransportParameters) {
	if err := s.checkTransportParameters(params); err != nil {
		s.closeLocal(&qerr.TransportError{
			ErrorCode:    qerr.TransportParameterError,
			ErrorMessage: err.Error(),
		})
	}
	s.peerParams = params
	// On the client side we have to wait for handshake completion.
	// During a 0-RTT connection, we are only allowed to use the new transport parameters for 1-RTT packets.
	if s.perspective == protocol.PerspectiveServer {
		s.applyTransportParameters()
		// On the server side, the early connection is ready as soon as we processed
		// the client's transport parameters.
		close(s.earlyConnReadyChan)
	}
}

func (s *connection) checkTransportParameters(params *wire.TransportParameters) error {
	if s.logger.Debug() {
		s.logger.Debugf("Processed Transport Parameters: %s", params)
	}
	if s.tracer != nil {
		s.tracer.ReceivedTransportParameters(params)
	}

	// check the initial_source_connection_id
	if params.InitialSourceConnectionID != s.handshakeDestConnID {
		return fmt.Errorf("expected initial_source_connection_id to equal %s, is %s", s.handshakeDestConnID, params.InitialSourceConnectionID)
	}

	if s.config.ExtraStreamEncryption == EnforceExtraStreamEncryption && !params.ExtraStreamEncryption {
		return fmt.Errorf("missing extra_stream_encryption; peer does not support the XSE-QUIC extension")
	}

	if s.perspective == protocol.PerspectiveServer {
		return nil
	}
	// check the original_destination_connection_id
	if params.OriginalDestinationConnectionID != s.origDestConnID {
		return fmt.Errorf("expected original_destination_connection_id to equal %s, is %s", s.origDestConnID, params.OriginalDestinationConnectionID)
	}
	if s.retrySrcConnID != nil { // a Retry was performed
		if params.RetrySourceConnectionID == nil {
			return errors.New("missing retry_source_connection_id")
		}
		if *params.RetrySourceConnectionID != *s.retrySrcConnID {
			return fmt.Errorf("expected retry_source_connection_id to equal %s, is %s", s.retrySrcConnID, *params.RetrySourceConnectionID)
		}
	} else if params.RetrySourceConnectionID != nil {
		return errors.New("received retry_source_connection_id, although no Retry was performed")
	}
	return nil
}

func (s *connection) applyTransportParameters() {
	params := s.peerParams
	// Our local idle timeout will always be > 0.
	s.idleTimeout = utils.MinNonZeroDuration(s.config.MaxIdleTimeout, params.MaxIdleTimeout)
	s.keepAliveInterval = utils.Min(s.config.KeepAlivePeriod, utils.Min(s.idleTimeout/2, protocol.MaxKeepAliveInterval))
	s.streamsMap.UpdateLimits(params)
	s.packer.HandleTransportParameters(params)
	s.frameParser.SetAckDelayExponent(params.AckDelayExponent)
	s.connFlowController.UpdateSendWindow(params.InitialMaxData)
	s.rttStats.SetMaxAckDelay(params.MaxAckDelay)
	s.connIDGenerator.SetMaxActiveConnIDs(params.ActiveConnectionIDLimit)
	if params.StatelessResetToken != nil && s.connIDManager.activeSequenceNumber == 0 {
		s.connIDManager.SetStatelessResetToken(*params.StatelessResetToken)
	}
	// We don't support connection migration yet, so we don't have any use for the preferred_address.
	if params.PreferredAddress != nil {
		// Retire the connection ID.
		s.connIDManager.AddFromPreferredAddress(params.PreferredAddress.ConnectionID, params.PreferredAddress.StatelessResetToken)
	}
	if s.config.ExtraStreamEncryption.enabled() && params.ExtraStreamEncryption {
		s.streamsMap.SetXseCryptoSetup(xse.NewCryptoSetupFromConn(s.cryptoStreamHandler.TlsConn(), s.perspective, s.tracer))
	}
}

func (s *connection) sendPackets() error {
	s.pacingDeadline = time.Time{}

	var sentPacket bool // only used in for packets sent in send mode SendAny
	for {
		sendMode := s.sentPacketHandler.SendMode()
		if sendMode == ackhandler.SendAny && s.handshakeComplete && !s.sentPacketHandler.HasPacingBudget() {
			deadline := s.sentPacketHandler.TimeUntilSend()
			if deadline.IsZero() {
				deadline = deadlineSendImmediately
			}
			s.pacingDeadline = deadline
			// Allow sending of an ACK if we're pacing limit (if we haven't sent out a packet yet).
			// This makes sure that a peer that is mostly receiving data (and thus has an inaccurate cwnd estimate)
			// sends enough ACKs to allow its peer to utilize the bandwidth.
			if sentPacket {
				return nil
			}
			sendMode = ackhandler.SendAck
		}

		err := s.maybeSendPathChallengeOnlyPacket()
		if err != nil {
			return err
		}

		switch sendMode {
		case ackhandler.SendNone:
			return nil
		case ackhandler.SendAck:
			// If we already sent packets, and the send mode switches to SendAck,
			// as we've just become congestion limited.
			// There's no need to try to send an ACK at this moment.
			if sentPacket {
				return nil
			}
			// We can at most send a single ACK only packet.
			// There will only be a new ACK after receiving new packets.
			// SendAck is only returned when we're congestion limited, so we don't need to set the pacinggs timer.
			return s.maybeSendAckOnlyPacket()
		case ackhandler.SendPTOInitial:
			if err := s.sendProbePacket(protocol.EncryptionInitial); err != nil {
				return err
			}
		case ackhandler.SendPTOHandshake:
			if err := s.sendProbePacket(protocol.EncryptionHandshake); err != nil {
				return err
			}
		case ackhandler.SendPTOAppData:
			if err := s.sendProbePacket(protocol.Encryption1RTT); err != nil {
				return err
			}
		case ackhandler.SendAny:
			sent, err := s.sendPacket()
			if err != nil || !sent {
				return err
			}
			sentPacket = true
		default:
			return fmt.Errorf("BUG: invalid send mode %d", sendMode)
		}
		// Prioritize receiving of packets over sending out more packets.
		if len(s.receivedPackets) > 0 {
			s.pacingDeadline = deadlineSendImmediately
			return nil
		}
		if s.sendQueue.WouldBlock() {
			return nil
		}
	}
}

func (s *connection) maybeSendPathChallengeOnlyPacket() error {
	path := s.pathManager.CurrentSendPath()
	if !path.PeerAddressValidated() && s.handshakeComplete {
		packet, err := s.packer.PackPathChallengeOnlyPacket()
		if err != nil {
			return err
		}
		if packet == nil {
			return nil
		}
		s.sendPackedPacket(packet, time.Now())
	}
	return nil
}

func (s *connection) maybeSendAckOnlyPacket() error {
	if !s.handshakeConfirmed {
		packet, err := s.packer.PackCoalescedPacket(true)
		if err != nil {
			return err
		}
		if packet == nil {
			return nil
		}
		s.logCoalescedPacket(packet)
		for _, p := range packet.packets {
			s.sentPacketHandler.SentPacket(p.ToAckHandlerPacket(time.Now(), s.retransmissionQueue))
		}
		s.connIDManager.SentPacket()
		s.sendQueue.Send(packet.buffer)
		return nil
	}

	packet, err := s.packer.PackPacket(true)
	if err != nil {
		return err
	}
	if packet == nil {
		return nil
	}
	s.sendPackedPacket(packet, time.Now())
	return nil
}

// TODO clarify because Ping is a non-probing frame
func (s *connection) sendProbePacket(encLevel protocol.EncryptionLevel) error {
	// Queue probe packets until we actually send out a packet,
	// or until there are no more packets to queue.
	var packet *packedPacket
	for {
		if wasQueued := s.sentPacketHandler.QueueProbePacket(encLevel); !wasQueued {
			break
		}
		var err error
		packet, err = s.packer.MaybePackProbePacket(encLevel)
		if err != nil {
			return err
		}
		if packet != nil {
			break
		}
	}
	if packet == nil {
		//nolint:exhaustive // Cannot send probe packets for 0-RTT.
		switch encLevel {
		case protocol.EncryptionInitial:
			s.retransmissionQueue.AddInitial(&wire.PingFrame{})
		case protocol.EncryptionHandshake:
			s.retransmissionQueue.AddHandshake(&wire.PingFrame{})
		case protocol.Encryption1RTT:
			s.retransmissionQueue.AddAppData(&wire.PingFrame{})
		default:
			panic("unexpected encryption level")
		}
		var err error
		packet, err = s.packer.MaybePackProbePacket(encLevel)
		if err != nil {
			return err
		}
	}
	if packet == nil || packet.packetContents == nil {
		return fmt.Errorf("connection BUG: couldn't pack %s probe packet", encLevel)
	}
	s.sendPackedPacket(packet, time.Now())
	return nil
}

func (s *connection) sendPacket() (bool, error) {
	if s.isRemoteAddressIgnored(s.conn.RemoteAddr()) && s.handshakeComplete {
		s.logger.Debugf("Ignore sending packets to this remote address")
		return false, nil
	}

	if isBlocked, offset := s.connFlowController.IsNewlyBlocked(); isBlocked {
		s.framer.QueueControlFrame(&wire.DataBlockedFrame{MaximumData: offset})
	}
	s.windowUpdateQueue.QueueAll()

	now := time.Now()
	if !s.handshakeConfirmed {
		packet, err := s.packer.PackCoalescedPacket(false)
		if err != nil || packet == nil {
			return false, err
		}
		s.sentFirstPacket = true
		s.logCoalescedPacket(packet)
		for _, p := range packet.packets {
			if s.firstAckElicitingPacketAfterIdleSentTime.IsZero() && p.IsAckEliciting() {
				s.firstAckElicitingPacketAfterIdleSentTime = now
			}
			s.sentPacketHandler.SentPacket(p.ToAckHandlerPacket(now, s.retransmissionQueue))
		}
		s.connIDManager.SentPacket()
		s.sendQueue.Send(packet.buffer)
		return true, nil
	}
	if !s.config.DisablePathMTUDiscovery && s.mtuDiscoverer.ShouldSendProbe(now) {
		packet, err := s.packer.PackMTUProbePacket(s.mtuDiscoverer.GetPing())
		if err != nil {
			return false, err
		}
		s.sendPackedPacket(packet, now)
		return true, nil
	}
	packet, err := s.packer.PackPacket(false)
	if err != nil || packet == nil {
		return false, err
	}
	s.sendPackedPacket(packet, now)
	return true, nil
}

func (s *connection) sendPackedPacket(packet *packedPacket, now time.Time) {
	if s.isRemoteAddressIgnored(s.conn.RemoteAddr()) {
		s.logger.Debugf("Ignore sending the following packet to this remote address:")
		s.logPacket(packet)
		return
	}
	//TODO this can be removed, when server ensures HandshakeDoneFrame is sent and stream can be already opened
	if s.isWaitingForHandoverMigration() && packet.EncryptionLevel() == protocol.Encryption1RTT && s.perspective == protocol.PerspectiveClient {
		s.logger.Debugf("Ignore sending 1RTT packets until path change")
		s.logPacket(packet)
		return
	}
	if s.firstAckElicitingPacketAfterIdleSentTime.IsZero() && packet.IsAckEliciting() {
		s.firstAckElicitingPacketAfterIdleSentTime = now
	}
	s.logPacket(packet)
	s.sentPacketHandler.SentPacket(packet.ToAckHandlerPacket(now, s.retransmissionQueue))
	s.connIDManager.SentPacket()
	s.sendQueue.Send(packet.buffer)
}

func (s *connection) sendConnectionClose(e error) ([]byte, error) {
	//TODO try to remove
	if s.isRemoteAddressIgnored(s.conn.RemoteAddr()) {
		return nil, nil
	}
	var packet *coalescedPacket
	var err error
	var transportErr *qerr.TransportError
	var applicationErr *qerr.ApplicationError
	if errors.As(e, &transportErr) {
		packet, err = s.packer.PackConnectionClose(transportErr)
	} else if errors.As(e, &applicationErr) {
		packet, err = s.packer.PackApplicationClose(applicationErr)
	} else {
		packet, err = s.packer.PackConnectionClose(&qerr.TransportError{
			ErrorCode:    qerr.InternalError,
			ErrorMessage: fmt.Sprintf("connection BUG: unspecified error type (msg: %s)", e.Error()),
		})
	}
	if err != nil {
		return nil, err
	}
	s.logCoalescedPacket(packet)
	return packet.buffer.Data, s.conn.Write(packet.buffer.Data)
}

func (s *connection) logPacketContents(p *packetContents) {
	// tracing
	if s.tracer != nil {
		frames := make([]logging.Frame, 0, len(p.frames))
		for _, f := range p.frames {
			frames = append(frames, logutils.ConvertFrame(f.Frame))
		}
		var ack *logging.AckFrame
		if p.ack != nil {
			ack = logutils.ConvertAckFrame(p.ack)
		}
		s.tracer.SentPacket(p.header, p.length, ack, frames)
	}

	// quic-go logging
	if !s.logger.Debug() {
		return
	}
	p.header.Log(s.logger)
	if p.ack != nil {
		wire.LogFrame(s.logger, p.ack, true)
	}
	for _, frame := range p.frames {
		wire.LogFrame(s.logger, frame.Frame, true)
	}
}

func (s *connection) logCoalescedPacket(packet *coalescedPacket) {
	if s.logger.Debug() {
		if len(packet.packets) > 1 {
			s.logger.Debugf("-> Sending coalesced packet (%d parts, %d bytes) for connection %s", len(packet.packets), packet.buffer.Len(), s.logID)
		} else {
			s.logger.Debugf("-> Sending packet %d (%d bytes) for connection %s, %s", packet.packets[0].header.PacketNumber, packet.buffer.Len(), s.logID, packet.packets[0].EncryptionLevel())
		}
	}
	for _, p := range packet.packets {
		s.logPacketContents(p)
	}
}

func (s *connection) logPacket(packet *packedPacket) {
	if s.logger.Debug() {
		s.logger.Debugf("-> Sending packet %d (%d bytes) for connection %s, %s", packet.header.PacketNumber, packet.buffer.Len(), s.logID, packet.EncryptionLevel())
	}
	s.logPacketContents(packet.packetContents)
}

// awaitMigration returns nil when no error occurred
func (s *connection) awaitHandoverMigration() error {
	select {
	case <-s.handoverMigrationCtx.Done():
		return nil
	case <-s.ctx.Done():
		return nil
	case <-time.After(s.config.MaxIdleTimeout):
		err := fmt.Errorf("failed to open stream: handover timeout")
		s.closeLocal(err)
		return err
	}
}

// AcceptStream returns the next stream openend by the peer
func (s *connection) AcceptStream(ctx context.Context) (Stream, error) {
	err := s.awaitHandoverMigration()
	if err != nil {
		return nil, err
	}
	return s.streamsMap.AcceptStream(ctx)
}

func (s *connection) AcceptUniStream(ctx context.Context) (ReceiveStream, error) {
	err := s.awaitHandoverMigration()
	if err != nil {
		return nil, err
	}
	return s.streamsMap.AcceptUniStream(ctx)
}

// OpenStream opens a stream
func (s *connection) OpenStream() (Stream, error) {
	err := s.awaitHandoverMigration()
	if err != nil {
		return nil, err
	}
	return s.streamsMap.OpenStream()
}

func (s *connection) OpenStreamSync(ctx context.Context) (Stream, error) {
	err := s.awaitHandoverMigration()
	if err != nil {
		return nil, err
	}
	return s.streamsMap.OpenStreamSync(ctx)
}

func (s *connection) OpenUniStream() (SendStream, error) {
	err := s.awaitHandoverMigration()
	if err != nil {
		return nil, err
	}
	return s.streamsMap.OpenUniStream()
}

func (s *connection) OpenUniStreamSync(ctx context.Context) (SendStream, error) {
	err := s.awaitHandoverMigration()
	if err != nil {
		return nil, err
	}
	return s.streamsMap.OpenUniStreamSync(ctx)
}

func (s *connection) newFlowController(id protocol.StreamID) flowcontrol.StreamFlowController {
	initialSendWindow := s.peerParams.InitialMaxStreamDataUni
	if id.Type() == protocol.StreamTypeBidi {
		if id.InitiatedBy() == s.perspective {
			initialSendWindow = s.peerParams.InitialMaxStreamDataBidiRemote
		} else {
			initialSendWindow = s.peerParams.InitialMaxStreamDataBidiLocal
		}
	}
	return flowcontrol.NewStreamFlowController(
		id,
		s.connFlowController,
		protocol.ByteCount(s.config.InitialStreamReceiveWindow),
		protocol.ByteCount(s.config.MaxStreamReceiveWindow),
		initialSendWindow,
		s.onHasStreamWindowUpdate,
		s.rttStats,
		s.logger,
	)
}

// scheduleSending signals that we have data for sending
func (s *connection) scheduleSending() {
	select {
	case s.sendingScheduled <- struct{}{}:
	default:
	}
}

// tryQueueingUndecryptablePacket queues a packet for which we're missing the decryption keys.
// The logging.PacketType is only used for logging purposes.
func (s *connection) tryQueueingUndecryptablePacket(p *receivedPacket, pt logging.PacketType) {
	if s.handshakeComplete {
		panic("shouldn't queue undecryptable packets after handshake completion")
	}
	if len(s.undecryptablePackets)+1 > protocol.MaxUndecryptablePackets {
		if s.tracer != nil {
			s.tracer.DroppedPacket(pt, p.Size(), logging.PacketDropDOSPrevention)
		}
		s.logger.Infof("Dropping undecryptable packet (%d bytes). Undecryptable packet queue full.", p.Size())
		return
	}
	s.logger.Infof("Queueing packet (%d bytes) for later decryption", p.Size())
	if s.tracer != nil {
		s.tracer.BufferedPacket(pt, p.Size())
	}
	s.undecryptablePackets = append(s.undecryptablePackets, p)
}

func (s *connection) queuePathChallengeFrame(path path.Path) {
	s.queueControlFrame(&wire.PathChallengeFrame{Data: path.ChallengeData()})
}

func (s *connection) queueControlFrame(f wire.Frame) {
	s.framer.QueueControlFrame(f)
	s.scheduleSending()
}

func (s *connection) onHasStreamWindowUpdate(id protocol.StreamID) {
	s.windowUpdateQueue.AddStream(id)
	s.scheduleSending()
}

func (s *connection) onHasConnectionWindowUpdate() {
	s.windowUpdateQueue.AddConnection()
	s.scheduleSending()
}

func (s *connection) onHasStreamData(id protocol.StreamID) {
	s.framer.AddActiveStream(id)
	s.scheduleSending()
}

func (s *connection) onStreamCompleted(id protocol.StreamID) {
	if err := s.streamsMap.DeleteStream(id); err != nil {
		s.closeLocal(err)
	}
}

func (s *connection) SendMessage(p []byte) error {
	if !s.supportsDatagrams() {
		return errors.New("datagram support disabled")
	}

	f := &wire.DatagramFrame{DataLenPresent: true}
	if protocol.ByteCount(len(p)) > f.MaxDataLen(s.peerParams.MaxDatagramFrameSize, s.version) {
		return errors.New("message too large")
	}
	f.Data = make([]byte, len(p))
	copy(f.Data, p)
	return s.datagramQueue.AddAndWait(f)
}

func (s *connection) ReceiveMessage() ([]byte, error) {
	if !s.config.EnableDatagrams {
		return nil, errors.New("datagram support disabled")
	}
	return s.datagramQueue.Receive()
}

func (s *connection) LocalAddr() net.Addr {
	return s.conn.LocalAddr()
}

func (s *connection) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *connection) getPerspective() protocol.Perspective {
	return s.perspective
}

func (s *connection) GetVersion() protocol.VersionNumber {
	return s.version
}

func (s *connection) NextConnection() Connection {
	<-s.HandshakeComplete().Done()
	s.streamsMap.UseResetMaps()
	return s
}

// Handover waits until handshake is confirmed
// if AllowEarlyHandover is set, waits until handshake is completed
func (s *connection) Handover(destroy bool, ignoreCurrentPath bool) (handover.State, error) {
	if s.config.AllowEarlyHandover {
		<-s.handshakeCtx.Done()
	} else {
		// wait until confirmed, as described in RFC 9000
		_, _ = <-s.handshakeConfirmedChan
	}

	return s.handover(destroy, ignoreCurrentPath)
}

// handover must be called after handshake is confirmed
// if AllowEarlyHandover is set, it can be called after handshake is completed
func (s *connection) handover(destroy bool, ignoreCurrentPath bool) (handover.State, error) {
	//TODO validate supported options
	if s.config.EnableDatagrams {
		panic("option is currently not supported")
	}
	if !s.handshakeComplete {
		panic("illegal handover state: handshake must be completed")
	}

	if !s.config.AllowEarlyHandover && !s.handshakeConfirmed {
		// must not initiate connection migration before handshake is confirmed, as described in RFC 9000
		panic("illegal handover state: handshake must be confirmed")
	}

	if destroy {
		s.destroy(errors.New("destroyed by handover"))
	}

	if ignoreCurrentPath {
		s.IgnoreCurrentRemoteAddress()
	}

	// TODO allow handover with open streams
	if len(s.streamsMap.(*streamsMap).incomingBidiStreams.streams) != 0 ||
		len(s.streamsMap.(*streamsMap).outgoingBidiStreams.streams) != 0 ||
		len(s.streamsMap.(*streamsMap).incomingUniStreams.streams) != 0 ||
		len(s.streamsMap.(*streamsMap).outgoingUniStreams.streams) != 0 {
		panic("illegal handover state")
	}

	state := handover.State{}

	state.Version = s.version
	logConnectionId, err := protocol.ParseConnectionIDHex(s.logID)
	if err != nil {
		return handover.State{}, err
	}
	state.OriginalDestinationConnectionID = logConnectionId.Bytes()
	state.SetRemoteAddress(s.perspective, *s.conn.RemoteAddr().(*net.UDPAddr))
	state.SetLocalAddress(s.perspective, *s.conn.LocalAddr().(*net.UDPAddr)) //TODO
	//state.SetHighest1RTTPacketNumber(s.perspective, s.sentPacketHandler.Highest1RTTPacketNumber()) //TODO
	s.cryptoStreamHandler.StoreHandoverState(&state, s.perspective)

	activeSrcConnIDs := make([]handover.ActiveConnectionID, 0)
	for sn, connID := range s.connIDGenerator.activeSrcConnIDs {
		statelessResetToken := s.connIDGenerator.getStatelessResetToken(connID) //TODO maybe store history
		activeSrcConnIDs = append(activeSrcConnIDs, handover.ActiveConnectionID{
			SequenceNumber:      sn,
			ConnectionID:        connID.Bytes(),
			StatelessResetToken: statelessResetToken[:],
		})
	}
	state.SetActiveSrcConnectionIDs(s.perspective, activeSrcConnIDs)

	// store dest conn ids
	var activeDestStatelessResetToken []byte = nil
	if s.connIDManager.activeStatelessResetToken != nil {
		activeDestStatelessResetToken = s.connIDManager.activeStatelessResetToken[:]
	}
	activeDestConnIDs := make([]handover.ActiveConnectionID, 0)
	activeDestConnIDs = append(activeDestConnIDs, handover.ActiveConnectionID{
		SequenceNumber:      s.connIDManager.activeSequenceNumber,
		ConnectionID:        s.connIDManager.activeConnectionID.Bytes(),
		StatelessResetToken: activeDestStatelessResetToken,
	})
	for destConnID := s.connIDManager.queue.Front(); destConnID != nil; destConnID = destConnID.Next() {
		activeDestConnIDs = append(activeDestConnIDs, handover.ActiveConnectionID{
			SequenceNumber:      destConnID.Value.SequenceNumber,
			ConnectionID:        destConnID.Value.ConnectionID.Bytes(),
			StatelessResetToken: destConnID.Value.StatelessResetToken[:],
		})
	}
	state.SetActiveDestConnectionIDs(s.perspective, activeDestConnIDs)

	if s.perspective == protocol.PerspectiveServer {
		//state.ServerHighestConnectionIDSequenceNumber = s.connIDGenerator.highestSeq //TODO
	}

	state.SetHighestSentPacketNumber(s.perspective, s.sentPacketHandler.Highest1RTTPacketNumber())
	state.SetHighestSentPacketNumber(s.perspective.Opposite(), s.receivedPacketHandler.Highest1RTTPacketNumber()) // this is an estimate

	if s.tracer != nil {
		s.tracer.Debug("hquic_handover_state_created", "")
	}

	return state, nil
}

// update path to peer
// migrate to new remote address
// is called by s.pathManager
func (s *connection) updateSendPath(path path.Path) {
	s.conn.SetRemoteAddr(path.Addr())
	s.rttStats.OnConnectionMigration()
	s.sentPacketHandler.UpdateSendPath(path) // and reset congestion control
	s.handoverMigrationCtxCancel()
}

// Returns nil when handshake is confirmed,
// otherwise return error
func (s *connection) awaitHandshakeConfirmed() error {
	for {
		if s.handshakeConfirmed {
			return nil // handshake is done
		}
		select {
		case <-s.ctx.Done():
			//TODO wrap cause if close is caused by error
			return fmt.Errorf("session closed before handshake confirmed")
		case <-time.After(time.Millisecond):
			continue
		}
	}
}

func (s *connection) MigrateUDPSocket() (*net.UDPAddr, error) {
	err := s.awaitHandshakeConfirmed()
	if err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	if s.peerParams.DisableActiveMigration {
		return nil, errors.New("peer disabled active migration")
	}

	switch conn := s.conn.(type) {
	case *sconn:
		switch conn := conn.rawConn.(type) {
		case *basicConn:
			switch conn := conn.PacketConn.(type) {
			case *MigratableUDPConn:
				return conn.MigrateUDPSocket()
			}
		}
	case *spconn:
		switch conn := conn.PacketConn.(type) {
		case *MigratableUDPConn:
			return conn.MigrateUDPSocket()
		}
	}

	return nil, errors.New("unexpected type")
}

// IgnoreCurrentRemoteAddress adds current remote address to ignore list.
// Session will drop all packets from and to the current destination.
func (s *connection) IgnoreCurrentRemoteAddress() {
	s.ignoredRemoteAddresses = append(s.ignoredRemoteAddresses, *s.RemoteAddr().(*net.UDPAddr))
}

func (s *connection) isRemoteAddressIgnored(addr net.Addr) bool {
	if addr == nil {
		return false
	}
	switch addr := addr.(type) {
	case *net.UDPAddr:
		for _, ira := range s.ignoredRemoteAddresses {
			if ira.IP.Equal(addr.IP) && ira.Port == addr.Port {
				return true
			}
		}
	}
	return false
}

// correct config after handover, based on own transport parameters
func correctConfig(conf *Config, ownParams wire.TransportParameters) {
	if conf == nil {
		return
	}

	conf.MaxStreamReceiveWindow = utils.MaxUint64V(
		conf.MaxStreamReceiveWindow,
		uint64(ownParams.InitialMaxStreamDataBidiLocal),
		uint64(ownParams.InitialMaxStreamDataBidiRemote),
		uint64(ownParams.InitialMaxStreamDataUni),
	)

	conf.MaxConnectionReceiveWindow = utils.MaxUint64V(
		conf.MaxConnectionReceiveWindow,
		conf.MaxStreamReceiveWindow,
		uint64(ownParams.InitialMaxData),
	)
}

// Restore session from H-QUIC state
func Restore(state handover.State, perspective protocol.Perspective, conf *Config) (Connection, error) {
	if conf == nil {
		conf = &Config{}
	}

	logger := utils.DefaultLogger.WithPrefix(conf.LoggerPrefix)

	pconn, err := ListenMigratableUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0}) //TODO make param
	if err != nil {
		return nil, err
	}
	remoteAddr := state.RemoteAddress(perspective)

	var conn sendConn
	if perspective == protocol.PerspectiveClient {
		conn = newSendPconn(pconn, remoteAddr)
	} else {
		sconn, err := wrapConn(pconn)
		if err != nil {
			return nil, err
		}
		conn = newSendConn(sconn, remoteAddr, nil)
	}

	// must be called before populateClientConfig or populateServerConfig
	conf.ConnectionIDLength = state.SrcConnectionIDLength(perspective)

	if perspective == protocol.PerspectiveClient {
		conf = populateClientConfig(conf, conf.ConnectionIDLength == 0)
	} else {
		conf = populateServerConfig(conf)
	}

	correctConfig(conf, state.OwnTransportParameters(perspective))

	// This should stay empty because this is no longer required after the Handshake
	tlsConf := &tls.Config{}

	tracingID := nextConnTracingID()
	var connTracer logging.ConnectionTracer
	if conf.Tracer != nil {
		connTracer = conf.Tracer.TracerForConnection(
			context.WithValue(context.Background(), ConnectionTracingKey, tracingID),
			perspective,
			protocol.ParseConnectionID(state.OriginalDestinationConnectionID), //TODO maybe add field with original id
		)
	}

	runner, err := getMultiplexer().AddConn(pconn, state.SrcConnectionIDLength(perspective), conf.StatelessResetKey, conf.Tracer)
	if err != nil {
		return nil, err
	}

	var tokenGenerator *handshake.TokenGenerator
	if perspective == protocol.PerspectiveServer {
		tokenGenerator, err = handshake.NewTokenGenerator(rand.Reader)
		if err != nil {
			return nil, err
		}
	}

	s := &connection{
		conn:                   conn,
		config:                 conf,
		srcConnIDLen:           state.SrcConnectionIDLength(perspective),
		perspective:            perspective,
		handshakeCompleteChan:  make(chan struct{}),
		handshakeConfirmedChan: make(chan struct{}),
		logID:                  protocol.ParseConnectionID(state.OriginalDestinationConnectionID).String(),
		logger:                 logger,
		tracer:                 connTracer,
		versionNegotiated:      false,
		version:                state.Version,
		runner:                 runner,
		tokenGenerator:         tokenGenerator,
		cloned:                 true,
	}

	s.connIDManager = newConnIDManager(
		*state.MinActiveDestConnectionID(perspective),
		func(token protocol.StatelessResetToken) { runner.AddResetToken(token, s) },
		runner.RemoveResetToken,
		s.queueControlFrame,
	)

	s.connIDGenerator = newConnIDGenerator(
		state.MinActiveSrcConnectionID(perspective),
		state.MinActiveDestConnectionID(perspective),
		func(connID protocol.ConnectionID) { runner.Add(connID, s) },
		runner.GetStatelessResetToken,
		runner.Remove,
		runner.Retire,
		runner.ReplaceWithClosed,
		s.queueControlFrame,
		s.config.ConnectionIDGenerator,
		s.version,
	)

	// this has to happen before the server calls applyTransportParameters
	{
		maxSN, _ := state.MaxActiveSrcConnectionID(perspective)
		activeScrConnIDMap := make(map[uint64]protocol.ConnectionID, 0)
		for _, activeConnID := range state.ActiveSrcConnectionIDs(perspective) {
			activeScrConnIDMap[activeConnID.SequenceNumber] = protocol.ParseConnectionID(activeConnID.ConnectionID)
		}
		s.connIDGenerator.activeSrcConnIDs = activeScrConnIDMap
		s.connIDGenerator.highestSeq = maxSN
	}

	s.ctx, s.ctxCancel = context.WithCancel(context.WithValue(context.Background(), ConnectionTracingKey, tracingID))
	s.preSetup()
	s.sentPacketHandler, s.receivedPacketHandler = ackhandler.NewAckHandler(
		s.pathManager.CurrentSendPath(),
		0,
		getMaxPacketSize(s.conn.RemoteAddr()),
		s.config.InitialCongestionWindow,
		s.config.MinCongestionWindow,
		s.config.MaxCongestionWindow,
		s.config.InitialSlowStartThreshold,
		s.config.MinSlowStartThreshold,
		s.config.MaxSlowStartThreshold,
		s.rttStats,
		true, // TODO path challenge
		s.perspective,
		s.config.HyblaWestwoodCongestionControl,
		s.config.FixedPTO,
		s.tracer,
		s.logger,
		s.version,
	)

	// skip some packets for two reasons:
	//  - this number might be an estimate from the opposite perspective
	//  - some packets might be sent during the handshake
	s.sentPacketHandler.SetHighest1RTTPacketNumber(state.HighestSentPacketNumber(perspective) + 10000)

	initialStream := newCryptoStream()
	handshakeStream := newCryptoStream()
	ownParams := state.OwnTransportParameters(perspective)

	if s.config.EnableDatagrams {
		ownParams.MaxDatagramFrameSize = protocol.MaxDatagramFrameSize
	}
	if s.tracer != nil {
		s.tracer.SentTransportParameters(&ownParams)
	}

	cs, err := handshake.RestoreCryptoSetupFromHandoverState(state, conn.LocalAddr(), perspective, connTracer, logger, s.rttStats, tlsConf)
	if err != nil {
		return nil, err
	}

	s.cryptoStreamHandler = cs
	s.cryptoStreamManager = newCryptoStreamManager(cs, initialStream, handshakeStream, newCryptoStream())
	s.unpacker = newPacketUnpacker(cs, s.srcConnIDLen, s.version)
	s.packer = newPacketPacker(
		state.MinActiveSrcConnectionID(perspective),
		s.connIDManager.Get,
		initialStream,
		handshakeStream,
		s.sentPacketHandler,
		s.retransmissionQueue,
		s.RemoteAddr(),
		cs,
		s.framer,
		s.receivedPacketHandler,
		s.datagramQueue,
		s.perspective,
		s.version,
	)
	if len(tlsConf.ServerName) > 0 {
		s.tokenStoreKey = tlsConf.ServerName
	} else {
		s.tokenStoreKey = conn.RemoteAddr().String()
	}
	if s.config.TokenStore != nil {
		if token := s.config.TokenStore.Pop(s.tokenStoreKey); token != nil {
			s.packer.SetToken(token.data)
		}
	}

	for _, activeConnID := range state.ActiveSrcConnectionIDs(perspective) {
		runner.Add(protocol.ParseConnectionID(activeConnID.ConnectionID), s)
		//TODO restore stateless reset tokens
	}

	peerParams := state.PeerTransportParameters(perspective)

	s.peerParams = &peerParams

	if s.tracer != nil {
		s.tracer.RestoredTransportParameters(&peerParams)
	}

	for _, activeConnID := range state.ActiveDestConnectionIDs(perspective) {
		if protocol.ParseConnectionID(activeConnID.ConnectionID).Equal(s.connIDManager.activeConnectionID) {
			s.connIDManager.activeSequenceNumber = activeConnID.SequenceNumber
			s.connIDManager.activeStatelessResetToken = new(protocol.StatelessResetToken)
			if activeConnID.StatelessResetToken != nil {
				s.connIDManager.packetsPerConnectionID = protocol.PacketsPerConnectionID/2 + uint32(s.connIDManager.rand.Int31n(protocol.PacketsPerConnectionID))
				copy(s.connIDManager.activeStatelessResetToken[:], activeConnID.StatelessResetToken[:16])
			}
			continue
		}
		var resetToken [16]byte
		copy(resetToken[:], activeConnID.StatelessResetToken[:16])
		err := s.connIDManager.Add(&wire.NewConnectionIDFrame{
			SequenceNumber:      activeConnID.SequenceNumber,
			ConnectionID:        protocol.ParseConnectionID(activeConnID.ConnectionID),
			StatelessResetToken: resetToken,
		})
		if err != nil {
			return nil, err
		}
	}

	if perspective == protocol.PerspectiveClient {
		s.handleHandshakeComplete()
		s.handleHandshakeConfirmed()
	} else {
		s.handshakeComplete = true
		s.handshakeCompleteChan = nil
		s.undecryptablePackets = nil
		s.connIDManager.SetHandshakeComplete()
		s.connIDGenerator.SetHandshakeComplete()
		s.handleHandshakeConfirmed()
		s.sentPacketHandler.SetHighest1RTTPacketNumber(0) //TODO use packet number from state
		s.applyTransportParameters()
		// On the server side, the early session is ready as soon as we processed
		// the client's transport parameters.
		close(s.earlyConnReadyChan)
	}

	s.receivedFirstPacket = true

	//TODO send PATH_CHALLENGE instead of PING
	err = s.sendProbePacket(protocol.Encryption1RTT)
	if err != nil {
		return nil, err
	}

	go func() {
		_ = s.run() // returns as soon as the session is closed
	}()
	return s, nil
}

func (s *connection) useProxy() error {
	if !s.config.EnableActiveMigration {
		return errors.New("active migration has to be enabled")
	}
	if s.peerParams.DisableActiveMigration {
		return errors.New("active migration has to be enabled by peer")
	}

	proxyControlSession, err := DialProxyAddr(s.config.ProxyConf.Addr,
		s.config.ProxyConf.TlsConf,
		s.config.ProxyConf.Config,
	)
	if err != nil {
		return err
	}
	state, err := s.handover(false, true)
	if err != nil {
		return err
	}
	if s.config.ProxyConf.ModifyState != nil {
		s.config.ProxyConf.ModifyState(&state)
	}

	marshalledState, err := json.Marshal(state)
	if err != nil {
		return err
	}
	if s.logger.Debug() {
		s.logger.Debugf("created handover state: %s", string(marshalledState))
	}

	err = proxyControlSession.SendHandover(&state)
	if err != nil {
		s.closeLocal(fmt.Errorf("failed to send handover state: %v", err))
	}

	return nil
}

func (s *connection) ExtraStreamEncrypted() bool {
	return s.config.ExtraStreamEncryption.enabled() && s.peerParams.ExtraStreamEncryption
}

// does not block
func (s *connection) isWaitingForHandoverMigration() bool {
	select {
	case _, _ = <-s.handoverMigrationCtx.Done():
		return false
	default:
		return true
	}
}

func (s *connection) QueueHandshakeDoneFrame() error {
	if s.perspective == PerspectiveClient {
		return fmt.Errorf("client cannot send HANDSHAKE_DONE frame")
	}
	s.queueControlFrame(&wire.HandshakeDoneFrame{})
	return nil
}

func (s *connection) OriginalDestinationConnectionID() ConnectionID {
	logID, err := protocol.ParseConnectionIDHex(s.logID)
	if err != nil {
		panic(fmt.Errorf("failed to parse odcid: %v", err))
	}
	return logID
}
