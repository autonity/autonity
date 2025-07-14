package helpers

func IsHeightExpired(coreHeight uint64, height uint64, heightRange uint64) bool {
	return height < MinNonExpiredHeight(coreHeight, heightRange)
}

func MinNonExpiredHeight(coreHeight uint64, heightRange uint64) uint64 {
	if coreHeight <= heightRange {
		return 0
	}
	return coreHeight - heightRange
}
