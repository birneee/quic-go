package qstate

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Stream) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "stream_id":
			z.StreamID, err = dc.ReadInt64()
			if err != nil {
				err = msgp.WrapError(err, "StreamID")
				return
			}
		case "write_offset":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					err = msgp.WrapError(err, "WriteOffset")
					return
				}
				z.WriteOffset = nil
			} else {
				if z.WriteOffset == nil {
					z.WriteOffset = new(int64)
				}
				*z.WriteOffset, err = dc.ReadInt64()
				if err != nil {
					err = msgp.WrapError(err, "WriteOffset")
					return
				}
			}
		case "write_fin":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					err = msgp.WrapError(err, "WriteFin")
					return
				}
				z.WriteFin = nil
			} else {
				if z.WriteFin == nil {
					z.WriteFin = new(int64)
				}
				*z.WriteFin, err = dc.ReadInt64()
				if err != nil {
					err = msgp.WrapError(err, "WriteFin")
					return
				}
			}
		case "write_max_data":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					err = msgp.WrapError(err, "WriteMaxData")
					return
				}
				z.WriteMaxData = nil
			} else {
				if z.WriteMaxData == nil {
					z.WriteMaxData = new(int64)
				}
				*z.WriteMaxData, err = dc.ReadInt64()
				if err != nil {
					err = msgp.WrapError(err, "WriteMaxData")
					return
				}
			}
		case "write_ack":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					err = msgp.WrapError(err, "WriteAck")
					return
				}
				z.WriteAck = nil
			} else {
				if z.WriteAck == nil {
					z.WriteAck = new(int64)
				}
				*z.WriteAck, err = dc.ReadInt64()
				if err != nil {
					err = msgp.WrapError(err, "WriteAck")
					return
				}
			}
		case "write_queue":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				err = msgp.WrapError(err, "WriteQueue")
				return
			}
			if cap(z.WriteQueue) >= int(zb0002) {
				z.WriteQueue = (z.WriteQueue)[:zb0002]
			} else {
				z.WriteQueue = make([]StreamRange, zb0002)
			}
			for za0001 := range z.WriteQueue {
				err = z.WriteQueue[za0001].DecodeMsg(dc)
				if err != nil {
					err = msgp.WrapError(err, "WriteQueue", za0001)
					return
				}
			}
		case "read_offset":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					err = msgp.WrapError(err, "ReadOffset")
					return
				}
				z.ReadOffset = nil
			} else {
				if z.ReadOffset == nil {
					z.ReadOffset = new(int64)
				}
				*z.ReadOffset, err = dc.ReadInt64()
				if err != nil {
					err = msgp.WrapError(err, "ReadOffset")
					return
				}
			}
		case "read_fin":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					err = msgp.WrapError(err, "ReadFin")
					return
				}
				z.ReadFin = nil
			} else {
				if z.ReadFin == nil {
					z.ReadFin = new(int64)
				}
				*z.ReadFin, err = dc.ReadInt64()
				if err != nil {
					err = msgp.WrapError(err, "ReadFin")
					return
				}
			}
		case "read_max_data":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					err = msgp.WrapError(err, "ReadMaxData")
					return
				}
				z.ReadMaxData = nil
			} else {
				if z.ReadMaxData == nil {
					z.ReadMaxData = new(int64)
				}
				*z.ReadMaxData, err = dc.ReadInt64()
				if err != nil {
					err = msgp.WrapError(err, "ReadMaxData")
					return
				}
			}
		case "read_queue":
			var zb0003 uint32
			zb0003, err = dc.ReadArrayHeader()
			if err != nil {
				err = msgp.WrapError(err, "ReadQueue")
				return
			}
			if cap(z.ReadQueue) >= int(zb0003) {
				z.ReadQueue = (z.ReadQueue)[:zb0003]
			} else {
				z.ReadQueue = make([]StreamRange, zb0003)
			}
			for za0002 := range z.ReadQueue {
				err = z.ReadQueue[za0002].DecodeMsg(dc)
				if err != nil {
					err = msgp.WrapError(err, "ReadQueue", za0002)
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
func (z *Stream) EncodeMsg(en *msgp.Writer) (err error) {
	// omitempty: check for empty values
	zb0001Len := uint32(10)
	var zb0001Mask uint16 /* 10 bits */
	_ = zb0001Mask
	if z.WriteOffset == nil {
		zb0001Len--
		zb0001Mask |= 0x2
	}
	if z.WriteFin == nil {
		zb0001Len--
		zb0001Mask |= 0x4
	}
	if z.WriteMaxData == nil {
		zb0001Len--
		zb0001Mask |= 0x8
	}
	if z.WriteAck == nil {
		zb0001Len--
		zb0001Mask |= 0x10
	}
	if z.WriteQueue == nil {
		zb0001Len--
		zb0001Mask |= 0x20
	}
	if z.ReadOffset == nil {
		zb0001Len--
		zb0001Mask |= 0x40
	}
	if z.ReadFin == nil {
		zb0001Len--
		zb0001Mask |= 0x80
	}
	if z.ReadMaxData == nil {
		zb0001Len--
		zb0001Mask |= 0x100
	}
	if z.ReadQueue == nil {
		zb0001Len--
		zb0001Mask |= 0x200
	}
	// variable map header, size zb0001Len
	err = en.Append(0x80 | uint8(zb0001Len))
	if err != nil {
		return
	}
	if zb0001Len == 0 {
		return
	}
	// write "stream_id"
	err = en.Append(0xa9, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x5f, 0x69, 0x64)
	if err != nil {
		return
	}
	err = en.WriteInt64(z.StreamID)
	if err != nil {
		err = msgp.WrapError(err, "StreamID")
		return
	}
	if (zb0001Mask & 0x2) == 0 { // if not empty
		// write "write_offset"
		err = en.Append(0xac, 0x77, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74)
		if err != nil {
			return
		}
		if z.WriteOffset == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			err = en.WriteInt64(*z.WriteOffset)
			if err != nil {
				err = msgp.WrapError(err, "WriteOffset")
				return
			}
		}
	}
	if (zb0001Mask & 0x4) == 0 { // if not empty
		// write "write_fin"
		err = en.Append(0xa9, 0x77, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x66, 0x69, 0x6e)
		if err != nil {
			return
		}
		if z.WriteFin == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			err = en.WriteInt64(*z.WriteFin)
			if err != nil {
				err = msgp.WrapError(err, "WriteFin")
				return
			}
		}
	}
	if (zb0001Mask & 0x8) == 0 { // if not empty
		// write "write_max_data"
		err = en.Append(0xae, 0x77, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x6d, 0x61, 0x78, 0x5f, 0x64, 0x61, 0x74, 0x61)
		if err != nil {
			return
		}
		if z.WriteMaxData == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			err = en.WriteInt64(*z.WriteMaxData)
			if err != nil {
				err = msgp.WrapError(err, "WriteMaxData")
				return
			}
		}
	}
	if (zb0001Mask & 0x10) == 0 { // if not empty
		// write "write_ack"
		err = en.Append(0xa9, 0x77, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x61, 0x63, 0x6b)
		if err != nil {
			return
		}
		if z.WriteAck == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			err = en.WriteInt64(*z.WriteAck)
			if err != nil {
				err = msgp.WrapError(err, "WriteAck")
				return
			}
		}
	}
	if (zb0001Mask & 0x20) == 0 { // if not empty
		// write "write_queue"
		err = en.Append(0xab, 0x77, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x71, 0x75, 0x65, 0x75, 0x65)
		if err != nil {
			return
		}
		err = en.WriteArrayHeader(uint32(len(z.WriteQueue)))
		if err != nil {
			err = msgp.WrapError(err, "WriteQueue")
			return
		}
		for za0001 := range z.WriteQueue {
			err = z.WriteQueue[za0001].EncodeMsg(en)
			if err != nil {
				err = msgp.WrapError(err, "WriteQueue", za0001)
				return
			}
		}
	}
	if (zb0001Mask & 0x40) == 0 { // if not empty
		// write "read_offset"
		err = en.Append(0xab, 0x72, 0x65, 0x61, 0x64, 0x5f, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74)
		if err != nil {
			return
		}
		if z.ReadOffset == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			err = en.WriteInt64(*z.ReadOffset)
			if err != nil {
				err = msgp.WrapError(err, "ReadOffset")
				return
			}
		}
	}
	if (zb0001Mask & 0x80) == 0 { // if not empty
		// write "read_fin"
		err = en.Append(0xa8, 0x72, 0x65, 0x61, 0x64, 0x5f, 0x66, 0x69, 0x6e)
		if err != nil {
			return
		}
		if z.ReadFin == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			err = en.WriteInt64(*z.ReadFin)
			if err != nil {
				err = msgp.WrapError(err, "ReadFin")
				return
			}
		}
	}
	if (zb0001Mask & 0x100) == 0 { // if not empty
		// write "read_max_data"
		err = en.Append(0xad, 0x72, 0x65, 0x61, 0x64, 0x5f, 0x6d, 0x61, 0x78, 0x5f, 0x64, 0x61, 0x74, 0x61)
		if err != nil {
			return
		}
		if z.ReadMaxData == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			err = en.WriteInt64(*z.ReadMaxData)
			if err != nil {
				err = msgp.WrapError(err, "ReadMaxData")
				return
			}
		}
	}
	if (zb0001Mask & 0x200) == 0 { // if not empty
		// write "read_queue"
		err = en.Append(0xaa, 0x72, 0x65, 0x61, 0x64, 0x5f, 0x71, 0x75, 0x65, 0x75, 0x65)
		if err != nil {
			return
		}
		err = en.WriteArrayHeader(uint32(len(z.ReadQueue)))
		if err != nil {
			err = msgp.WrapError(err, "ReadQueue")
			return
		}
		for za0002 := range z.ReadQueue {
			err = z.ReadQueue[za0002].EncodeMsg(en)
			if err != nil {
				err = msgp.WrapError(err, "ReadQueue", za0002)
				return
			}
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Stream) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// omitempty: check for empty values
	zb0001Len := uint32(10)
	var zb0001Mask uint16 /* 10 bits */
	_ = zb0001Mask
	if z.WriteOffset == nil {
		zb0001Len--
		zb0001Mask |= 0x2
	}
	if z.WriteFin == nil {
		zb0001Len--
		zb0001Mask |= 0x4
	}
	if z.WriteMaxData == nil {
		zb0001Len--
		zb0001Mask |= 0x8
	}
	if z.WriteAck == nil {
		zb0001Len--
		zb0001Mask |= 0x10
	}
	if z.WriteQueue == nil {
		zb0001Len--
		zb0001Mask |= 0x20
	}
	if z.ReadOffset == nil {
		zb0001Len--
		zb0001Mask |= 0x40
	}
	if z.ReadFin == nil {
		zb0001Len--
		zb0001Mask |= 0x80
	}
	if z.ReadMaxData == nil {
		zb0001Len--
		zb0001Mask |= 0x100
	}
	if z.ReadQueue == nil {
		zb0001Len--
		zb0001Mask |= 0x200
	}
	// variable map header, size zb0001Len
	o = append(o, 0x80|uint8(zb0001Len))
	if zb0001Len == 0 {
		return
	}
	// string "stream_id"
	o = append(o, 0xa9, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x5f, 0x69, 0x64)
	o = msgp.AppendInt64(o, z.StreamID)
	if (zb0001Mask & 0x2) == 0 { // if not empty
		// string "write_offset"
		o = append(o, 0xac, 0x77, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74)
		if z.WriteOffset == nil {
			o = msgp.AppendNil(o)
		} else {
			o = msgp.AppendInt64(o, *z.WriteOffset)
		}
	}
	if (zb0001Mask & 0x4) == 0 { // if not empty
		// string "write_fin"
		o = append(o, 0xa9, 0x77, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x66, 0x69, 0x6e)
		if z.WriteFin == nil {
			o = msgp.AppendNil(o)
		} else {
			o = msgp.AppendInt64(o, *z.WriteFin)
		}
	}
	if (zb0001Mask & 0x8) == 0 { // if not empty
		// string "write_max_data"
		o = append(o, 0xae, 0x77, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x6d, 0x61, 0x78, 0x5f, 0x64, 0x61, 0x74, 0x61)
		if z.WriteMaxData == nil {
			o = msgp.AppendNil(o)
		} else {
			o = msgp.AppendInt64(o, *z.WriteMaxData)
		}
	}
	if (zb0001Mask & 0x10) == 0 { // if not empty
		// string "write_ack"
		o = append(o, 0xa9, 0x77, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x61, 0x63, 0x6b)
		if z.WriteAck == nil {
			o = msgp.AppendNil(o)
		} else {
			o = msgp.AppendInt64(o, *z.WriteAck)
		}
	}
	if (zb0001Mask & 0x20) == 0 { // if not empty
		// string "write_queue"
		o = append(o, 0xab, 0x77, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x71, 0x75, 0x65, 0x75, 0x65)
		o = msgp.AppendArrayHeader(o, uint32(len(z.WriteQueue)))
		for za0001 := range z.WriteQueue {
			o, err = z.WriteQueue[za0001].MarshalMsg(o)
			if err != nil {
				err = msgp.WrapError(err, "WriteQueue", za0001)
				return
			}
		}
	}
	if (zb0001Mask & 0x40) == 0 { // if not empty
		// string "read_offset"
		o = append(o, 0xab, 0x72, 0x65, 0x61, 0x64, 0x5f, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74)
		if z.ReadOffset == nil {
			o = msgp.AppendNil(o)
		} else {
			o = msgp.AppendInt64(o, *z.ReadOffset)
		}
	}
	if (zb0001Mask & 0x80) == 0 { // if not empty
		// string "read_fin"
		o = append(o, 0xa8, 0x72, 0x65, 0x61, 0x64, 0x5f, 0x66, 0x69, 0x6e)
		if z.ReadFin == nil {
			o = msgp.AppendNil(o)
		} else {
			o = msgp.AppendInt64(o, *z.ReadFin)
		}
	}
	if (zb0001Mask & 0x100) == 0 { // if not empty
		// string "read_max_data"
		o = append(o, 0xad, 0x72, 0x65, 0x61, 0x64, 0x5f, 0x6d, 0x61, 0x78, 0x5f, 0x64, 0x61, 0x74, 0x61)
		if z.ReadMaxData == nil {
			o = msgp.AppendNil(o)
		} else {
			o = msgp.AppendInt64(o, *z.ReadMaxData)
		}
	}
	if (zb0001Mask & 0x200) == 0 { // if not empty
		// string "read_queue"
		o = append(o, 0xaa, 0x72, 0x65, 0x61, 0x64, 0x5f, 0x71, 0x75, 0x65, 0x75, 0x65)
		o = msgp.AppendArrayHeader(o, uint32(len(z.ReadQueue)))
		for za0002 := range z.ReadQueue {
			o, err = z.ReadQueue[za0002].MarshalMsg(o)
			if err != nil {
				err = msgp.WrapError(err, "ReadQueue", za0002)
				return
			}
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Stream) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "stream_id":
			z.StreamID, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "StreamID")
				return
			}
		case "write_offset":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.WriteOffset = nil
			} else {
				if z.WriteOffset == nil {
					z.WriteOffset = new(int64)
				}
				*z.WriteOffset, bts, err = msgp.ReadInt64Bytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "WriteOffset")
					return
				}
			}
		case "write_fin":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.WriteFin = nil
			} else {
				if z.WriteFin == nil {
					z.WriteFin = new(int64)
				}
				*z.WriteFin, bts, err = msgp.ReadInt64Bytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "WriteFin")
					return
				}
			}
		case "write_max_data":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.WriteMaxData = nil
			} else {
				if z.WriteMaxData == nil {
					z.WriteMaxData = new(int64)
				}
				*z.WriteMaxData, bts, err = msgp.ReadInt64Bytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "WriteMaxData")
					return
				}
			}
		case "write_ack":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.WriteAck = nil
			} else {
				if z.WriteAck == nil {
					z.WriteAck = new(int64)
				}
				*z.WriteAck, bts, err = msgp.ReadInt64Bytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "WriteAck")
					return
				}
			}
		case "write_queue":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "WriteQueue")
				return
			}
			if cap(z.WriteQueue) >= int(zb0002) {
				z.WriteQueue = (z.WriteQueue)[:zb0002]
			} else {
				z.WriteQueue = make([]StreamRange, zb0002)
			}
			for za0001 := range z.WriteQueue {
				bts, err = z.WriteQueue[za0001].UnmarshalMsg(bts)
				if err != nil {
					err = msgp.WrapError(err, "WriteQueue", za0001)
					return
				}
			}
		case "read_offset":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.ReadOffset = nil
			} else {
				if z.ReadOffset == nil {
					z.ReadOffset = new(int64)
				}
				*z.ReadOffset, bts, err = msgp.ReadInt64Bytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "ReadOffset")
					return
				}
			}
		case "read_fin":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.ReadFin = nil
			} else {
				if z.ReadFin == nil {
					z.ReadFin = new(int64)
				}
				*z.ReadFin, bts, err = msgp.ReadInt64Bytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "ReadFin")
					return
				}
			}
		case "read_max_data":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.ReadMaxData = nil
			} else {
				if z.ReadMaxData == nil {
					z.ReadMaxData = new(int64)
				}
				*z.ReadMaxData, bts, err = msgp.ReadInt64Bytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "ReadMaxData")
					return
				}
			}
		case "read_queue":
			var zb0003 uint32
			zb0003, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "ReadQueue")
				return
			}
			if cap(z.ReadQueue) >= int(zb0003) {
				z.ReadQueue = (z.ReadQueue)[:zb0003]
			} else {
				z.ReadQueue = make([]StreamRange, zb0003)
			}
			for za0002 := range z.ReadQueue {
				bts, err = z.ReadQueue[za0002].UnmarshalMsg(bts)
				if err != nil {
					err = msgp.WrapError(err, "ReadQueue", za0002)
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
func (z *Stream) Msgsize() (s int) {
	s = 1 + 10 + msgp.Int64Size + 13
	if z.WriteOffset == nil {
		s += msgp.NilSize
	} else {
		s += msgp.Int64Size
	}
	s += 10
	if z.WriteFin == nil {
		s += msgp.NilSize
	} else {
		s += msgp.Int64Size
	}
	s += 15
	if z.WriteMaxData == nil {
		s += msgp.NilSize
	} else {
		s += msgp.Int64Size
	}
	s += 10
	if z.WriteAck == nil {
		s += msgp.NilSize
	} else {
		s += msgp.Int64Size
	}
	s += 12 + msgp.ArrayHeaderSize
	for za0001 := range z.WriteQueue {
		s += z.WriteQueue[za0001].Msgsize()
	}
	s += 12
	if z.ReadOffset == nil {
		s += msgp.NilSize
	} else {
		s += msgp.Int64Size
	}
	s += 9
	if z.ReadFin == nil {
		s += msgp.NilSize
	} else {
		s += msgp.Int64Size
	}
	s += 14
	if z.ReadMaxData == nil {
		s += msgp.NilSize
	} else {
		s += msgp.Int64Size
	}
	s += 11 + msgp.ArrayHeaderSize
	for za0002 := range z.ReadQueue {
		s += z.ReadQueue[za0002].Msgsize()
	}
	return
}
