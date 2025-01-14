package quic

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/lucas-clemente/quic-go/handover"
	"io"
	"net"
	"time"

	"github.com/lucas-clemente/quic-go/internal/handshake"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/logging"
)

// The StreamID is the ID of a QUIC stream.
type StreamID = protocol.StreamID

// A VersionNumber is a QUIC version number.
type VersionNumber = protocol.VersionNumber

type ByteCount = protocol.ByteCount

type Perspective = protocol.Perspective

const PerspectiveClient = protocol.PerspectiveClient
const PerspectiveServer = protocol.PerspectiveServer

// ConnectionFlowControlMultiplier determines how much larger the connection flow control windows needs to be relative to any stream's flow control window
// This is the value that Chromium is using
const ConnectionFlowControlMultiplier = protocol.ConnectionFlowControlMultiplier

// DefaultInitialCongestionWindow is the default initial congestion window in number of packets.
const DefaultInitialCongestionWindow = protocol.DefaultInitialCongestionWindow

// DefaultMinCongestionWindow is the default minimum congestion window in number of packets.
const DefaultMinCongestionWindow = protocol.DefaultMinCongestionWindow

// DefaultMaxCongestionWindow is the maximum congestion window in packet.
const DefaultMaxCongestionWindow = protocol.DefaultMaxCongestionWindow

const DefaultInitialSlowStartThreshold = protocol.DefaultInitialSlowStartThreshold
const DefaultMinSlowStartThreshold = protocol.DefaultMinSlowStartThreshold
const DefaultMaxSlowStartThreshold = protocol.DefaultMaxSlowStartThreshold

const (
	// VersionDraft29 is IETF QUIC draft-29
	VersionDraft29 = protocol.VersionDraft29
	// Version1 is RFC 9000
	Version1 = protocol.Version1
	Version2 = protocol.Version2
)

type ExtraStreamEncryptionMode uint

const (
	// DisableExtraStreamEncryption disables the use of the XSE-QUIC extension (default)
	DisableExtraStreamEncryption ExtraStreamEncryptionMode = 0
	// PreferExtraStreamEncryption enables XSE-QUIC only when the peer supports it
	PreferExtraStreamEncryption = 1
	// EnforceExtraStreamEncryption enables XSE-QUIC, but terminates the connection with a TRANSPORT_PARAMETER_ERROR when the peer does not support it
	EnforceExtraStreamEncryption = 2
)

func (a ExtraStreamEncryptionMode) enabled() bool {
	return a != DisableExtraStreamEncryption
}

// A ClientToken is a token received by the client.
// It can be used to skip address validation on future connection attempts.
type ClientToken struct {
	data []byte
}

type TokenStore interface {
	// Pop searches for a ClientToken associated with the given key.
	// Since tokens are not supposed to be reused, it must remove the token from the cache.
	// It returns nil when no token is found.
	Pop(key string) (token *ClientToken)

	// Put adds a token to the cache with the given key. It might get called
	// multiple times in a connection.
	Put(key string, token *ClientToken)
}

// Err0RTTRejected is the returned from:
// * Open{Uni}Stream{Sync}
// * Accept{Uni}Stream
// * Stream.Read and Stream.Write
// when the server rejects a 0-RTT connection attempt.
var Err0RTTRejected = errors.New("0-RTT rejected")

// ConnectionTracingKey can be used to associate a ConnectionTracer with a Connection.
// It is set on the Connection.Context() context,
// as well as on the context passed to logging.Tracer.NewConnectionTracer.
var ConnectionTracingKey = connTracingCtxKey{}

type connTracingCtxKey struct{}

// Stream is the interface implemented by QUIC streams
// In addition to the errors listed on the Connection,
// calls to stream functions can return a StreamError if the stream is canceled.
type Stream interface {
	ReceiveStream
	SendStream
	// SetDeadline sets the read and write deadlines associated
	// with the connection. It is equivalent to calling both
	// SetReadDeadline and SetWriteDeadline.
	SetDeadline(t time.Time) error
}

