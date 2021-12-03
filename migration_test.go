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
		Eventually(areSessionsRunning).Should(BeFalse())
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
	})

	It("server migration", func() {
		server, err := ListenAddr("127.0.0.1:0", serverTlsConf, &Config{MaxIdleTimeout: time.Second})
		defer server.Close()
		originalServerAddr := server.Addr()
		Expect(err).ToNot(HaveOccurred())
		//start client
		go func() {
			defer GinkgoRecover()
			client, err := DialAddr(originalServerAddr.String(), clientTlsConf, &Config{EnableActiveMigration: true, MaxIdleTimeout: time.Second})
			Expect(err).ToNot(HaveOccurred())
			stream, err := client.AcceptStream(context.Background())
			Expect(err).ToNot(HaveOccurred())
			buf := make([]byte, len(message))
			n, err := stream.Read(buf[:])
			Expect(err).ToNot(HaveOccurred())
			Expect(n).To(BeEquivalentTo(len(message)))
			n, err = stream.Write(buf[:])
			Expect(err).ToNot(HaveOccurred())
			Expect(n).To(BeEquivalentTo(len(message)))
			<-stream.Context().Done()
			err = client.CloseWithError(ApplicationErrorCode(0), "client_close")
			Expect(err).ToNot(HaveOccurred())
		}()
		//run server
		session, err := server.Accept(context.Background())
		Expect(err).ToNot(HaveOccurred())
		migratedServerAddr, err := session.MigrateUDPSocket()
		Expect(err).ToNot(HaveOccurred())
		Expect(originalServerAddr.String()).ToNot(Equal(migratedServerAddr.String()))
		Expect(server.Addr().String()).To(Equal(migratedServerAddr.String()))
		stream, err := session.OpenStream()
		Expect(err).ToNot(HaveOccurred())
		n, err := stream.Write(message)
		Expect(err).ToNot(HaveOccurred())
		Expect(n).To(BeEquivalentTo(len(message)))
		buf := make([]byte, len(message))
		n, err = stream.Read(buf[:])
		Expect(err).ToNot(HaveOccurred())
		Expect(n).To(Equal(len(message)))
		Expect(buf[:]).To(Equal(message))
		err = session.CloseWithError(ApplicationErrorCode(0), "server_close")
		Expect(err).ToNot(HaveOccurred())
	})

	It("client migration", func() {
		server, err := ListenAddr("127.0.0.1:0", serverTlsConf, &Config{EnableActiveMigration: true, MaxIdleTimeout: time.Second})
		Expect(err).ToNot(HaveOccurred())
		//run server
		go func() {
			defer GinkgoRecover()
			defer server.Close()
			session, err := server.Accept(context.Background())
			Expect(err).ToNot(HaveOccurred())
			stream, err := session.AcceptStream(context.Background())
			Expect(err).ToNot(HaveOccurred())
			buf := make([]byte, len(message))
			n, err := stream.Read(buf[:])
			Expect(err).ToNot(HaveOccurred())
			Expect(n).To(BeEquivalentTo(len(message)))
			n, err = stream.Write(buf[:])
			Expect(err).ToNot(HaveOccurred())
			Expect(n).To(BeEquivalentTo(len(message)))
			<-stream.Context().Done()
			<-session.Context().Done()
			err = session.CloseWithError(ApplicationErrorCode(0), "server_close")
			Expect(err).ToNot(HaveOccurred())
		}()
		//start client
		client, err := DialAddr(server.Addr().String(), clientTlsConf, &Config{MaxIdleTimeout: time.Second})
		Expect(err).ToNot(HaveOccurred())
		originalClientAddr := client.LocalAddr()
		migratedClientAddr, err := client.MigrateUDPSocket()
		Expect(err).ToNot(HaveOccurred())
		Expect(originalClientAddr.String()).ToNot(Equal(migratedClientAddr.String()))
		Expect(client.LocalAddr().String()).To(Equal(migratedClientAddr.String()))
		stream, err := client.OpenStream()
		Expect(err).ToNot(HaveOccurred())
		n, err := stream.Write(message)
		Expect(err).ToNot(HaveOccurred())
		Expect(n).To(BeEquivalentTo(len(message)))
		buf := make([]byte, len(message))
		n, err = stream.Read(buf[:])
		Expect(err).ToNot(HaveOccurred())
		Expect(n).To(Equal(len(message)))
		Expect(buf[:]).To(Equal(message))
		err = client.CloseWithError(ApplicationErrorCode(0), "client_close")
		Expect(err).ToNot(HaveOccurred())
	})

})
