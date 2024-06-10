package qstate

var Codecs = map[string]Codec[*Connection]{
	"std_json":       StdJsonCodec[*Connection]{},
	"msgp":           MsgpCodec[*Connection]{},
	"msgp_zstd":      NewZstdCodec[*Connection](MsgpCodec[*Connection]{}),
	"std_json_zstd":  NewZstdCodec[*Connection](StdJsonCodec[*Connection]{}),
	"easy_json":      EasyJsonCodec[*Connection]{},
	"easy_json_zstd": NewZstdCodec[*Connection](EasyJsonCodec[*Connection]{}),
	"cbor":           CborCodec[*Connection]{},
	"cbor_zstd":      NewZstdCodec[*Connection](CborCodec[*Connection]{}),
}
