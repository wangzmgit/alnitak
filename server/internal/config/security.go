package config

type Security struct {
	AccessJwtSecret          string `mapstructure:"access_jwt_secret" json:"access_jwt_secret" yaml:"access_jwt_secret"`
	RefreshJwtSecret         string `mapstructure:"refresh_jwt_secret" json:"refresh_jwt_secret" yaml:"refresh_jwt_secret"`
	UserIdSalt               string `mapstructure:"user_id_salt" json:"user_id_salt" yaml:"user_id_salt"`
	CloseRecordUserOperation bool   `mapstructure:"close_record_user_operation" json:"close_record_user_operation" yaml:"close_record_user_operation"`
}
