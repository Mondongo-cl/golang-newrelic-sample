package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
)

type Configuration struct {
	username string
	password string
	port     int
	host     string
	database string
}

var (
	cfg Configuration
)

func (c Configuration) Parse() (string, error) {
	if c.username == "" {
		return "", errors.New("username is required")
	}
	if c.port <= 0 {
		return "", errors.New("port number must be greather than 0")
	}
	if c.host == "" {
		return "", errors.New("hostname can not be empty")
	}
	if c.database == "" {
		return "", errors.New("database name can not be null")
	}

	cnnstr := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", c.username, c.password, c.host, c.port, c.database)
	return cnnstr, nil
}

func Handle(c context.Context, e events.APIGatewayV2HTTPResponse) (events.APIGatewayV2HTTPResponse, error) {
	connectionstring, err := cfg.Parse()
	if err != nil {
		log.Panic(err)
	}
	db, err := sql.Open("mysql", connectionstring)

	if err != nil {
		log.Panic(err)
	}
	mysqldb := db.QueryRow("select version()")
	if mysqldb.Err() != nil {
		log.Panic(mysqldb.Err())
	}
	var s string
	mysqldb.Scan(&s)
	db.Close()
	return events.APIGatewayV2HTTPResponse{
		Headers:         e.Headers,
		Body:            s,
		StatusCode:      200,
		IsBase64Encoded: false,
	}, nil
}

func main() {
	port, err := strconv.Atoi(os.Getenv("MYSQL_PORT"))
	if err != nil {
		port = 3306
	}
	cfg = Configuration{
		username: os.Getenv("MYSQL_USERNAME"),
		password: os.Getenv("MYSQL_PASSWORD"),
		port:     int(port),
		host:     os.Getenv("MYSQL_HOST"),
		database: os.Getenv("MYSQL_DATABASE"),
	}

	lambda.Start(Handle)
}
