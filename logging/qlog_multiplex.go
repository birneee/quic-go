package logging

import "time"

type qlogWriterMultiplexer struct {
	writers []QlogWriter
}

var _ QlogWriter = &qlogWriterMultiplexer{}

func (t *qlogWriterMultiplexer) RecordEvent(eventTime time.Time, details QlogEventDetails) {
	for _, writer := range t.writers {
		writer.RecordEvent(eventTime, details)
	}
}
