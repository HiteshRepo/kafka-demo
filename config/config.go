package config

import (
	"io"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type ProducerConfig struct {
	Brokers       string `mapstructure:"brokers"`
	ReturnSuccess bool   `mapstructure:"returnSuccesses"`
	RetryCount    int    `mapstructure:"retryCount"`
	RequiredAcks  int    `mapstructure:"requiredAcks"`
}

type TopicsConfig struct {
	Name            string `mapstructure:"name"`
	KeySerializer   string `mapstructure:"keySerializer"`
	ValueSerializer string `mapstructure:"valueSerializer"`
}

type AppConfig interface {
	GetProducerConfig() *ProducerConfig
	GetTopicsConfig() *TopicsConfig
}

type appConfig struct {
	ProducerConfig ProducerConfig `mapstructure:"producer"`
	TopicsConfig   TopicsConfig   `mapstructure:"topic"`
}

func (a *appConfig) GetProducerConfig() *ProducerConfig {
	return &a.ProducerConfig
}

func (a *appConfig) GetTopicsConfig() *TopicsConfig {
	return &a.TopicsConfig
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
