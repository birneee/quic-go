package qstate

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Parameters) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "initial_max_stream_data_bidi_local":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					err = msgp.WrapError(err, "InitialMaxStreamDataBidiLocal")
					return
				}
				z.InitialMaxStreamDataBidiLocal = nil
			} else {
				if z.InitialMaxStreamDataBidiLocal == nil {
					z.InitialMaxStreamDataBidiLocal = new(int64)
				}
				*z.InitialMaxStreamDataBidiLocal, err = dc.ReadInt64()
				if err != nil {
					err = msgp.WrapError(err, "InitialMaxStreamDataBidiLocal")
					return
				}
			}
		case "initial_max_stream_data_bidi_remote":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					err = msgp.WrapError(err, "InitialMaxStreamDataBidiRemote")
					return
				}
				z.InitialMaxStreamDataBidiRemote = nil
			} else {
				if z.InitialMaxStreamDataBidiRemote == nil {
					z.InitialMaxStreamDataBidiRemote = new(int64)
				}
				*z.InitialMaxStreamDataBidiRemote, err = dc.ReadInt64()
				if err != nil {
					err = msgp.WrapError(err, "InitialMaxStreamDataBidiRemote")
					return
				}
			}
		case "initial_max_stream_data_uni":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					err = msgp.WrapError(err, "InitialMaxStreamDataUni")
					return
				}
				z.InitialMaxStreamDataUni = nil
			} else {
				if z.InitialMaxStreamDataUni == nil {
					z.InitialMaxStreamDataUni = new(int64)
				}
				*z.InitialMaxStreamDataUni, err = dc.ReadInt64()
				if err != nil {
					err = msgp.WrapError(err, "InitialMaxStreamDataUni")
					return
				}
			}
		case "max_ack_delay":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					err = msgp.WrapError(err, "MaxAckDelay")
					return
				}
				z.MaxAckDelay = nil
			} else {
				if z.MaxAckDelay == nil {
					z.MaxAckDelay = new(int64)
				}
				*z.MaxAckDelay, err = dc.ReadInt64()
				if err != nil {
					err = msgp.WrapError(err, "MaxAckDelay")
					return
				}
			}
		case "ack_delay_exponent":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					err = msgp.WrapError(err, "AckDelayExponent")
					return
				}
				z.AckDelayExponent = nil
			} else {
				if z.AckDelayExponent == nil {
					z.AckDelayExponent = new(uint8)
				}
				*z.AckDelayExponent, err = dc.ReadUint8()
				if err != nil {
					err = msgp.WrapError(err, "AckDelayExponent")
					return
				}
			}
		case "disable_active_migration":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					err = msgp.WrapError(err, "DisableActiveMigration")
					return
				}
				z.DisableActiveMigration = nil
			} else {
				if z.DisableActiveMigration == nil {
					z.DisableActiveMigration = new(bool)
				}
				*z.DisableActiveMigration, err = dc.ReadBool()
				if err != nil {
					err = msgp.WrapError(err, "DisableActiveMigration")
					return
				}
			}
		case "max_udp_payload_size":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					err = msgp.WrapError(err, "MaxUDPPayloadSize")
					return
				}
				z.MaxUDPPayloadSize = nil
			} else {
				if z.MaxUDPPayloadSize == nil {
					z.MaxUDPPayloadSize = new(int64)
				}
				*z.MaxUDPPayloadSize, err = dc.ReadInt64()
				if err != nil {
					err = msgp.WrapError(err, "MaxUDPPayloadSize")
					return
				}
			}
		case "max_idle_timeout":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					err = msgp.WrapError(err, "MaxIdleTimeout")
					return
				}
				z.MaxIdleTimeout = nil
			} else {
				if z.MaxIdleTimeout == nil {
					z.MaxIdleTimeout = new(int64)
				}
				*z.MaxIdleTimeout, err = dc.ReadInt64()
				if err != nil {
					err = msgp.WrapError(err, "MaxIdleTimeout")
					return
				}
			}
		case "OriginalDestinationConnectionID":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					err = msgp.WrapError(err, "OriginalDestinationConnectionID")
					return
				}
				z.OriginalDestinationConnectionID = nil
			} else {
				if z.OriginalDestinationConnectionID == nil {
					z.OriginalDestinationConnectionID = new([]byte)
				}
				*z.OriginalDestinationConnectionID, err = dc.ReadBytes(*z.OriginalDestinationConnectionID)
				if err != nil {
					err = msgp.WrapError(err, "OriginalDestinationConnectionID")
					return
				}
			}
		case "active_connection_id_limit":
			z.ActiveConnectionIDLimit, err = dc.ReadUint64()
			if err != nil {
				err = msgp.WrapError(err, "ActiveConnectionIDLimit")
				return
			}
		case "max_datagram_frame_size":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					err = msgp.WrapError(err, "MaxDatagramFrameSize")
					return
				}
				z.MaxDatagramFrameSize = nil
			} else {
				if z.MaxDatagramFrameSize == nil {
					z.MaxDatagramFrameSize = new(int64)
				}
				*z.MaxDatagramFrameSize, err = dc.ReadInt64()
				if err != nil {
					err = msgp.WrapError(err, "MaxDatagramFrameSize")
					return
				}
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *Parameters) EncodeMsg(en *msgp.Writer) (err error) {
	// omitempty: check for empty values
	zb0001Len := uint32(11)
	var zb0001Mask uint16 /* 11 bits */
	_ = zb0001Mask
	if z.InitialMaxStreamDataBidiLocal == nil {
		zb0001Len--
		zb0001Mask |= 0x1
	}
	if z.InitialMaxStreamDataBidiRemote == nil {
		zb0001Len--
		zb0001Mask |= 0x2
	}
	if z.InitialMaxStreamDataUni == nil {
		zb0001Len--
		zb0001Mask |= 0x4
	}
	if z.MaxAckDelay == nil {
		zb0001Len--
		zb0001Mask |= 0x8
	}
	if z.AckDelayExponent == nil {
		zb0001Len--
		zb0001Mask |= 0x10
	}
	if z.DisableActiveMigration == nil {
		zb0001Len--
		zb0001Mask |= 0x20
	}
	if z.MaxUDPPayloadSize == nil {
		zb0001Len--
		zb0001Mask |= 0x40
	}
	if z.MaxIdleTimeout == nil {
		zb0001Len--
		zb0001Mask |= 0x80
	}
	if z.ActiveConnectionIDLimit == 0 {
		zb0001Len--
		zb0001Mask |= 0x200
	}
	if z.MaxDatagramFrameSize == nil {
		zb0001Len--
		zb0001Mask |= 0x400
	}
	// variable map header, size zb0001Len
	err = en.Append(0x80 | uint8(zb0001Len))
	if err != nil {
		return
	}
	if zb0001Len == 0 {
		return
	}
	if (zb0001Mask & 0x1) == 0 { // if not empty
		// write "initial_max_stream_data_bidi_local"
		err = en.Append(0xd9, 0x22, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x5f, 0x6d, 0x61, 0x78, 0x5f, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x62, 0x69, 0x64, 0x69, 0x5f, 0x6c, 0x6f, 0x63, 0x61, 0x6c)
		if err != nil {
			return
		}
		if z.InitialMaxStreamDataBidiLocal == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			err = en.WriteInt64(*z.InitialMaxStreamDataBidiLocal)
			if err != nil {
				err = msgp.WrapError(err, "InitialMaxStreamDataBidiLocal")
				return
			}
		}
	}
	if (zb0001Mask & 0x2) == 0 { // if not empty
		// write "initial_max_stream_data_bidi_remote"
		err = en.Append(0xd9, 0x23, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x5f, 0x6d, 0x61, 0x78, 0x5f, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x62, 0x69, 0x64, 0x69, 0x5f, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65)
		if err != nil {
			return
		}
		if z.InitialMaxStreamDataBidiRemote == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			err = en.WriteInt64(*z.InitialMaxStreamDataBidiRemote)
			if err != nil {
				err = msgp.WrapError(err, "InitialMaxStreamDataBidiRemote")
				return
			}
		}
	}
	if (zb0001Mask & 0x4) == 0 { // if not empty
		// write "initial_max_stream_data_uni"
		err = en.Append(0xbb, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x5f, 0x6d, 0x61, 0x78, 0x5f, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x75, 0x6e, 0x69)
		if err != nil {
			return
		}
		if z.InitialMaxStreamDataUni == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			err = en.WriteInt64(*z.InitialMaxStreamDataUni)
			if err != nil {
				err = msgp.WrapError(err, "InitialMaxStreamDataUni")
				return
			}
		}
	}
	if (zb0001Mask & 0x8) == 0 { // if not empty
		// write "max_ack_delay"
		err = en.Append(0xad, 0x6d, 0x61, 0x78, 0x5f, 0x61, 0x63, 0x6b, 0x5f, 0x64, 0x65, 0x6c, 0x61, 0x79)
		if err != nil {
			return
		}
		if z.MaxAckDelay == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			err = en.WriteInt64(*z.MaxAckDelay)
			if err != nil {
				err = msgp.WrapError(err, "MaxAckDelay")
				return
			}
		}
	}
	if (zb0001Mask & 0x10) == 0 { // if not empty
		// write "ack_delay_exponent"
		err = en.Append(0xb2, 0x61, 0x63, 0x6b, 0x5f, 0x64, 0x65, 0x6c, 0x61, 0x79, 0x5f, 0x65, 0x78, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74)
		if err != nil {
			return
		}
		if z.AckDelayExponent == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			err = en.WriteUint8(*z.AckDelayExponent)
			if err != nil {
				err = msgp.WrapError(err, "AckDelayExponent")
				return
			}
		}
	}
	if (zb0001Mask & 0x20) == 0 { // if not empty
		// write "disable_active_migration"
		err = en.Append(0xb8, 0x64, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x6d, 0x69, 0x67, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e)
		if err != nil {
			return
		}
		if z.DisableActiveMigration == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			err = en.WriteBool(*z.DisableActiveMigration)
			if err != nil {
				err = msgp.WrapError(err, "DisableActiveMigration")
				return
			}
		}
	}
	if (zb0001Mask & 0x40) == 0 { // if not empty
		// write "max_udp_payload_size"
		err = en.Append(0xb4, 0x6d, 0x61, 0x78, 0x5f, 0x75, 0x64, 0x70, 0x5f, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x73, 0x69, 0x7a, 0x65)
		if err != nil {
			return
		}
		if z.MaxUDPPayloadSize == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			err = en.WriteInt64(*z.MaxUDPPayloadSize)
			if err != nil {
				err = msgp.WrapError(err, "MaxUDPPayloadSize")
				return
			}
		}
	}
	if (zb0001Mask & 0x80) == 0 { // if not empty
		// write "max_idle_timeout"
		err = en.Append(0xb0, 0x6d, 0x61, 0x78, 0x5f, 0x69, 0x64, 0x6c, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74)
		if err != nil {
			return
		}
		if z.MaxIdleTimeout == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			err = en.WriteInt64(*z.MaxIdleTimeout)
			if err != nil {
				err = msgp.WrapError(err, "MaxIdleTimeout")
				return
			}
		}
	}
	// write "OriginalDestinationConnectionID"
	err = en.Append(0xbf, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x44, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x44)
	if err != nil {
		return
	}
	if z.OriginalDestinationConnectionID == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = en.WriteBytes(*z.OriginalDestinationConnectionID)
		if err != nil {
			err = msgp.WrapError(err, "OriginalDestinationConnectionID")
			return
		}
	}
	if (zb0001Mask & 0x200) == 0 { // if not empty
		// write "active_connection_id_limit"
		err = en.Append(0xba, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74)
		if err != nil {
			return
		}
		err = en.WriteUint64(z.ActiveConnectionIDLimit)
		if err != nil {
			err = msgp.WrapError(err, "ActiveConnectionIDLimit")
			return
		}
	}
	if (zb0001Mask & 0x400) == 0 { // if not empty
		// write "max_datagram_frame_size"
		err = en.Append(0xb7, 0x6d, 0x61, 0x78, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x67, 0x72, 0x61, 0x6d, 0x5f, 0x66, 0x72, 0x61, 0x6d, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65)
		if err != nil {
			return
		}
		if z.MaxDatagramFrameSize == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			err = en.WriteInt64(*z.MaxDatagramFrameSize)
			if err != nil {
				err = msgp.WrapError(err, "MaxDatagramFrameSize")
				return
			}
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Parameters) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// omitempty: check for empty values
	zb0001Len := uint32(11)
	var zb0001Mask uint16 /* 11 bits */
	_ = zb0001Mask
	if z.InitialMaxStreamDataBidiLocal == nil {
		zb0001Len--
		zb0001Mask |= 0x1
	}
	if z.InitialMaxStreamDataBidiRemote == nil {
		zb0001Len--
		zb0001Mask |= 0x2
	}
	if z.InitialMaxStreamDataUni == nil {
		zb0001Len--
		zb0001Mask |= 0x4
	}
	if z.MaxAckDelay == nil {
		zb0001Len--
		zb0001Mask |= 0x8
	}
	if z.AckDelayExponent == nil {
		zb0001Len--
		zb0001Mask |= 0x10
	}
	if z.DisableActiveMigration == nil {
		zb0001Len--
		zb0001Mask |= 0x20
	}
	if z.MaxUDPPayloadSize == nil {
		zb0001Len--
		zb0001Mask |= 0x40
	}
	if z.MaxIdleTimeout == nil {
		zb0001Len--
		zb0001Mask |= 0x80
	}
	if z.ActiveConnectionIDLimit == 0 {
		zb0001Len--
		zb0001Mask |= 0x200
	}
	if z.MaxDatagramFrameSize == nil {
		zb0001Len--
		zb0001Mask |= 0x400
	}
	// variable map header, size zb0001Len
	o = append(o, 0x80|uint8(zb0001Len))
	if zb0001Len == 0 {
		return
	}
	if (zb0001Mask & 0x1) == 0 { // if not empty
		// string "initial_max_stream_data_bidi_local"
		o = append(o, 0xd9, 0x22, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x5f, 0x6d, 0x61, 0x78, 0x5f, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x62, 0x69, 0x64, 0x69, 0x5f, 0x6c, 0x6f, 0x63, 0x61, 0x6c)
		if z.InitialMaxStreamDataBidiLocal == nil {
			o = msgp.AppendNil(o)
		} else {
			o = msgp.AppendInt64(o, *z.InitialMaxStreamDataBidiLocal)
		}
	}
	if (zb0001Mask & 0x2) == 0 { // if not empty
		// string "initial_max_stream_data_bidi_remote"
		o = append(o, 0xd9, 0x23, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x5f, 0x6d, 0x61, 0x78, 0x5f, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x62, 0x69, 0x64, 0x69, 0x5f, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65)
		if z.InitialMaxStreamDataBidiRemote == nil {
			o = msgp.AppendNil(o)
		} else {
			o = msgp.AppendInt64(o, *z.InitialMaxStreamDataBidiRemote)
		}
	}
	if (zb0001Mask & 0x4) == 0 { // if not empty
		// string "initial_max_stream_data_uni"
		o = append(o, 0xbb, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x5f, 0x6d, 0x61, 0x78, 0x5f, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x75, 0x6e, 0x69)
		if z.InitialMaxStreamDataUni == nil {
			o = msgp.AppendNil(o)
		} else {
			o = msgp.AppendInt64(o, *z.InitialMaxStreamDataUni)
		}
	}
	if (zb0001Mask & 0x8) == 0 { // if not empty
		// string "max_ack_delay"
		o = append(o, 0xad, 0x6d, 0x61, 0x78, 0x5f, 0x61, 0x63, 0x6b, 0x5f, 0x64, 0x65, 0x6c, 0x61, 0x79)
		if z.MaxAckDelay == nil {
			o = msgp.AppendNil(o)
		} else {
			o = msgp.AppendInt64(o, *z.MaxAckDelay)
		}
	}
	if (zb0001Mask & 0x10) == 0 { // if not empty
		// string "ack_delay_exponent"
		o = append(o, 0xb2, 0x61, 0x63, 0x6b, 0x5f, 0x64, 0x65, 0x6c, 0x61, 0x79, 0x5f, 0x65, 0x78, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74)
		if z.AckDelayExponent == nil {
			o = msgp.AppendNil(o)
		} else {
			o = msgp.AppendUint8(o, *z.AckDelayExponent)
		}
	}
	if (zb0001Mask & 0x20) == 0 { // if not empty
		// string "disable_active_migration"
		o = append(o, 0xb8, 0x64, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x6d, 0x69, 0x67, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e)
		if z.DisableActiveMigration == nil {
			o = msgp.AppendNil(o)
		} else {
			o = msgp.AppendBool(o, *z.DisableActiveMigration)
		}
	}
	if (zb0001Mask & 0x40) == 0 { // if not empty
		// string "max_udp_payload_size"
		o = append(o, 0xb4, 0x6d, 0x61, 0x78, 0x5f, 0x75, 0x64, 0x70, 0x5f, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x73, 0x69, 0x7a, 0x65)
		if z.MaxUDPPayloadSize == nil {
			o = msgp.AppendNil(o)
		} else {
			o = msgp.AppendInt64(o, *z.MaxUDPPayloadSize)
		}
	}
	if (zb0001Mask & 0x80) == 0 { // if not empty
		// string "max_idle_timeout"
		o = append(o, 0xb0, 0x6d, 0x61, 0x78, 0x5f, 0x69, 0x64, 0x6c, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74)
		if z.MaxIdleTimeout == nil {
			o = msgp.AppendNil(o)
		} else {
			o = msgp.AppendInt64(o, *z.MaxIdleTimeout)
		}
	}
	// string "OriginalDestinationConnectionID"
	o = append(o, 0xbf, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x44, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x44)
	if z.OriginalDestinationConnectionID == nil {
		o = msgp.AppendNil(o)
	} else {
		o = msgp.AppendBytes(o, *z.OriginalDestinationConnectionID)
	}
	if (zb0001Mask & 0x200) == 0 { // if not empty
		// string "active_connection_id_limit"
		o = append(o, 0xba, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74)
		o = msgp.AppendUint64(o, z.ActiveConnectionIDLimit)
	}
	if (zb0001Mask & 0x400) == 0 { // if not empty
		// string "max_datagram_frame_size"
		o = append(o, 0xb7, 0x6d, 0x61, 0x78, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x67, 0x72, 0x61, 0x6d, 0x5f, 0x66, 0x72, 0x61, 0x6d, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65)
		if z.MaxDatagramFrameSize == nil {
			o = msgp.AppendNil(o)
		} else {
			o = msgp.AppendInt64(o, *z.MaxDatagramFrameSize)
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Parameters) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "initial_max_stream_data_bidi_local":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.InitialMaxStreamDataBidiLocal = nil
			} else {
				if z.InitialMaxStreamDataBidiLocal == nil {
					z.InitialMaxStreamDataBidiLocal = new(int64)
				}
				*z.InitialMaxStreamDataBidiLocal, bts, err = msgp.ReadInt64Bytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "InitialMaxStreamDataBidiLocal")
					return
				}
			}
		case "initial_max_stream_data_bidi_remote":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.InitialMaxStreamDataBidiRemote = nil
			} else {
				if z.InitialMaxStreamDataBidiRemote == nil {
					z.InitialMaxStreamDataBidiRemote = new(int64)
				}
				*z.InitialMaxStreamDataBidiRemote, bts, err = msgp.ReadInt64Bytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "InitialMaxStreamDataBidiRemote")
					return
				}
			}
		case "initial_max_stream_data_uni":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.InitialMaxStreamDataUni = nil
			} else {
				if z.InitialMaxStreamDataUni == nil {
					z.InitialMaxStreamDataUni = new(int64)
				}
				*z.InitialMaxStreamDataUni, bts, err = msgp.ReadInt64Bytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "InitialMaxStreamDataUni")
					return
				}
			}
		case "max_ack_delay":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.MaxAckDelay = nil
			} else {
				if z.MaxAckDelay == nil {
					z.MaxAckDelay = new(int64)
				}
				*z.MaxAckDelay, bts, err = msgp.ReadInt64Bytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "MaxAckDelay")
					return
				}
			}
		case "ack_delay_exponent":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.AckDelayExponent = nil
			} else {
				if z.AckDelayExponent == nil {
					z.AckDelayExponent = new(uint8)
				}
				*z.AckDelayExponent, bts, err = msgp.ReadUint8Bytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "AckDelayExponent")
					return
				}
			}
		case "disable_active_migration":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.DisableActiveMigration = nil
			} else {
				if z.DisableActiveMigration == nil {
					z.DisableActiveMigration = new(bool)
				}
				*z.DisableActiveMigration, bts, err = msgp.ReadBoolBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "DisableActiveMigration")
					return
				}
			}
		case "max_udp_payload_size":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.MaxUDPPayloadSize = nil
			} else {
				if z.MaxUDPPayloadSize == nil {
					z.MaxUDPPayloadSize = new(int64)
				}
				*z.MaxUDPPayloadSize, bts, err = msgp.ReadInt64Bytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "MaxUDPPayloadSize")
					return
				}
			}
		case "max_idle_timeout":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.MaxIdleTimeout = nil
			} else {
				if z.MaxIdleTimeout == nil {
					z.MaxIdleTimeout = new(int64)
				}
				*z.MaxIdleTimeout, bts, err = msgp.ReadInt64Bytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "MaxIdleTimeout")
					return
				}
			}
		case "OriginalDestinationConnectionID":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.OriginalDestinationConnectionID = nil
			} else {
				if z.OriginalDestinationConnectionID == nil {
					z.OriginalDestinationConnectionID = new([]byte)
				}
				*z.OriginalDestinationConnectionID, bts, err = msgp.ReadBytesBytes(bts, *z.OriginalDestinationConnectionID)
				if err != nil {
					err = msgp.WrapError(err, "OriginalDestinationConnectionID")
					return
				}
			}
		case "active_connection_id_limit":
			z.ActiveConnectionIDLimit, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "ActiveConnectionIDLimit")
				return
			}
		case "max_datagram_frame_size":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.MaxDatagramFrameSize = nil
			} else {
				if z.MaxDatagramFrameSize == nil {
					z.MaxDatagramFrameSize = new(int64)
				}
				*z.MaxDatagramFrameSize, bts, err = msgp.ReadInt64Bytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "MaxDatagramFrameSize")
					return
				}
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *Parameters) Msgsize() (s int) {
	s = 1 + 36
	if z.InitialMaxStreamDataBidiLocal == nil {
		s += msgp.NilSize
	} else {
		s += msgp.Int64Size
	}
	s += 37
	if z.InitialMaxStreamDataBidiRemote == nil {
		s += msgp.NilSize
	} else {
		s += msgp.Int64Size
	}
	s += 28
	if z.InitialMaxStreamDataUni == nil {
		s += msgp.NilSize
	} else {
		s += msgp.Int64Size
	}
	s += 14
	if z.MaxAckDelay == nil {
		s += msgp.NilSize
	} else {
		s += msgp.Int64Size
	}
	s += 19
	if z.AckDelayExponent == nil {
		s += msgp.NilSize
	} else {
		s += msgp.Uint8Size
	}
	s += 25
	if z.DisableActiveMigration == nil {
		s += msgp.NilSize
	} else {
		s += msgp.BoolSize
	}
	s += 21
	if z.MaxUDPPayloadSize == nil {
		s += msgp.NilSize
	} else {
		s += msgp.Int64Size
	}
	s += 17
	if z.MaxIdleTimeout == nil {
		s += msgp.NilSize
	} else {
		s += msgp.Int64Size
	}
	s += 32
	if z.OriginalDestinationConnectionID == nil {
		s += msgp.NilSize
	} else {
		s += msgp.BytesPrefixSize + len(*z.OriginalDestinationConnectionID)
	}
	s += 27 + msgp.Uint64Size + 24
	if z.MaxDatagramFrameSize == nil {
		s += msgp.NilSize
	} else {
		s += msgp.Int64Size
	}
	return
}