// A ReceiveStream is a unidirectional Receive Stream.
type ReceiveStream interface {
	// StreamID returns the stream ID.
	StreamID() StreamID
	// Read reads data from the stream.
	// Read can be made to time out and return a net.Error with Timeout() == true
	// after a fixed time limit; see SetDeadline and SetReadDeadline.
	// If the stream was canceled by the peer, the error implements the StreamError
	// interface, and Canceled() == true.
	// If the connection was closed due to a timeout, the error satisfies
	// the net.Error interface, and Timeout() will be true.
	io.Reader
	// CancelRead aborts receiving on this stream.
	// It will ask the peer to stop transmitting stream data.
	// Read will unblock immediately, and future Read calls will fail.
	// When called multiple times or after reading the io.EOF it is a no-op.
	CancelRead(StreamErrorCode)
	// SetReadDeadline sets the deadline for future Read calls and
	// any currently-blocked Read call.
	// A zero value for t means Read will not time out.

	SetReadDeadline(t time.Time) error
	ReadOffset() ByteCount
}

// A SendStream is a unidirectional Send Stream.
type SendStream interface {
	// StreamID returns the stream ID.
	StreamID() StreamID
	// Write writes data to the stream.
	// Write can be made to time out and return a net.Error with Timeout() == true
	// after a fixed time limit; see SetDeadline and SetWriteDeadline.
	// If the stream was canceled by the peer, the error implements the StreamError
	// interface, and Canceled() == true.
	// If the connection was closed due to a timeout, the error satisfies
	// the net.Error interface, and Timeout() will be true.
	io.Writer
	// Close closes the write-direction of the stream.
	// Future calls to Write are not permitted after calling Close.
	// It must not be called concurrently with Write.
	// It must not be called after calling CancelWrite.
	io.Closer
	// CancelWrite aborts sending on this stream.
	// Data already written, but not yet delivered to the peer is not guaranteed to be delivered reliably.
	// Write will unblock immediately, and future calls to Write will fail.
	// When called multiple times or after closing the stream it is a no-op.
	CancelWrite(StreamErrorCode)
	// The Context is canceled as soon as the write-side of the stream is closed.
	// This happens when Close() or CancelWrite() is called, or when the peer
	// cancels the read-side of their stream.
	Context() context.Context
	// SetWriteDeadline sets the deadline for future Write calls
	// and any currently-blocked Write call.
	// Even if write times out, it may return n > 0, indicating that
	// some data was successfully written.
	// A zero value for t means Write will not time out.
	SetWriteDeadline(t time.Time) error
	WriteOffset() ByteCount
}

