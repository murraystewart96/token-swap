package config

import (
	"github.com/spf13/viper"
)

type Migration struct {
	DB   DB     `mapstructure:"db"       validate:"required"`
	Path string `mapstructure:"path"     validate:"required"`
}

type DB struct {
	Host     string `mapstructure:"host"     validate:"required"`
	Port     string `mapstructure:"port"     validate:"required"`
	Name     string `mapstructure:"name"     validate:"required"`
	User     string `mapstructure:"user"     validate:"required"`
	Password string `mapstructure:"password" validate:"omitempty"`
}

type Redis struct {
	Addr string `mapstructure:"addr" validate:"required"`
}

func (m *Migration) Defaults() {
	viper.SetDefault("path", "./internal/storage/postgres/migrations/sql")

	viper.SetDefault("db.host", "localhost")
	viper.SetDefault("db.port", "5432")
	viper.SetDefault("db.name", "tokenswap")
	viper.SetDefault("db.user", "tokenswap")
	viper.SetDefault("db.password", "password")
}
