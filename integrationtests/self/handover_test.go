package self

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/handover"
	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/testdata"
	"github.com/quic-go/quic-go/internal/utils"
	"github.com/quic-go/quic-go/qlog"
	"github.com/quic-go/quic-go/qstate"
	"github.com/stretchr/testify/require"
	"io"
	"net"
	"sync"
	"testing"
	"time"
)

func tlsConf() (client *tls.Config, server *tls.Config) {
	protos := []string{"proto1"}
	server = testdata.GetTLSConfig()
	server.NextProtos = protos
	certPool := x509.NewCertPool()
	testdata.AddRootCA(certPool)
	client = &tls.Config{
		RootCAs:            certPool,
		InsecureSkipVerify: true, // used because of missing IP SAN
	}
	client.NextProtos = protos
	return
}

func TestServerToServerHandoverMidStream(t *testing.T) {
	MaxIdleTimeout := 100 * time.Millisecond
	server1AddrChan := make(chan net.Addr, 1)
	server2AddrChan := make(chan net.Addr, 1)
	serverStateChan := make(chan qstate.Connection, 1)
	server1Ctx, server1CtxCancel := context.WithCancelCause(context.Background())
	server2Ctx, server2CtxCancel := context.WithCancelCause(context.Background())
	clientTlsConf, serverTlsConf := tlsConf()
	message1 := "hello"
	message2 := "hquic"

	go func() { // server 1
		server, err := quic.ListenAddr("127.0.0.1:0", serverTlsConf, &quic.Config{MaxIdleTimeout: MaxIdleTimeout})
		require.NoError(t, err)
		defer server.Close()
		server1AddrChan <- server.Addr()
		conn, err := server.Accept(context.Background())
		require.NoError(t, err)
		defer conn.Destroy()
		// transfer data
		err = acceptAndReceive(conn, message1, false)
		require.NoError(t, err)
		err = openAndSend(conn, message1, false)
		require.NoError(t, err)
		// state handover
		res := conn.Handover(true, &handover.ConnectionStateStoreConf{
			IncludePendingOutgoingFrames: true,
			IncludePendingIncomingFrames: true,
		})
		serverState := res.State
		err = res.Error
		require.NoError(t, err)
		serverStateChan <- serverState
		<-conn.Context().Done()
		server1CtxCancel(nil)
	}()
	go func() { // server 2
		defer GinkgoRecover()
		serverState := <-serverStateChan
		require.Equal(t, "server", serverState.Transport.VantagePoint)
		conn, restoredStreams, err := quic.Restore(nil, &serverState, &quic.ConnectionRestoreConfig{
			MaxIdleTimeout: MaxIdleTimeout,
		})
		bidiStreams := restoredStreams.BidiStreams
		sendStreams := restoredStreams.SendStreams
		receiveStreams := restoredStreams.ReceiveStreams
		require.NoError(t, err)
		require.Equal(t, 2, len(bidiStreams))
		require.Empty(t, sendStreams)
		require.Empty(t, receiveStreams)
		defer conn.Destroy()
		server2AddrChan <- conn.LocalAddr()
		// transfer
		err = receiveMessage(bidiStreams[0], message2, true)
		require.NoError(t, err)
		err = sendMessage(bidiStreams[1], message2, true)
		require.NoError(t, err)
		<-conn.Context().Done()
		server2CtxCancel(nil)
	}()

	originalServerAddr := <-server1AddrChan
	clientConn, err := quic.DialAddr(context.Background(), originalServerAddr.String(), clientTlsConf, &quic.Config{MaxIdleTimeout: MaxIdleTimeout})
	require.NoError(t, err)
	// transfer
	err = openAndSend(clientConn, message1+message2, true)
	require.NoError(t, err)
	err = acceptAndReceive(clientConn, message1+message2, true)
	require.NoError(t, err)
	migratedServerAddr := clientConn.RemoteAddr()
	require.NotEqual(t, originalServerAddr.String(), migratedServerAddr.String()) // check if migrated

	require.Eventually(t, func() bool {
		<-server1Ctx.Done()
		return true
	}, time.Second, 100*time.Millisecond)
	require.Eventually(t, func() bool {
		<-server2Ctx.Done()
		return true
	}, time.Second, 100*time.Millisecond)
	// destroy sessions
	clientConn.Destroy()
}

