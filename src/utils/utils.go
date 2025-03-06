package utils

import "time"

var StartTime time.Time

func SetStartTime() {
	StartTime = time.Now()
}

func GetUptime() time.Duration {
	return time.Since(StartTime)
}

func BToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
