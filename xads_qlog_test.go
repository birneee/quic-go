package quic

import (
	"context"
	"crypto/tls"
	"github.com/golang/mock/gomock"
	mocklogging "github.com/quic-go/quic-go/internal/mocks/logging"
	"github.com/quic-go/quic-go/internal/testdata"
	"github.com/quic-go/quic-go/logging"
	"io"
	"testing"
)

func xadsConnection(ctx context.Context, tracer func(context.Context, logging.Perspective, ConnectionID) logging.ConnectionTracer) (Connection, Connection, error) {
	serverListener, err := ListenAddr("127.0.0.0:0", &tls.Config{NextProtos: []string{"test"}, Certificates: testdata.GetTLSConfig().Certificates}, &Config{EnableDatagrams: true, Experimental: ExperimentalConfig{ExtraApplicationDataSecurity: PreferExtraApplicationDataSecurity}, Tracer: tracer})
	if err != nil {
		return nil, nil, err
	}
	client, err := DialAddr(ctx, serverListener.Addr().String(), &tls.Config{NextProtos: []string{"test"}, InsecureSkipVerify: true}, &Config{EnableDatagrams: true, Experimental: ExperimentalConfig{ExtraApplicationDataSecurity: PreferExtraApplicationDataSecurity}, Tracer: tracer})
	if err != nil {
		return nil, nil, err
	}
	server, err := serverListener.Accept(ctx)
	if err != nil {
		return nil, nil, err
	}
	return client, server, nil
}

func connectionTracerWithoutXads(ctrl *gomock.Controller) *mocklogging.MockConnectionTracer {
	t := mocklogging.NewMockConnectionTracer(ctrl)
	t.EXPECT().StartedConnection(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	t.EXPECT().UpdatedCongestionState(gomock.Any()).AnyTimes()
	t.EXPECT().SentTransportParameters(gomock.Any()).AnyTimes()
	t.EXPECT().UpdatedKeyFromTLS(gomock.Any(), gomock.Any()).AnyTimes()
	t.EXPECT().SentLongHeaderPacket(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	t.EXPECT().UpdatedMetrics(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	t.EXPECT().SetLossTimer(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	t.EXPECT().NegotiatedVersion(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	t.EXPECT().ReceivedLongHeaderPacket(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	t.EXPECT().AcknowledgedPacket(gomock.Any(), gomock.Any()).AnyTimes()
	t.EXPECT().ReceivedTransportParameters(gomock.Any()).AnyTimes()
	t.EXPECT().DroppedEncryptionLevel(gomock.Any()).AnyTimes()
	t.EXPECT().ReceivedShortHeaderPacket(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	t.EXPECT().SentShortHeaderPacket(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	t.EXPECT().LossTimerCanceled().AnyTimes()
	t.EXPECT().ClosedConnection(gomock.Any()).AnyTimes()
	t.EXPECT().Close().AnyTimes()
	return t
}

func TestQlog(t *testing.T) {
	const expectedOverhead = 22
	ctrl := gomock.NewController(t)
	msg := "hello"
	connectionTracer := connectionTracerWithoutXads(ctrl)
	connectionTracer.EXPECT().XadsReceiveRecord(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Do(func(streamID StreamID, rawLength int, dataLength int) {
		if dataLength != len(msg) {
			t.Error("wrong payload size")
		}
		if rawLength != len(msg)+expectedOverhead {
			t.Error("wrong record size")
		}
	})
	tracer := func(_ context.Context, _ logging.Perspective, _ ConnectionID) logging.ConnectionTracer {
		return connectionTracer
	}
	client, server, err := xadsConnection(context.Background(), tracer)
	if err != nil {
		t.Error(err)
	}
	if !client.ExtraApplicationDataSecurity() {
		t.Errorf("XADS negotiation failed")
	}
	if !server.ExtraApplicationDataSecurity() {
		t.Errorf("XADS negotiation failed")
	}
	stream, err := client.OpenStream()
	if err != nil {
		t.Error(err)
	}
	_, err = stream.Write([]byte(msg))
	if err != nil {
		t.Error(err)
	}
	serverStream, err := server.AcceptStream(context.Background())
	if err != nil {
		t.Error(err)
	}
	_, err = io.CopyN(io.Discard, serverStream, int64(len(msg)))
	if err != nil {
		t.Error(err)
	}
	_ = client.CloseWithError(ApplicationErrorCode(0), "done")
	_ = server.CloseWithError(ApplicationErrorCode(0), "done")
}
