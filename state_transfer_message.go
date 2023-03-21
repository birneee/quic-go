package quic

import (
	"bytes"
	"github.com/lucas-clemente/quic-go/handover"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"io"
)

type StateTransferMessageType = uint8

const TransferMessageTypeInvalid = 0
const TransferMessageTypeState = 1
const TransferMessageTypeRequest = 2

type StateTransferMessage interface {
	Serialize() ([]byte, error)
}

type DataStateTransferMessage struct {
	state *handover.State
}

func (s DataStateTransferMessage) Serialize() ([]byte, error) {
	b := make([]byte, 0, 1)
	b = append(b, TransferMessageTypeState)
	serializedState, err := s.state.Serialize()
	if err != nil {
		return nil, err
	}
	b = append(b, serializedState...)
	return b, nil
}

var _ StateTransferMessage = &DataStateTransferMessage{}

type RequestStateTransferMessage struct {
	connectionID protocol.ConnectionID
}

func (r RequestStateTransferMessage) Serialize() ([]byte, error) {
	b := make([]byte, 0, 1)
	b = append(b, TransferMessageTypeRequest)
	b = append(b, r.connectionID.Bytes()...)
	return b, nil
}

var _ StateTransferMessage = &RequestStateTransferMessage{}

func parseRequestTransferMessage(r *bytes.Reader) (*RequestStateTransferMessage, error) {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	connID := protocol.ParseConnectionID(bytes)
	return &RequestStateTransferMessage{
		connectionID: connID,
	}, nil
}

func parseStateTransferMessage(r *bytes.Reader) (*DataStateTransferMessage, error) {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	state, err := handover.Parse(bytes)
	if err != nil {
		return nil, err
	}
	return &DataStateTransferMessage{
		state: state,
	}, nil
}
