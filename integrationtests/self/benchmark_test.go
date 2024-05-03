package self_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/integrationtests/tools"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
)

func BenchmarkHandshakeTCP(b *testing.B) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		b.Fatal(err)
	}
	defer ln.Close()

	connChan := make(chan net.Conn, 1)
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			connChan <- conn
			require.NoError(b, err)
		}
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c, err := net.DialTCP("tcp", nil, ln.Addr().(*net.TCPAddr))
		if err != nil {
			b.Fatal(err)
		}
		<-connChan
		c.Close()
	}
}

func BenchmarkHandshakeTLS12(b *testing.B) {
	tlsConfig := getTLSConfig()
	tlsClientConfig := getTLSClientConfig()
	tlsConfig.MinVersion = tls.VersionTLS12
	tlsConfig.MaxVersion = tls.VersionTLS12
	tlsClientConfig.MinVersion = tls.VersionTLS12
	tlsClientConfig.MaxVersion = tls.VersionTLS12

	ln, err := tls.Listen("tcp", "127.0.0.1:8080", tlsConfig)
	require.NoError(b, err)
	defer ln.Close()

	connChan := make(chan net.Conn, 1)
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			connChan <- c
			err = c.(*tls.Conn).Handshake()
			require.NoError(b, err)
		}
	}()

	addr := ln.Addr().String()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c, err := tls.Dial("tcp", addr, tlsClientConfig)
		require.NoError(b, err)
		err = c.Handshake()
		require.NoError(b, err)
		<-connChan
		c.Close()
	}
}

func BenchmarkHandshakeTLS13(b *testing.B) {
	tlsConfig := getTLSConfig()
	tlsClientConfig := getTLSClientConfig()
	tlsConfig.MinVersion = tls.VersionTLS13
	tlsConfig.MaxVersion = tls.VersionTLS13
	tlsClientConfig.MinVersion = tls.VersionTLS13
	tlsClientConfig.MaxVersion = tls.VersionTLS13

	ln, err := tls.Listen("tcp", "127.0.0.1:8080", tlsConfig)
	require.NoError(b, err)
	defer ln.Close()

	connChan := make(chan net.Conn, 1)
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			connChan <- c
			err = c.(*tls.Conn).Handshake()
			require.NoError(b, err)
		}
	}()

	addr := ln.Addr().String()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c, err := tls.Dial("tcp", addr, tlsClientConfig)
		require.NoError(b, err)
		err = c.Handshake()
		require.NoError(b, err)
		<-connChan
		c.Close()
	}
}

func BenchmarkHandshake0RTT(b *testing.B) {
	//b.ReportAllocs()
	tokenStore := tools.NewSingleTokenStore()
	sessionCache := tools.NewSingleSessionCache()
	tlsClientConfig := getTLSClientConfig()
	tlsClientConfig.ClientSessionCache = sessionCache
	tlsConfig := getTLSConfig()
	ln, err := quic.ListenAddrEarly("localhost:0", tlsConfig, &quic.Config{Allow0RTT: true})
	require.NoError(b, err)
	defer ln.Close()

	go func() {
		for {
			_, err := ln.Accept(context.Background())
			if err != nil {
				return
			}
		}
	}()

	quicClientConfig := &quic.Config{
		TokenStore: tokenStore,
	}

	// fetch 0-rtt info
	{
		_, err = quic.DialAddr(context.Background(), ln.Addr().String(), tlsClientConfig, quicClientConfig)
		require.NoError(b, err)
		sessionCache.Await()
		tokenStore.Await()
	}

	addr := ln.Addr().String()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c, err := quic.DialAddrEarly(context.Background(), addr, tlsClientConfig, quicClientConfig)
		require.NoError(b, err)
		c.CloseWithError(0, "")
	}
}

func BenchmarkHandshake(b *testing.B) {
	//b.ReportAllocs()

	ln, err := quic.ListenAddrEarly("localhost:0", tlsConfig, nil)
	if err != nil {
		b.Fatal(err)
	}
	defer ln.Close()

	connChan := make(chan quic.Connection, 1)
	go func() {
		for {
			conn, err := ln.Accept(context.Background())
			if err != nil {
				return
			}
			connChan <- conn
		}
	}()

	addr := ln.Addr().String()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c, err := quic.DialAddr(context.Background(), addr, tlsClientConfig, nil)
		if err != nil {
			b.Fatal(err)
		}
		<-connChan
		c.CloseWithError(0, "")
	}
}

func BenchmarkFirstSentByte(b *testing.B) {
	b.ReportAllocs()

	ln, err := quic.ListenAddrEarly("localhost:0", tlsConfig, nil)
	if err != nil {
		b.Fatal(err)
	}
	defer ln.Close()

	connChan := make(chan quic.Connection, 1)
	go func() {
		for {
			conn, err := ln.Accept(context.Background())
			if err != nil {
				return
			}
			stream, err := conn.AcceptStream(context.Background())
			if err != nil {
				return
			}
			var b [1]byte
			_, err = stream.Read(b[:])
			if err != nil {
				return
			}
			connChan <- conn
		}
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c, err := quic.DialAddr(context.Background(), ln.Addr().String(), tlsClientConfig, nil)
		if err != nil {
			b.Fatal(err)
		}
		s, err := c.OpenStream()
		if err != nil {
			b.Fatal(err)
		}
		_, err = s.Write([]byte{'a'})
		if err != nil {
			b.Fatal(err)
		}
		<-connChan
		c.CloseWithError(0, "")
	}
}

func BenchmarkStreamChurn(b *testing.B) {
	b.ReportAllocs()

	ln, err := quic.ListenAddr("localhost:0", tlsConfig, &quic.Config{MaxIncomingStreams: 1e10})
	if err != nil {
		b.Fatal(err)
	}
	defer ln.Close()

	errChan := make(chan error, 1)
	go func() {
		conn, err := ln.Accept(context.Background())
		if err != nil {
			errChan <- err
			return
		}
		close(errChan)
		for {
			str, err := conn.AcceptStream(context.Background())
			if err != nil {
				return
			}
			str.Close()
		}
	}()

	c, err := quic.DialAddr(context.Background(), fmt.Sprintf("localhost:%d", ln.Addr().(*net.UDPAddr).Port), tlsClientConfig, nil)
	if err != nil {
		b.Fatal(err)
	}
	if err := <-errChan; err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		str, err := c.OpenStreamSync(context.Background())
		if err != nil {
			b.Fatal(err)
		}
		if err := str.Close(); err != nil {
			b.Fatal(err)
		}
	}
}
