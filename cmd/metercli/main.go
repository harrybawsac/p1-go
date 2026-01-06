package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/harrybawsac/p1-go/src/app"
	"github.com/harrybawsac/p1-go/src/buffer"
	"github.com/harrybawsac/p1-go/src/config"
	"github.com/harrybawsac/p1-go/src/models"
	"github.com/harrybawsac/p1-go/src/scheduler"
	"github.com/harrybawsac/p1-go/src/services/csvloader"
	"github.com/harrybawsac/p1-go/src/services/db"
	"github.com/harrybawsac/p1-go/src/services/parser"
	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()
	cfgPath := flag.String("config", "./config.json", "path to JSON config file")
	loop := flag.Bool("loop", false, "run in loop mode (use scheduler)")
	interval := flag.Int("interval", 60, "interval in seconds when running in loop mode")
	drain := flag.Bool("drain-buffer", false, "drain local buffer and attempt to persist entries")
	dryRun := flag.Bool("dry-run", false, "fetch and log data without inserting into database")
	importCSV := flag.Bool("import", false, "import CSV files from data directory")
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

	if *importCSV {
		if err := importCSVData(ctx, cfg, adapter, *dryRun); err != nil {
			log.Fatalf("import CSV failed: %v", err)
		}
		log.Println("import completed")
		return
	}

	if *drain {
		if err := buf.Drain(ctx, func(ctx context.Context, raw json.RawMessage) error {
			// attempt to parse and insert
			r, err := parser.ParseFullReading([]byte(raw))
			if err != nil {
				return err
			}
			return adapter.InsertReading(ctx, r)
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
		if err := s.Run(ctx, func(ctx context.Context) error { return app.RunOnceWithDeps(ctx, adapter, buf, *dryRun) }); err != nil {
			log.Fatalf("scheduler failed: %v", err)
		}
	} else {
		if err := app.RunOnceWithDeps(ctx, adapter, buf, *dryRun); err != nil {
			log.Fatalf("run failed: %v", err)
		}
	}
	log.Println("run completed")
}

// importCSVData loads CSV files and imports them into the database day-by-day
func importCSVData(ctx context.Context, cfg config.Config, adapter *db.PostgresAdapter, dryRun bool) error {
	if cfg.DataDir == "" {
		return fmt.Errorf("data_dir not configured")
	}

	loader := &csvloader.CSVLoader{DataDir: cfg.DataDir}

	log.Printf("Loading CSV files from %s...\n", cfg.DataDir)
	merged, err := loader.LoadAndMerge()
	if err != nil {
		return fmt.Errorf("load and merge CSV: %w", err)
	}
	log.Printf("Loaded %d records\n", len(merged))

	// Group by day
	dayMap := csvloader.GroupByDay(merged)

	// Sort days for consistent processing
	var days []string
	for day := range dayMap {
		days = append(days, day)
	}
	sort.Strings(days)

	log.Printf("Found %d days of data\n", len(days))

	if dryRun {
		// Only process first two days for dry-run
		processCount := 2
		if len(days) < 2 {
			processCount = len(days)
		}

		log.Printf("Dry-run mode: generating SQL for first %d days\n", processCount)
		for i := 0; i < processCount; i++ {
			day := days[i]
			dayReadings := dayMap[day]

			log.Printf("\n--- Day %d: %s (%d readings) ---\n", i+1, day, len(dayReadings))

			// Convert to models.Reading
			readings := make([]models.Reading, len(dayReadings))
			for j, mr := range dayReadings {
				readings[j] = mr.ToReading()
			}

			// Generate SQL (single statement for all readings of the day)
			sqlStatement := adapter.GenerateInsertSQL(readings)
			fmt.Println(sqlStatement)
		}
		return nil
	}

	// Insert day by day
	for i, day := range days {
		dayReadings := dayMap[day]
		log.Printf("Processing day %d/%d: %s (%d readings)\n", i+1, len(days), day, len(dayReadings))

		// Convert to models.Reading
		readings := make([]models.Reading, len(dayReadings))
		for j, mr := range dayReadings {
			readings[j] = mr.ToReading()
		}

		// Insert batch for this day
		if err := adapter.InsertReadingsBatch(ctx, readings); err != nil {
			return fmt.Errorf("insert day %s: %w", day, err)
		}
		log.Printf("Inserted %d readings for %s\n", len(readings), day)
	}

	return nil
}

// runOnceWithDeps performs a single fetch -> parse -> persist cycle using injected dependencies.
// run logic moved to src/app/runner.go
