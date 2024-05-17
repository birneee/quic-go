// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package qstate

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson28548e40DecodeGithubComQuicGoQuicGoQstate(in *jlexer.Lexer, out *Connection) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "state":
			out.State = ConnectionState(in.String())
		case "transport":
			easyjson28548e40DecodeGithubComQuicGoQuicGoQstate1(in, &out.Transport)
		case "crypto":
			easyjson28548e40DecodeGithubComQuicGoQuicGoQstate2(in, &out.Crypto)
		case "metrics":
			easyjson28548e40DecodeGithubComQuicGoQuicGoQstate3(in, &out.Metrics)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson28548e40EncodeGithubComQuicGoQuicGoQstate(out *jwriter.Writer, in Connection) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"state\":"
		out.RawString(prefix[1:])
		out.String(string(in.State))
	}
	{
		const prefix string = ",\"transport\":"
		out.RawString(prefix)
		easyjson28548e40EncodeGithubComQuicGoQuicGoQstate1(out, in.Transport)
	}
	{
		const prefix string = ",\"crypto\":"
		out.RawString(prefix)
		easyjson28548e40EncodeGithubComQuicGoQuicGoQstate2(out, in.Crypto)
	}
	{
		const prefix string = ",\"metrics\":"
		out.RawString(prefix)
		easyjson28548e40EncodeGithubComQuicGoQuicGoQstate3(out, in.Metrics)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Connection) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson28548e40EncodeGithubComQuicGoQuicGoQstate(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Connection) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson28548e40EncodeGithubComQuicGoQuicGoQstate(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Connection) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson28548e40DecodeGithubComQuicGoQuicGoQstate(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Connection) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson28548e40DecodeGithubComQuicGoQuicGoQstate(l, v)
}
func easyjson28548e40DecodeGithubComQuicGoQuicGoQstate3(in *jlexer.Lexer, out *Metrics) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "congestion_window":
			if in.IsNull() {
				in.Skip()
				out.CongestionWindow = nil
			} else {
				if out.CongestionWindow == nil {
					out.CongestionWindow = new(int64)
				}
				*out.CongestionWindow = int64(in.Int64())
			}
		case "smoothed_rtt":
			if in.IsNull() {
				in.Skip()
				out.SmoothedRTT = nil
			} else {
				if out.SmoothedRTT == nil {
					out.SmoothedRTT = new(int64)
				}
				*out.SmoothedRTT = int64(in.Int64())
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson28548e40EncodeGithubComQuicGoQuicGoQstate3(out *jwriter.Writer, in Metrics) {
	out.RawByte('{')
	first := true
	_ = first
	if in.CongestionWindow != nil {
		const prefix string = ",\"congestion_window\":"
		first = false
		out.RawString(prefix[1:])
		out.Int64(int64(*in.CongestionWindow))
	}
	if in.SmoothedRTT != nil {
		const prefix string = ",\"smoothed_rtt\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(*in.SmoothedRTT))
	}
	out.RawByte('}')
}
func easyjson28548e40DecodeGithubComQuicGoQuicGoQstate2(in *jlexer.Lexer, out *Crypto) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "key_phase":
			out.KeyPhase = uint64(in.Uint64())
		case "tls_cipher":
			out.TlsCipher = string(in.String())
		case "remote_header_protection_key":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.RemoteHeaderProtectionKey).UnmarshalJSON(data))
			}
		case "header_protection_key":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.HeaderProtectionKey).UnmarshalJSON(data))
			}
		case "remote_traffic_secret":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.RemoteTrafficSecret).UnmarshalJSON(data))
			}
		case "traffic_secret":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.TrafficSecret).UnmarshalJSON(data))
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson28548e40EncodeGithubComQuicGoQuicGoQstate2(out *jwriter.Writer, in Crypto) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"key_phase\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.KeyPhase))
	}
	{
		const prefix string = ",\"tls_cipher\":"
		out.RawString(prefix)
		out.String(string(in.TlsCipher))
	}
	{
		const prefix string = ",\"remote_header_protection_key\":"
		out.RawString(prefix)
		out.Raw((in.RemoteHeaderProtectionKey).MarshalJSON())
	}
	{
		const prefix string = ",\"header_protection_key\":"
		out.RawString(prefix)
		out.Raw((in.HeaderProtectionKey).MarshalJSON())
	}
	{
		const prefix string = ",\"remote_traffic_secret\":"
		out.RawString(prefix)
		out.Raw((in.RemoteTrafficSecret).MarshalJSON())
	}
	{
		const prefix string = ",\"traffic_secret\":"
		out.RawString(prefix)
		out.Raw((in.TrafficSecret).MarshalJSON())
	}
	out.RawByte('}')
}
func easyjson28548e40DecodeGithubComQuicGoQuicGoQstate1(in *jlexer.Lexer, out *Transport) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "version":
			out.Version = uint32(in.Uint32())
		case "chosen_alpn":
			out.ChosenALPN = string(in.String())
		case "vantage_point":
			out.VantagePoint = string(in.String())
		case "connection_ids":
			if in.IsNull() {
				in.Skip()
				out.ConnectionIDs = nil
			} else {
				in.Delim('[')
				if out.ConnectionIDs == nil {
					if !in.IsDelim(']') {
						out.ConnectionIDs = make([]ConnectionID, 0, 1)
					} else {
						out.ConnectionIDs = []ConnectionID{}
					}
				} else {
					out.ConnectionIDs = (out.ConnectionIDs)[:0]
				}
				for !in.IsDelim(']') {
					var v1 ConnectionID
					easyjson28548e40DecodeGithubComQuicGoQuicGoQstate4(in, &v1)
					out.ConnectionIDs = append(out.ConnectionIDs, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "remote_connection_ids":
			if in.IsNull() {
				in.Skip()
				out.RemoteConnectionIDs = nil
			} else {
				in.Delim('[')
				if out.RemoteConnectionIDs == nil {
					if !in.IsDelim(']') {
						out.RemoteConnectionIDs = make([]ConnectionID, 0, 1)
					} else {
						out.RemoteConnectionIDs = []ConnectionID{}
					}
				} else {
					out.RemoteConnectionIDs = (out.RemoteConnectionIDs)[:0]
				}
				for !in.IsDelim(']') {
					var v2 ConnectionID
					easyjson28548e40DecodeGithubComQuicGoQuicGoQstate4(in, &v2)
					out.RemoteConnectionIDs = append(out.RemoteConnectionIDs, v2)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "dst_ip":
			out.DestinationIP = string(in.String())
		case "dst_port":
			out.DestinationPort = uint16(in.Uint16())
		case "parameters":
			easyjson28548e40DecodeGithubComQuicGoQuicGoQstate5(in, &out.Parameters)
		case "remote_parameters":
			easyjson28548e40DecodeGithubComQuicGoQuicGoQstate5(in, &out.RemoteParameters)
		case "idle_timeout":
			out.IdleTimeout = int64(in.Int64())
		case "max_data":
			out.MaxData = int64(in.Int64())
		case "remote_max_data":
			out.RemoteMaxData = int64(in.Int64())
		case "sent_data":
			out.SentData = int64(in.Int64())
		case "received_data":
			out.ReceivedData = int64(in.Int64())
		case "max_bidirectional_streams":
			out.MaxBidirectionalStreams = int64(in.Int64())
		case "max_unidirectional_streams":
			out.MaxUnidirectionalStreams = int64(in.Int64())
		case "remote_max_bidirectional_streams":
			out.RemoteMaxBidirectionalStreams = int64(in.Int64())
		case "remote_max_unidirectional_streams":
			out.RemoteMaxUnidirectionalStreams = int64(in.Int64())
		case "next_unidirectional_stream":
			out.NextUnidirectionalStream = int64(in.Int64())
		case "next_bidirectional_stream":
			out.NextBidirectionalStream = int64(in.Int64())
		case "remote_next_unidirectional_stream":
			out.RemoteNextUnidirectionalStream = int64(in.Int64())
		case "remote_next_bidirectional_stream":
			out.RemoteNextBidirectionalStream = int64(in.Int64())
		case "streams":
			if in.IsNull() {
				in.Skip()
				out.Streams = nil
			} else {
				in.Delim('[')
				if out.Streams == nil {
					if !in.IsDelim(']') {
						out.Streams = make([]Stream, 0, 0)
					} else {
						out.Streams = []Stream{}
					}
				} else {
					out.Streams = (out.Streams)[:0]
				}
				for !in.IsDelim(']') {
					var v3 Stream
					easyjson28548e40DecodeGithubComQuicGoQuicGoQstate6(in, &v3)
					out.Streams = append(out.Streams, v3)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "next_packet_number":
			out.NextPacketNumber = int64(in.Int64())
		case "highest_observed_packet_number":
			out.HighestObservedPacketNumber = int64(in.Int64())
		case "ack_ranges":
			if in.IsNull() {
				in.Skip()
				out.AckRanges = nil
			} else {
				in.Delim('[')
				if out.AckRanges == nil {
					if !in.IsDelim(']') {
						out.AckRanges = make([][2]int64, 0, 4)
					} else {
						out.AckRanges = [][2]int64{}
					}
				} else {
					out.AckRanges = (out.AckRanges)[:0]
				}
				for !in.IsDelim(']') {
					var v4 [2]int64
					if in.IsNull() {
						in.Skip()
					} else {
						in.Delim('[')
						v5 := 0
						for !in.IsDelim(']') {
							if v5 < 2 {
								(v4)[v5] = int64(in.Int64())
								v5++
							} else {
								in.SkipRecursive()
							}
							in.WantComma()
						}
						in.Delim(']')
					}
					out.AckRanges = append(out.AckRanges, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "remote_ack_ranges":
			if in.IsNull() {
				in.Skip()
				out.RemoteAckRanges = nil
			} else {
				in.Delim('[')
				if out.RemoteAckRanges == nil {
					if !in.IsDelim(']') {
						out.RemoteAckRanges = make([][2]int64, 0, 4)
					} else {
						out.RemoteAckRanges = [][2]int64{}
					}
				} else {
					out.RemoteAckRanges = (out.RemoteAckRanges)[:0]
				}
				for !in.IsDelim(']') {
					var v6 [2]int64
					if in.IsNull() {
						in.Skip()
					} else {
						in.Delim('[')
						v7 := 0
						for !in.IsDelim(']') {
							if v7 < 2 {
								(v6)[v7] = int64(in.Int64())
								v7++
							} else {
								in.SkipRecursive()
							}
							in.WantComma()
						}
						in.Delim(']')
					}
					out.RemoteAckRanges = append(out.RemoteAckRanges, v6)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "pending_acks":
			if in.IsNull() {
				in.Skip()
				out.PendingAcks = nil
			} else {
				in.Delim('[')
				if out.PendingAcks == nil {
					if !in.IsDelim(']') {
						out.PendingAcks = make([]Packet, 0, 2)
					} else {
						out.PendingAcks = []Packet{}
					}
				} else {
					out.PendingAcks = (out.PendingAcks)[:0]
				}
				for !in.IsDelim(']') {
					var v8 Packet
					easyjson28548e40DecodeGithubComQuicGoQuicGoQstate7(in, &v8)
					out.PendingAcks = append(out.PendingAcks, v8)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson28548e40EncodeGithubComQuicGoQuicGoQstate1(out *jwriter.Writer, in Transport) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"version\":"
		out.RawString(prefix[1:])
		out.Uint32(uint32(in.Version))
	}
	{
		const prefix string = ",\"chosen_alpn\":"
		out.RawString(prefix)
		out.String(string(in.ChosenALPN))
	}
	{
		const prefix string = ",\"vantage_point\":"
		out.RawString(prefix)
		out.String(string(in.VantagePoint))
	}
	{
		const prefix string = ",\"connection_ids\":"
		out.RawString(prefix)
		if in.ConnectionIDs == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v9, v10 := range in.ConnectionIDs {
				if v9 > 0 {
					out.RawByte(',')
				}
				easyjson28548e40EncodeGithubComQuicGoQuicGoQstate4(out, v10)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"remote_connection_ids\":"
		out.RawString(prefix)
		if in.RemoteConnectionIDs == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v11, v12 := range in.RemoteConnectionIDs {
				if v11 > 0 {
					out.RawByte(',')
				}
				easyjson28548e40EncodeGithubComQuicGoQuicGoQstate4(out, v12)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"dst_ip\":"
		out.RawString(prefix)
		out.String(string(in.DestinationIP))
	}
	{
		const prefix string = ",\"dst_port\":"
		out.RawString(prefix)
		out.Uint16(uint16(in.DestinationPort))
	}
	{
		const prefix string = ",\"parameters\":"
		out.RawString(prefix)
		easyjson28548e40EncodeGithubComQuicGoQuicGoQstate5(out, in.Parameters)
	}
	{
		const prefix string = ",\"remote_parameters\":"
		out.RawString(prefix)
		easyjson28548e40EncodeGithubComQuicGoQuicGoQstate5(out, in.RemoteParameters)
	}
	{
		const prefix string = ",\"idle_timeout\":"
		out.RawString(prefix)
		out.Int64(int64(in.IdleTimeout))
	}
	{
		const prefix string = ",\"max_data\":"
		out.RawString(prefix)
		out.Int64(int64(in.MaxData))
	}
	{
		const prefix string = ",\"remote_max_data\":"
		out.RawString(prefix)
		out.Int64(int64(in.RemoteMaxData))
	}
	{
		const prefix string = ",\"sent_data\":"
		out.RawString(prefix)
		out.Int64(int64(in.SentData))
	}
	{
		const prefix string = ",\"received_data\":"
		out.RawString(prefix)
		out.Int64(int64(in.ReceivedData))
	}
	{
		const prefix string = ",\"max_bidirectional_streams\":"
		out.RawString(prefix)
		out.Int64(int64(in.MaxBidirectionalStreams))
	}
	{
		const prefix string = ",\"max_unidirectional_streams\":"
		out.RawString(prefix)
		out.Int64(int64(in.MaxUnidirectionalStreams))
	}
	{
		const prefix string = ",\"remote_max_bidirectional_streams\":"
		out.RawString(prefix)
		out.Int64(int64(in.RemoteMaxBidirectionalStreams))
	}
	{
		const prefix string = ",\"remote_max_unidirectional_streams\":"
		out.RawString(prefix)
		out.Int64(int64(in.RemoteMaxUnidirectionalStreams))
	}
	{
		const prefix string = ",\"next_unidirectional_stream\":"
		out.RawString(prefix)
		out.Int64(int64(in.NextUnidirectionalStream))
	}
	{
		const prefix string = ",\"next_bidirectional_stream\":"
		out.RawString(prefix)
		out.Int64(int64(in.NextBidirectionalStream))
	}
	{
		const prefix string = ",\"remote_next_unidirectional_stream\":"
		out.RawString(prefix)
		out.Int64(int64(in.RemoteNextUnidirectionalStream))
	}
	{
		const prefix string = ",\"remote_next_bidirectional_stream\":"
		out.RawString(prefix)
		out.Int64(int64(in.RemoteNextBidirectionalStream))
	}
	{
		const prefix string = ",\"streams\":"
		out.RawString(prefix)
		if in.Streams == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v13, v14 := range in.Streams {
				if v13 > 0 {
					out.RawByte(',')
				}
				easyjson28548e40EncodeGithubComQuicGoQuicGoQstate6(out, v14)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"next_packet_number\":"
		out.RawString(prefix)
		out.Int64(int64(in.NextPacketNumber))
	}
	{
		const prefix string = ",\"highest_observed_packet_number\":"
		out.RawString(prefix)
		out.Int64(int64(in.HighestObservedPacketNumber))
	}
	{
		const prefix string = ",\"ack_ranges\":"
		out.RawString(prefix)
		if in.AckRanges == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v15, v16 := range in.AckRanges {
				if v15 > 0 {
					out.RawByte(',')
				}
				out.RawByte('[')
				for v17 := range v16 {
					if v17 > 0 {
						out.RawByte(',')
					}
					out.Int64(int64((v16)[v17]))
				}
				out.RawByte(']')
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"remote_ack_ranges\":"
		out.RawString(prefix)
		if in.RemoteAckRanges == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v18, v19 := range in.RemoteAckRanges {
				if v18 > 0 {
					out.RawByte(',')
				}
				out.RawByte('[')
				for v20 := range v19 {
					if v20 > 0 {
						out.RawByte(',')
					}
					out.Int64(int64((v19)[v20]))
				}
				out.RawByte(']')
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"pending_acks\":"
		out.RawString(prefix)
		if in.PendingAcks == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v21, v22 := range in.PendingAcks {
				if v21 > 0 {
					out.RawByte(',')
				}
				easyjson28548e40EncodeGithubComQuicGoQuicGoQstate7(out, v22)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}
func easyjson28548e40DecodeGithubComQuicGoQuicGoQstate7(in *jlexer.Lexer, out *Packet) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "packet_number":
			out.PacketNumber = int64(in.Int64())
		case "frames":
			if in.IsNull() {
				in.Skip()
				out.Frames = nil
			} else {
				in.Delim('[')
				if out.Frames == nil {
					if !in.IsDelim(']') {
						out.Frames = make([]Frame, 0, 0)
					} else {
						out.Frames = []Frame{}
					}
				} else {
					out.Frames = (out.Frames)[:0]
				}
				for !in.IsDelim(']') {
					var v23 Frame
					easyjson28548e40DecodeGithubComQuicGoQuicGoQstate8(in, &v23)
					out.Frames = append(out.Frames, v23)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson28548e40EncodeGithubComQuicGoQuicGoQstate7(out *jwriter.Writer, in Packet) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"packet_number\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.PacketNumber))
	}
	{
		const prefix string = ",\"frames\":"
		out.RawString(prefix)
		if in.Frames == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v24, v25 := range in.Frames {
				if v24 > 0 {
					out.RawByte(',')
				}
				easyjson28548e40EncodeGithubComQuicGoQuicGoQstate8(out, v25)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}
func easyjson28548e40DecodeGithubComQuicGoQuicGoQstate8(in *jlexer.Lexer, out *Frame) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "frame_type":
			out.Type = string(in.String())
		case "stream_id":
			if in.IsNull() {
				in.Skip()
				out.StreamID = nil
			} else {
				if out.StreamID == nil {
					out.StreamID = new(int64)
				}
				*out.StreamID = int64(in.Int64())
			}
		case "offset":
			if in.IsNull() {
				in.Skip()
				out.Offset = nil
			} else {
				if out.Offset == nil {
					out.Offset = new(int64)
				}
				*out.Offset = int64(in.Int64())
			}
		case "length":
			if in.IsNull() {
				in.Skip()
				out.Length = nil
			} else {
				if out.Length == nil {
					out.Length = new(int64)
				}
				*out.Length = int64(in.Int64())
			}
		case "token":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Token).UnmarshalJSON(data))
			}
		case "data":
			if in.IsNull() {
				in.Skip()
				out.Data = nil
			} else {
				out.Data = in.Bytes()
			}
		case "sequence_number":
			if in.IsNull() {
				in.Skip()
				out.SequenceNumber = nil
			} else {
				if out.SequenceNumber == nil {
					out.SequenceNumber = new(uint64)
				}
				*out.SequenceNumber = uint64(in.Uint64())
			}
		case "stream_type":
			out.StreamType = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson28548e40EncodeGithubComQuicGoQuicGoQstate8(out *jwriter.Writer, in Frame) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"frame_type\":"
		out.RawString(prefix[1:])
		out.String(string(in.Type))
	}
	if in.StreamID != nil {
		const prefix string = ",\"stream_id\":"
		out.RawString(prefix)
		out.Int64(int64(*in.StreamID))
	}
	if in.Offset != nil {
		const prefix string = ",\"offset\":"
		out.RawString(prefix)
		out.Int64(int64(*in.Offset))
	}
	if in.Length != nil {
		const prefix string = ",\"length\":"
		out.RawString(prefix)
		out.Int64(int64(*in.Length))
	}
	if len(in.Token) != 0 {
		const prefix string = ",\"token\":"
		out.RawString(prefix)
		out.Raw((in.Token).MarshalJSON())
	}
	if len(in.Data) != 0 {
		const prefix string = ",\"data\":"
		out.RawString(prefix)
		out.Base64Bytes(in.Data)
	}
	if in.SequenceNumber != nil {
		const prefix string = ",\"sequence_number\":"
		out.RawString(prefix)
		out.Uint64(uint64(*in.SequenceNumber))
	}
	if in.StreamType != "" {
		const prefix string = ",\"stream_type\":"
		out.RawString(prefix)
		out.String(string(in.StreamType))
	}
	out.RawByte('}')
}
func easyjson28548e40DecodeGithubComQuicGoQuicGoQstate6(in *jlexer.Lexer, out *Stream) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "stream_id":
			out.StreamID = int64(in.Int64())
		case "write_offset":
			if in.IsNull() {
				in.Skip()
				out.WriteOffset = nil
			} else {
				if out.WriteOffset == nil {
					out.WriteOffset = new(int64)
				}
				*out.WriteOffset = int64(in.Int64())
			}
		case "write_fin":
			if in.IsNull() {
				in.Skip()
				out.WriteFin = nil
			} else {
				if out.WriteFin == nil {
					out.WriteFin = new(int64)
				}
				*out.WriteFin = int64(in.Int64())
			}
		case "write_max_data":
			if in.IsNull() {
				in.Skip()
				out.WriteMaxData = nil
			} else {
				if out.WriteMaxData == nil {
					out.WriteMaxData = new(int64)
				}
				*out.WriteMaxData = int64(in.Int64())
			}
		case "write_ack":
			if in.IsNull() {
				in.Skip()
				out.WriteAck = nil
			} else {
				if out.WriteAck == nil {
					out.WriteAck = new(int64)
				}
				*out.WriteAck = int64(in.Int64())
			}
		case "write_queue":
			if in.IsNull() {
				in.Skip()
				out.WriteQueue = nil
			} else {
				in.Delim('[')
				if out.WriteQueue == nil {
					if !in.IsDelim(']') {
						out.WriteQueue = make([]StreamRange, 0, 2)
					} else {
						out.WriteQueue = []StreamRange{}
					}
				} else {
					out.WriteQueue = (out.WriteQueue)[:0]
				}
				for !in.IsDelim(']') {
					var v29 StreamRange
					easyjson28548e40DecodeGithubComQuicGoQuicGoQstate9(in, &v29)
					out.WriteQueue = append(out.WriteQueue, v29)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "read_offset":
			if in.IsNull() {
				in.Skip()
				out.ReadOffset = nil
			} else {
				if out.ReadOffset == nil {
					out.ReadOffset = new(int64)
				}
				*out.ReadOffset = int64(in.Int64())
			}
		case "read_fin":
			if in.IsNull() {
				in.Skip()
				out.ReadFin = nil
			} else {
				if out.ReadFin == nil {
					out.ReadFin = new(int64)
				}
				*out.ReadFin = int64(in.Int64())
			}
		case "read_max_data":
			if in.IsNull() {
				in.Skip()
				out.ReadMaxData = nil
			} else {
				if out.ReadMaxData == nil {
					out.ReadMaxData = new(int64)
				}
				*out.ReadMaxData = int64(in.Int64())
			}
		case "read_queue":
			if in.IsNull() {
				in.Skip()
				out.ReadQueue = nil
			} else {
				in.Delim('[')
				if out.ReadQueue == nil {
					if !in.IsDelim(']') {
						out.ReadQueue = make([]StreamRange, 0, 2)
					} else {
						out.ReadQueue = []StreamRange{}
					}
				} else {
					out.ReadQueue = (out.ReadQueue)[:0]
				}
				for !in.IsDelim(']') {
					var v30 StreamRange
					easyjson28548e40DecodeGithubComQuicGoQuicGoQstate9(in, &v30)
					out.ReadQueue = append(out.ReadQueue, v30)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson28548e40EncodeGithubComQuicGoQuicGoQstate6(out *jwriter.Writer, in Stream) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"stream_id\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.StreamID))
	}
	if in.WriteOffset != nil {
		const prefix string = ",\"write_offset\":"
		out.RawString(prefix)
		out.Int64(int64(*in.WriteOffset))
	}
	if in.WriteFin != nil {
		const prefix string = ",\"write_fin\":"
		out.RawString(prefix)
		out.Int64(int64(*in.WriteFin))
	}
	if in.WriteMaxData != nil {
		const prefix string = ",\"write_max_data\":"
		out.RawString(prefix)
		out.Int64(int64(*in.WriteMaxData))
	}
	if in.WriteAck != nil {
		const prefix string = ",\"write_ack\":"
		out.RawString(prefix)
		out.Int64(int64(*in.WriteAck))
	}
	if len(in.WriteQueue) != 0 {
		const prefix string = ",\"write_queue\":"
		out.RawString(prefix)
		{
			out.RawByte('[')
			for v31, v32 := range in.WriteQueue {
				if v31 > 0 {
					out.RawByte(',')
				}
				easyjson28548e40EncodeGithubComQuicGoQuicGoQstate9(out, v32)
			}
			out.RawByte(']')
		}
	}
	if in.ReadOffset != nil {
		const prefix string = ",\"read_offset\":"
		out.RawString(prefix)
		out.Int64(int64(*in.ReadOffset))
	}
	if in.ReadFin != nil {
		const prefix string = ",\"read_fin\":"
		out.RawString(prefix)
		out.Int64(int64(*in.ReadFin))
	}
	if in.ReadMaxData != nil {
		const prefix string = ",\"read_max_data\":"
		out.RawString(prefix)
		out.Int64(int64(*in.ReadMaxData))
	}
	if len(in.ReadQueue) != 0 {
		const prefix string = ",\"read_queue\":"
		out.RawString(prefix)
		{
			out.RawByte('[')
			for v33, v34 := range in.ReadQueue {
				if v33 > 0 {
					out.RawByte(',')
				}
				easyjson28548e40EncodeGithubComQuicGoQuicGoQstate9(out, v34)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}
func easyjson28548e40DecodeGithubComQuicGoQuicGoQstate9(in *jlexer.Lexer, out *StreamRange) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "offset":
			out.Offset = int64(in.Int64())
		case "data":
			if in.IsNull() {
				in.Skip()
				out.Data = nil
			} else {
				out.Data = in.Bytes()
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson28548e40EncodeGithubComQuicGoQuicGoQstate9(out *jwriter.Writer, in StreamRange) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"offset\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.Offset))
	}
	{
		const prefix string = ",\"data\":"
		out.RawString(prefix)
		out.Base64Bytes(in.Data)
	}
	out.RawByte('}')
}
func easyjson28548e40DecodeGithubComQuicGoQuicGoQstate5(in *jlexer.Lexer, out *Parameters) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "initial_max_stream_data_bidi_local":
			if in.IsNull() {
				in.Skip()
				out.InitialMaxStreamDataBidiLocal = nil
			} else {
				if out.InitialMaxStreamDataBidiLocal == nil {
					out.InitialMaxStreamDataBidiLocal = new(int64)
				}
				*out.InitialMaxStreamDataBidiLocal = int64(in.Int64())
			}
		case "initial_max_stream_data_bidi_remote":
			if in.IsNull() {
				in.Skip()
				out.InitialMaxStreamDataBidiRemote = nil
			} else {
				if out.InitialMaxStreamDataBidiRemote == nil {
					out.InitialMaxStreamDataBidiRemote = new(int64)
				}
				*out.InitialMaxStreamDataBidiRemote = int64(in.Int64())
			}
		case "initial_max_stream_data_uni":
			if in.IsNull() {
				in.Skip()
				out.InitialMaxStreamDataUni = nil
			} else {
				if out.InitialMaxStreamDataUni == nil {
					out.InitialMaxStreamDataUni = new(int64)
				}
				*out.InitialMaxStreamDataUni = int64(in.Int64())
			}
		case "max_ack_delay":
			if in.IsNull() {
				in.Skip()
				out.MaxAckDelay = nil
			} else {
				if out.MaxAckDelay == nil {
					out.MaxAckDelay = new(int64)
				}
				*out.MaxAckDelay = int64(in.Int64())
			}
		case "ack_delay_exponent":
			if in.IsNull() {
				in.Skip()
				out.AckDelayExponent = nil
			} else {
				if out.AckDelayExponent == nil {
					out.AckDelayExponent = new(uint8)
				}
				*out.AckDelayExponent = uint8(in.Uint8())
			}
		case "disable_active_migration":
			if in.IsNull() {
				in.Skip()
				out.DisableActiveMigration = nil
			} else {
				if out.DisableActiveMigration == nil {
					out.DisableActiveMigration = new(bool)
				}
				*out.DisableActiveMigration = bool(in.Bool())
			}
		case "max_udp_payload_size":
			if in.IsNull() {
				in.Skip()
				out.MaxUDPPayloadSize = nil
			} else {
				if out.MaxUDPPayloadSize == nil {
					out.MaxUDPPayloadSize = new(int64)
				}
				*out.MaxUDPPayloadSize = int64(in.Int64())
			}
		case "original_destination_connection_id":
			if in.IsNull() {
				in.Skip()
				out.OriginalDestinationConnectionID = nil
			} else {
				if out.OriginalDestinationConnectionID == nil {
					out.OriginalDestinationConnectionID = new(HexByteSlice)
				}
				if data := in.Raw(); in.Ok() {
					in.AddError((*out.OriginalDestinationConnectionID).UnmarshalJSON(data))
				}
			}
		case "active_connection_id_limit":
			out.ActiveConnectionIDLimit = uint64(in.Uint64())
		case "max_datagram_frame_size":
			if in.IsNull() {
				in.Skip()
				out.MaxDatagramFrameSize = nil
			} else {
				if out.MaxDatagramFrameSize == nil {
					out.MaxDatagramFrameSize = new(int64)
				}
				*out.MaxDatagramFrameSize = int64(in.Int64())
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson28548e40EncodeGithubComQuicGoQuicGoQstate5(out *jwriter.Writer, in Parameters) {
	out.RawByte('{')
	first := true
	_ = first
	if in.InitialMaxStreamDataBidiLocal != nil {
		const prefix string = ",\"initial_max_stream_data_bidi_local\":"
		first = false
		out.RawString(prefix[1:])
		out.Int64(int64(*in.InitialMaxStreamDataBidiLocal))
	}
	if in.InitialMaxStreamDataBidiRemote != nil {
		const prefix string = ",\"initial_max_stream_data_bidi_remote\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(*in.InitialMaxStreamDataBidiRemote))
	}
	if in.InitialMaxStreamDataUni != nil {
		const prefix string = ",\"initial_max_stream_data_uni\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(*in.InitialMaxStreamDataUni))
	}
	if in.MaxAckDelay != nil {
		const prefix string = ",\"max_ack_delay\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(*in.MaxAckDelay))
	}
	if in.AckDelayExponent != nil {
		const prefix string = ",\"ack_delay_exponent\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint8(uint8(*in.AckDelayExponent))
	}
	if in.DisableActiveMigration != nil {
		const prefix string = ",\"disable_active_migration\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(*in.DisableActiveMigration))
	}
	if in.MaxUDPPayloadSize != nil {
		const prefix string = ",\"max_udp_payload_size\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(*in.MaxUDPPayloadSize))
	}
	if in.OriginalDestinationConnectionID != nil {
		const prefix string = ",\"original_destination_connection_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Raw((*in.OriginalDestinationConnectionID).MarshalJSON())
	}
	if in.ActiveConnectionIDLimit != 0 {
		const prefix string = ",\"active_connection_id_limit\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint64(uint64(in.ActiveConnectionIDLimit))
	}
	if in.MaxDatagramFrameSize != nil {
		const prefix string = ",\"max_datagram_frame_size\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(*in.MaxDatagramFrameSize))
	}
	out.RawByte('}')
}
func easyjson28548e40DecodeGithubComQuicGoQuicGoQstate4(in *jlexer.Lexer, out *ConnectionID) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "sequence_number":
			out.SequenceNumber = uint64(in.Uint64())
		case "connection_id":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.ConnectionID).UnmarshalJSON(data))
			}
		case "stateless_reset_token":
			if in.IsNull() {
				in.Skip()
				out.StatelessResetToken = nil
			} else {
				if out.StatelessResetToken == nil {
					out.StatelessResetToken = new(StatelessResetToken)
				}
				if data := in.Raw(); in.Ok() {
					in.AddError((*out.StatelessResetToken).UnmarshalJSON(data))
				}
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson28548e40EncodeGithubComQuicGoQuicGoQstate4(out *jwriter.Writer, in ConnectionID) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"sequence_number\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.SequenceNumber))
	}
	{
		const prefix string = ",\"connection_id\":"
		out.RawString(prefix)
		out.Raw((in.ConnectionID).MarshalJSON())
	}
	{
		const prefix string = ",\"stateless_reset_token\":"
		out.RawString(prefix)
		if in.StatelessResetToken == nil {
			out.RawString("null")
		} else {
			out.Raw((*in.StatelessResetToken).MarshalJSON())
		}
	}
	out.RawByte('}')
}
