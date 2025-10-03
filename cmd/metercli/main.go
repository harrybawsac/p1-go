package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"

	"github.com/harrybawsac/p1-go/src/app"
	"github.com/harrybawsac/p1-go/src/buffer"
	"github.com/harrybawsac/p1-go/src/config"
	"github.com/harrybawsac/p1-go/src/scheduler"
	"github.com/harrybawsac/p1-go/src/services/db"
	"github.com/harrybawsac/p1-go/src/services/parser"
	_ "github.com/lib/pq"
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

	loop := flag.Bool("loop", false, "run in loop mode (use scheduler)")
	interval := flag.Int("interval", 60, "interval in seconds when running in loop mode")
	drain := flag.Bool("drain-buffer", false, "drain local buffer and attempt to persist entries")

	// parse flags (note: we already parsed earlier; reparse to pick new flags)
	flag.Parse()

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatalf("DB_DSN not set; provide in config or environment")
	}
	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer dbConn.Close()

	adapter := &db.PostgresAdapter{DB: dbConn}
	buf := buffer.New("/tmp/p1-buffer.jsonl")

	if *drain {
		if err := buf.Drain(ctx, func(ctx context.Context, raw json.RawMessage) error {
			// attempt to parse and insert
			r, externals, err := parser.ParseFullReading([]byte(raw))
			if err != nil {
				return err
			}
			return adapter.InsertReading(ctx, r, externals)
		}); err != nil {
			log.Fatalf("drain buffer failed: %v", err)
		}
		log.Println("drain completed")
		return
	}

	if *loop {
		s := &scheduler.Scheduler{
			DB:       dbConn,
			LockKey:  42,
			Interval: time.Duration(*interval) * time.Second,
		}
		if err := s.Run(ctx, func(ctx context.Context) error { return app.RunOnceWithDeps(ctx, adapter, buf) }); err != nil {
			log.Fatalf("scheduler failed: %v", err)
		}
	} else {
		if err := app.RunOnceWithDeps(ctx, adapter, buf); err != nil {
			log.Fatalf("run failed: %v", err)
		}
	}
	log.Println("run completed")
}

// runOnceWithDeps performs a single fetch -> parse -> persist cycle using injected dependencies.
// run logic moved to src/app/runner.go
