package handover

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	jsoniter "github.com/json-iterator/go"
	"github.com/klauspost/compress/zstd"
	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/stretchr/testify/require"
	"github.com/tinylib/msgp/msgp"
	"testing"
)

func nonDefaultState() State {
	s := State{
		ClientConnectionIDs: map[ConnectionIDSequenceNumber]*ConnectionIDWithResetToken{
			0: {
				ConnectionID:        []byte{1, 2, 3},
				StatelessResetToken: []byte{4, 5, 6},
			},
		},
		Version:                   1,
		ServerHeaderProtectionKey: []byte{1, 2, 3},
		BidiStreams: map[protocol.StreamID]*BidiStreamState{
			0: {
				ServerDirectionMaxData: 10_000,
			},
		},
	}
	// append pending stream frames
	for i := 0; i < 20; i++ {
		s.ClientAckPending = append(s.ClientAckPending, PacketState{
			PacketNumber: int64(100 + i),
			Frames: []Frame{
				{Type: "stream", StreamID: 0, Offset: protocol.ByteCount(i * 1000), Length: 1000},
			},
		})
	}
	return s
}

func benchmarkBaseSerialize(b *testing.B, serialize func(State) ([]byte, error)) {
	s := nonDefaultState()
	var buf []byte
	var err error
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf, err = serialize(s)
		if err != nil {
			b.Error(err)
		}
	}
	b.StopTimer()
	b.ReportMetric(float64(len(buf)), "bytes")
}

func benchmarkBaseParse(b *testing.B, serialize func(State) ([]byte, error), parse func([]byte) (State, error)) {
	s := nonDefaultState()
	serialized, err := serialize(s)
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
	benchmarkBaseSerialize(b,
		func(s State) ([]byte, error) {
			return json.Marshal(s)
		},
	)
}

func BenchmarkJsonParse(b *testing.B) {
	benchmarkBaseParse(b,
		func(s State) ([]byte, error) {
			return json.Marshal(s)
		},
		func(buf []byte) (State, error) {
			s := State{}
			err := json.Unmarshal(buf, &s)
			return s, err
		},
	)
}

func BenchmarkJsoniterSerialize(b *testing.B) {
	benchmarkBaseSerialize(b,
		func(s State) ([]byte, error) {
			return jsoniter.Marshal(s)
		},
	)
}

func BenchmarkJsoniterParse(b *testing.B) {
	benchmarkBaseParse(b,
		func(s State) ([]byte, error) {
			return jsoniter.Marshal(s)
		},
		func(buf []byte) (State, error) {
			s := State{}
			err := jsoniter.Unmarshal(buf, &s)
			return s, err
		},
	)
}

func BenchmarkGobSerialize(b *testing.B) {
	benchmarkBaseSerialize(b,
		func(s State) ([]byte, error) {
			buf := bytes.NewBuffer(nil)
			encoder := gob.NewEncoder(buf)
			err := encoder.Encode(s)
			return buf.Bytes(), err
		},
	)
}

func BenchmarkGobParse(b *testing.B) {
	benchmarkBaseParse(b,
		func(s State) ([]byte, error) {
			buf := bytes.NewBuffer(nil)
			encoder := gob.NewEncoder(buf)
			err := encoder.Encode(s)
			return buf.Bytes(), err
		},
		func(buf []byte) (State, error) {
			decoder := gob.NewDecoder(bytes.NewReader(buf))
			s := State{}
			err := decoder.Decode(&s)
			return s, err
		},
	)
}

func BenchmarkMsgpSerialize(b *testing.B) {
	buf := make([]byte, 0, 100_000)
	benchmarkBaseSerialize(b,
		func(s State) ([]byte, error) {
			buf, err := s.MarshalMsg(buf[:0])
			return buf, err
		},
	)
}

func BenchmarkMsgpParse(b *testing.B) {
	benchmarkBaseParse(b,
		func(s State) ([]byte, error) {
			buf, err := s.MarshalMsg(nil)
			return buf, err
		},
		func(buf []byte) (State, error) {
			s := State{}
			_, err := s.UnmarshalMsg(buf)
			return s, err
		},
	)
}

func BenchmarkMsgpJsonSerialize(b *testing.B) {
	msgpBuf := make([]byte, 0, 100_000)
	jsonBuf := bytes.NewBuffer(make([]byte, 0, 100_000))
	benchmarkBaseSerialize(b,
		func(s State) ([]byte, error) {
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
		func(s State) ([]byte, error) {
			msgpBuf, err = s.MarshalMsg(msgpBuf[:0])
			zstdBuf = zstdWriter.EncodeAll(msgpBuf, zstdBuf[:0])
			return zstdBuf, err
		},
	)
}