// A Connection is a QUIC connection between two peers.
// Calls to the connection (and to streams) can return the following types of errors:
// * ApplicationError: for errors triggered by the application running on top of QUIC
// * TransportError: for errors triggered by the QUIC transport (in many cases a misbehaving peer)
// * IdleTimeoutError: when the peer goes away unexpectedly (this is a net.Error timeout error)
// * HandshakeTimeoutError: when the cryptographic handshake takes too long (this is a net.Error timeout error)
// * StatelessResetError: when we receive a stateless reset (this is a net.Error temporary error)
// * VersionNegotiationError: returned by the client, when there's no version overlap between the peers
type Connection interface {
	// AcceptStream returns the next stream opened by the peer, blocking until one is available.
	// If the connection was closed due to a timeout, the error satisfies
	// the net.Error interface, and Timeout() will be true.
	AcceptStream(context.Context) (Stream, error)
	// AcceptUniStream returns the next unidirectional stream opened by the peer, blocking until one is available.
	// If the connection was closed due to a timeout, the error satisfies
	// the net.Error interface, and Timeout() will be true.
	AcceptUniStream(context.Context) (ReceiveStream, error)
	// OpenStream opens a new bidirectional QUIC stream.
	// There is no signaling to the peer about new streams:
	// The peer can only accept the stream after data has been sent on the stream.
	// If the error is non-nil, it satisfies the net.Error interface.
	// When reaching the peer's stream limit, err.Temporary() will be true.
	// If the connection was closed due to a timeout, Timeout() will be true.
	OpenStream() (Stream, error)
	// OpenStreamSync opens a new bidirectional QUIC stream.
	// It blocks until a new stream can be opened.
	// If the error is non-nil, it satisfies the net.Error interface.
	// If the connection was closed due to a timeout, Timeout() will be true.
	OpenStreamSync(context.Context) (Stream, error)
	// OpenUniStream opens a new outgoing unidirectional QUIC stream.
	// If the error is non-nil, it satisfies the net.Error interface.
	// When reaching the peer's stream limit, Temporary() will be true.
	// If the connection was closed due to a timeout, Timeout() will be true.
	OpenUniStream() (SendStream, error)
	// OpenUniStreamSync opens a new outgoing unidirectional QUIC stream.
	// It blocks until a new stream can be opened.
	// If the error is non-nil, it satisfies the net.Error interface.
	// If the connection was closed due to a timeout, Timeout() will be true.
	OpenUniStreamSync(context.Context) (SendStream, error)
	// LocalAddr returns the local address.
	LocalAddr() net.Addr
	// RemoteAddr returns the address of the peer.
	RemoteAddr() net.Addr
	// CloseWithError closes the connection with an error.
	// The error string will be sent to the peer.
	CloseWithError(ApplicationErrorCode, string) error
	// Context returns a context that is cancelled when the connection is closed.
	Context() context.Context
	// ConnectionState returns basic details about the QUIC connection.
	// It blocks until the handshake completes.
	// Warning: This API should not be considered stable and might change soon.
	ConnectionState() ConnectionState

	// SendMessage sends a message as a datagram, as specified in RFC 9221.
	SendMessage([]byte) error
	SendMessageWithoutWaitForDequeue([]byte) error
	// ReceiveMessage gets a message received in a datagram, as specified in RFC 9221.
	ReceiveMessage() ([]byte, error)
	// Handover creates H-QUIC state.
	// Session is silently destroyed when destroy is set.
	// Session no longer sends and ignores incoming packets from the current path when ignoreCurrentPath is set.
	Handover(destroy bool, config *ConnectionStateStoreConf) HandoverStateResponse
	// MigrateUDPSocket migrates connection to a new UDP socket.
	// Returns new UDP address.
	MigrateUDPSocket() (*net.UDPAddr, error)
	// ExtraStreamEncrypted returns true if XSE-QUIC is used
	ExtraStreamEncrypted() bool
	// QueueHandshakeDoneFrame is required by H-QUIC proxy
	QueueHandshakeDoneFrame() error
	// OriginalDestinationConnectionID is used for connection identification e.g. for qlog.
	// Remains unchanged over the whole connection time
	OriginalDestinationConnectionID() ConnectionID
	AwaitPathUpdate() <-chan struct{}
	// get already opened stream.
	// returns error if stream is not yet opened
	OpenedBidiStream(id StreamID) (Stream, error)
	AddProxy(conf *ProxyConfig) ProxySetupResponse
	UpdateRemoteAddr(addr net.UDPAddr, ignoreReceivedPacketsFromCurrentPath bool, ignoreMigrationToCurrentPath bool) error
	QlogWriter() logging.QlogWriter
}

// An EarlyConnection is a connection that is handshaking.
// Data sent during the handshake is encrypted using the forward secure keys.
// When using client certificates, the client's identity is only verified
// after completion of the handshake.
type EarlyConnection interface {
	Connection

	// HandshakeComplete blocks until the handshake completes (or fails).
	// Data sent before completion of the handshake is encrypted with 1-RTT keys.
	// Note that the client's identity hasn't been verified yet.
	HandshakeComplete() context.Context

	NextConnection() Connection
}

