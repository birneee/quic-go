package qstate

var Codecs = map[string]Codec{
	"std_json":      StdJsonCodec{},
	"msgp":          MsgpCodec{},
	"msgp_zstd":     NewMsgpZstdCodec(),
	"std_json_zstd": NewStdJsonZstdCodec(),
}
