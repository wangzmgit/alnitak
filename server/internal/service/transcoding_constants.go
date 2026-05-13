package service

import "time"

const (
	maxCPUConcurrentTranscoding = 2
	maxGPUConcurrentTranscoding = 2

	maxGpuFailCountThreshold = 3

	ossUploadMaxConcurrency = 10
	ossUploadRetryDelay     = 500 * time.Millisecond
)

