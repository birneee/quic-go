package qstate

import (
	"bytes"
	"encoding/json"
	"github.com/klauspost/compress/zstd"
)

type Codec interface {
	Encode(dst []byte, connection *Connection) ([]byte, error)
	Decode(connection *Connection, src []byte) error
}

type StdJsonCodec struct{}

var _ Codec = &StdJsonCodec{}

func (s StdJsonCodec) Encode(dst []byte, connection *Connection) ([]byte, error) {
	buf := bytes.NewBuffer(dst)
	enc := json.NewEncoder(buf)
	err := enc.Encode(connection)
	return buf.Bytes(), err
}

func (s StdJsonCodec) Decode(connection *Connection, src []byte) error {
	return json.Unmarshal(src, connection)
}

type MsgpCodec struct{}

var _ Codec = &MsgpCodec{}

func (m MsgpCodec) Encode(dst []byte, connection *Connection) ([]byte, error) {
	return connection.MarshalMsg(dst)
}

func (m MsgpCodec) Decode(connection *Connection, src []byte) error {
	_, err := connection.UnmarshalMsg(src)
	return err
}

type MsgpZstdCodec struct {
	zstdWriter *zstd.Encoder
	zstdReader *zstd.Decoder
	reusedBuf  [100_000]byte
}

var _ Codec = &MsgpZstdCodec{}

func NewMsgpZstdCodec() Codec {
	writer, err := zstd.NewWriter(nil, zstd.WithEncoderConcurrency(1), zstd.WithEncoderLevel(zstd.SpeedFastest))
	if err != nil {
		panic(err)
	}
	reader, err := zstd.NewReader(nil, zstd.WithDecoderConcurrency(1))
	if err != nil {
		panic(err)
	}
	return &MsgpZstdCodec{
		zstdWriter: writer,
		zstdReader: reader,
	}
}

func (m *MsgpZstdCodec) Encode(dst []byte, connection *Connection) ([]byte, error) {
	msgp, err := connection.MarshalMsg(m.reusedBuf[:0])
	if err != nil {
		return nil, err
	}
	dst = m.zstdWriter.EncodeAll(msgp, dst[:0])
	return dst, nil
}

func (m *MsgpZstdCodec) Decode(connection *Connection, src []byte) error {
	msgp, err := m.zstdReader.DecodeAll(src, m.reusedBuf[:0])
	if err != nil {
		return err
	}
	_, err = connection.UnmarshalMsg(msgp)
	return err
}

type StdJsonZstdCodec struct {
	zstdWriter *zstd.Encoder
	zstdReader *zstd.Decoder
	reusedBuf  [100_000]byte
}

var _ Codec = &StdJsonZstdCodec{}

func NewStdJsonZstdCodec() Codec {
	writer, err := zstd.NewWriter(nil, zstd.WithEncoderConcurrency(1), zstd.WithEncoderLevel(zstd.SpeedFastest))
	if err != nil {
		panic(err)
	}
	reader, err := zstd.NewReader(nil, zstd.WithDecoderConcurrency(1))
	if err != nil {
		panic(err)
	}
	return &StdJsonZstdCodec{
		zstdWriter: writer,
		zstdReader: reader,
	}
}

func (m *StdJsonZstdCodec) Encode(dst []byte, connection *Connection) ([]byte, error) {
	jsonBuf, err := json.Marshal(connection)
	if err != nil {
		return nil, err
	}
	dst = m.zstdWriter.EncodeAll(jsonBuf, dst[:0])
	return dst, nil
}

func (m *StdJsonZstdCodec) Decode(connection *Connection, src []byte) error {
	jsonBuf, err := m.zstdReader.DecodeAll(src, m.reusedBuf[:0])
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonBuf, connection)
}
