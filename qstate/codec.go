package qstate

import (
	"bytes"
	"encoding/json"
	"github.com/fxamacker/cbor/v2"
	"github.com/klauspost/compress/zstd"
	"github.com/mailru/easyjson"
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

type EasyJsonCodec[T easyjson.MarshalerUnmarshaler] struct{}

var _ Codec[*Connection] = &EasyJsonCodec[*Connection]{}

func (e EasyJsonCodec[T]) Encode(dst []byte, connection T) ([]byte, error) {
	w := bytes.NewBuffer(dst)
	_, err := easyjson.MarshalToWriter(connection, w)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (e EasyJsonCodec[T]) Decode(connection T, src []byte) error {
	err := easyjson.Unmarshal(src, connection)
	if err != nil {
		return err
	}
	return nil
}

type CborCodec[T any] struct{}

func (c CborCodec[T]) Encode(dst []byte, connection T) ([]byte, error) {
	return cbor.Marshal(connection)
}

func (c CborCodec[T]) Decode(connection T, src []byte) error {
	return cbor.Unmarshal(src, connection)
}

var _ Codec[*Connection] = &CborCodec[*Connection]{}

type zstdCodec[T any] struct {
	inner      Codec[T]
	zstdWriter *zstd.Encoder
	zstdReader *zstd.Decoder
	reusedBuf  [100_000]byte
}

func NewZstdCodec[T any](inner Codec[T]) Codec[T] {
	zstdWriter, err := zstd.NewWriter(nil, zstd.WithEncoderConcurrency(1), zstd.WithEncoderLevel(zstd.SpeedFastest))
	if err != nil {
		panic(err)
	}
	zstdReader, err := zstd.NewReader(nil, zstd.WithDecoderConcurrency(1))
	if err != nil {
		panic(err)
	}
	return &zstdCodec[T]{
		inner:      inner,
		zstdReader: zstdReader,
		zstdWriter: zstdWriter,
	}
}

var _ Codec[*Connection] = &zstdCodec[*Connection]{}

func (z *zstdCodec[T]) Encode(dst []byte, connection T) ([]byte, error) {
	plain, _ := z.inner.Encode(z.reusedBuf[:0], connection)
	dst = z.zstdWriter.EncodeAll(plain, dst[:0])
	return dst, nil
}

func (z *zstdCodec[T]) Decode(connection T, src []byte) error {
	plain, err := z.zstdReader.DecodeAll(src, z.reusedBuf[:0])
	if err != nil {
		return err
	}
	return z.inner.Decode(connection, plain)
}
