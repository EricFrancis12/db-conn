package main

import (
	"context"
	"log"
	"os"

	"db-conn/pkg"
)

func main() {
	err := pkg.Run(
		&pkg.Config{
			Ctx:    context.Background(),
			Logger: pkg.NewLogger(),
			Args:   os.Args,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}
