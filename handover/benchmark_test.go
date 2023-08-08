package handover

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	jsoniter "github.com/json-iterator/go"
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
	}
	return s
}

func benchmarkBaseSerialize(b *testing.B, serialize func(State) ([]byte, error)) {
	s := nonDefaultState()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := serialize(s)
		if err != nil {
			b.Error(err)
		}
	}
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
	benchmarkBaseSerialize(b,
		func(s State) ([]byte, error) {
			buf, err := s.MarshalMsg(nil)
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
