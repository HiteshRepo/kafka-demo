package producer

import (
	"github.com/demos/kafka/config"
	"gopkg.in/Shopify/sarama.v1"
	"strings"
)

func ConnectProducer(producerConfig *config.ProducerConfig) (sarama.SyncProducer,error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = producerConfig.ReturnSuccess
	config.Producer.RequiredAcks = producerAcks[producerConfig.RequiredAcks]
	config.Producer.Retry.Max = producerConfig.RetryCount
	// NewSyncProducer creates a new SyncProducer using the given broker addresses and configuration.
	brokers := strings.Split(producerConfig.Brokers, ";")
	conn, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
