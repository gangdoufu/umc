package config

type Redis struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host" `
	Password string `mapstructure:"password" json:"password" yaml:"password" `
}
