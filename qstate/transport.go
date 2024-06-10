//go:generate msgp
package qstate

import "github.com/quic-go/quic-go/internal/protocol"

type Transport struct {
	Version    uint32 `msg:"version" json:"version" cbor:"1,keyasint"`
	ChosenALPN string `msg:"chosen_alpn" json:"chosen_alpn" cbor:"2,keyasint"`
	// client or server
	VantagePoint string `msg:"vantage_point" json:"vantage_point" cbor:"3,keyasint"`
	// active connection IDs;
	// must be sorted ascending by sequence number;
	ConnectionIDs []ConnectionID `msg:"connection_ids" json:"connection_ids" cbor:"4,keyasint"`
	// active peer connection IDs;
	// must be sorted ascending by sequence number;
	RemoteConnectionIDs []ConnectionID `msg:"remote_connection_ids" json:"remote_connection_ids" cbor:"5,keyasint"`
	DestinationIP       string         `msg:"dst_ip" json:"dst_ip" cbor:"6,keyasint"`
	DestinationPort     uint16         `msg:"dst_port" json:"dst_port" cbor:"7,keyasint"`
	// TODO only include non-default parameters
	Parameters Parameters `msg:"parameters" json:"parameters" cbor:"8,keyasint"`
	// TODO only include non-default parameters
	RemoteParameters Parameters `msg:"remote_parameters" json:"remote_parameters" cbor:"9,keyasint"`
	// minimum of max_idle_timout transport parameter advertised by both endpoints;
	// 0 if default
	IdleTimeout int64 `msg:"idle_timeout,omitzero" json:"idle_timeout,omitzero" cbor:"10,keyasint,omitzero"`
	// in byte;
	// max data that can be received;
	MaxData int64 `msg:"max_data" json:"max_data" cbor:"11,keyasint"`
	// in byte;
	// max data that can be sent;
	RemoteMaxData int64 `msg:"remote_max_data" json:"remote_max_data" cbor:"12,keyasint"`
	// in byte
	SentData int64 `msg:"sent_data" json:"sent_data" cbor:"13,keyasint"`
	// in byte
	ReceivedData                   int64 `msg:"received_data" json:"received_data" cbor:"14,keyasint"`
	MaxBidirectionalStreams        int64 `msg:"max_bidirectional_streams" json:"max_bidirectional_streams" cbor:"15,keyasint"`
	MaxUnidirectionalStreams       int64 `msg:"max_unidirectional_streams" json:"max_unidirectional_streams" cbor:"16,keyasint"`
	RemoteMaxBidirectionalStreams  int64 `msg:"remote_max_bidirectional_streams" json:"remote_max_bidirectional_streams" cbor:"17,keyasint"`
	RemoteMaxUnidirectionalStreams int64 `msg:"remote_max_unidirectional_streams" json:"remote_max_unidirectional_streams" cbor:"18,keyasint"`
	NextUnidirectionalStream       int64 `msg:"next_unidirectional_stream" json:"next_unidirectional_stream" cbor:"19,keyasint"`
	NextBidirectionalStream        int64 `msg:"next_bidirectional_stream" json:"next_bidirectional_stream" cbor:"20,keyasint"`
	// next unidirectional stream to accept from remote
	RemoteNextUnidirectionalStream int64 `msg:"remote_next_unidirectional_stream" json:"remote_next_unidirectional_stream" cbor:"21,keyasint"`
	// next bidirectional stream to accept from remote
	RemoteNextBidirectionalStream int64    `msg:"remote_next_bidirectional_stream" json:"remote_next_bidirectional_stream" cbor:"22,keyasint"`
	Streams                       []Stream `msg:"streams" json:"streams" cbor:"23,keyasint"`
	NextPacketNumber              int64    `msg:"next_packet_number" json:"next_packet_number" cbor:"24,keyasint"`
	HighestObservedPacketNumber   int64    `msg:"highest_observed_packet_number" json:"highest_observed_packet_number" cbor:"25,keyasint"`
	// received packet numbers
	AckRanges [][2]int64 `msg:"ack_ranges" json:"ack_ranges" cbor:"26,keyasint"`
	// acknowledged packets by peer
	RemoteAckRanges [][2]int64 `msg:"remote_ack_ranges" json:"remote_ack_ranges" cbor:"27,keyasint"`
	PendingAcks     []Packet   `msg:"pending_acks" json:"pending_acks" cbor:"28,keyasint"`
}

