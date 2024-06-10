//go:generate msgp
package qstate

import "github.com/quic-go/quic-go/internal/utils"

type Stream struct {
	StreamID int64 `msg:"stream_id" json:"stream_id" cbor:"1,keyasint"`
	// in byte;
	// only set for write streams;
	// offset for writing to write queue;
	WriteOffset *int64 `msg:"write_offset,omitempty" json:"write_offset,omitempty" cbor:"2,keyasint,omitempty"`
	// in byte;
	// only set for write streams if fin is written;
	WriteFin *int64 `msg:"write_fin,omitempty" json:"write_fin,omitempty" cbor:"3,keyasint,omitempty"`
	// in byte;
	// only set for write streams;
	// also required for sending RESET_STREAM frames;
	WriteMaxData *int64 `msg:"write_max_data,omitempty" json:"write_max_data,omitempty" cbor:"4,keyasint,omitempty"`
	// in byte;
	// only set for write streams;
	// offset until stream data is acknowledged
	WriteAck *int64 `msg:"write_ack,omitempty" json:"write_ack,omitempty" cbor:"5,keyasint,omitempty"`
	// only set for write streams;
	// stream data written by the application but not yet sent on the network;
	// must be sorted ascending by offset;
	WriteQueue []StreamRange `msg:"write_queue,omitempty" json:"write_queue,omitempty" cbor:"6,keyasint,omitempty"`
	// in byte;
	// only set for read streams
	ReadOffset *int64 `msg:"read_offset,omitempty" json:"read_offset,omitempty" cbor:"7,keyasint,omitempty"`
	// in byte;
	// only set for read streams if fin was read;
	ReadFin *int64 `msg:"read_fin,omitempty" json:"read_fin,omitempty" cbor:"8,keyasint,omitempty"`
	// in byte;
	// only set for read streams;
	ReadMaxData *int64 `msg:"read_max_data,omitempty" json:"read_max_data,omitempty" cbor:"9,keyasint,omitempty"`
	// only set for read streams;
	// stream data received from network but not yet read by application
	// must be sorted ascending by offset;
	ReadQueue []StreamRange `msg:"read_queue,omitempty" json:"read_queue,omitempty" cbor:"10,keyasint,omitempty"`
}

func (s *Stream) ChangeVantagePoint() Stream {
	f := Stream{
		StreamID:     s.StreamID,
		WriteMaxData: utils.New(*s.ReadMaxData),
		WriteQueue:   nil,
		ReadMaxData:  utils.New(*s.WriteMaxData),
	}

	if s.ReadFin != nil {
		f.WriteFin = utils.New(*s.ReadFin)
	}

	if s.WriteFin != nil {
		f.ReadFin = utils.New(*s.WriteFin)
	}

	f.ReadQueue = make([]StreamRange, len(s.WriteQueue))
	copy(f.ReadQueue, s.WriteQueue)

	if s.ReadQueue != nil {
		lastRange := s.ReadQueue[len(s.ReadQueue)-1]
		f.WriteOffset = utils.New(lastRange.Offset + int64(len(lastRange.Data)))
	} else {
		f.WriteOffset = utils.New(*s.ReadOffset)
	}
	f.WriteAck = utils.New(*f.WriteOffset)

	if s.WriteQueue != nil {
		f.ReadOffset = utils.New(s.WriteQueue[0].Offset)
	} else {
		f.ReadOffset = utils.New(*s.WriteOffset)
	}

	return f
}

func (s *Stream) PutBack(offset int64, data []byte) {
	if len(s.ReadQueue) != 0 && offset+int64(len(data)) > s.ReadQueue[0].Offset {
		panic("unsorted or overlapping")
	}
	s.ReadQueue = append([]StreamRange{{Offset: offset, Data: data}}, s.ReadQueue...)
	*s.ReadOffset = min(*s.ReadOffset, offset)
}
