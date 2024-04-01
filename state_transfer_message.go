package quic

import (
	"bytes"
	"github.com/quic-go/quic-go/internal/protocol"
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
	State []byte
}

func (s DataStateTransferMessage) Serialize() ([]byte, error) {
	b := make([]byte, 0, 1)
	b = append(b, TransferMessageTypeState)
	b = append(b, s.State...)
	return b, nil
}

var _ StateTransferMessage = &DataStateTransferMessage{}

type RequestStateTransferMessage struct {
	ConnectionID protocol.ConnectionID
}

func (r RequestStateTransferMessage) Serialize() ([]byte, error) {
	b := make([]byte, 0, 1)
	b = append(b, TransferMessageTypeRequest)
	b = append(b, r.ConnectionID.Bytes()...)
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
		ConnectionID: connID,
	}, nil
}

func parseStateTransferMessage(r *bytes.Reader) (*DataStateTransferMessage, error) {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return &DataStateTransferMessage{
		State: bytes,
	}, nil
}
