package quic

type ConnectionStateStoreConf struct {
	// base stream state without pending frames
	IncludeStreamState           bool
	IncludePendingOutgoingFrames bool
	IncludePendingIncomingFrames bool
	IgnoreCurrentPath            bool
	IncludeCongestionState       bool
}