func (c *Transport) Perspective() protocol.Perspective {
	if c.VantagePoint == "client" {
		return protocol.PerspectiveClient
	} else {
		return protocol.PerspectiveServer
	}
}

func (s *Transport) ConnectionIDLength() int {
	for _, connectionID := range s.ConnectionIDs {
		return len(connectionID.ConnectionID)
	}
	panic("unexpected empty set")
}

func (s *Transport) OriginalDestinationConnectionID() []byte {
	if s.VantagePoint == "client" {
		return *s.RemoteParameters.OriginalDestinationConnectionID
	} else {
		return *s.Parameters.OriginalDestinationConnectionID
	}
}

func (s *Transport) GetStream(streamID int64) *Stream {
	for i := range s.Streams {
		if s.Streams[i].StreamID == streamID {
			return &s.Streams[i]
		}
	}
	panic("no such stream")
}

func (s *Transport) PutBack(streamID int64, offset int64, data []byte) {
	s.GetStream(streamID).PutBack(offset, data)
}

func (t *Transport) ChangeVantagePoint(DestinationIP string, DestinationPort uint16) Transport {
	f := Transport{
		Version:                        t.Version,
		ChosenALPN:                     t.ChosenALPN,
		ConnectionIDs:                  t.RemoteConnectionIDs,
		RemoteConnectionIDs:            t.ConnectionIDs,
		DestinationIP:                  DestinationIP,
		DestinationPort:                DestinationPort,
		Parameters:                     t.RemoteParameters,
		IdleTimeout:                    t.IdleTimeout,
		RemoteParameters:               t.Parameters,
		MaxData:                        t.RemoteMaxData,
		RemoteMaxData:                  t.MaxData,
		SentData:                       t.ReceivedData,
		ReceivedData:                   t.SentData,
		MaxBidirectionalStreams:        t.RemoteMaxBidirectionalStreams,
		MaxUnidirectionalStreams:       t.RemoteMaxUnidirectionalStreams,
		RemoteMaxBidirectionalStreams:  t.MaxBidirectionalStreams,
		RemoteMaxUnidirectionalStreams: t.MaxUnidirectionalStreams,
		NextUnidirectionalStream:       t.RemoteNextUnidirectionalStream,
		NextBidirectionalStream:        t.RemoteNextBidirectionalStream,
		RemoteNextUnidirectionalStream: t.NextUnidirectionalStream,
		RemoteNextBidirectionalStream:  t.NextBidirectionalStream,
		NextPacketNumber:               t.HighestObservedPacketNumber,
		HighestObservedPacketNumber:    t.NextPacketNumber,
		AckRanges:                      t.RemoteAckRanges,
		RemoteAckRanges:                t.AckRanges,
		PendingAcks:                    nil, // TODO apply acks to state
	}

	if t.VantagePoint == "client" {
		f.VantagePoint = "server"
	} else {
		f.VantagePoint = "client"
	}

	for _, stream := range t.Streams {
		f.Streams = append(f.Streams, stream.ChangeVantagePoint())
	}

	return f
}

// ResetStreamWriteOffsetsToAck reset not acknowledged stream write offsets
func (c *Transport) ResetStreamWriteOffsetsToAck() {
	for _, stream := range c.Streams {
		if stream.WriteOffset == nil {
			continue
		}
		reducedOffset := *stream.WriteAck + 1
		diff := *stream.WriteOffset - reducedOffset
		*stream.WriteOffset = reducedOffset
		c.SentData -= diff
	}
}
