package producer

import "gopkg.in/Shopify/sarama.v1"

var producerAcks = map[int]sarama.RequiredAcks {
	0: sarama.NoResponse,
	1: sarama.WaitForLocal,
	-1: sarama.WaitForAll,
}
