package config

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"os"
)

const (
	defaultConfig = "./config/umc.yaml"
)

type Config struct {
	Server Server `mapstructure:"server" json:"server" yaml:"server"`
	Grpc   Grpc   `mapstructure:"grpc" json:"grpc" yaml:"grpc"`
	Http   Http   `mapstructure:"http" json:"http" yaml:"http"`
	SSL    SSL    `mapstructure:"ssl" json:"ssl" yaml:"ssl"`
	Mysql  Mysql  `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis  Redis  `mapstructure:"redis" json:"redis" yaml:"redis"`
	Jwt    Jwt    `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Log    Log    `mapstructure:"log" json:"log" yaml:"log"`
	Queue  Queue  `mapstructure:"queue" json:"queue" yaml:"queue"`
	Kafka  Kafka  `mapstructure:"kafka" json:"kafka" yaml:"kafka"`
}

// 命令行 > 环境变量 > 默认值
func LoadConfig() *Config {
	var configPath string
	flag.StringVar(&configPath, "c", "", "choose config file.")
	flag.Parse()
	if configPath == "" {
		if configEnv := os.Getenv("UMC_CONFIG"); configEnv != "" {
			configPath = configEnv
		} else {
			configPath = defaultConfig
		}
	}
	if configPath == "" {
		panic("need config path")
	}
	fmt.Println("config path-->", configPath)
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("read config file error:%s\n", err))
	}
	var configs = &Config{}
	if err = v.Unmarshal(configs); err != nil {
		panic(fmt.Errorf("read config file error:%s\n", err))
	}
	return configs
}
