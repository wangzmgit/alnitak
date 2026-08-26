package config

type Transcoding struct {
	// Mode 转码执行模式：local（本地进程）| remote（远程 Worker 池）。
	// 默认 local，切换 remote 前必须部署 transcoder-worker。
	Mode string `mapstructure:"mode" json:"mode" yaml:"mode"`

	UseGpu             bool   `mapstructure:"use_gpu" json:"use_gpu" yaml:"use_gpu"`
	UseH265            bool   `mapstructure:"use_h265" json:"use_h265" yaml:"use_h265"`
	UseAv1             bool   `mapstructure:"use_av1" json:"use_av1" yaml:"use_av1"`
	Generate1080p60    bool   `mapstructure:"generate_1080p60" json:"generate_1080p60" yaml:"generate_1080p60"`
	MaxCpuConcurrency  int    `mapstructure:"max_cpu_concurrency" json:"max_cpu_concurrency" yaml:"max_cpu_concurrency"`
	MaxGpuConcurrency  int    `mapstructure:"max_gpu_concurrency" json:"max_gpu_concurrency" yaml:"max_gpu_concurrency"`
	WorkerConcurrency  int    `mapstructure:"worker_concurrency" json:"worker_concurrency" yaml:"worker_concurrency"`
	EncodingConcurrency int   `mapstructure:"encoding_concurrency" json:"encoding_concurrency" yaml:"encoding_concurrency"`
	MaxQueueDepth       int   `mapstructure:"max_queue_depth" json:"max_queue_depth" yaml:"max_queue_depth"`
	WorkDir            string `mapstructure:"work_dir" json:"work_dir" yaml:"work_dir"`
}
