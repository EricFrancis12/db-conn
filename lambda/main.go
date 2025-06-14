package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"

	"db-conn/pkg"
)

func lambdaHandler(ctx context.Context, _ any) error {
	logger := pkg.NewLogger()

	err := pkg.Run(
		&pkg.Config{
			Ctx:            ctx,
			Logger:         logger,
			ConnStrs:       targets[:],
			ListConnErrors: true,
		},
	)
	if err != nil {
		logger.Printf("%v\n", err)
	}

	return err
}

func main() {
	lambda.Start(lambdaHandler)
}
