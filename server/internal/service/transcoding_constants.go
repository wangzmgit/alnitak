package service

import "time"

const (
	maxCPUConcurrentTranscoding = 2
	maxGPUConcurrentTranscoding = 2

	maxGpuFailCountThreshold = 3

	ossUploadMaxConcurrency = 10
	ossUploadMaxRetries     = 3
)

// ossUploadBackoff 指数退避：1s, 2s, 4s
var ossUploadBackoff = []time.Duration{1 * time.Second, 2 * time.Second, 4 * time.Second}

