package handover

import (
	"bytes"
	"crypto/rand"
	"encoding/gob"
	"encoding/json"
	jsoniter "github.com/json-iterator/go"
	"github.com/klauspost/compress/zstd"
	"github.com/quic-go/quic-go/internal/utils"
	"github.com/quic-go/quic-go/qstate"
	"github.com/stretchr/testify/require"
	"github.com/tinylib/msgp/msgp"
	"testing"
)

func randomByteSlice(length int) []byte {
	b := make([]byte, length)
	rand.Read(b)
	return b
}

func nonDefaultState() qstate.Connection {
	s := qstate.Connection{
		Transport: qstate.Transport{
			Version:      1,
			ChosenALPN:   "proto1",
			VantagePoint: "client",

			DestinationIP:   "127.0.0.1",
			DestinationPort: 6000,
			Parameters: qstate.Parameters{
				ActiveConnectionIDLimit: 4,
			},
			RemoteParameters: qstate.Parameters{
				OriginalDestinationConnectionID: utils.New(randomByteSlice(20)),
			},
			MaxData:       100_000,
			RemoteMaxData: 20_000,
			//TODO complete
		},
		Crypto: qstate.Crypto{
			KeyPhase:                  2,
			TlsCipher:                 "TLS_AES_128_GCM_SHA256",
			HeaderProtectionKey:       randomByteSlice(16),
			RemoteHeaderProtectionKey: randomByteSlice(16),
			TrafficSecret:             randomByteSlice(32),
			RemoteTrafficSecret:       randomByteSlice(32),
		},
		Metrics: qstate.Metrics{},
	}

	// add connection IDs
	for i := 0; i < 4; i++ {
		s.Transport.ConnectionIDs = append(s.Transport.ConnectionIDs, qstate.ConnectionID{
			SequenceNumber:      uint64(i),
			ConnectionID:        randomByteSlice(4),
			StatelessResetToken: (*[16]byte)(randomByteSlice(16)),
		})
	}

	// add remote connection IDs
	for i := 0; i < 4; i++ {
		s.Transport.RemoteConnectionIDs = append(s.Transport.RemoteConnectionIDs, qstate.ConnectionID{
			SequenceNumber:      uint64(i),
			ConnectionID:        randomByteSlice(20),
			StatelessResetToken: (*[16]byte)(randomByteSlice(16)),
		})
	}

	// add streams
	for i := 0; i < 3; i++ {
		s.Transport.Streams = append(s.Transport.Streams, qstate.Stream{
			StreamID:     int64(i * 4),
			WriteMaxData: utils.New(int64(10_000)),
		})
	}

	// add pending stream frames
	for i := 0; i < 20; i++ {
		s.Transport.PendingAcks = append(s.Transport.PendingAcks, qstate.Packet{
			PacketNumber: int64(100 + i),
			Frames: []qstate.Frame{
				{
					Type:     "stream",
					StreamID: utils.New(int64(0)),
					Offset:   utils.New(int64(i * 1000)),
					Length:   utils.New(int64(1000)),
				},
			},
		})
	}
	return s
}

func benchmarkBaseSerialize(b *testing.B, serialize func(*qstate.Connection) ([]byte, error)) {
	s := nonDefaultState()
	var buf []byte
	var err error
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf, err = serialize(&s)
		if err != nil {
			b.Error(err)
		}
	}
	b.StopTimer()
	b.ReportMetric(float64(len(buf)), "bytes")
}

