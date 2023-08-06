package config

type App struct {
	Name string `mapstructure:"name" json:"name" yaml:"name"`
	Desc string `mapstructure:"desc" json:"desc" yaml:"desc"`
}
