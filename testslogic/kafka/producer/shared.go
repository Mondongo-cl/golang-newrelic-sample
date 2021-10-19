package producer

import (
	"bltrain/heathcheck/internal/configuration"

	"github.com/Shopify/sarama"
)

type IProducer interface {
	Connect()
	Disconnect()
	Produce(message interface{}, brokers []string, topic string) error
}

type producer struct {
	configuration.Settings
	syncProducer sarama.SyncProducer
}

type ProduceKafkaRequest struct {
	Topics []string
	Data   interface{}
}
