package config

import "github.com/spf13/viper"

type Server struct {
	Addr  string `mapstructure:"addr"    validate:"required"`
	Redis `mapstructure:"redis"    validate:"required"`
	DB    `mapstructure:"db"       validate:"required"`
}

func (s *Server) Defaults() {
	viper.SetDefault("addr", "localhost:2025")
	viper.SetDefault("redis.addr", "localhost:6379")
}
