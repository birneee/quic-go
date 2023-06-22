package self_test

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/quic-go/quic-go"
)

func BenchmarkHandshake(b *testing.B) {
	b.ReportAllocs()

	ln, err := quic.ListenAddr("localhost:0", tlsConfig, nil)
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c, err := quic.DialAddr(context.Background(), ln.Addr().String(), tlsClientConfig, nil)
		if err != nil {
			b.Fatal(err)
		}
		<-connChan
		c.CloseWithError(0, "")
	}
}

func BenchmarkFirstSentByte(b *testing.B) {
	b.ReportAllocs()

	ln, err := quic.ListenAddr("localhost:0", tlsConfig, nil)
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
