package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

func run(c *Config) bool {
	if err := c.Init(); err != nil {
		c.Logger.Printf("%w\n", err)
		return false
	}

	connCount := len(c.ConnStrs)
	errChan := make(chan error, connCount)
	var wg sync.WaitGroup

	c.Logger.Printf("Starting (%d) connection(s)\n", connCount)
	now := time.Now()

	for _, connStr := range c.ConnStrs {
		wg.Add(1)
		go func(cs string) {
			err := connect(c.Ctx, c.DriverName, cs, c.ConnTimeout)
			if err != nil {
				errChan <- fmt.Errorf("error connecting to '%s': %v", connStr, err)
			}
			wg.Done()
		}(connStr)
	}

	wg.Wait()
	took := time.Since(now)
	c.Logger.Printf("Finished (%d) connection(s) in '%s'\n", connCount, took.String())
	close(errChan)

	if len(errChan) > 0 {
		c.Logger.Println("--- START ERRORS ---")
		for err := range errChan {
			c.Logger.Printf("%w", err)
		}
		c.Logger.Println("---  END ERRORS  ---")
	}

	errorCount := len(errChan)
	successCount := connCount - errorCount
	successRate := float64(successCount) / float64(connCount) * 100

	c.Logger.Println("--- START RESULTS ---")
	c.Logger.Printf("Sent %d connections in %s\n", connCount, took.String())
	c.Logger.Printf("%d/%d connections were successful\n", successCount, connCount)
	c.Logger.Printf("%d/%d connections failed\n", errorCount, connCount)
	c.Logger.Printf("Success Rate: %.2f%%\n", successRate)
	c.Logger.Println("---  END RESULTS  ---")

	return errorCount == 0
}

type ConnectFunc func(
	ctx context.Context,
	driverName string,
	connStr string,
	timeout time.Duration,
) error

func connect(
	ctx context.Context,
	driverName string,
	connStr string,
	timeout time.Duration,
) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	errChan := make(chan error)

	go func() {
		db, err := sql.Open(driverName, connStr)
		if err != nil {
			errChan <- err
			return
		}
		defer db.Close()
		errChan <- db.PingContext(ctx)
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("reached timeout of '%s'", timeout.String())
	case err := <-errChan:
		return err
	}
}

const defaultFilePath string = "targets.txt"

var defaultFileContent string

func init() {
	var once sync.Once
	once.Do(func() {
		b, err := os.ReadFile(defaultFilePath)
		if err == nil {
			defaultFileContent = string(b)
		}
	})
}

func main() {
	var (
		filePath    = flag.String("f", defaultFilePath, "Path to targets file. Each target should be on a new line")
		driverName  = flag.String("d", defaultDriverName, "SQL driver name")
		connTimeout = flag.Duration("t", defaultConnTimeout, "DB connection timeout")
		cleanRead   = flag.Bool("clean-read", false, "If true, always read the targets file from disk, ignoring any cached/default content")
	)
	flag.Parse()

	stat, err := os.Stat(*filePath)
	if os.IsNotExist(err) {
		log.Fatalf("file '%s' not found", *filePath)
	} else if err != nil {
		log.Fatalf("error reading file '%s': %s", *filePath, err.Error())
	} else if stat.IsDir() {
		log.Fatalf("'%s' needs to be a path to a file (found dir)", *filePath)
	}

	var content string
	if !*cleanRead && *filePath == defaultFilePath {
		content = defaultFileContent
	} else {
		b, err := os.ReadFile(*filePath)
		if err != nil {
			log.Fatalf("error reading file '%s': %s", *filePath, err.Error())
		}
		content = string(b)
	}

	success := run(
		&Config{
			Ctx:         context.Background(),
			Logger:      logger{},
			ConnStrs:    strings.Split(content, "\n"),
			DriverName:  *driverName,
			ConnTimeout: *connTimeout,
			Connect:     connect,
		},
	)
	if !success {
		os.Exit(1)
	}
	os.Exit(0)
}
