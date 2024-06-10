package qstate

import (
	"github.com/stretchr/testify/require"
	"gonum.org/v1/gonum/stat"
	"testing"
)

func BenchmarkEncode(b *testing.B) {
	c := nonDefaultState()
	for codecName, codec := range Codecs {
		b.Run(codecName, func(b *testing.B) {
			buf := make([]byte, 0, 100_000)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var err error
				_, err = codec.Encode(buf[:0], &c)
				require.NoError(b, err)
			}
		})
	}
}

func BenchmarkDecode(b *testing.B) {
	originalState := nonDefaultState()
	decodedState := Connection{}
	for codecName, codec := range Codecs {
		b.Run(codecName, func(b *testing.B) {
			encodedState, err := codec.Encode(nil, &originalState)
			require.NoError(b, err)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				err := codec.Decode(&decodedState, encodedState)
				require.NoError(b, err)
			}
		})
	}
}

func BenchmarkSize(b *testing.B) {
	for codecName, codec := range Codecs {
		b.Run(codecName, func(b *testing.B) {
			sizes := make([]float64, 0, b.N)
			for i := 0; i < b.N; i++ {
				state := nonDefaultState()
				encodedState, err := codec.Encode(nil, &state)
				require.NoError(b, err)
				//tmp, err := os.CreateTemp("", "*.cbor")
				//require.NoError(b, err)
				//_, err = tmp.Write(encodedState)
				//require.NoError(b, err)
				//fmt.Printf("stored encoded state: %s\n", tmp.Name())
				sizes = append(sizes, float64(len(encodedState)))
			}
			mean, stdDev := stat.MeanStdDev(sizes, nil)
			b.ReportMetric(mean, "mean")
			b.ReportMetric(stdDev, "std_dev")
		})
	}
}
