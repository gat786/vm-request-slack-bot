package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

type NotifyEvent struct{}
type Result struct{}

func Handler(ctx context.Context, event NotifyEvent) (Result, error) {
	return Result{}, nil
}

func main() {
	lambda.Start(Handler)
}
