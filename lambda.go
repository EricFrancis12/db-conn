package main

import (
	"context"
	"log"
)

type lambdaResponse struct {
	Error error `json:"error"`
}

func lambdaHandler(ctx context.Context, _ any) lambdaResponse {
	log.Println("HELLO")
	logger := logger{}

	err := run(
		&Config{
			Ctx:    ctx,
			Logger: logger,
		},
	)
	if err != nil {
		logger.Printf("%v\n", err)
	}

	return lambdaResponse{
		Error: nil,
	}
}
