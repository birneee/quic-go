package handover

type ConnectionStateStoreConf struct {
	IncludePendingOutgoingFrames bool
	IncludePendingIncomingFrames bool
	IgnoreCurrentPath            bool
	IncludeCongestionState       bool
}