package config

import "github.com/spf13/viper"

type Events struct {
	RPCUrl       string `mapstructure:"rpc_url"          validate:"required"`
	ContractAddr string `mapstructure:"contract_addr"    validate:"required"`
}

func (e *Events) Defaults() {
	// Binds ENV vars to struct
	// ENV vars, if defined, take precedence over defaults and config.yaml
	viper.SetDefault("rpc_url", "127.0.0.1:8545")
	viper.SetDefault("contract_addr", "vault")
}
