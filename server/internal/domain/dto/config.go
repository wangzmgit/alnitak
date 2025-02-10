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
	MaxImgSize   int64
	MaxVideoSize int64

	Type      string
	KeyID     string
	KeySecret string
	Bucket    string
	Endpoint  string
	AppID     string
	Region    string
	Domain    string
	Private   bool

	UploadMp4File bool
}

type OtherConfigReq struct {
	AllowOrigin     string
	Prefix          string
	Generate1080p60 bool
	UseGpu          bool
}