// StatelessResetKey is a key used to derive stateless reset tokens.
type StatelessResetKey [32]byte

// A ConnectionID is a QUIC Connection ID, as defined in RFC 9000.
// It is not able to handle QUIC Connection IDs longer than 20 bytes,
// as they are allowed by RFC 8999.
type ConnectionID = protocol.ConnectionID

// ConnectionIDFromBytes interprets b as a Connection ID. It panics if b is
// longer than 20 bytes.
func ConnectionIDFromBytes(b []byte) ConnectionID {
	return protocol.ParseConnectionID(b)
}

// A ConnectionIDGenerator is an interface that allows clients to implement their own format
// for the Connection IDs that servers/clients use as SrcConnectionID in QUIC packets.
//
// Connection IDs generated by an implementation should always produce IDs of constant size.
type ConnectionIDGenerator interface {
	// GenerateConnectionID generates a new ConnectionID.
	// Generated ConnectionIDs should be unique and observers should not be able to correlate two ConnectionIDs.
	GenerateConnectionID() (ConnectionID, error)

	// ConnectionIDLen tells what is the length of the ConnectionIDs generated by the implementation of
	// this interface.
	// Effectively, this means that implementations of ConnectionIDGenerator must always return constant-size
	// connection IDs. Valid lengths are between 0 and 20 and calls to GenerateConnectionID.
	// 0-length ConnectionsIDs can be used when an endpoint (server or client) does not require multiplexing connections
	// in the presence of a connection migration environment.
	ConnectionIDLen() int
}

