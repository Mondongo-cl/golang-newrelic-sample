package producer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func init() {
	brokers := getBrokersFromEnv()
	partition := getPartitionFromEnv()
	handler = NewProducer(brokers, partition)
}

func ProducerPostHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("content-type", "application/json")
	if request.Method != "POST" {
		log.Printf("Method %s not allowed", request.Method)
		writer.WriteHeader(http.StatusMethodNotAllowed)
		writer.Write([]byte("method not allowed"))
		return
	}

	err := handler.Connect()
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}

	defer handler.Disconnect()
	var msg ProduceKafkaRequest
	data := make([]byte, request.ContentLength)
	request.Body.Read(data)
	err = json.Unmarshal(data, &msg)
	if err != nil || msg.Topics == nil {
		sample, _ := json.Marshal(&ProduceKafkaRequest{Topics: []string{"topic-1", "topic-2", "topic-nth"}, Data: map[string]string{"key": "value"}})
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(fmt.Sprintf("{\"message\":\"body format is not allowed error %v, use correct one:\":\"%v\"}", err, string(sample))))
		return
	}
	topicsList := msg.Topics
	createTopic(handler.Brokers, topicsList)
	for i, topicName := range topicsList {
		log.Printf("prossesing element number %v - topic name %v", i, topicName)
		if msg.Data != nil {
			data, err := json.Marshal(msg.Data)
			if err != nil {
				writer.WriteHeader(http.StatusInternalServerError)
				writer.Write([]byte(fmt.Sprintf("\"message:\"\"Cant serialize payload error is :%v\"", err)))
				return

			}
			err = handler.Produce(string(data), topicName)
			if err != nil {
				writer.WriteHeader(http.StatusInternalServerError)
				writer.Write([]byte(fmt.Sprintf("{\"message\":\"%v\"}", err)))
				return
			}
		} else {
			writer.WriteHeader(http.StatusOK)
			writer.Write([]byte("{\"message\":\"Kafka server is running, no topics was created\"}"))
			return
		}
	}

	if err == nil {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(fmt.Sprintf("{\"message\":\"Kafka server is running the all topics are written %v\"}", topicsList)))
		return
	}
	writer.WriteHeader(http.StatusInternalServerError)
	writer.Write([]byte(fmt.Sprintf("{\"message\":\"%v\"}", err)))
}

func getBrokersFromEnv() []string {
	brokers := os.Getenv(KafkaHostConfigKey)
	if brokers == "" {
		brokers = "localhost:9092"
	}
	brokersList := strings.Split(brokers, ",")
	return brokersList
}

func getTopicsFromEnv() []string {
	topics := os.Getenv(KafkaTopicConfigKey)
	if topics == "" {
		topics = "default"
	}
	brokersList := strings.Split(topics, ",")
	return brokersList
}

func getPartitionFromEnv() int {
	partition := os.Getenv(KafkaPartionConfigKey)
	result, err := strconv.Atoi(partition)
	if err != nil {
		result = 0
	}
	return result
}
