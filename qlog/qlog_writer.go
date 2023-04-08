package qlog

import (
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/logging"
	"io"
)

func NewQlogWriter(writer io.WriteCloser, perspective protocol.Perspective, odcid protocol.ConnectionID, config *Config) logging.QlogWriter {
	return NewConnectionTracer(writer, perspective, odcid, config).(*connectionTracer)
}
