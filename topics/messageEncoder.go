package topics

import "gopkg.in/Shopify/sarama.v1"

func ValueEncoder(message []byte, valueSerializer string) sarama.Encoder {
	switch valueSerializer {
	case "string": return sarama.StringEncoder(message)
	case "byte": return sarama.ByteEncoder(message)
	default: return nil
	}
}

func KeyEncoder(message []byte, keySerializer string) sarama.Encoder {
	switch keySerializer {
	case "string": return sarama.StringEncoder(message)
	case "byte": return sarama.ByteEncoder(message)
	default: return nil
	}
}
