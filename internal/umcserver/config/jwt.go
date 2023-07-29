package config

type Jwt struct {
	Issuer     string `mapstructure:"issuer" json:"issuer" yaml:"issuer"`
	Key        string `mapstructure:"key" json:"key" yaml:"key" `
	Timeout    int    `mapstructure:"timeout" json:"timeout" yaml:"timeout" `
	BufferTime int    `mapstructure:"buffer-time" json:"buffer-time" yaml:"buffer-time"`
}
