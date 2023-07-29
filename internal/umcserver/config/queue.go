package config

type Queue struct {
	Model string `mapstructure:"model" json:"model" yaml:"model"`
}

type Kafka struct {
	Servers string `mapstructure:"servers" json:"servers" yaml:"servers"`
}
