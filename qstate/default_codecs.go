package qstate

var Codecs = map[string]Codec[*Connection]{
	"std_json":      StdJsonCodec[*Connection]{},
	"msgp":          MsgpCodec[*Connection]{},
	"msgp_zstd":     NewMsgpZstdCodec[*Connection](),
	"std_json_zstd": NewStdJsonZstdCodec[*Connection](),
}