var _ = Describe("Handover", func() {
	var (
		message1      = "hello"
		idleTimeout   = 100 * time.Millisecond
		serverTlsConf *tls.Config
		clientTlsConf *tls.Config
	)

	BeforeEach(func() {
		//utils.DefaultLogger.SetLogLevel(utils.LogLevelDebug)
		//log.SetOutput(os.Stdout)
		protos := []string{"proto1"}
		serverTlsConf = testdata.GetTLSConfig()
		serverTlsConf.NextProtos = protos
		certPool := x509.NewCertPool()
		testdata.AddRootCA(certPool)
		clientTlsConf = &tls.Config{
			RootCAs:            certPool,
			InsecureSkipVerify: true, // used because of missing IP SAN
		}
		clientTlsConf.NextProtos = protos
	})

	AfterEach(func() {
		//Eventually(testutils.AreConnsRunning).Should(BeFalse())
		//Eventually(testutils.AreServersRunning).Should(BeFalse())
	})

	It("server to server handover", func() {
		MaxIdleTimeout := 100 * time.Millisecond
		server1AddrChan := make(chan net.Addr, 1)
		server2AddrChan := make(chan net.Addr, 1)
		serverStateChan := make(chan qstate.Connection, 1)
		server2DoneChan := make(chan struct{})

		go func() { // server 1
			defer GinkgoRecover()
			server, err := quic.ListenAddr("127.0.0.1:0", serverTlsConf, &quic.Config{
				MaxIdleTimeout: MaxIdleTimeout,
				Tracer:         qlog.DefaultTracerWithLabel("server"),
			})
			Expect(err).ToNot(HaveOccurred())
			defer server.Close()
			server1AddrChan <- server.Addr()
			conn, err := server.Accept(context.Background())
			Expect(err).ToNot(HaveOccurred())
			defer conn.Destroy()
			res := conn.Handover(true, &handover.ConnectionStateStoreConf{})
			serverState := res.State
			err = res.Error
			Expect(err).ToNot(HaveOccurred())
			serverStateChan <- serverState
			<-conn.Context().Done()
		}()
		go func() { // server 2
			defer GinkgoRecover()
			serverState := <-serverStateChan
			Expect(serverState.Transport.VantagePoint).To(Equal("server"))
			conn, _, err := quic.Restore(nil, &serverState, &quic.ConnectionRestoreConfig{
				MaxIdleTimeout: MaxIdleTimeout,
				Tracer:         qlog.DefaultTracerWithLabel("restored_server"),
			})
			Expect(err).ToNot(HaveOccurred())
			defer conn.Destroy()
			server2AddrChan <- conn.LocalAddr()
			// transfer
			err = openAndSend(conn, message1, true)
			Expect(err).ToNot(HaveOccurred())
			err = acceptAndReceive(conn, message1, true)
			Expect(err).ToNot(HaveOccurred())
			close(server2DoneChan)
		}()

		originalServerAddr := <-server1AddrChan
		clientConn, err := quic.DialAddr(context.Background(), originalServerAddr.String(), clientTlsConf, &quic.Config{
			MaxIdleTimeout: MaxIdleTimeout,
			Tracer:         qlog.DefaultTracerWithLabel("client"),
		})
		Expect(err).ToNot(HaveOccurred())
		// transfer
		err = acceptAndReceive(clientConn, message1, true)
		Expect(err).ToNot(HaveOccurred())
		err = openAndSend(clientConn, message1, true)
		Expect(err).ToNot(HaveOccurred())
		migratedServerAddr := clientConn.RemoteAddr()
		Expect(originalServerAddr.String()).ToNot(Equal(migratedServerAddr.String())) // check if migrated
		<-server2DoneChan
		clientConn.CloseWithError(0, "")
	})

	It("client to client handover", func() {
		server, err := quic.ListenAddr("127.0.0.1:0", serverTlsConf, &quic.Config{MaxIdleTimeout: 100 * time.Millisecond})
		var serverSession quic.Connection
		Expect(err).ToNot(HaveOccurred())
		go func() {
			defer GinkgoRecover()
			serverSession, err = server.Accept(context.Background())
			Expect(err).ToNot(HaveOccurred())
			originalClientAddr := serverSession.RemoteAddr()

			// transfer
			err := acceptAndReceive(serverSession, message1, true)
			Expect(err).ToNot(HaveOccurred())
			err = openAndSend(serverSession, message1, true)
			Expect(err).ToNot(HaveOccurred())
			<-serverSession.Context().Done()

			migratedClientAddr := serverSession.RemoteAddr()
			Expect(originalClientAddr.(*net.UDPAddr).Port).ToNot(Equal(migratedClientAddr.(*net.UDPAddr).Port)) // check if migrated
		}()
		clientSession, err := quic.DialAddr(context.Background(), server.Addr().String(), clientTlsConf, &quic.Config{MaxIdleTimeout: 100 * time.Millisecond})
		Expect(err).ToNot(HaveOccurred())
		res := clientSession.Handover(true, &handover.ConnectionStateStoreConf{
			IncludePendingOutgoingFrames: false,
		})
		handoverState := res.State
		err = res.Error
		Expect(err).ToNot(HaveOccurred())
		Expect(handoverState.Transport.VantagePoint).To(Equal("client"))
		migratedClientSession, _, err := quic.Restore(nil, &handoverState, &quic.ConnectionRestoreConfig{})
		Expect(err).ToNot(HaveOccurred())

		// transfer
		err = openAndSend(migratedClientSession, message1, true)
		Expect(err).ToNot(HaveOccurred())
		err = acceptAndReceive(migratedClientSession, message1, true)
		Expect(err).ToNot(HaveOccurred())

		// destroy sessions
		clientSession.Destroy()
		migratedClientSession.Destroy()
		serverSession.Destroy()
		Expect(server.Close()).ToNot(HaveOccurred())
	})

	It("restore server from client", func() {
		originalServerAddrChan := make(chan net.Addr, 1)
		go func() {
			defer GinkgoRecover()
			server, err := quic.ListenAddr("127.0.0.1:0", serverTlsConf, &quic.Config{MaxIdleTimeout: idleTimeout})
			Expect(err).ToNot(HaveOccurred())
			defer server.Close()
			originalServerAddrChan <- server.Addr()
			serverConn, err := server.Accept(context.Background())
			Expect(err).ToNot(HaveOccurred())
			<-serverConn.Context().Done()
		}()

		originalServerAddr := <-originalServerAddrChan
		clientConn, err := quic.DialAddr(context.Background(), originalServerAddr.String(), clientTlsConf, &quic.Config{MaxIdleTimeout: idleTimeout})
		Expect(err).ToNot(HaveOccurred())
		res := clientConn.Handover(false, &handover.ConnectionStateStoreConf{
			IgnoreCurrentPath: true,
		})
		handoverState := res.State
		err = res.Error
		Expect(err).ToNot(HaveOccurred())

		restoredServerAddrChan := make(chan net.Addr, 1)
		go func() {
			defer GinkgoRecover()
			handoverState = handoverState.ChangeVantagePoint(
				clientConn.LocalAddr().(*net.UDPAddr).IP.String(),
				uint16(clientConn.LocalAddr().(*net.UDPAddr).Port),
			)
			Expect(handoverState.Transport.VantagePoint).To(Equal("server"))
			restoredServerConn, _, err := quic.Restore(nil, &handoverState, &quic.ConnectionRestoreConfig{
				MaxIdleTimeout: idleTimeout,
			})
			Expect(err).ToNot(HaveOccurred())
			restoredServerAddrChan <- restoredServerConn.LocalAddr()
			err = acceptAndReceive(restoredServerConn, message1, true)
			Expect(err).ToNot(HaveOccurred())
			err = openAndSend(restoredServerConn, message1, true)
			Expect(err).ToNot(HaveOccurred())
			<-restoredServerConn.Context().Done()
		}()

		// transfer
		err = openAndSend(clientConn, message1, true)
		Expect(err).ToNot(HaveOccurred())
		err = acceptAndReceive(clientConn, message1, true)
		Expect(err).ToNot(HaveOccurred())

		restoredServerAddr := <-restoredServerAddrChan
		Expect(clientConn.RemoteAddr().(*net.UDPAddr).Port).ToNot(Equal(originalServerAddr.(*net.UDPAddr).Port))
		Expect(clientConn.RemoteAddr().(*net.UDPAddr).Port).To(Equal(restoredServerAddr.(*net.UDPAddr).Port))

		// destroy sessions
		clientConn.Destroy()
	})

	It("client handover twice", func() {
		format.MaxLength = 8000
		var wg sync.WaitGroup
		serverAddrChan := make(chan net.Addr, 1)
		serverRemoteAddrChan := make(chan net.Addr, 1)
		wg.Add(1)
		go func() {
			defer GinkgoRecover()
			defer wg.Done()
			server, err := quic.ListenAddr("127.0.0.1:0", serverTlsConf, &quic.Config{MaxIdleTimeout: 100 * time.Millisecond})
			Expect(err).ToNot(HaveOccurred())
			defer server.Close()
			serverAddrChan <- server.Addr()
			serverConn, err := server.Accept(context.Background())
			Expect(err).ToNot(HaveOccurred())
			defer serverConn.Destroy()
			// transfer
			err = acceptAndReceive(serverConn, message1, true)
			Expect(err).ToNot(HaveOccurred())
			err = openAndSend(serverConn, message1, true)
			Expect(err).ToNot(HaveOccurred())
			serverRemoteAddrChan <- serverConn.RemoteAddr()
			<-serverConn.Context().Done()
		}()
		serverAddr := <-serverAddrChan
		clientConn1, err := quic.DialAddr(context.Background(), serverAddr.String(), clientTlsConf, &quic.Config{MaxIdleTimeout: protocol.MinRemoteIdleTimeout, Tracer: qlog.DefaultTracerWithLabel("client")})
		Expect(err).ToNot(HaveOccurred())
		res := clientConn1.Handover(true, &handover.ConnectionStateStoreConf{})
		clientState1 := res.State
		err = res.Error
		Expect(err).ToNot(HaveOccurred())
		Expect(clientState1.Transport.VantagePoint).To(Equal("client"))
		clientConn2, _, err := quic.Restore(nil, &clientState1, &quic.ConnectionRestoreConfig{
			MaxIdleTimeout: protocol.MinRemoteIdleTimeout,
			Tracer:         qlog.DefaultTracerWithLabel("restored_client"),
			DefaultRTT:     utils.New(time.Second),
		})
		Expect(err).ToNot(HaveOccurred())
		res = clientConn2.Handover(true, &handover.ConnectionStateStoreConf{})
		clientState2 := res.State
		err = res.Error
		Expect(err).ToNot(HaveOccurred())
		Expect(clientState2.Transport.VantagePoint).To(Equal("client"))
		clientConn3, _, err := quic.Restore(nil, &clientState2, &quic.ConnectionRestoreConfig{
			MaxIdleTimeout: protocol.MinRemoteIdleTimeout,
			Tracer:         qlog.DefaultTracerWithLabel("restored_client_2"),
			DefaultRTT:     utils.New(time.Second),
		})
		Expect(err).ToNot(HaveOccurred())
		defer clientConn3.Destroy()
		// compare handover states
		clientState1.Transport.ConnectionIDs[0].StatelessResetToken = nil // clients are not configured to send stateless resets
		clientState2.Transport.ConnectionIDs[0].StatelessResetToken = nil // clients are not configured to send stateless resets
		Expect(clientState2.Transport.NextPacketNumber >= clientState1.Transport.NextPacketNumber).To(BeTrue())
		clientState2.Transport.NextPacketNumber = -1
		clientState1.Transport.NextPacketNumber = -1
		Expect(clientState1).To(BeEquivalentTo(clientState2))
		// transmit
		err = openAndSend(clientConn3, message1, true)
		Expect(err).ToNot(HaveOccurred())
		err = acceptAndReceive(clientConn3, message1, true)
		Expect(err).ToNot(HaveOccurred())
		// check if migrated
		serverRemoteAddr := <-serverRemoteAddrChan
		Expect(serverRemoteAddr.(*net.UDPAddr).Port).ToNot(Equal(clientConn1.LocalAddr().(*net.UDPAddr).Port))
		Expect(serverRemoteAddr.(*net.UDPAddr).Port).ToNot(Equal(clientConn2.LocalAddr().(*net.UDPAddr).Port))
		Expect(serverRemoteAddr.(*net.UDPAddr).Port).To(Equal(clientConn3.LocalAddr().(*net.UDPAddr).Port))
		Expect(err).ToNot(HaveOccurred())
		wg.Wait()
	})
})

