package config

import "github.com/spf13/viper"

type Worker struct {
	Kafka  KafkaConsumer `mapstructure:"kafka"    validate:"required"`
	Redis  `mapstructure:"redis"    validate:"required"`
	DB     `mapstructure:"db"       validate:"required"`
	Topics []string `mapstructure:"topics" validate:"required"`
}

func (w *Worker) Defaults() {
	viper.SetDefault("kafka.bootstrap_servers", "localhost:9092")
	viper.SetDefault("kafka.offset_reset", "earliest")
	viper.SetDefault("redis.addr", "localhost:6379")
}
