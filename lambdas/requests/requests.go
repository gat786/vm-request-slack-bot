package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

type RequestEvent struct{}
type Result struct{}

func Handler(ctx context.Context, event RequestEvent) (Result, error) {
	return Result{}, nil
}

func main() {
	lambda.Start(Handler)
}
