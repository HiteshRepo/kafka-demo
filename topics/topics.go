package topics

import (
	"fmt"
	"github.com/demos/kafka/config"
	"github.com/demos/kafka/producer"
	"gopkg.in/Shopify/sarama.v1"
	"log"
	"strings"
)

type topic struct {
	cnf   *config.TopicsConfig
	prCnf *config.ProducerConfig
}

func GetNewTopic(cnf *config.TopicsConfig, prCnf *config.ProducerConfig) topic {
	return topic{
		cnf:   cnf,
		prCnf: prCnf,
	}
}

func (t topic) IsTopicAvailable() bool {
	client, err := sarama.NewClient(strings.Split(t.prCnf.Brokers, ";"), sarama.NewConfig())
	defer client.Close()

	if err != nil {
		log.Printf("error while creating client: %v", err)
		return false
	}

	availableTopics, err := client.Topics()
	if err != nil {
		log.Printf("error while getting topics from kafka: %v", err)
		return false
	}

	for _, topic := range availableTopics {
		if topic == t.cnf.Name {
			return true
		}
	}

	return false
}

func (t topic) Publish(message []byte) error {
	p, err := producer.ConnectProducer(t.prCnf)
	if err != nil {
		return err
	}
	defer p.Close()

	msg := &sarama.ProducerMessage{
		Topic: t.cnf.Name,
		Value: ValueEncoder(message, t.cnf.ValueSerializer),
	}

	partition, offset, err := p.SendMessage(msg)
	if err != nil {
		return err
	}
	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", t.cnf.Name, partition, offset)
	return nil
}
