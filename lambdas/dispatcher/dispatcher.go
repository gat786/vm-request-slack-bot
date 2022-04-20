package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

type DispatchEvent struct{}
type Result struct{}

func Handler(ctx context.Context, event DispatchEvent) (Result, error) {
	return Result{}, nil
}

func main() {
	lambda.Start(Handler)
}
