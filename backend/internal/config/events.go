package config

type Events struct {
	RPCUrl       string `mapstructure:"rpc_url"          validate:"required"`
	ContractAddr string `mapstructure:"contract_addr"     validate:"required"`
}
