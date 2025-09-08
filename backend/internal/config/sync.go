package config

import "github.com/spf13/viper"

type Sync struct {
	Listener Listener `mapstructure:"listener" validate:"required"`
	Redis    `mapstructure:"redis"    validate:"required"`
	DB       `mapstructure:"db"       validate:"required"`
}

func (s *Sync) Defaults() {
	// Binds ENV vars to struct
	// ENV vars, if defined, take precedence over defaults and config.yaml
	viper.SetDefault("listener.rpc_url", "127.0.0.1:8545")
	viper.SetDefault("listener.contract_addr", "")
	viper.SetDefault("redis.addr", "localhost:6379")
}
