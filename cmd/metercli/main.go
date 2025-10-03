package main

import (
	"context"
	"log"
	"time"
)

func main() {
	ctx := context.Background()
	log.Println("metercli starting")

	// TODO: load config
	// TODO: call run once (scheduler runs externally via cron)
	if err := runOnce(ctx); err != nil {
		log.Fatalf("run failed: %v", err)
	}
	log.Println("run completed")
}

func runOnce(ctx context.Context) error {
	// Placeholder for orchestration: call meter clients, parse, and insert
	// This will be implemented in later tasks.
	time.Sleep(100 * time.Millisecond)
	return nil
}
