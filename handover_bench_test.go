package quic

import (
	"context"
	"crypto/tls"
	"github.com/quic-go/quic-go/handover"
	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/testdata"
	"github.com/stretchr/testify/require"
	"io"
	"net"
	"testing"
	"time"
)

func clientAndServer(t require.TestingT) (transport *Transport, newClient func() (Connection, error), serverAccept func() (Connection, error)) {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	require.NoError(t, err)
	transport = &Transport{
		Conn:               conn,
		ConnectionIDLength: 20,
	}
	newClient = func() (Connection, error) {
		return transport.Dial(context.Background(), conn.LocalAddr(),
			&tls.Config{
				NextProtos:         []string{"test"},
				InsecureSkipVerify: true,
			}, &Config{
				MaxIdleTimeout: 100 * time.Millisecond,
			})
	}
	serverAccept = func() (Connection, error) {
		tlsConf := testdata.GetTLSConfig()
		tlsConf.NextProtos = []string{"test"}
		listener, err := transport.Listen(tlsConf, &Config{
			MaxIdleTimeout: 100 * time.Millisecond,
		})
		if err != nil {
			return nil, err
		}
		server, err := listener.Accept(context.Background())
		if err != nil {
			return nil, err
		}
		return server, nil
	}
	return transport, newClient, serverAccept
}

func runClient(b *testing.B, newClient func() (Connection, error), clientCtxCancel context.CancelFunc, msgLen int) {
	client, err := newClient()
	require.NoError(b, err)
	stream, err := client.OpenStream()
	require.NoError(b, err)
	stream.Write([]byte{1})
	msg, err := io.ReadAll(stream)
	require.NoError(b, err)
	require.Equal(b, len(msg), msgLen)
	for i := 0; i < len(msg); i++ {
		require.Equal(b, byte(i), msg[i])
	}
	clientCtxCancel()
}

func BenchmarkHandover(b *testing.B) {
	transport, newClient, serverAccept := clientAndServer(b)
	clientCtx, clientCtxCancel := context.WithCancel(context.Background())
	_, serverCtxCancel := context.WithCancel(context.Background())

	go runClient(b, newClient, clientCtxCancel, b.N+1)

	server, err := serverAccept()
	require.NoError(b, err)
	stream, err := server.AcceptStream(context.Background())
	require.NoError(b, err)
	stream.Write([]byte{0})

	b.ResetTimer()
	for n := 1; n <= b.N; n++ {
		resp := server.Handover(true, &handover.ConnectionStateStoreConf{IncludePendingOutgoingFrames: true})
		require.NoError(b, resp.Error)
		serializedState, err := resp.State.SerializeMsgp()
		require.NoError(b, err)
		state, err := (&handover.State{}).ParseMsgp(serializedState)
		require.NoError(b, err)
		newServer, streams, err := Restore(transport, state, &ConnectionRestoreConfig{
			Perspective: protocol.PerspectiveServer,
		})
		server = newServer
		require.NoError(b, err)
		require.Equal(b, 1, len(streams.BidiStreams))
		stream = streams.BidiStreams[0]
		stream.Write([]byte{byte(n)})
	}
	stream.Close()
	select {
	case <-clientCtx.Done():
	case <-time.After(2 * time.Second):
		require.Fail(b, "timeout")
	}
	b.StopTimer()
	serverCtxCancel()
	err = transport.Close()
	require.NoError(b, err)
}
