package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"ramenman/handler"
)

func main() {
	lambda.Start(handler.HandleRequest)
}