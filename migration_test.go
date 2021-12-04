package quic

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/lucas-clemente/quic-go/internal/testdata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("Migration", func() {
	var (
		message       []byte
		serverTlsConf *tls.Config
		clientTlsConf *tls.Config
	)

	BeforeEach(func() {
		message = []byte("hello")
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
		Eventually(areSessionsRunning).Should(BeFalse())
		Eventually(areServersRunning).Should(BeFalse())
		Eventually(areClosedSessionsRunning).Should(BeFalse())
	})

	It("server migration", func() {
		server, err := ListenAddr("127.0.0.1:0", serverTlsConf, &Config{MaxIdleTimeout: time.Second})
		originalServerAddr := server.Addr()
		Expect(err).ToNot(HaveOccurred())
		var clientSession Session
		//start client
		go func() {
			defer GinkgoRecover()
			clientSession, err = DialAddr(originalServerAddr.String(), clientTlsConf, &Config{EnableActiveMigration: true, MaxIdleTimeout: time.Second})
			Expect(err).ToNot(HaveOccurred())
			stream, err := clientSession.AcceptStream(context.Background())
			Expect(err).ToNot(HaveOccurred())
			buf := make([]byte, len(message))
			n, err := stream.Read(buf[:])
			Expect(err).ToNot(HaveOccurred())
			Expect(n).To(Equal(len(message)))
			n, err = stream.Write(buf[:])
			Expect(err).ToNot(HaveOccurred())
			Expect(n).To(Equal(len(message)))
		}()
		//run server
		serverSession, err := server.Accept(context.Background())
		Expect(err).ToNot(HaveOccurred())
		migratedServerAddr, err := serverSession.MigrateUDPSocket()
		Expect(err).ToNot(HaveOccurred())
		Expect(originalServerAddr.String()).ToNot(Equal(migratedServerAddr.String()))
		Expect(server.Addr().String()).To(Equal(migratedServerAddr.String()))
		stream, err := serverSession.OpenStream()
		Expect(err).ToNot(HaveOccurred())
		n, err := stream.Write(message)
		Expect(err).ToNot(HaveOccurred())
		Expect(n).To(Equal(len(message)))
		buf := make([]byte, len(message))
		n, err = stream.Read(buf[:])
		Expect(err).ToNot(HaveOccurred())
		Expect(n).To(Equal(len(message)))
		Expect(buf[:]).To(Equal(message))
		// destroy sessions
		serverSession.(*session).destroyImpl(nil)
		clientSession.(*session).destroyImpl(nil)
		Expect(server.Close()).ToNot(HaveOccurred())
	})

	It("client migration", func() {
		server, err := ListenAddr("127.0.0.1:0", serverTlsConf, &Config{EnableActiveMigration: true, MaxIdleTimeout: time.Second})
		Expect(err).ToNot(HaveOccurred())
		var serverSession Session
		//run server
		go func() {
			defer GinkgoRecover()
			serverSession, err = server.Accept(context.Background())
			Expect(err).ToNot(HaveOccurred())
			stream, err := serverSession.AcceptStream(context.Background())
			Expect(err).ToNot(HaveOccurred())
			buf := make([]byte, len(message))
			n, err := stream.Read(buf[:])
			Expect(err).ToNot(HaveOccurred())
			Expect(n).To(Equal(len(message)))
			n, err = stream.Write(buf[:])
			Expect(err).ToNot(HaveOccurred())
			Expect(n).To(Equal(len(message)))
		}()
		//start client
		clientSession, err := DialAddr(server.Addr().String(), clientTlsConf, &Config{MaxIdleTimeout: time.Second})
		Expect(err).ToNot(HaveOccurred())
		originalClientAddr := clientSession.LocalAddr()
		migratedClientAddr, err := clientSession.MigrateUDPSocket()
		Expect(err).ToNot(HaveOccurred())
		Expect(originalClientAddr.String()).ToNot(Equal(migratedClientAddr.String()))
		Expect(clientSession.LocalAddr().String()).To(Equal(migratedClientAddr.String()))
		stream, err := clientSession.OpenStream()
		Expect(err).ToNot(HaveOccurred())
		n, err := stream.Write(message)
		Expect(err).ToNot(HaveOccurred())
		Expect(n).To(Equal(len(message)))
		buf := make([]byte, len(message))
		n, err = stream.Read(buf[:])
		Expect(err).ToNot(HaveOccurred())
		Expect(n).To(Equal(len(message)))
		Expect(buf[:]).To(Equal(message))
		// destroy sessions
		serverSession.(*session).destroyImpl(nil)
		clientSession.(*session).destroyImpl(nil)
		Expect(server.Close()).ToNot(HaveOccurred())
	})

})
