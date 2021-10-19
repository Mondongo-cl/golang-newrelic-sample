package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	MongoDbHostKey       = "MONGO_DB"
	MongoDbNameKey       = "MONGO_DBNAME"
	MonngoCollectionKey  = "MONGO_COLLECTION"
	MongoUserNameKey     = "MONGO_USERNAME"
	MongoPasswordHostKey = "MONGO_PASSWORD"

	readPreference           = "secondaryPreferred"
	connectionStringTemplate = "mongodb://%s:%s@%s/%s?retryWrites=false&tls=true&replicaSet=rs0&readpreference=%s&retryWrites=false"
)

var MONGO_DB *mongo.Database

type configuration struct {
	MongoDbHost string
	Database    string
	Collection  string
	Username    string
	Password    string
	Certfile    string
}

var cfg configuration

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	c := ConnectMongoDB(ctx, cfg)

	coll := c.Collection(cfg.Collection)
	idValue := request.PathParameters["id"]
	var result string
	err := coll.FindOne(ctx, bson.D{{"id", idValue}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return events.APIGatewayV2HTTPResponse{
				StatusCode:      403,
				Body:            "{ \"code\":\"403\", \"error\":\"record not found\"}",
				IsBase64Encoded: false,
				Headers:         request.Headers,
			}, nil
		}
		return events.APIGatewayV2HTTPResponse{
			StatusCode:      503,
			Body:            fmt.Sprintf("{ \"code\":\"503\", \"error\":\"%v\"}", request),
			IsBase64Encoded: false,
			Headers:         request.Headers,
		}, nil
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode:      200,
		Body:            fmt.Sprintf("{ \"code\":\"200\", \"result\":\"%v\"}", result),
		IsBase64Encoded: false,
		Headers:         request.Headers,
	}, nil
}

func ConnectMongoDB(ctx context.Context, cfg configuration) *mongo.Database {

	var client *mongo.Client

	connectionURI := fmt.Sprintf(connectionStringTemplate, cfg.Username, cfg.Password, cfg.MongoDbHost, cfg.Database, readPreference)
	tlsConfig, err := getCustomTLSConfig(cfg.Certfile)
	if err != nil {
		log.Fatal(err)
	}
	_client, err := mongo.NewClient(options.Client().ApplyURI(connectionURI).SetTLSConfig(tlsConfig))
	if err != nil {
		log.Fatal(err)
	}
	client = _client

	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer cancel()

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Couldn't connect to the database", err)
	} else {
		log.Println("MDB Connected!")
	}

	// Connect to the database
	db := client.Database(cfg.Database)
	return db
}

func getCustomTLSConfig(caFile string) (*tls.Config, error) {
	tlsConfig := new(tls.Config)
	certs, err := ioutil.ReadFile(caFile)

	if err != nil {
		return tlsConfig, err
	}

	tlsConfig.RootCAs = x509.NewCertPool()
	ok := tlsConfig.RootCAs.AppendCertsFromPEM(certs)

	if !ok {
		return tlsConfig, errors.New("failed parsing pem file")
	}

	return tlsConfig, nil
}

func main() {
	cfg = configuration{
		MongoDbHost: os.Getenv(MongoDbHostKey),
		Database:    os.Getenv(MongoDbNameKey),
		Collection:  os.Getenv(MonngoCollectionKey),
		Username:    os.Getenv(MongoUserNameKey),
		Password:    os.Getenv(MongoPasswordHostKey),
		Certfile:    "/opt/rds-combined-ca-bundle.pem",
	}
	lambda.Start(handler)
}
