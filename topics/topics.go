package topics

import (
	"fmt"
	"github.com/demos/kafka/config"
	"github.com/demos/kafka/producer"
	"gopkg.in/Shopify/sarama.v1"
	"log"
	"os"
	"time"
)

type Topic struct {
	cnf   *config.TopicsConfig
	prCnf *config.ProducerConfig
}

func GetNewTopic(cnf *config.TopicsConfig, prCnf *config.ProducerConfig) Topic {
	return Topic{
		cnf:   cnf,
		prCnf: prCnf,
	}
}

func (t Topic) IsTopicAvailable(brokers []string) bool {
	client, err := sarama.NewClient(brokers, sarama.NewConfig())
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

func (t Topic) Publish(message []byte, brokers []string) error {
	p, err := producer.ConnectProducer(t.prCnf, brokers)
	if err != nil {
		return err
	}
	defer p.Close()

	msg := &sarama.ProducerMessage{
		Topic: t.cnf.Name,
		Key: KeyEncoder([]byte("id_1"),t.cnf.KeySerializer), // not required since single comment is published
		Value: ValueEncoder(message, t.cnf.ValueSerializer),
	}

	partition, offset, err := p.SendMessage(msg)
	if err != nil {
		return err
	}
	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", t.cnf.Name, partition, offset)
	return nil
}

func (t Topic) PublishAsync(messages [][]byte, brokers []string, signals chan os.Signal) {
	p, err := producer.ConnectAsyncProducer(t.prCnf, brokers)
	if err != nil {
		log.Printf("unable to send message asynchronously: %v", err)
	}
	defer p.Close()

	for i, message := range messages {
		time.Sleep(time.Second)
		msg := &sarama.ProducerMessage{
			Topic: t.cnf.Name,
			Key: KeyEncoder([]byte(fmt.Sprintf("id_%s", i)), t.cnf.KeySerializer), // to send data to same partition based on key, in here all keys are different,so no message will go to same partition
			Value: ValueEncoder(message, t.cnf.ValueSerializer)}
		select {
		case p.Input() <- msg:
			log.Println("New Message produced")
		case resMsg, ok := <-p.Successes():
			if ok {
				fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", t.cnf.Name, resMsg.Partition, resMsg.Partition)
			}
		case errMsg, ok := <-p.Errors():
			if ok {
				fmt.Printf("Failed to publish message(%s) in topic(%s) : %v\n", t.cnf.Name, errMsg.Msg.Value, errMsg.Err)
			}
		case <-signals:
			p.AsyncClose()
			return
		}
	}
}
