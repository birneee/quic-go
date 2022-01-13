package test

import (
	"bytes"
	"crypto/tls"
	"github.com/golang/mock/gomock"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/qtls"
	"github.com/lucas-clemente/quic-go/internal/xse"
	"github.com/lucas-clemente/quic-go/internal/xse/mocks"
	"io"
	"math/rand"
	"reflect"
	"testing"
)

func mockSendReceiveQueue(ctrl *gomock.Controller) (xse.SendStream, xse.ReceiveStream) {
	var buf bytes.Buffer
	streamId := protocol.StreamID(0)
	sendStream := mocks.NewMockSendStream(ctrl)
	sendStream.EXPECT().StreamID().AnyTimes().Return(streamId)
	sendStream.EXPECT().Write(gomock.Any()).AnyTimes().DoAndReturn(buf.Write)
	receiveStream := mocks.NewMockReceiveStream(ctrl)
	receiveStream.EXPECT().StreamID().AnyTimes().Return(streamId)
	receiveStream.EXPECT().Read(gomock.Any()).AnyTimes().DoAndReturn(buf.Read)
	return sendStream, receiveStream
}

func mockCryptoSetup() xse.CryptoSetup {
	suite := qtls.CipherSuiteTLS13ByID(tls.TLS_AES_128_GCM_SHA256)
	secret := make([]byte, suite.Hash.Size())
	rand.Read(secret)
	return xse.NewCryptoSetup(secret, secret, suite)
}

func TestShortMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	baseSendStream, baseReceiveStream := mockSendReceiveQueue(ctrl)
	cryptoSetup := mockCryptoSetup()
	sendStream := xse.NewSendStream(baseSendStream, cryptoSetup)
	receiveStream := xse.NewReceiveStream(baseReceiveStream, cryptoSetup)
	sendBuf := []byte("hello")
	_, err := sendStream.Write(sendBuf)
	if err != nil {
		t.Errorf(err.Error())
	}
	receiveBuf, err := io.ReadAll(receiveStream)
	if err != nil {
		t.Errorf(err.Error())
	}
	if !reflect.DeepEqual(sendBuf, receiveBuf) {
		t.Errorf("got %q, wanted %q", receiveBuf, sendBuf)
	}
}

func TestLongMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	baseSendStream, baseReceiveStream := mockSendReceiveQueue(ctrl)
	cryptoSetup := mockCryptoSetup()
	sendStream := xse.NewSendStream(baseSendStream, cryptoSetup)
	receiveStream := xse.NewReceiveStream(baseReceiveStream, cryptoSetup)
	sendBuf := make([]byte, 1e6)
	rand.Read(sendBuf)
	_, err := sendStream.Write(sendBuf)
	if err != nil {
		t.Errorf(err.Error())
	}
	receiveBuf, err := io.ReadAll(receiveStream)
	if err != nil {
		t.Errorf(err.Error())
	}
	if !reflect.DeepEqual(sendBuf, receiveBuf) {
		t.Errorf("sent and received data are different")
	}
}
