package qstate

import (
	"bytes"
	"encoding/json"
	"github.com/klauspost/compress/zstd"
	"github.com/tinylib/msgp/msgp"
)

type Codec[T any] interface {
	Encode(dst []byte, connection T) ([]byte, error)
	Decode(connection T, src []byte) error
}

type StdJsonCodec[T any] struct{}

var _ Codec[*Connection] = &StdJsonCodec[*Connection]{}

func (s StdJsonCodec[T]) Encode(dst []byte, connection T) ([]byte, error) {
	buf := bytes.NewBuffer(dst)
	enc := json.NewEncoder(buf)
	err := enc.Encode(connection)
	return buf.Bytes(), err
}

func (s StdJsonCodec[T]) Decode(connection T, src []byte) error {
	return json.Unmarshal(src, connection)
}

type MarshlerAndUnmarshler interface {
	msgp.Marshaler
	msgp.Unmarshaler
}

type MsgpCodec[T MarshlerAndUnmarshler] struct{}

var _ Codec[*Connection] = &MsgpCodec[*Connection]{}

func (m MsgpCodec[T]) Encode(dst []byte, connection T) ([]byte, error) {
	return connection.MarshalMsg(dst)
}

func (m MsgpCodec[T]) Decode(connection T, src []byte) error {
	_, err := connection.UnmarshalMsg(src)
	return err
}

type MsgpZstdCodec[T MarshlerAndUnmarshler] struct {
	zstdWriter *zstd.Encoder
	zstdReader *zstd.Decoder
	reusedBuf  [100_000]byte
}

var _ Codec[*Connection] = &MsgpZstdCodec[*Connection]{}

func NewMsgpZstdCodec[T MarshlerAndUnmarshler]() Codec[T] {
	writer, err := zstd.NewWriter(nil, zstd.WithEncoderConcurrency(1), zstd.WithEncoderLevel(zstd.SpeedFastest))
	if err != nil {
		panic(err)
	}
	reader, err := zstd.NewReader(nil, zstd.WithDecoderConcurrency(1))
	if err != nil {
		panic(err)
	}
	return &MsgpZstdCodec[T]{
		zstdWriter: writer,
		zstdReader: reader,
	}
}

func (m *MsgpZstdCodec[T]) Encode(dst []byte, connection T) ([]byte, error) {
	msgp, err := connection.MarshalMsg(m.reusedBuf[:0])
	if err != nil {
		return nil, err
	}
	dst = m.zstdWriter.EncodeAll(msgp, dst[:0])
	return dst, nil
}

func (m *MsgpZstdCodec[T]) Decode(connection T, src []byte) error {
	msgp, err := m.zstdReader.DecodeAll(src, m.reusedBuf[:0])
	if err != nil {
		return err
	}
	_, err = connection.UnmarshalMsg(msgp)
	return err
}

type StdJsonZstdCodec[T any] struct {
	zstdWriter *zstd.Encoder
	zstdReader *zstd.Decoder
	reusedBuf  [100_000]byte
}

var _ Codec[*Connection] = &StdJsonZstdCodec[*Connection]{}

func NewStdJsonZstdCodec[T any]() Codec[T] {
	writer, err := zstd.NewWriter(nil, zstd.WithEncoderConcurrency(1), zstd.WithEncoderLevel(zstd.SpeedFastest))
	if err != nil {
		panic(err)
	}
	reader, err := zstd.NewReader(nil, zstd.WithDecoderConcurrency(1))
	if err != nil {
		panic(err)
	}
	return &StdJsonZstdCodec[T]{
		zstdWriter: writer,
		zstdReader: reader,
	}
}

func (m *StdJsonZstdCodec[T]) Encode(dst []byte, connection T) ([]byte, error) {
	jsonBuf, err := json.Marshal(connection)
	if err != nil {
		return nil, err
	}
	dst = m.zstdWriter.EncodeAll(jsonBuf, dst[:0])
	return dst, nil
}

func (m *StdJsonZstdCodec[T]) Decode(connection T, src []byte) error {
	jsonBuf, err := m.zstdReader.DecodeAll(src, m.reusedBuf[:0])
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonBuf, connection)
}
