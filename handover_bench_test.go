package quic

import (
	"context"
	"crypto/tls"
	"github.com/quic-go/quic-go/handover"
	"github.com/quic-go/quic-go/internal/testdata"
	"github.com/quic-go/quic-go/qstate"
	"github.com/stretchr/testify/require"
	"gonum.org/v1/gonum/stat"
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
				MaxIdleTimeout: 10 * time.Second,
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

func BenchmarkHandoverContinuousStream(b *testing.B) {
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
		serializedState, err := resp.State.MarshalMsg(nil)
		require.NoError(b, err)
		state := qstate.Connection{}
		_, err = state.UnmarshalMsg(serializedState)
		require.NoError(b, err)
		newServer, streams, err := Restore(transport, &state, &ConnectionRestoreConfig{})
		server = newServer
		require.NoError(b, err)
		require.Equal(b, 1, len(streams.BidiStreams))
		stream = streams.BidiStreams[0]
		stream.Write([]byte{byte(n)})
	}
	stream.Close()
	select {
	case <-clientCtx.Done():
	case <-time.After(10 * time.Second):
		require.Fail(b, "timeout")
	}
	b.StopTimer()
	serverCtxCancel()
	err = transport.Close()
	require.NoError(b, err)
}

func BenchmarkHandoverFinOnly(b *testing.B) {
	var stateBuf [100_000]byte
	for codecName, codec := range qstate.Codecs {
		b.Run(codecName, func(b *testing.B) {
			transport, newClient, serverAccept := clientAndServer(b)
			clientCtx, clientCtxCancel := context.WithCancel(context.Background())
			_, serverCtxCancel := context.WithCancel(context.Background())

			go runClient(b, newClient, clientCtxCancel, 1)

			server, err := serverAccept()
			require.NoError(b, err)
			stream, err := server.AcceptStream(context.Background())
			require.NoError(b, err)
			stream.Write([]byte{0})

			storeDurations := make([]float64, 0, b.N)
			restoreDurations := make([]float64, 0, b.N)
			encodeSizes := make([]float64, 0, b.N)
			b.ResetTimer()
			for n := 1; n <= b.N; n++ {
				startStore := time.Now()
				resp := server.Handover(true, &handover.ConnectionStateStoreConf{IncludePendingOutgoingFrames: true})
				require.NoError(b, resp.Error)
				serializedState, err := codec.Encode(stateBuf[:0], &resp.State)
				require.NoError(b, err)
				encodeSizes = append(encodeSizes, float64(len(serializedState)))
				storeDurations = append(storeDurations, time.Since(startStore).Seconds())
				startRestore := time.Now()
				state := qstate.Connection{}
				err = codec.Decode(&state, serializedState)
				require.NoError(b, err)
				newServer, streams, err := Restore(transport, &state, &ConnectionRestoreConfig{})
				server = newServer
				require.NoError(b, err)
				require.Equal(b, 1, len(streams.BidiStreams))
				stream = streams.BidiStreams[0]
				require.NotNil(b, stream)
				restoreDurations = append(restoreDurations, time.Since(startRestore).Seconds())
			}
			stream.Close()
			select {
			case <-clientCtx.Done():
			case <-time.After(10 * time.Second):
				require.Fail(b, "timeout")
			}
			b.StopTimer()
			serverCtxCancel()
			err = transport.Close()
			require.NoError(b, err)
			size_mean, size_std_dev := stat.MeanStdDev(encodeSizes, nil)
			b.ReportMetric(size_mean, "mean_B")
			b.ReportMetric(size_std_dev, "std_dev_B")
			store_mean, store_std_dev := stat.MeanStdDev(storeDurations, nil)
			b.ReportMetric(store_mean, "mean_store_s")
			b.ReportMetric(store_std_dev, "std_dev_store_s")
			restore_mean, restore_std_dev := stat.MeanStdDev(restoreDurations, nil)
			b.ReportMetric(restore_mean, "mean_restore_s")
			b.ReportMetric(restore_std_dev, "std_dev_restore_s")
		})
	}
}