func receiveMessage(stream quic.Stream, msg string, checkEOF bool) error {
	buf := make([]byte, len(msg))
	n, err := io.ReadAtLeast(stream, buf, len(msg))
	if err != nil && err != io.EOF {
		return err
	}
	if string(buf[:n]) != msg {
		return fmt.Errorf("failed to read message: expected \"%s\" but received \"%s\"", msg, buf[:n])
	}
	if checkEOF {
		err := checkStreamEOF(stream)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkStreamEOF(stream quic.ReceiveStream) error {
	buf := make([]byte, 1)
	n, err := stream.Read(buf)
	if err != io.EOF || n != 0 {
		return fmt.Errorf("not at EOF")
	}
	return nil
}

func sendMessage(stream quic.Stream, msg string, closeStreamAfterWrite bool) error {
	buf := []byte(msg)
	n, err := stream.Write(buf)
	if err != nil {
		return err
	}
	if n != len(buf) {
		return fmt.Errorf("failed to write all")
	}
	if closeStreamAfterWrite {
		err := stream.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func openAndSend(conn quic.Connection, msg string, closeStreamAfterWrite bool) error {
	stream, err := conn.OpenStream()
	if err != nil {
		return err
	}
	err = sendMessage(stream, msg, closeStreamAfterWrite)
	if err != nil {
		return err
	}
	return nil
}

func acceptAndReceive(conn quic.Connection, msg string, checkEOF bool) error {
	stream, err := conn.AcceptStream(context.Background())
	if err != nil {
		return err
	}
	err = receiveMessage(stream, msg, checkEOF)
	if err != nil {
		return err
	}
	return nil
}
