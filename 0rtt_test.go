package quic

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/golang/mock/gomock"
	mocktls "github.com/quic-go/quic-go/internal/mocks/tls"
	"github.com/quic-go/quic-go/internal/testdata"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
	"time"
)

func get0RttConn(t assert.TestingT) (newClient func() EarlyConnection, listener *EarlyListener) {
	maxIdleTimeout := 10000 * time.Millisecond
	protos := []string{"proto1"}
	serverTlsConf := testdata.GetTLSConfig()
	serverTlsConf.NextProtos = protos
	certPool := x509.NewCertPool()
	testdata.AddRootCA(certPool)
	serverConfig := &Config{
		Allow0RTT:      true,
		MaxIdleTimeout: maxIdleTimeout,
	}
	tokenStore := NewMockTokenStore(mockCtrl)
	var token *ClientToken
	tokenSet := make(chan struct{})
	tokenStore.EXPECT().Pop(gomock.Any()).
		Return(nil).
		DoAndReturn(func(key string) *ClientToken {
			return token
		})
	tokenStore.EXPECT().Put(gomock.Any(), gomock.Any()).
		Do(func(key string, _token *ClientToken) {
			token = _token
			close(tokenSet)
		})
	tokenStore.EXPECT().Pop(gomock.Any()).AnyTimes().DoAndReturn(func(key string) *ClientToken {
		return token
	})
	tokenStore.EXPECT().Put(gomock.Any(), gomock.Any()).AnyTimes().Do(func(key string, _token *ClientToken) {
		token = _token
	})
	clientConf := &Config{
		MaxIdleTimeout: maxIdleTimeout,
		TokenStore:     tokenStore,
	}
	var ticket *tls.ClientSessionState
	ticketSet := make(chan struct{})
	sessionCache := mocktls.NewMockClientSessionCache(mockCtrl)
	sessionCache.EXPECT().Get(gomock.Any()).Return(nil, false)
	sessionCache.EXPECT().Put(gomock.Any(), gomock.Any()).
		Do(func(key string, _ticket *tls.ClientSessionState) {
			ticket = _ticket
			close(ticketSet)
		})
	sessionCache.EXPECT().Get(gomock.Any()).AnyTimes().DoAndReturn(func(key string) (*tls.ClientSessionState, bool) {
		return ticket, true
	})
	sessionCache.EXPECT().Put(gomock.Any(), gomock.Any()).AnyTimes().Do(func(key string, _ticket *tls.ClientSessionState) {
		ticket = _ticket
	})
	clientTlsConf := &tls.Config{
		RootCAs:            certPool,
		InsecureSkipVerify: true,
		ClientSessionCache: sessionCache,
	}
	clientTlsConf.NextProtos = protos
	var err error
	listener, err = ListenAddrEarly("localhost:0", serverTlsConf, serverConfig)
	assert.NoError(t, err)

	fetch0RttClient, err := DialAddrEarly(context.Background(), listener.Addr().String(), clientTlsConf, clientConf)
	assert.NoError(t, err)

	fetch0RttServer, err := listener.Accept(context.Background())
	assert.NoError(t, err)

	_ = fetch0RttClient
	_ = fetch0RttServer

	<-tokenSet
	<-ticketSet

	err = fetch0RttClient.CloseWithError(0, "0rtt state fetched")
	assert.NoError(t, err)

	err = fetch0RttServer.CloseWithError(0, "0rtt state fetched")
	assert.NoError(t, err)

	newClient = func() EarlyConnection {
		client, err := DialAddrEarly(context.Background(), listener.Addr().String(), clientTlsConf, clientConf)
		assert.NoError(t, err)
		return client
	}

	return
}

func Benchmark0RTTUpload(b *testing.B) {
	//utils.DefaultLogger.SetLogLevel(utils.LogLevelError)
	//utils.DefaultLogger.SetLogTimeFormat("15:04:05.9999")
	mockCtrl = gomock.NewController(b)
	newClient, serverListener := get0RttConn(b)

	go func() {
		for {
			server, err := serverListener.Accept(context.Background())
			assert.NoError(b, err)
			go func() {
				stream, err := server.AcceptStream(context.Background())
				assert.NoError(b, err)
				_, err = io.ReadAll(stream)
				assert.NoError(b, err)
				err = server.CloseWithError(0, "done")
				assert.NoError(b, err)
			}()
		}
	}()

	var buf [1e5]byte
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client := newClient()
		stream, err := client.OpenStream()
		assert.NoError(b, err)
		_, err = stream.Write(buf[:])
		assert.NoError(b, err)
		err = stream.Close()
		assert.NoError(b, err)
		<-client.Context().Done()
		assert.True(b, client.ConnectionState().TLS.Used0RTT)
	}
}
