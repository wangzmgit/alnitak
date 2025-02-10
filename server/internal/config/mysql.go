package config

type Mysql struct {
	Datasource string `mapstructure:"datasource" json:"datasource" yaml:"datasource"`
	Host       string `mapstructure:"host" json:"host" yaml:"host"`
	Port       int    `mapstructure:"port" json:"port" yaml:"port"`
	Username   string `mapstructure:"username" json:"username" yaml:"username"`
	Password   string `mapstructure:"password" json:"password" yaml:"password"`
	Param      string `mapstructure:"param" json:"param" yaml:"param"`
}
