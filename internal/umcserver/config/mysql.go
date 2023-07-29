package config

type Mysql struct {
	Host                  string `mapstructure:"host" json:"host" yaml:"host" `
	Username              string `mapstructure:"username" json:"username" yaml:"username" `
	Password              string `mapstructure:"password" json:"password" yaml:"password" `
	Database              string `mapstructure:"database" json:"database" yaml:"database" `
	MaxIdleConnections    int    `mapstructure:"max-idle-connections" json:"max-idle-connections" yaml:"max-idle-connections" `
	MaxOpenConnections    int    `mapstructure:"max-open-connections" json:"max-open-connections" yaml:"max-open-connections" `
	MaxConnectionLifeTime int    `mapstructure:"max-connection-life-time" json:"max-connection-life-time" yaml:"max-connection-life-time" `
	LogLevel              string `mapstructure:"log-level" json:"log-level" yaml:"log-level" `
}
