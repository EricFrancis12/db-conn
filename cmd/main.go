package main

import (
	"context"
	"flag"
	"log"

	"db-conn/pkg"
)

func main() {
	var (
		filePath       = flag.String("f", pkg.DefaultFilePath, "Path to targets file. Each target should be on a new line")
		driverName     = flag.String("d", pkg.DefaultDriverName, "SQL driver name")
		connTimeout    = flag.Duration("t", pkg.DefaultConnTimeout, "DB connection timeout")
		listConnErrors = flag.Bool("errors", true, "Prints connection errors to terminal")
	)
	flag.Parse()

	connStrs, err := pkg.ReadToConnStrs(*filePath)
	if err != nil {
		log.Fatal(err)
	}

	err = pkg.Run(
		&pkg.Config{
			Ctx:            context.Background(),
			Logger:         pkg.NewLogger(),
			ConnStrs:       connStrs,
			DriverName:     *driverName,
			ConnTimeout:    *connTimeout,
			ListConnErrors: *listConnErrors,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}
