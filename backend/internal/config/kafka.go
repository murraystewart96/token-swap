package config

type KafkaProducer struct {
	BootstrapServers string `mapstructure:"bootstrap_servers" validate:"required"`
	Acks             string `mapstructure:"acks" validate:"required"`
}

type KafkaConsumer struct {
	BootstrapServers string `mapstructure:"bootstrap_servers" validate:"required"`
	GroupID          string `mapstructure:"group_id" validate:"required"`
	OffsetReset      string `mapstructure:"offset_reset" validate:"required"`
}
