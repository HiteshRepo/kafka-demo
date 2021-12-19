package consumer

import "gopkg.in/Shopify/sarama.v1"

var consumerOffsets = map[int]int64 {
	-1: sarama.OffsetNewest,
	-2: sarama.OffsetOldest,
}

