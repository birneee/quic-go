package qlog

import (
	"bufio"
	"fmt"
	"github.com/lucas-clemente/quic-go/logging"
	"io"
	"os"
	"strings"
)

func NewStdoutQlogTracer(config *Config) logging.Tracer {
	return NewTracer(func(p logging.Perspective, connectionID []byte) io.WriteCloser {
		return os.Stdout
	}, config)
}

func NewFileQlogTracer(filepath string, config *Config) logging.Tracer {
	return NewTracer(func(p logging.Perspective, connectionID []byte) io.WriteCloser {
		filepath := strings.ReplaceAll(filepath, "{odcid}", fmt.Sprintf("%x", connectionID))
		f, err := os.Create(filepath)
		if err != nil {
			panic(err)
		}
		return NewBufferedWriteCloser(bufio.NewWriter(f), f)
	}, config)
}

type bufferedWriteCloser struct {
	*bufio.Writer
	io.Closer
}

// NewBufferedWriteCloser creates an io.WriteCloser from a bufio.Writer and an io.Closer
func NewBufferedWriteCloser(writer *bufio.Writer, closer io.Closer) io.WriteCloser {
	return &bufferedWriteCloser{
		Writer: writer,
		Closer: closer,
	}
}

func (h bufferedWriteCloser) Close() error {
	if err := h.Writer.Flush(); err != nil {
		return err
	}
	return h.Closer.Close()
}
