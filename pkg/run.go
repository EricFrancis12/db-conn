package pkg

import (
	"errors"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func Run(c *Config) error {
	if err := c.Init(); err != nil {
		return err
	}

	connCount := len(c.ConnStrs)
	errChan := make(chan error, connCount)
	var wg sync.WaitGroup

	c.Logger.Printf("Starting (%d) connection(s)\n", connCount)
	now := time.Now()

	for _, connStr := range c.ConnStrs {
		wg.Add(1)
		go func(cs string) {
			err := c.Connect(c.Ctx, c.DriverName, cs, c.ConnTimeout)
			if err != nil {
				errChan <- fmt.Errorf("error connecting to '%s': %v", cs, err)
			}
			wg.Done()
		}(connStr)
	}

	wg.Wait()
	took := time.Since(now)
	c.Logger.Printf("Finished (%d) connection(s) in %s\n", connCount, took.String())
	close(errChan)

	// len() measures unread elements in a channel,
	// so we must take the len before ranging over
	errorCount := len(errChan)

	if c.ListConnErrors && errorCount > 0 {
		c.Logger.Println("--- START ERRORS ---")
		for err := range errChan {
			c.Logger.Printf("%v\n", err)
		}
		c.Logger.Println("---  END ERRORS  ---")
	}

	successCount := connCount - errorCount
	successRate := float64(successCount) / float64(connCount) * 100

	msg := fmt.Sprintf(
		"%d/%d connection attempt(s) succeeded (%.2f%% success rate)",
		successCount,
		connCount,
		successRate,
	)

	if errorCount > 0 {
		return errors.New(msg)
	}

	c.Logger.Println(msg)
	return nil
}
