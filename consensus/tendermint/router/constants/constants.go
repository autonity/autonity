package constants

import "time"

const (
	ScaleThresholdForClustering = 21
	DefaultLatency              = uint(132)
	DefaultNearThreshold        = 50
	LatencyMeasurementDelayCap  = 2000
	MaxLatencyCapFactor         = 0.75
	LatencyDataExpiry           = 5 * time.Minute
	RetryLatencyTimeout         = 50 * time.Second
	CacheCleanupInterval        = 10 * time.Minute
)
