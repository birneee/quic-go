package xse

import (
	"bytes"
	"crypto/tls"
	"github.com/golang/mock/gomock"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/qtls"
	"io"
	"math/rand"
	"reflect"
	"testing"
)

func mockSendReceiveQueue(ctrl *gomock.Controller) (SendStream, ReceiveStream) {
	var buf bytes.Buffer
	streamId := protocol.StreamID(4)
	sendStream := NewMockSendStream(ctrl)
	sendStream.EXPECT().StreamID().AnyTimes().Return(streamId)
	sendStream.EXPECT().Write(gomock.Any()).AnyTimes().DoAndReturn(buf.Write)
	receiveStream := NewMockReceiveStream(ctrl)
	receiveStream.EXPECT().StreamID().AnyTimes().Return(streamId)
	receiveStream.EXPECT().Read(gomock.Any()).AnyTimes().DoAndReturn(buf.Read)
	return sendStream, receiveStream
}

func mockCryptoSetup() (client CryptoSetup, server CryptoSetup) {
	suite := qtls.CipherSuiteTLS13ByID(tls.TLS_AES_128_GCM_SHA256)
	secret := make([]byte, suite.Hash.Size())
	rand.Read(secret)
	client = NewCryptoSetup(suite, secret, protocol.PerspectiveClient, nil)
	server = NewCryptoSetup(suite, secret, protocol.PerspectiveServer, nil)
	return
}

func TestShortMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	baseSendStream, baseReceiveStream := mockSendReceiveQueue(ctrl)
	clientCryptoSetup, serverCryptoSetup := mockCryptoSetup()
	sendStream := serverCryptoSetup.NewSendStream(baseSendStream)
	receiveStream := clientCryptoSetup.NewReceiveStream(baseReceiveStream)
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
	clientCryptoSetup, serverCryptoSetup := mockCryptoSetup()
	sendStream := serverCryptoSetup.NewSendStream(baseSendStream)
	receiveStream := clientCryptoSetup.NewReceiveStream(baseReceiveStream)
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
