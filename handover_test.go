package quic

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/testdata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("Handover", func() {
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

	It("client handover", func() {
		server, err := ListenAddr("127.0.0.1:0", serverTlsConf, &Config{EnableActiveMigration: true, MaxIdleTimeout: 100 * time.Millisecond})
		var serverSession Session
		Expect(err).ToNot(HaveOccurred())
		go func() {
			defer GinkgoRecover()
			serverSession, err = server.Accept(context.Background())
			Expect(err).ToNot(HaveOccurred())
			originalClientAddr := serverSession.RemoteAddr()
			stream, err := serverSession.AcceptStream(context.Background())
			Expect(err).ToNot(HaveOccurred())
			migratedClientAddr := serverSession.RemoteAddr()
			Expect(originalClientAddr.String()).ToNot(Equal(migratedClientAddr.String())) // check if migrated
			buf := make([]byte, len(message))
			n, err := stream.Read(buf)
			Expect(err).ToNot(HaveOccurred())
			Expect(n).To(Equal(len(message)))
			n, err = stream.Write(buf)
			Expect(err).ToNot(HaveOccurred())
			Expect(n).To(Equal(len(message)))
		}()
		clientSession, err := DialAddr(server.Addr().String(), clientTlsConf, &Config{IgnoreReceived1RTTPacketsUntilFirstPathMigration: true, MaxIdleTimeout: 100 * time.Millisecond})
		Expect(err).ToNot(HaveOccurred())
		handoverState, err := clientSession.Handover(true, false)
		Expect(err).ToNot(HaveOccurred())
		migratedClientSession, err := Restore(handoverState, protocol.PerspectiveClient, &Config{LoggerPrefix: "cloned"})
		Expect(err).ToNot(HaveOccurred())
		stream, err := migratedClientSession.OpenStream()
		Expect(err).ToNot(HaveOccurred())
		n, err := stream.Write(message)
		Expect(err).ToNot(HaveOccurred())
		Expect(n).To(Equal(len(message)))
		buf := make([]byte, len(message))
		n, err = stream.Read(buf)
		Expect(err).ToNot(HaveOccurred())
		Expect(n).To(Equal(len(message)))
		Expect(buf).To(Equal(message))
		// destroy sessions
		clientSession.(*session).destroyImpl(nil)
		migratedClientSession.(*session).destroyImpl(nil)
		serverSession.(*session).destroyImpl(nil)
		Expect(server.Close()).ToNot(HaveOccurred())
	})

})
