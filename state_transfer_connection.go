package quic

import (
	"bytes"
	"context"
	"fmt"
	"github.com/quic-go/quic-go/internal/protocol"
	"io"
	"net"
)

type StateTransferConnection interface {
	ReceiveState() ([]byte, error)
	SendState(state []byte) error
	SendRequest(connectionID protocol.ConnectionID) error
	ReceiveRequest() (protocol.ConnectionID, error)
	CloseWithError(ApplicationErrorCode, string) error
	RemoteAddr() net.Addr
	QuicConn() Connection
	ReceiveMessage() (StateTransferMessage, error)
}

type transferConnection struct {
	quicConn         EarlyConnection
	ctx              context.Context
	ctxCancel        context.CancelFunc
	receivedStates   chan *DataStateTransferMessage
	receivedRequests chan *RequestStateTransferMessage
	closeError       error
}

var _ StateTransferConnection = &transferConnection{}

func DialStateTransfer(transport *Transport, ctx context.Context, addr string, config *StateTransferConfig) (StateTransferConnection, error) {
	config = config.Populate()
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	quicConn, err := transport.DialEarly(ctx, udpAddr, config.TlsConfig, config.QuicConfig)
	if err != nil {
		return nil, err
	}
	return NewStateTransferConnection(quicConn), nil
}

func NewStateTransferConnection(quicConn EarlyConnection) StateTransferConnection {
	c := &transferConnection{
		quicConn:         quicConn,
		receivedStates:   make(chan *DataStateTransferMessage, 1),
		receivedRequests: make(chan *RequestStateTransferMessage, 1),
	}
	c.ctx, c.ctxCancel = context.WithCancel(context.Background())

	go c.runReceiveLoop()
	return c
}

// should run in goroutine
func (c *transferConnection) runReceiveLoop() {
	for {
		stream, err := c.quicConn.AcceptUniStream(c.ctx)
		if err != nil {
			c.CloseWithError(ApplicationErrorCode(0), err.Error())
			break
		}
		b, err := io.ReadAll(stream)
		if err != nil {
			c.CloseWithError(ApplicationErrorCode(0), fmt.Sprintf("failed to read: %v", err))
			break
		}
		msg, err := c.parseMessage(bytes.NewReader(b))
		if err != nil {
			c.CloseWithError(ApplicationErrorCode(0), fmt.Sprintf("failed to parse: %v", err))
			break
		}
		switch msg := msg.(type) {
		case *DataStateTransferMessage:
			select {
			case c.receivedStates <- msg:
			default:
				c.CloseWithError(ApplicationErrorCode(0), "blocked")
				break
			}
		case *RequestStateTransferMessage:
			select {
			case c.receivedRequests <- msg:
			default:
				c.CloseWithError(ApplicationErrorCode(0), "blocked")
				break
			}
		default:
			panic("unexpected")
		}
	}
}

func (c *transferConnection) sendMessage(message StateTransferMessage) error {
	stream, err := c.quicConn.OpenUniStream()
	if err != nil {
		return err
	}
	ss, err := message.Serialize()
	if err != nil {
		return err
	}
	_, err = io.Copy(stream, bytes.NewReader(ss))
	if err != nil {
		return err
	}
	err = stream.Close()
	if err != nil {
		return err
	}
	return nil
}

func (c *transferConnection) SendState(state []byte) error {
	message := &DataStateTransferMessage{
		State: state,
	}
	err := c.sendMessage(message)
	if err != nil {
		return err
	}
	return nil
}

func (c *transferConnection) SendRequest(connectionID protocol.ConnectionID) error {
	message := &RequestStateTransferMessage{
		ConnectionID: connectionID,
	}
	err := c.sendMessage(message)
	if err != nil {
		return err
	}
	return nil
}

func (c *transferConnection) ReceiveRequest() (protocol.ConnectionID, error) {
	select {
	case <-c.ctx.Done():
		return ConnectionID{}, c.closeError
	case message := <-c.receivedRequests:
		return message.ConnectionID, nil
	}
}

func (c *transferConnection) ReceiveState() ([]byte, error) {
	select {
	case <-c.ctx.Done():
		return nil, c.closeError
	case message := <-c.receivedStates:
		return message.State, nil
	}
}

func (c *transferConnection) ReceiveMessage() (StateTransferMessage, error) {
	select {
	case <-c.ctx.Done():
		return nil, c.closeError
	case m := <-c.receivedRequests:
		return m, nil
	case m := <-c.receivedStates:
		return m, nil
	}
}

func (c *transferConnection) parseMessage(reader *bytes.Reader) (StateTransferMessage, error) {
	messageType, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	switch messageType {
	case TransferMessageTypeState:
		return parseStateTransferMessage(reader)
	case TransferMessageTypeRequest:
		return parseRequestTransferMessage(reader)
	default:
		return nil, fmt.Errorf("invalid type")
	}
}

func (c *transferConnection) CloseWithError(err ApplicationErrorCode, s string) error {
	c.quicConn.CloseWithError(err, s)
	c.closeError = fmt.Errorf(s)
	c.ctxCancel()
	return nil
}

func (c *transferConnection) RemoteAddr() net.Addr {
	return c.quicConn.RemoteAddr()
}

func (c *transferConnection) QuicConn() Connection {
	return c.quicConn
}
