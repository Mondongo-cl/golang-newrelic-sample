package producer

import (
	"bltrain/heathcheck/internal/configuration"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/Shopify/sarama"
)

const (
	KafkaHostConfigKey    = "KAFKA_SERVER"
	KafkaTopicConfigKey   = "KAFKA_TOPIC"
	KafkaPartionConfigKey = "KAFKA_PARTITION"
)

var (
	handler producer
)

func NewProducer(brokers []string, defaultPartition int) producer {
	return producer{
		Settings: configuration.Settings{
			Brokers:          brokers,
			DefaultPartition: defaultPartition,
		},
		syncProducer: nil,
	}
}

func (p *producer) Disconnect() error {
	if p.syncProducer != nil {
		producer := p.syncProducer
		if producer != nil {
			if err := producer.Close(); err != nil {
				return err
			}
		} else {
			return errors.New("producer connection is null")
		}
	}
	return nil
}
func (p *producer) Connect() error {
	if p.Brokers == nil || len(p.Brokers) == 0 {
		return errors.New("brokers config is empty")
	}
	log.Print("Starting connection to kafka")
	client, error := getSyncProducer(p.Brokers)
	if error != nil {
		log.Printf("failed to connect to kafka %v", error)
		return error
	}
	log.Print("End connection to kafka")
	p.syncProducer = client
	return nil
}

func getSyncProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Version = sarama.MaxVersion
	return sarama.NewSyncProducer(brokers, config)
}

func (p *producer) Produce(message interface{}, topic string) error {
	if p.syncProducer != nil {
		producer := p.syncProducer
		out, err := json.Marshal(message)
		if err != nil {
			return err
		}

		msg := &sarama.ProducerMessage{
			Topic:     topic,
			Value:     sarama.ByteEncoder(out),
			Timestamp: time.Now(),
		}
		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			return err
		}
		log.Printf("Message Writtern on patition:%v offset:%v", partition, offset)
		return nil
	}
	return errors.New("send Message Failed")
}

func createTopic(brokers []string, topics []string) {
	topicdetails := &sarama.TopicDetail{
		NumPartitions:     2,
		ReplicationFactor: 2,
		ConfigEntries:     map[string]*string{},
	}

	cfg := sarama.NewConfig()

	cluster, err := sarama.NewClusterAdmin(brokers, cfg)

	if err == nil {
		for _, topicName := range topics {
			err := cluster.CreateTopic(topicName, topicdetails, true)
			if err == nil {
				log.Printf("Creating topic %v", topicName)
				cluster.CreateTopic(topicName, topicdetails, false)
			}
		}
	} else {
		log.Panicf("error while create the clusterAdmin object %v", err)
	}
}
