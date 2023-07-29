package config

type Log struct {
	Path          string `mapstructure:"path" json:"path" yaml:"path" `
	MaxAge        int    `mapstructure:"max-age" json:"max-age" yaml:"max-age"`
	MaxSize       int64  `mapstructure:"max-size" json:"max-size" yaml:"max-size" `
	Level         string `mapstructure:"level" json:"level" yaml:"level" `
	Format        string `mapstructure:"format" json:"format" yaml:"format" `
	LogIncConsole bool   `mapstructure:"log-inc-console" json:"log-inc-console" yaml:"log-inc-console" `
	LevelEncode   string `mapstructure:"level-encode" json:"level-encode" yaml:"level-encode" `
}
