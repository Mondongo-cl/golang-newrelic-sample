package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	// KafkaHostConfigKey is the env.os value name
	KafkaHostConfigKey = "KAFKA_SERVER"
	// KafkaTopicConfigKey is the env.os value name
	KafkaTopicConfigKey = "KAFKA_TOPIC"
)

var (
	brokersList []string
	topicName   string
)

func init() {
	brokers := os.Getenv(KafkaHostConfigKey)
	topicName = os.Getenv(KafkaTopicConfigKey)
	brokersList = strings.Split(brokers, ",")

}
func describeTopic() []byte {

	cfg := sarama.NewConfig()

	cluster, err := sarama.NewClusterAdmin(brokersList, cfg)

	if err == nil {
		topics, err := cluster.ListTopics()
		if err != nil {
			log.Printf("listing topics error  %v", err)
		}
		response, err := json.Marshal(topics)
		if err != nil {
			log.Panicf("error while parse topics descriptions, error %v", err)
		}
		return response
	} else {
		log.Panicf("error while create the clusterAdmin object %v", err)
	}
	return nil
}

// ConsumerHandler CONSUME MESSAGES
func ConsumerHandler(e interface{}) {

	data, err := json.Marshal(e)
	if err != nil {
		log.Fatalf("Error while decode to json kafka event error is\n%v", err)
	}
	log.Printf("Message Received from data stream message:\n%v", string(data))

}

func main() {
	lambda.Start(ConsumerHandler)
}
