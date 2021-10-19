package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(c context.Context, e events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	return events.APIGatewayV2HTTPResponse{}, nil
}

func main() {

	lambda.Start(handler)
}
