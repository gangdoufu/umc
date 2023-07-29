package redis

import (
	"github.com/go-redis/redis/v8"
	"time"
)

var client redis.UniversalClient

const (
	defaultStep = "-"
)

type Option struct {
	Addr                  string   `mapstructure:"addr" json:"addr" yaml:"addr" `
	DB                    int      `mapstructure:"db" json:"db" yaml:"db" `
	UserName              string   `mapstructure:"user-name" json:"user-name" yaml:"user-name" `
	Password              string   `mapstructure:"password" json:"password" yaml:"password" `
	EnableCluster         bool     `mapstructure:"is-cluster" json:"is-cluster" yaml:"is-cluster" `
	Addrs                 []string `mapstructure:"addrs" json:"addrs" yaml:"addrs" `
	MasterName            string   `mapstructure:"master-name" json:"master-name" yaml:"master-name" `
	MaxIdle               int      `mapstructure:"max-idle" json:"max-idle" yaml:"max-idle" `
	MaxActive             int      `mapstructure:"max-active" json:"max-active" yaml:"max-active" `
	Timeout               int      `mapstructure:"timeout" json:"timeout" yaml:"timeout" `
	UseSSL                bool     `mapstructure:"use-ssl" json:"use-ssl" yaml:"use-ssl" `
	SSLInsecureSkipVerify bool     `mapstructure:"ssl-insecure-skip-verify" json:"ssl-insecure-skip-verify" yaml:"ssl-insecure-skip-verify" `
}

func (o Option) getUniversalClientOption() *redis.UniversalOptions {
	if o.MaxActive <= 0 {
		o.MaxActive = 500
	}
	timeout := time.Duration(o.Timeout) * time.Second
	op := &redis.UniversalOptions{
		DB:           o.DB,
		Username:     o.UserName,
		Password:     o.Password,
		PoolSize:     o.MaxActive,
		IdleTimeout:  240 * time.Second,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		DialTimeout:  timeout,
		MasterName:   o.MasterName,
	}
	if o.EnableCluster {
		op.Addrs = o.Addrs
	} else {
		op.Addrs = []string{o.Addr}
	}
	return op
}

func NewRedisCache(op *Option) {
	universalOp := op.getUniversalClientOption()
	client = redis.NewUniversalClient(universalOp)
}
