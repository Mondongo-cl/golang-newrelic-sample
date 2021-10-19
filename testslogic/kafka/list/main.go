package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/aws/aws-lambda-go/events"
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

// ProducerHandler Produce Messages
func ProducerHandler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	data := describeTopic()

	if data != nil {
		return events.APIGatewayV2HTTPResponse{StatusCode: 200, Body: fmt.Sprintf("{\"Sucess\":\"%v\"}", string(data))}, nil
	}
	return events.APIGatewayV2HTTPResponse{StatusCode: 500, Body: "{\"failed\":\"\"}"}, nil
}

func main() {
	lambda.Start(ProducerHandler)
}
