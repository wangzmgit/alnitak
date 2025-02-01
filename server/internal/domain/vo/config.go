package vo

type EmailConfigResp struct {
	// Debug     bool   `json:"debug"`
	User string `json:"user"`
	// Pass      string `json:"pass"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Addresser string `json:"addresser"`
}

type StorageConfigResp struct {
	MaxImgSize   int `json:"maxImgSize"`
	MaxVideoSize int `json:"maxVideoSize"`

	Type     string `json:"type"`
	KeyID    string `json:"keyId"`
	Bucket   string `json:"bucket"`
	Endpoint string `json:"endpoint"`
	AppID    string `json:"appId"`
	Region   string `json:"region"`
	Domain   string `json:"domain"`
	Private  bool   `json:"private"`

	UploadMp4File bool `json:"uploadMp4File"`
}

type OtherConfigResp struct {
	AllowOrigin string `json:"allowOrigin"`
	Prefix      string `json:"prefix"`
}
