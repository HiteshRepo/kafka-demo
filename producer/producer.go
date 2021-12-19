package producer

import (
	"github.com/demos/kafka/config"
	"gopkg.in/Shopify/sarama.v1"
)

func ConnectProducer(producerConfig *config.ProducerConfig, brokers []string) (sarama.SyncProducer,error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = producerConfig.ReturnSuccess
	config.Producer.RequiredAcks = producerAcks[producerConfig.RequiredAcks]
	config.Producer.Retry.Max = producerConfig.RetryCount
	// NewSyncProducer creates a new SyncProducer using the given broker addresses and configuration.
	conn, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func ConnectAsyncProducer(producerConfig *config.ProducerConfig, brokers []string) (sarama.AsyncProducer,error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = producerConfig.ReturnSuccess
	config.Producer.RequiredAcks = producerAcks[producerConfig.RequiredAcks]
	config.Producer.Retry.Max = producerConfig.RetryCount
	// NewSyncProducer creates a new SyncProducer using the given broker addresses and configuration.
	conn, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
