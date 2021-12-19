package config

import (
	"io"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Brokers string `mapstructure:"brokers"`
}

type ProducerConfig struct {
	ReturnSuccess bool `mapstructure:"returnSuccesses"`
	RetryCount    int  `mapstructure:"retryCount"`
	RequiredAcks  int  `mapstructure:"requiredAcks"`
}

type TopicsConfig struct {
	Name            string `mapstructure:"name"`
	KeySerializer   string `mapstructure:"keySerializer"`
	ValueSerializer string `mapstructure:"valueSerializer"`
}

type ConsumerConfig struct {
	Offset int `mapstructure:"offset"`
}

type ConsumerGroupConfig struct {
	Name     string `mapstructure:"name"`
	Version  string `mapstructure:"version"`
	Assignor string `mapstructure:"assignor"`
	Oldest   bool   `mapstructure:"oldest"`
}

type AppConfig interface {
	GetProducerConfig() *ProducerConfig
	GetTopicsConfig() *TopicsConfig
	GetServerConfig() *ServerConfig
	GetConsumerConfig() *ConsumerConfig
	GetConsumerGroupConfig() *ConsumerGroupConfig
}

type appConfig struct {
	ProducerConfig      ProducerConfig      `mapstructure:"producer"`
	TopicsConfig        TopicsConfig        `mapstructure:"topic"`
	ServerConfig        ServerConfig        `mapstructure:"app"`
	ConsumerConfig      ConsumerConfig      `mapstructure:"consumer"`
	ConsumerGroupConfig ConsumerGroupConfig `mapstructure:"group"`
}

func (a *appConfig) GetProducerConfig() *ProducerConfig {
	return &a.ProducerConfig
}

func (a *appConfig) GetTopicsConfig() *TopicsConfig {
	return &a.TopicsConfig
}

func (a *appConfig) GetServerConfig() *ServerConfig {
	return &a.ServerConfig
}

func (a *appConfig) GetConsumerConfig() *ConsumerConfig {
	return &a.ConsumerConfig
}

func (a *appConfig) GetConsumerGroupConfig() *ConsumerGroupConfig {
	return &a.ConsumerGroupConfig
}

func LoadConfig(reader io.Reader) (AppConfig, error) {
	var appConfig appConfig

	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	if err := viper.ReadConfig(reader); err != nil {
		return nil, errors.Wrap(err, "Failed to load app config file")
	}

	if err := viper.Unmarshal(&appConfig); err != nil {
		return nil, errors.Wrap(err, "Unable to parse app config file")
	}

	return &appConfig, nil
}
