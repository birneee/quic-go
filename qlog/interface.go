package qlog

import (
	"github.com/lucas-clemente/quic-go/logging"
	"time"
)

type ConnectionTracer = connectionTracer

var _ logging.QlogWriter = &ConnectionTracer{}

func (t *ConnectionTracer) RecordEvent(eventTime time.Time, details logging.QlogEventDetails) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.recordEvent(eventTime, details)
}

type Config struct {
	ExcludeEventsByDefault bool
	// keys in form "<category>:<name>"
	// e.g. "transport:packet_received"
	IncludedEvents map[string]bool
}

func (c *Config) Included(event string) bool {
	if included, ok := c.IncludedEvents[event]; ok {
		return included
	}
	return !c.ExcludeEventsByDefault
}
