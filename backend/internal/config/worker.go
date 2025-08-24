package config

import "github.com/spf13/viper"

type Worker struct {
	Kafka KafkaConsumer `mapstructure:"kafka"    validate:"required"`
}

func (w *Worker) Defaults() {
	viper.SetDefault("kafka.bootstrap_servers", "localhost:9092")
	viper.SetDefault("kafka.offset_reset", "earliest")
}