// Config contains all configuration data needed for a QUIC server or client.
type Config struct {
	// The QUIC versions that can be negotiated.
	// If not set, it uses all versions available.
	Versions []VersionNumber
	// The length of the connection ID in bytes.
	// It can be 0, or any value between 4 and 18.
	// If not set, the interpretation depends on where the Config is used:
	// If used for dialing an address, a 0 byte connection ID will be used.
	// If used for a server, or dialing on a packet conn, a 4 byte connection ID will be used.
	// When dialing on a packet conn, the ConnectionIDLength value must be the same for every Dial call.
	ConnectionIDLength int
	// An optional ConnectionIDGenerator to be used for ConnectionIDs generated during the lifecycle of a QUIC connection.
	// The goal is to give some control on how connection IDs, which can be useful in some scenarios, in particular for servers.
	// By default, if not provided, random connection IDs with the length given by ConnectionIDLength is used.
	// Otherwise, if one is provided, then ConnectionIDLength is ignored.
	ConnectionIDGenerator ConnectionIDGenerator
	// HandshakeIdleTimeout is the idle timeout before completion of the handshake.
	// Specifically, if we don't receive any packet from the peer within this time, the connection attempt is aborted.
	// If this value is zero, the timeout is set to 5 seconds.
	HandshakeIdleTimeout time.Duration
	// MaxIdleTimeout is the maximum duration that may pass without any incoming network activity.
	// The actual value for the idle timeout is the minimum of this value and the peer's.
	// This value only applies after the handshake has completed.
	// If the timeout is exceeded, the connection is closed.
	// If this value is zero, the timeout is set to 30 seconds.
	MaxIdleTimeout time.Duration
	// RequireAddressValidation determines if a QUIC Retry packet is sent.
	// This allows the server to verify the client's address, at the cost of increasing the handshake latency by 1 RTT.
	// See https://datatracker.ietf.org/doc/html/rfc9000#section-8 for details.
	// If not set, every client is forced to prove its remote address.
	RequireAddressValidation func(net.Addr) bool
	// MaxRetryTokenAge is the maximum age of a Retry token.
	// If not set, it defaults to 5 seconds. Only valid for a server.
	MaxRetryTokenAge time.Duration
	// MaxTokenAge is the maximum age of the token presented during the handshake,
	// for tokens that were issued on a previous connection.
	// If not set, it defaults to 24 hours. Only valid for a server.
	MaxTokenAge time.Duration
	// The TokenStore stores tokens received from the server.
	// Tokens are used to skip address validation on future connection attempts.
	// The key used to store tokens is the ServerName from the tls.Config, if set
	// otherwise the token is associated with the server's IP address.
	TokenStore TokenStore
	// InitialStreamReceiveWindow is the initial size of the stream-level flow control window for receiving data.
	// If the application is consuming data quickly enough, the flow control auto-tuning algorithm
	// will increase the window up to MaxStreamReceiveWindow.
	// If this value is zero, it will default to 512 KB.
	InitialStreamReceiveWindow uint64
	// MaxStreamReceiveWindow is the maximum stream-level flow control window for receiving data.
	// If this value is zero, it will default to 6 MB.
	MaxStreamReceiveWindow uint64
	// InitialConnectionReceiveWindow is the initial size of the stream-level flow control window for receiving data.
	// If the application is consuming data quickly enough, the flow control auto-tuning algorithm
	// will increase the window up to MaxConnectionReceiveWindow.
	// If this value is zero, it will default to 512 KB.
	InitialConnectionReceiveWindow uint64
	// MaxConnectionReceiveWindow is the connection-level flow control window for receiving data.
	// If this value is zero, it will default to 15 MB.
	MaxConnectionReceiveWindow uint64
	// AllowConnectionWindowIncrease is called every time the connection flow controller attempts
	// to increase the connection flow control window.
	// If set, the caller can prevent an increase of the window. Typically, it would do so to
	// limit the memory usage.
	// To avoid deadlocks, it is not valid to call other functions on the connection or on streams
	// in this callback.
	AllowConnectionWindowIncrease func(conn Connection, delta uint64) bool
	// MaxIncomingStreams is the maximum number of concurrent bidirectional streams that a peer is allowed to open.
	// Values above 2^60 are invalid.
	// If not set, it will default to 100.
	// If set to a negative value, it doesn't allow any bidirectional streams.
	MaxIncomingStreams int64
	// MaxIncomingUniStreams is the maximum number of concurrent unidirectional streams that a peer is allowed to open.
	// Values above 2^60 are invalid.
	// If not set, it will default to 100.
	// If set to a negative value, it doesn't allow any unidirectional streams.
	MaxIncomingUniStreams int64
	// The StatelessResetKey is used to generate stateless reset tokens.
	// If no key is configured, sending of stateless resets is disabled.
	StatelessResetKey *StatelessResetKey
	// KeepAlivePeriod defines whether this peer will periodically send a packet to keep the connection alive.
	// If set to 0, then no keep alive is sent. Otherwise, the keep alive is sent on that period (or at most
	// every half of MaxIdleTimeout, whichever is smaller).
	KeepAlivePeriod time.Duration
	// DisablePathMTUDiscovery disables Path MTU Discovery (RFC 8899).
	// Packets will then be at most 1252 (IPv4) / 1232 (IPv6) bytes in size.
	// Note that if Path MTU discovery is causing issues on your system, please open a new issue
	DisablePathMTUDiscovery bool
	// DisableVersionNegotiationPackets disables the sending of Version Negotiation packets.
	// This can be useful if version information is exchanged out-of-band.
	// It has no effect for a client.
	DisableVersionNegotiationPackets bool
	// Enable QUIC datagram support (RFC 9221).
	EnableDatagrams bool
	Tracer          logging.Tracer
	// IgnoreReceived1RTTPacketsUntilFirstPathMigration pauses the session directly after the handshake
	// this is e.g. used for connection handover
	// TODO this is no longer necessary when handover is supported at an arbitrary state
	IgnoreReceived1RTTPacketsUntilFirstPathMigration bool
	// LoggerPrefix add prefix to every log line.
	// if nil, "client" or "server" are used as prefix
	LoggerPrefix string
	// EnableActiveMigration oppositely sets the disable_active_migration transport parameter.
	// If not set, it will default to false.
	EnableActiveMigration bool
	// The Proxy to use
	// if nil, no proxy is used
	ProxyConf *ProxyConfig
	// The InitialCongestionWindow to use, in number of packets
	InitialCongestionWindow uint32
	// in number of packets
	MinCongestionWindow uint32
	// in number of packets
	MaxCongestionWindow uint32
	// Initial slow start congestion threshold in bytes, aka initial ssthresh.
	InitialSlowStartThreshold ByteCount
	MinSlowStartThreshold     ByteCount
	MaxSlowStartThreshold     ByteCount
	// Use XSE-QUIC extension.
	// Derive an additional key in the handshake,
	// to additionally encrypt the stream payload,
	// before the QUIC Packet is encrypted.
	ExtraStreamEncryption          ExtraStreamEncryptionMode
	HyblaWestwoodCongestionControl bool
	// Allows H-QUIC state creation when handshake is completed but not yet confirmed.
	// Only affects client behaviour.
	// Optimistic approach! Success is not guaranteed due to race conditions.
	AllowEarlyHandover bool
	// if not 0 use this PTO instead of the calculated PTO
	FixedPTO time.Duration
	// Handler for short header packets with an unknown connection id
	HandleUnknownConnectionPacket func(ConnectionID, *receivedPacket)
}

