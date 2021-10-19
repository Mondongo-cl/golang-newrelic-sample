package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-redis/redis/v8"
)

type configuration struct {
	Address  string
	Password string
	DB       int
}

const (
	RedisAddressConfigKey  = "REDIS_HOST"
	RedisPassWordConfigKey = "REDIS_PASSWORD"
	RedisDBConfigKey       = "REDIS_DB"
)

var (
	cfg configuration
)

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.Printf("Processing message With cfg: %v", cfg)
	options := &redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			id := cn.ClientID(ctx)
			log.Printf("New connection for %v", id)
			return nil
		},
	}

	log.Printf("redis options: %v", options)
	rdb := redis.NewClient(options)
	key := request.PathParameters["id"]
	result := rdb.Get(ctx, key)
	err := result.Err()
	log.Printf("get error : %v", err)
	log.Printf("get result : %v", result)

	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode:      503,
			Body:            "{ \"code\":\"503\", \"key\":\"nil\"}",
			IsBase64Encoded: false,
			Headers:         request.Headers,
		}, err
	}
	data, err := result.Result()
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode:      503,
			Body:            fmt.Sprintf("{ \"code\":\"503\", \"key\":\"%v\"}", key),
			IsBase64Encoded: false,
			Headers:         request.Headers,
		}, err
	}
	return events.APIGatewayV2HTTPResponse{
		StatusCode:      200,
		Body:            fmt.Sprintf("{ \"code\":\"200\", \"key\":\"%v\", \"value\":\"%v\"}", key, data),
		IsBase64Encoded: false,
		Headers:         request.Headers,
	}, nil
}

func main() {
	db, err := strconv.Atoi(os.Getenv(RedisDBConfigKey))
	if err != nil {
		db = 0
	}

	cfg = configuration{
		Address:  os.Getenv(RedisAddressConfigKey),
		Password: os.Getenv(RedisPassWordConfigKey),
		DB:       db,
	}
	log.Printf("Starting lambda with followin configuratuib %v", cfg)
	lambda.Start(handler)
}
