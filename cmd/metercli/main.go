package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/harrybawsac/p1-go/src/config"
)

func main() {
	ctx := context.Background()
	cfgPath := flag.String("config", "./config.json", "path to JSON config file")
	flag.Parse()

	log.Println("metercli starting")

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	// expose config via env for other packages that might expect env vars
	if cfg.MeterEndpoint != "" {
		os.Setenv("METER_ENDPOINT", cfg.MeterEndpoint)
	}
	if cfg.DBDSN != "" {
		os.Setenv("DB_DSN", cfg.DBDSN)
	}

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
