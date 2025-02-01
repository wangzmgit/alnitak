package oss

type Config struct {
	KeyID     string //(必填，对应阿里云access_key_id、腾讯云secret_id，七牛云accessKey)
	KeySecret string //(必填，对应阿里云access_key_secret、腾讯云secret_key，七牛云secretKey)
	Bucket    string

	Endpoint string //阿里云必填

	AppID  string //腾讯云必填
	Region string //腾讯云必填

	Domain string //域名

	Private bool //是否私有 (仅七牛云)
}
