package quic

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/lucas-clemente/quic-go/internal/testdata"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("Migration", func() {
	var (
		message       = "hello"
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

	It("server migration", func() {
		server, err := ListenAddr("127.0.0.1:0", serverTlsConf, &Config{MaxIdleTimeout: time.Second})
		originalServerAddr := server.Addr()
		Expect(err).ToNot(HaveOccurred())
		var clientConn Connection
		//start client
		go func() {
			defer GinkgoRecover()
			clientConn, err = DialAddr(originalServerAddr.String(), clientTlsConf, &Config{EnableActiveMigration: true, MaxIdleTimeout: time.Second})
			Expect(err).ToNot(HaveOccurred())
			// transfer
			err = openAndSend(clientConn, message)
			Expect(err).ToNot(HaveOccurred())
			err = acceptAndReceive(clientConn, message)
			Expect(err).ToNot(HaveOccurred())
		}()
		//run server
		serverConn, err := server.Accept(context.Background())
		Expect(err).ToNot(HaveOccurred())
		migratedServerAddr, err := serverConn.MigrateUDPSocket()
		Expect(err).ToNot(HaveOccurred())
		Expect(originalServerAddr.String()).ToNot(Equal(migratedServerAddr.String()))
		Expect(server.Addr().String()).To(Equal(migratedServerAddr.String()))
		// transfer
		err = acceptAndReceive(serverConn, message)
		Expect(err).ToNot(HaveOccurred())
		err = openAndSend(serverConn, message)
		Expect(err).ToNot(HaveOccurred())
		<-serverConn.Context().Done()
		<-clientConn.Context().Done()
		// destroy sessions
		serverConn.(*connection).destroyImpl(nil)
		clientConn.(*connection).destroyImpl(nil)
		Expect(server.Close()).ToNot(HaveOccurred())
	})

	It("client migration", func() {
		server, err := ListenAddr("127.0.0.1:0", serverTlsConf, &Config{EnableActiveMigration: true, MaxIdleTimeout: time.Second})
		Expect(err).ToNot(HaveOccurred())
		var serverConn Connection
		//run server
		go func() {
			defer GinkgoRecover()
			serverConn, err = server.Accept(context.Background())
			Expect(err).ToNot(HaveOccurred())
			// transfer
			err = acceptAndReceive(serverConn, message)
			Expect(err).ToNot(HaveOccurred())
			err = openAndSend(serverConn, message)
			Expect(err).ToNot(HaveOccurred())
		}()
		//start client
		clientConn, err := DialAddr(server.Addr().String(), clientTlsConf, &Config{MaxIdleTimeout: time.Second})
		Expect(err).ToNot(HaveOccurred())
		originalClientAddr := clientConn.LocalAddr()
		migratedClientAddr, err := clientConn.MigrateUDPSocket()
		Expect(err).ToNot(HaveOccurred())
		Expect(originalClientAddr.String()).ToNot(Equal(migratedClientAddr.String()))
		Expect(clientConn.LocalAddr().String()).To(Equal(migratedClientAddr.String()))
		// transfer
		err = openAndSend(clientConn, message)
		Expect(err).ToNot(HaveOccurred())
		err = acceptAndReceive(clientConn, message)
		Expect(err).ToNot(HaveOccurred())
		<-clientConn.Context().Done()
		<-serverConn.Context().Done()
		// destroy sessions
		serverConn.(*connection).destroyImpl(nil)
		clientConn.(*connection).destroyImpl(nil)
		Expect(server.Close()).ToNot(HaveOccurred())
	})

})
