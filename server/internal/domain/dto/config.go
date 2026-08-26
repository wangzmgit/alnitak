package dto

type EmailConfigReq struct {
	Debug     bool
	User      string
	Pass      string
	Host      string
	Port      int
	Addresser string
}

type StorageConfigReq struct {
	MaxImgSize   int64 `json:"maxImgSize"`
	MaxVideoSize int64 `json:"maxVideoSize"`

	Type      string `json:"type"`
	KeyID     string `json:"keyId"`
	KeySecret string `json:"keySecret"`
	Bucket    string `json:"bucket"`
	Endpoint  string `json:"endpoint"`
	AppID     string `json:"appId"`
	Region    string `json:"region"`
	Domain    string `json:"domain"`
	Private   bool   `json:"private"`
	UseSSL    bool   `json:"useSSL"`

	UploadMp4File bool `json:"uploadMp4File"`
}

type OtherConfigReq struct {
	AllowOrigin string `json:"allowOrigin"`
	Prefix      string `json:"prefix"`

	// 服务器配置
	ServerPort  string `json:"serverPort"`
	SslEnabled  bool   `json:"sslEnabled"`
	SslPort     string `json:"sslPort"`
	SslCertFile string `json:"sslCertFile"`
	SslKeyFile  string `json:"sslKeyFile"`
}

type TranscodingConfigReq struct {
	Mode                string `json:"mode"`
	UseGpu              bool   `json:"useGpu"`
	UseH265             bool   `json:"useH265"`
	UseAv1              bool   `json:"useAv1"`
	Generate1080p60     bool   `json:"generate1080p60"`
	MaxCpuConcurrency   int    `json:"maxCpuConcurrency"`
	MaxGpuConcurrency   int    `json:"maxGpuConcurrency"`
	WorkerConcurrency   int    `json:"workerConcurrency"`
	EncodingConcurrency int    `json:"encodingConcurrency"`
	MaxQueueDepth       int    `json:"maxQueueDepth"`
	WorkDir             string `json:"workDir"`
}