func benchmarkBaseParse(b *testing.B, serialize func(*qstate.Connection) ([]byte, error), parse func([]byte) (qstate.Connection, error)) {
	s := nonDefaultState()
	serialized, err := serialize(&s)
	if err != nil {
		b.Error(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parse(serialized)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkJsonSerialize(b *testing.B) {
	buf := make([]byte, 0, 100_000)
	wr := bytes.NewBuffer(buf)
	en := json.NewEncoder(wr)
	benchmarkBaseSerialize(b,
		func(s *qstate.Connection) ([]byte, error) {
			wr.Reset()
			err := en.Encode(s)
			return wr.Bytes(), err
		},
	)
}

func BenchmarkJsonParse(b *testing.B) {
	s := qstate.Connection{}
	benchmarkBaseParse(b,
		func(s *qstate.Connection) ([]byte, error) {
			return json.Marshal(s)
		},
		func(buf []byte) (qstate.Connection, error) {
			err := json.Unmarshal(buf, &s)
			return s, err
		},
	)
}

func BenchmarkJsoniterSerialize(b *testing.B) {
	buf := make([]byte, 0, 100_000)
	wr := bytes.NewBuffer(buf)
	en := jsoniter.NewEncoder(wr)
	benchmarkBaseSerialize(b,
		func(s *qstate.Connection) ([]byte, error) {
			wr.Reset()
			err := en.Encode(s)
			return wr.Bytes(), err
		},
	)
}

func BenchmarkJsoniterParse(b *testing.B) {
	s := qstate.Connection{}
	benchmarkBaseParse(b,
		func(s *qstate.Connection) ([]byte, error) {
			return jsoniter.Marshal(s)
		},
		func(buf []byte) (qstate.Connection, error) {
			err := jsoniter.Unmarshal(buf, &s)
			return s, err
		},
	)
}

func BenchmarkGobSerialize(b *testing.B) {
	buf := make([]byte, 0, 100_000)
	wr := bytes.NewBuffer(buf)
	en := gob.NewEncoder(wr)
	benchmarkBaseSerialize(b,
		func(s *qstate.Connection) ([]byte, error) {
			wr.Reset()
			err := en.Encode(s)
			return wr.Bytes(), err
		},
	)
}

func BenchmarkGobParse(b *testing.B) {
	s := qstate.Connection{}
	benchmarkBaseParse(b,
		func(s *qstate.Connection) ([]byte, error) {
			buf := bytes.NewBuffer(nil)
			encoder := gob.NewEncoder(buf)
			err := encoder.Encode(s)
			return buf.Bytes(), err
		},
		func(buf []byte) (qstate.Connection, error) {
			decoder := gob.NewDecoder(bytes.NewReader(buf))
			err := decoder.Decode(&s)
			return s, err
		},
	)
}

func BenchmarkMsgpSerialize(b *testing.B) {
	buf := make([]byte, 0, 100_000)
	benchmarkBaseSerialize(b,
		func(s *qstate.Connection) ([]byte, error) {
			buf, err := s.MarshalMsg(buf[:0])
			return buf, err
		},
	)
}

func BenchmarkMsgpParse(b *testing.B) {
	s := qstate.Connection{}
	benchmarkBaseParse(b,
		func(s *qstate.Connection) ([]byte, error) {
			buf, err := s.MarshalMsg(nil)
			return buf, err
		},
		func(buf []byte) (qstate.Connection, error) {
			_, err := s.UnmarshalMsg(buf)
			return s, err
		},
	)
}

func BenchmarkMsgpJsonSerialize(b *testing.B) {
	msgpBuf := make([]byte, 0, 100_000)
	jsonBuf := bytes.NewBuffer(make([]byte, 0, 100_000))
	benchmarkBaseSerialize(b,
		func(s *qstate.Connection) ([]byte, error) {
			msgpBuf, err := s.MarshalMsg(msgpBuf[:0])
			jsonBuf.Reset()
			_, err = msgp.UnmarshalAsJSON(jsonBuf, msgpBuf)
			return jsonBuf.Bytes(), err
		},
	)
}

func BenchmarkMsgpZstdSerialize(b *testing.B) {
	msgpBuf := make([]byte, 0, 100_000)
	zstdBuf := make([]byte, 0, 100_000)
	zstdWriter, err := zstd.NewWriter(nil, zstd.WithEncoderConcurrency(1), zstd.WithEncoderLevel(zstd.SpeedFastest))
	require.NoError(b, err)
	benchmarkBaseSerialize(b,
		func(s *qstate.Connection) ([]byte, error) {
			msgpBuf, err = s.MarshalMsg(msgpBuf[:0])
			zstdBuf = zstdWriter.EncodeAll(msgpBuf, zstdBuf[:0])
			return zstdBuf, err
		},
	)
}

func BenchmarkMsgpZstdParse(b *testing.B) {
	msgpBuf := make([]byte, 0, 100_000)
	zstdReader, err := zstd.NewReader(nil, zstd.WithDecoderConcurrency(1))
	require.NoError(b, err)
	s := qstate.Connection{}
	benchmarkBaseParse(b,
		func(s *qstate.Connection) ([]byte, error) {
			zstdWriter, err := zstd.NewWriter(nil, zstd.WithEncoderConcurrency(1), zstd.WithEncoderLevel(zstd.SpeedFastest))
			require.NoError(b, err)
			msgpBuf, err := s.MarshalMsg(nil)
			zstdBuf := zstdWriter.EncodeAll(msgpBuf, nil)
			return zstdBuf, err
		},
		func(b []byte) (qstate.Connection, error) {
			msgpBuf, err = zstdReader.DecodeAll(b, msgpBuf[:0])
			_, err = s.UnmarshalMsg(msgpBuf)
			return s, err
		},
	)
}

func BenchmarkJsoniterZstdSerialize(b *testing.B) {
	jsonBuf := bytes.NewBuffer(make([]byte, 0, 100_000))
	jsonEncoder := jsoniter.NewEncoder(jsonBuf)
	zstdBuf := make([]byte, 0, 100_000)
	zstdWriter, err := zstd.NewWriter(nil, zstd.WithEncoderConcurrency(1), zstd.WithEncoderLevel(zstd.SpeedFastest))
	require.NoError(b, err)
	benchmarkBaseSerialize(b,
		func(s *qstate.Connection) ([]byte, error) {
			jsonBuf.Reset()
			err := jsonEncoder.Encode(s)
			zstdBuf := zstdWriter.EncodeAll(jsonBuf.Bytes(), zstdBuf[:0])
			return zstdBuf, err
		},
	)
}

func BenchmarkJsoniterZstdParse(b *testing.B) {
	jsonEncodeBuf := bytes.NewBuffer(make([]byte, 0, 100_000))
	jsonDecodeBuf := make([]byte, 0, 100_000)
	jsonEncoder := jsoniter.NewEncoder(jsonEncodeBuf)
	zstdWriter, err := zstd.NewWriter(nil, zstd.WithEncoderConcurrency(1), zstd.WithEncoderLevel(zstd.SpeedFastest))
	require.NoError(b, err)
	zstdReader, err := zstd.NewReader(nil, zstd.WithDecoderConcurrency(1))
	require.NoError(b, err)
	s := qstate.Connection{}
	benchmarkBaseParse(b,
		func(s *qstate.Connection) ([]byte, error) {
			jsonEncodeBuf.Reset()
			err := jsonEncoder.Encode(s)
			zstdBuf := zstdWriter.EncodeAll(jsonEncodeBuf.Bytes(), nil)
			return zstdBuf, err
		},
		func(zstdBuf []byte) (qstate.Connection, error) {
			jsonDecodeBuf, err := zstdReader.DecodeAll(zstdBuf, jsonDecodeBuf[:0])
			err = jsoniter.Unmarshal(jsonDecodeBuf, &s)
			return s, err
		},
	)
}
