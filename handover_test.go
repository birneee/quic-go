package quic

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/testdata"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"net"
	"time"
)

var _ = Describe("Handover", func() {
	var (
		message       = "hello"
		idleTimeout   = 100 * time.Millisecond
		serverTlsConf *tls.Config
		clientTlsConf *tls.Config
	)

	BeforeEach(func() {
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

	It("client to client handover", func() {
		server, err := ListenAddr("127.0.0.1:0", serverTlsConf, &Config{EnableActiveMigration: true, MaxIdleTimeout: 100 * time.Millisecond})
		var serverSession Connection
		Expect(err).ToNot(HaveOccurred())
		go func() {
			defer GinkgoRecover()
			serverSession, err = server.Accept(context.Background())
			Expect(err).ToNot(HaveOccurred())
			originalClientAddr := serverSession.RemoteAddr()

			// transfer
			err := acceptAndReceive(serverSession, message)
			Expect(err).ToNot(HaveOccurred())
			err = openAndSend(serverSession, message)
			Expect(err).ToNot(HaveOccurred())
			<-serverSession.Context().Done()

			migratedClientAddr := serverSession.RemoteAddr()
			Expect(originalClientAddr.(*net.UDPAddr).Port).ToNot(Equal(migratedClientAddr.(*net.UDPAddr).Port)) // check if migrated
		}()
		clientSession, err := DialAddr(server.Addr().String(), clientTlsConf, &Config{IgnoreReceived1RTTPacketsUntilFirstPathMigration: true, MaxIdleTimeout: 100 * time.Millisecond})
		Expect(err).ToNot(HaveOccurred())
		handoverState, err := clientSession.Handover(true, false)
		Expect(err).ToNot(HaveOccurred())
		migratedClientSession, err := Restore(handoverState, protocol.PerspectiveClient, &Config{LoggerPrefix: "cloned"})
		Expect(err).ToNot(HaveOccurred())

		// transfer
		err = openAndSend(migratedClientSession, message)
		Expect(err).ToNot(HaveOccurred())
		err = acceptAndReceive(migratedClientSession, message)
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
			server, err := ListenAddr("127.0.0.1:0", serverTlsConf, &Config{EnableActiveMigration: true, MaxIdleTimeout: idleTimeout})
			Expect(err).ToNot(HaveOccurred())
			defer server.Close()
			originalServerAddrChan <- server.Addr()
			serverConn, err := server.Accept(context.Background())
			Expect(err).ToNot(HaveOccurred())
			<-serverConn.Context().Done()
		}()

		originalServerAddr := <-originalServerAddrChan
		clientConn, err := DialAddr(originalServerAddr.String(), clientTlsConf, &Config{EnableActiveMigration: true, IgnoreReceived1RTTPacketsUntilFirstPathMigration: true, MaxIdleTimeout: idleTimeout})
		Expect(err).ToNot(HaveOccurred())
		handoverState, err := clientConn.Handover(false, true)
		Expect(err).ToNot(HaveOccurred())

		restoredServerAddrChan := make(chan net.Addr, 1)
		go func() {
			defer GinkgoRecover()
			restoredServerConn, err := Restore(handoverState, PerspectiveServer, &Config{MaxIdleTimeout: idleTimeout})
			Expect(err).ToNot(HaveOccurred())
			restoredServerAddrChan <- restoredServerConn.LocalAddr()
			err = acceptAndReceive(restoredServerConn, message)
			Expect(err).ToNot(HaveOccurred())
			err = openAndSend(restoredServerConn, message)
			Expect(err).ToNot(HaveOccurred())
			<-restoredServerConn.Context().Done()
		}()

		// transfer
		err = openAndSend(clientConn, message)
		Expect(err).ToNot(HaveOccurred())
		err = acceptAndReceive(clientConn, message)
		Expect(err).ToNot(HaveOccurred())

		restoredServerAddr := <-restoredServerAddrChan
		Expect(clientConn.RemoteAddr().(*net.UDPAddr).Port).ToNot(Equal(originalServerAddr.(*net.UDPAddr).Port))
		Expect(clientConn.RemoteAddr().(*net.UDPAddr).Port).To(Equal(restoredServerAddr.(*net.UDPAddr).Port))

		// destroy sessions
		clientConn.(*connection).destroyImpl(nil)
	})

	It("client handover twice", func() {
		format.MaxLength = 8000
		serverAddrChan := make(chan net.Addr, 1)
		serverRemoteAddrChan := make(chan net.Addr, 1)
		go func() {
			defer GinkgoRecover()
			server, err := ListenAddr("127.0.0.1:0", serverTlsConf, &Config{EnableActiveMigration: true, MaxIdleTimeout: 100 * time.Millisecond})
			Expect(err).ToNot(HaveOccurred())
			defer server.Close()
			serverAddrChan <- server.Addr()
			serverConn, err := server.Accept(context.Background())
			Expect(err).ToNot(HaveOccurred())
			defer serverConn.(*connection).destroy(nil)
			// transfer
			err = acceptAndReceive(serverConn, message)
			Expect(err).ToNot(HaveOccurred())
			err = openAndSend(serverConn, message)
			Expect(err).ToNot(HaveOccurred())
			serverRemoteAddrChan <- serverConn.RemoteAddr()
			<-serverConn.Context().Done()
		}()
		serverAddr := <-serverAddrChan
		clientConn1, err := DialAddr(serverAddr.String(), clientTlsConf, &Config{IgnoreReceived1RTTPacketsUntilFirstPathMigration: true, MaxIdleTimeout: 100 * time.Millisecond, EnableActiveMigration: true, LoggerPrefix: "client1"})
		Expect(err).ToNot(HaveOccurred())
		clientState1, err := clientConn1.Handover(true, true)
		Expect(err).ToNot(HaveOccurred())
		clientConn2, err := Restore(clientState1, PerspectiveClient, &Config{IgnoreReceived1RTTPacketsUntilFirstPathMigration: true, LoggerPrefix: "client2"})
		Expect(err).ToNot(HaveOccurred())
		clientState2, err := clientConn2.Handover(true, true)
		Expect(err).ToNot(HaveOccurred())
		clientConn3, err := Restore(clientState2, PerspectiveClient, &Config{MaxIdleTimeout: 100 * time.Millisecond, LoggerPrefix: "client3"})
		Expect(err).ToNot(HaveOccurred())
		defer clientConn3.(*connection).destroyImpl(nil)
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
		err = openAndSend(clientConn3, message)
		Expect(err).ToNot(HaveOccurred())
		err = acceptAndReceive(clientConn3, message)
		Expect(err).ToNot(HaveOccurred())
		// check if migrated
		serverRemoteAddr := <-serverRemoteAddrChan
		Expect(serverRemoteAddr.(*net.UDPAddr).Port).ToNot(Equal(clientConn1.LocalAddr().(*net.UDPAddr).Port))
		Expect(serverRemoteAddr.(*net.UDPAddr).Port).ToNot(Equal(clientConn2.LocalAddr().(*net.UDPAddr).Port))
		Expect(serverRemoteAddr.(*net.UDPAddr).Port).To(Equal(clientConn3.LocalAddr().(*net.UDPAddr).Port))
		Expect(err).ToNot(HaveOccurred())
	})
})

func receiveMessage(stream Stream, msg string) error {
	buf := make([]byte, 2*len(msg))
	n, err := stream.Read(buf)
	if err != nil {
		return err
	}
	if string(buf[:n]) != msg {
		return fmt.Errorf("failed to read message")
	}
	return nil
}

func sendMessage(stream Stream, msg string) error {
	buf := []byte(msg)
	n, err := stream.Write(buf)
	if err != nil {
		return err
	}
	if n != len(buf) {
		return fmt.Errorf("failed to write all")
	}
	return nil
}

func openAndSend(conn Connection, msg string) error {
	stream, err := conn.OpenStream()
	if err != nil {
		return err
	}
	err = sendMessage(stream, msg)
	if err != nil {
		return err
	}
	return nil
}

func acceptAndReceive(conn Connection, msg string) error {
	stream, err := conn.AcceptStream(context.Background())
	if err != nil {
		return err
	}
	err = receiveMessage(stream, msg)
	if err != nil {
		return err
	}
	return nil
}
