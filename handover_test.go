package quic

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"github.com/quic-go/quic-go/handover"
	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/testdata"
	"github.com/quic-go/quic-go/logging"
	"io"
	"net"
	"sync"
	"time"
)

var _ = Describe("Handover", func() {
	var (
		message1      = "hello"
		message2      = "hquic"
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
		Eventually(areConnsRunning).Should(BeFalse())
		Eventually(areServersRunning).Should(BeFalse())
	})

	It("server to server handover mid-stream", func() {
		MaxIdleTimeout := 100 * time.Millisecond
		server1AddrChan := make(chan net.Addr, 1)
		server2AddrChan := make(chan net.Addr, 1)
		serverStateChan := make(chan handover.State, 1)

		go func() { // server 1
			defer GinkgoRecover()
			server, err := ListenAddr("127.0.0.1:0", serverTlsConf, &Config{MaxIdleTimeout: MaxIdleTimeout})
			Expect(err).ToNot(HaveOccurred())
			defer server.Close()
			server1AddrChan <- server.Addr()
			conn, err := server.Accept(context.Background())
			Expect(err).ToNot(HaveOccurred())
			defer conn.(*connection).destroyImpl(nil)
			// transfer data
			err = acceptAndReceive(conn, message1, false)
			Expect(err).ToNot(HaveOccurred())
			err = openAndSend(conn, message1, false)
			Expect(err).ToNot(HaveOccurred())
			// state handover
			res := conn.Handover(true, &ConnectionStateStoreConf{
				IncludeStreamState:           true,
				IncludePendingOutgoingFrames: true,
				IncludePendingIncomingFrames: true,
			})
			serverState := res.State
			err = res.Error
			Expect(err).ToNot(HaveOccurred())
			serverStateChan <- serverState
			<-conn.Context().Done()
		}()
		go func() { // server 2
			defer GinkgoRecover()
			serverState := <-serverStateChan
			conn, restoredStreams, err := Restore(serverState, &ConnectionRestoreConfig{
				Perspective: logging.PerspectiveServer,
				QuicConf:    &Config{MaxIdleTimeout: MaxIdleTimeout},
			})
			bidiStreams := restoredStreams.BidiStreams
			sendStreams := restoredStreams.SendStreams
			receiveStreams := restoredStreams.ReceiveStreams
			Expect(err).ToNot(HaveOccurred())
			Expect(len(bidiStreams)).To(BeEquivalentTo(2))
			Expect(sendStreams).To(BeEmpty())
			Expect(receiveStreams).To(BeEmpty())
			defer conn.(*connection).destroyImpl(nil)
			server2AddrChan <- conn.LocalAddr()
			// transfer
			err = receiveMessage(bidiStreams[0], message2, true)
			Expect(err).ToNot(HaveOccurred())
			err = sendMessage(bidiStreams[1], message2, true)
			Expect(err).ToNot(HaveOccurred())
			<-conn.Context().Done()
		}()

		originalServerAddr := <-server1AddrChan
		clientConn, err := DialAddr(context.Background(), originalServerAddr.String(), clientTlsConf, &Config{MaxIdleTimeout: MaxIdleTimeout})
		Expect(err).ToNot(HaveOccurred())
		// transfer
		err = openAndSend(clientConn, message1+message2, true)
		Expect(err).ToNot(HaveOccurred())
		err = acceptAndReceive(clientConn, message1+message2, true)
		Expect(err).ToNot(HaveOccurred())
		migratedServerAddr := clientConn.RemoteAddr()
		Expect(originalServerAddr.String()).ToNot(Equal(migratedServerAddr.String())) // check if migrated

		// destroy sessions
		clientConn.(*connection).destroyImpl(nil)
	})

	It("server to server handover", func() {
		MaxIdleTimeout := 100 * time.Millisecond
		server1AddrChan := make(chan net.Addr, 1)
		server2AddrChan := make(chan net.Addr, 1)
		serverStateChan := make(chan handover.State, 1)

		go func() { // server 1
			defer GinkgoRecover()
			server, err := ListenAddr("127.0.0.1:0", serverTlsConf, &Config{MaxIdleTimeout: MaxIdleTimeout})
			Expect(err).ToNot(HaveOccurred())
			defer server.Close()
			server1AddrChan <- server.Addr()
			conn, err := server.Accept(context.Background())
			Expect(err).ToNot(HaveOccurred())
			defer conn.(*connection).destroyImpl(nil)
			res := conn.Handover(true, &ConnectionStateStoreConf{})
			serverState := res.State
			err = res.Error
			Expect(err).ToNot(HaveOccurred())
			serverStateChan <- serverState
			<-conn.Context().Done()
		}()
		go func() { // server 2
			defer GinkgoRecover()
			serverState := <-serverStateChan
			conn, _, err := Restore(serverState, &ConnectionRestoreConfig{
				Perspective: logging.PerspectiveServer,
				QuicConf:    &Config{MaxIdleTimeout: MaxIdleTimeout},
			})
			Expect(err).ToNot(HaveOccurred())
			defer conn.(*connection).destroyImpl(nil)
			server2AddrChan <- conn.LocalAddr()
			// transfer
			err = openAndSend(conn, message1, true)
			Expect(err).ToNot(HaveOccurred())
			err = acceptAndReceive(conn, message1, true)
			Expect(err).ToNot(HaveOccurred())
			<-conn.Context().Done()
		}()

		originalServerAddr := <-server1AddrChan
		clientConn, err := DialAddr(context.Background(), originalServerAddr.String(), clientTlsConf, &Config{MaxIdleTimeout: MaxIdleTimeout})
		Expect(err).ToNot(HaveOccurred())
		// transfer
		err = acceptAndReceive(clientConn, message1, true)
		Expect(err).ToNot(HaveOccurred())
		err = openAndSend(clientConn, message1, true)
		Expect(err).ToNot(HaveOccurred())
		migratedServerAddr := clientConn.RemoteAddr()
		Expect(originalServerAddr.String()).ToNot(Equal(migratedServerAddr.String())) // check if migrated
		<-clientConn.Context().Done()

		// destroy sessions
		clientConn.(*connection).destroyImpl(nil)
	})

	It("client to client handover", func() {
		server, err := ListenAddr("127.0.0.1:0", serverTlsConf, &Config{MaxIdleTimeout: 100 * time.Millisecond})
		var serverSession Connection
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
		clientSession, err := DialAddr(context.Background(), server.Addr().String(), clientTlsConf, &Config{MaxIdleTimeout: 100 * time.Millisecond})
		Expect(err).ToNot(HaveOccurred())
		res := clientSession.Handover(true, &ConnectionStateStoreConf{
			IncludePendingOutgoingFrames: false,
		})
		handoverState := res.State
		err = res.Error
		Expect(err).ToNot(HaveOccurred())
		migratedClientSession, _, err := Restore(handoverState, &ConnectionRestoreConfig{
			Perspective: protocol.PerspectiveClient,
			QuicConf:    &Config{},
		})
		Expect(err).ToNot(HaveOccurred())

		// transfer
		err = openAndSend(migratedClientSession, message1, true)
		Expect(err).ToNot(HaveOccurred())
		err = acceptAndReceive(migratedClientSession, message1, true)
		Expect(err).ToNot(HaveOccurred())

		// destroy sessions
		clientSession.(*connection).destroyImpl(nil)
		migratedClientSession.(*connection).destroyImpl(nil)
		serverSession.(*connection).destroyImpl(nil)
		Expect(server.Close()).ToNot(HaveOccurred())
	})

	It("restore server from client", func() {
		originalServerAddrChan := make(chan net.Addr, 1)
		go func() {
			defer GinkgoRecover()
			server, err := ListenAddr("127.0.0.1:0", serverTlsConf, &Config{MaxIdleTimeout: idleTimeout})
			Expect(err).ToNot(HaveOccurred())
			defer server.Close()
			originalServerAddrChan <- server.Addr()
			serverConn, err := server.Accept(context.Background())
			Expect(err).ToNot(HaveOccurred())
			<-serverConn.Context().Done()
		}()

		originalServerAddr := <-originalServerAddrChan
		clientConn, err := DialAddr(context.Background(), originalServerAddr.String(), clientTlsConf, &Config{MaxIdleTimeout: idleTimeout})
		Expect(err).ToNot(HaveOccurred())
		res := clientConn.Handover(false, &ConnectionStateStoreConf{
			IgnoreCurrentPath: true,
		})
		handoverState := res.State
		err = res.Error
		Expect(err).ToNot(HaveOccurred())

		restoredServerAddrChan := make(chan net.Addr, 1)
		go func() {
			defer GinkgoRecover()
			restoredServerConn, _, err := Restore(handoverState, &ConnectionRestoreConfig{
				Perspective: logging.PerspectiveServer,
				QuicConf:    &Config{MaxIdleTimeout: idleTimeout},
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
		clientConn.(*connection).destroyImpl(nil)
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
			server, err := ListenAddr("127.0.0.1:0", serverTlsConf, &Config{MaxIdleTimeout: 100 * time.Millisecond})
			Expect(err).ToNot(HaveOccurred())
			defer server.Close()
			serverAddrChan <- server.Addr()
			serverConn, err := server.Accept(context.Background())
			Expect(err).ToNot(HaveOccurred())
			defer serverConn.(*connection).destroy(nil)
			// transfer
			err = acceptAndReceive(serverConn, message1, true)
			Expect(err).ToNot(HaveOccurred())
			err = openAndSend(serverConn, message1, true)
			Expect(err).ToNot(HaveOccurred())
			serverRemoteAddrChan <- serverConn.RemoteAddr()
			<-serverConn.Context().Done()
		}()
		serverAddr := <-serverAddrChan
		clientConn1, err := DialAddr(context.Background(), serverAddr.String(), clientTlsConf, &Config{MaxIdleTimeout: protocol.MinRemoteIdleTimeout})
		Expect(err).ToNot(HaveOccurred())
		res := clientConn1.Handover(true, &ConnectionStateStoreConf{})
		clientState1 := res.State
		err = res.Error
		Expect(err).ToNot(HaveOccurred())
		clientConn2, _, err := Restore(clientState1, &ConnectionRestoreConfig{
			Perspective: logging.PerspectiveClient,
			QuicConf:    &Config{MaxIdleTimeout: protocol.MinRemoteIdleTimeout},
		})
		Expect(err).ToNot(HaveOccurred())
		res = clientConn2.Handover(true, &ConnectionStateStoreConf{})
		clientState2 := res.State
		err = res.Error
		Expect(err).ToNot(HaveOccurred())
		clientConn3, _, err := Restore(clientState2, &ConnectionRestoreConfig{
			Perspective: logging.PerspectiveClient,
			QuicConf:    &Config{MaxIdleTimeout: protocol.MinRemoteIdleTimeout},
		})
		Expect(err).ToNot(HaveOccurred())
		defer clientConn3.(*connection).destroy(nil)
		// compare handover states
		clientState1.ClientAddress = ""                               // ignore changed client address
		clientState2.ClientAddress = ""                               // ignore changed client address
		clientState1.ClientConnectionIDs[0].StatelessResetToken = nil // TODO check statelessResetEnabled before comparing
		clientState2.ClientConnectionIDs[0].StatelessResetToken = nil // TODO check statelessResetEnabled before comparing
		clientState1.ClientHighestSentPacketNumber = 0                // ignore
		clientState1.ServerHighestSentPacketNumber = 0                // ignore
		clientState2.ClientHighestSentPacketNumber = 0                // ignore
		clientState2.ServerHighestSentPacketNumber = 0                // ignore
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

func receiveMessage(stream Stream, msg string, checkEOF bool) error {
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

func checkStreamEOF(stream ReceiveStream) error {
	buf := make([]byte, 1)
	n, err := stream.Read(buf)
	if err != io.EOF || n != 0 {
		return fmt.Errorf("not at EOF")
	}
	return nil
}

func sendMessage(stream Stream, msg string, closeStreamAfterWrite bool) error {
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

func openAndSend(conn Connection, msg string, closeStreamAfterWrite bool) error {
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

func acceptAndReceive(conn Connection, msg string, checkEOF bool) error {
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
