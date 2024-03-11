//go:generate msgp
package qstate

import "github.com/quic-go/quic-go/internal/protocol"

type Transport struct {
	Version    uint32 `msg:"version" json:"version"`
	KeyPhase   uint64 `msg:"key_phase" json:"key_phase"`
	ChosenALPN string `msg:"chosen_alpn" json:"chosen_alpn"`
	// client or server
	VantagePoint string `msg:"vantage_point" json:"vantage_point"`
	// active connection IDs;
	// must be sorted ascending by sequence number;
	ConnectionIDs []ConnectionID `msg:"connection_ids" json:"connection_ids"`
	// active peer connection IDs;
	// must be sorted ascending by sequence number;
	RemoteConnectionIDs []ConnectionID `msg:"remote_connection_ids" json:"remote_connection_ids"`
	// id of the used TLS 1.3 cipher suites.
	// see RFC 8446 Appendix B.4. Cipher Suites.
	// e.g."AES_128_GCM_SHA256"
	TlsCipher string `msg:"tls_cipher" json:"tls_cipher"`
	// used for header protection sent by peer.
	// see RFC 9001 Section 5.4 Header Protection.
	RemoteHeaderProtectionKey []byte `msg:"remote_header_protection_key" json:"remote_header_protection_key"`
	// used for header protection sent to peer.
	// see RFC 9001 Section 5.4 Header Protection.
	HeaderProtectionKey []byte `msg:"header_protection_key" json:"header_protection_key"`
	// secret used on packets sent from peer.
	RemoteTrafficSecret []byte `msg:"remote_traffic_secret" json:"remote_traffic_secret"`
	// secret used on packets sent to peer.
	TrafficSecret   []byte `msg:"traffic_secret" json:"traffic_secret"`
	DestinationIP   string `msg:"dst_ip" json:"dst_ip"`
	DestinationPort uint16 `msg:"dst_port" json:"dst_port"`
	// TODO only include non-default parameters
	Parameters Parameters `msg:"parameters" json:"parameters"`
	// TODO only include non-default parameters
	RemoteParameters Parameters `msg:"remote_parameters" json:"remote_parameters"`
	// in byte;
	// max data that can be received;
	MaxData int64 `msg:"max_data" json:"max_data"`
	// in byte;
	// max data that can be sent;
	RemoteMaxData int64 `msg:"remote_max_data" json:"remote_max_data"`
	// in byte
	SentData int64 `msg:"sent_data" json:"sent_data"`
	// in byte
	ReceivedData                   int64 `msg:"received_data" json:"received_data"`
	MaxBidirectionalStreams        int64 `msg:"max_bidirectional_streams" json:"max_bidirectional_streams"`
	MaxUnidirectionalStreams       int64 `msg:"max_unidirectional_streams" json:"max_unidirectional_streams"`
	RemoteMaxBidirectionalStreams  int64 `msg:"remote_max_bidirectional_streams" json:"remote_max_bidirectional_streams"`
	RemoteMaxUnidirectionalStreams int64 `msg:"remote_max_unidirectional_streams" json:"remote_max_unidirectional_streams"`
	NextUnidirectionalStream       int64 `msg:"next_unidirectional_stream" json:"next_unidirectional_stream"`
	NextBidirectionalStream        int64 `msg:"next_bidirectional_stream" json:"next_bidirectional_stream"`
	// next unidirectional stream to accept from remote
	RemoteNextUnidirectionalStream int64 `msg:"remote_next_unidirectional_stream" json:"remote_next_unidirectional_stream"`
	// next bidirectional stream to accept from remote
	RemoteNextBidirectionalStream int64    `msg:"remote_next_bidirectional_stream" json:"remote_next_bidirectional_stream"`
	Streams                       []Stream `msg:"streams" json:"streams"`
	NextPacketNumber              int64    `msg:"next_packet_number" json:"next_packet_number"`
	HighestObservedPacketNumber   int64    `msg:"highest_observed_packet_number" json:"highest_observed_packet_number"`
	// received packet numbers
	AckRanges [][2]int64 `msg:"ack_ranges" json:"ack_ranges"`
	// acknowledged packets by peer
	RemoteAckRanges [][2]int64 `msg:"remote_ack_ranges" json:"remote_ack_ranges"`
	PendingAcks     []Packet   `msg:"pending_acks" json:"pending_acks"`
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

func (s *Transport) PutBack(streamID int64, offset int64, data []byte) {
	for i := range s.Streams {
		if s.Streams[i].StreamID == streamID {
			s.Streams[i].PutBack(offset, data)
			return
		}
	}
	panic("no such stream")
}

func (t *Transport) ChangeVantagePoint(DestinationIP string, DestinationPort uint16) Transport {
	f := Transport{
		Version:                        t.Version,
		KeyPhase:                       t.KeyPhase,
		ChosenALPN:                     t.ChosenALPN,
		ConnectionIDs:                  t.RemoteConnectionIDs,
		RemoteConnectionIDs:            t.ConnectionIDs,
		TlsCipher:                      t.TlsCipher,
		RemoteHeaderProtectionKey:      t.HeaderProtectionKey,
		HeaderProtectionKey:            t.RemoteHeaderProtectionKey,
		RemoteTrafficSecret:            t.TrafficSecret,
		TrafficSecret:                  t.RemoteTrafficSecret,
		DestinationIP:                  DestinationIP,
		DestinationPort:                DestinationPort,
		Parameters:                     t.RemoteParameters,
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
