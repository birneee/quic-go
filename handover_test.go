package quic

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/testdata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
	"net"
)

var _ = Describe("Handover", func() {
	var (
		serverTlsConf *tls.Config
		clientTlsConf *tls.Config
	)

	//TODO server handover does not work yet
	//var runEchoClient = func(serverAddr net.Addr) {
	//	session, err := DialAddr(serverAddr.String(), clientTlsConf, &Config{})
	//	Expect(err).ToNot(HaveOccurred())
	//	stream, err := session.OpenStream()
	//	Expect(err).ToNot(HaveOccurred())
	//	message := []byte("hello")
	//	n, err := stream.Write(message)
	//	Expect(err).ToNot(HaveOccurred())
	//	err = stream.Close()
	//	Expect(err).ToNot(HaveOccurred())
	//	Expect(n).To(BeEquivalentTo(len(message)))
	//	echo, err := io.ReadAll(stream)
	//	Expect(err).ToNot(HaveOccurred())
	//	Expect(echo).To(BeEquivalentTo(message))
	//	err = session.CloseWithError(ApplicationErrorCode(0), "client_close")
	//	Expect(err).ToNot(HaveOccurred())
	//}

	var runEchoServer = func() net.Addr {
		server, err := ListenAddr("127.0.0.1:0", serverTlsConf, &Config{})
		Expect(err).ToNot(HaveOccurred())
		go func() {
			defer GinkgoRecover()
			session, err := server.Accept(context.Background())
			Expect(err).ToNot(HaveOccurred())
			stream, err := session.AcceptStream(context.Background())
			Expect(err).ToNot(HaveOccurred())
			data, err := io.ReadAll(stream)
			Expect(err).ToNot(HaveOccurred())
			n, err := stream.Write(data)
			Expect(err).ToNot(HaveOccurred())
			err = stream.Close()
			Expect(err).ToNot(HaveOccurred())
			Expect(n).To(BeEquivalentTo(len(data)))
			<-session.Context().Done()
			err = server.Close()
			Expect(err).ToNot(HaveOccurred())
		}()
		return server.Addr()
	}

	BeforeEach(func() {
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

	It("client handover", func() {
		serverAddr := runEchoServer()

		session, err := DialAddr(serverAddr.String(), clientTlsConf, &Config{IgnoreReceived1RTTPacketsUntilFirstPathMigration: true})
		Expect(err).ToNot(HaveOccurred())
		handoverState, err := session.Handover(true)
		Expect(err).ToNot(HaveOccurred())
		session, err = RestoreSessionFromHandoverState(handoverState, protocol.PerspectiveClient, nil, "cloned")
		Expect(err).ToNot(HaveOccurred())
		stream, err := session.OpenStream()
		Expect(err).ToNot(HaveOccurred())
		message := []byte("hello")
		n, err := stream.Write(message)
		Expect(err).ToNot(HaveOccurred())
		err = stream.Close()
		Expect(err).ToNot(HaveOccurred())
		Expect(n).To(BeEquivalentTo(len(message)))
		echo, err := io.ReadAll(stream)
		Expect(err).ToNot(HaveOccurred())
		Expect(echo).To(BeEquivalentTo(message))
		err = session.CloseWithError(ApplicationErrorCode(0), "client_close")
		Expect(err).ToNot(HaveOccurred())
	})

	//TODO server handover does not work yet
	//It("server handover", func(){
	//	server, err := ListenAddr("127.0.0.1:0", serverTlsConf, &Config{IgnoreReceived1RTTPacketsUntilFirstPathMigration: true})
	//	Expect(err).ToNot(HaveOccurred())
	//	go func() {
	//		defer GinkgoRecover()
	//		runEchoClient(server.Addr())
	//	}()
	//	session, err := server.Accept(context.Background())
	//	Expect(err).ToNot(HaveOccurred())
	//	handoverState, err := session.Handover(true)
	//	Expect(err).ToNot(HaveOccurred())
	//	session2, err := RestoreSessionFromHandoverState(handoverState, protocol.PerspectiveServer, nil, "cloned")
	//	Expect(err).ToNot(HaveOccurred())
	//	stream, err := session2.AcceptStream(context.Background())
	//	Expect(err).ToNot(HaveOccurred())
	//	data, err := io.ReadAll(stream)
	//	Expect(err).ToNot(HaveOccurred())
	//	n, err := stream.Write(data)
	//	Expect(err).ToNot(HaveOccurred())
	//	err = stream.Close()
	//	Expect(err).ToNot(HaveOccurred())
	//	Expect(n).To(BeEquivalentTo(len(data)))
	//	<- session2.Context().Done()
	//	err = server.Close()
	//	Expect(err).ToNot(HaveOccurred())
	//})
})