// ConnectionState records basic details about a QUIC connection
type ConnectionState struct {
	TLS               handshake.ConnectionState
	SupportsDatagrams bool
	Version           VersionNumber
}

// A Listener for incoming QUIC connections
type Listener interface {
	// Close the server. All active connections will be closed.
	Close() error
	// Addr returns the local network addr that the server is listening on.
	Addr() net.Addr
	// Accept returns new connections. It should be called in a loop.
	Accept(context.Context) (Connection, error)
	// MigrateUDPSocket migrates connection to a new UDP socket.
	// Returns new UDP address.
	MigrateUDPSocket() (*net.UDPAddr, error)
}

// An EarlyListener listens for incoming QUIC connections,
// and returns them before the handshake completes.
type EarlyListener interface {
	// Close the server. All active connections will be closed.
	Close() error
	// Addr returns the local network addr that the server is listening on.
	Addr() net.Addr
	// Accept returns new early connections. It should be called in a loop.
	Accept(context.Context) (EarlyConnection, error)
	// MigrateUDPSocket migrates connection to a new UDP socket.
	// Returns new UDP address.
	MigrateUDPSocket() (*net.UDPAddr, error)
	PacketHandlerManager() PacketHandlerManager
}

type ProxyConfig struct {
	// the proxy address to use
	Addr string
	// used for proxy control connection
	Config *Config
	// used for proxy control connection
	TlsConf *tls.Config
	// before sending the handover state to the H-QUIC proxy, this function can be used to modify the state
	ModifyState func(state *handover.State)
}

func (c *ProxyConfig) Clone() *ProxyConfig {
	return &ProxyConfig{
		Addr:        c.Addr,
		Config:      c.Config.Clone(),
		TlsConf:     c.TlsConf.Clone(),
		ModifyState: c.ModifyState,
	}
}

type ReceivedPacket = receivedPacket

type PacketHandlerManager interface {
	AddWithConnID(protocol.ConnectionID, protocol.ConnectionID, func() packetHandler) bool
	Destroy() error
	connRunner
	SetServer(unknownPacketHandler)
	CloseServer()
	SetUnknownConnectionHandler(func(ConnectionID, *receivedPacket))
	ConnIDLength() int
	GetConnectionByID(id protocol.ConnectionID) Connection
	PacketConn() net.PacketConn
}
