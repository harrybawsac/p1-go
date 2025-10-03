package integration

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/lib/pq"

	dbpkg "github.com/harrybawsac/p1-go/src/services/db"
	"github.com/harrybawsac/p1-go/src/services/parser"
)

func must(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("err: %v", err)
	}
}

func TestIntegration_InsertReading(t *testing.T) {
	// Expect a local postgres running on localhost:5432 (docker-compose up -d)
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		t.Skip("skipping integration test: set TEST_DATABASE_DSN and start docker compose to run")
	}

	// wait a bit for db to come up (simple retry)
	var db *sql.DB
	var err error
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			err = db.Ping()
		}
		if err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	must(t, err)
	defer db.Close()

	// apply migration
	migPath := filepath.Join("..", "..", "migrations", "001_create_tables.sql")
	mig, err := os.ReadFile(migPath)
	must(t, err)
	if _, err := db.Exec(string(mig)); err != nil {
		t.Fatalf("apply migration: %v", err)
	}

	adapter := &dbpkg.PostgresAdapter{DB: db}

	// read sample JSON
	sample := filepath.Join("..", "..", "specs", "001-build-a-cli", "contracts", "meter_sample.json")
	data, err := os.ReadFile(sample)
	must(t, err)

	r, externals, err := parser.ParseFullReading(data)
	must(t, err)

	// insert reading
	ctx := context.Background()
	if err := adapter.InsertReading(ctx, r, externals); err != nil {
		t.Fatalf("insert reading: %v", err)
	}

	// verify
	var count int
	row := db.QueryRow("SELECT count(1) FROM p1.meter_readings WHERE unique_id = $1", r.UniqueID)
	must(t, row.Scan(&count))
	if count != 1 {
		t.Fatalf("expected 1 row, got %d", count)
	}

	// verify external readings if present
	if len(externals) > 0 {
		row = db.QueryRow("SELECT count(1) FROM p1.external_readings WHERE meter_reading_unique_id = $1", r.UniqueID)
		must(t, row.Scan(&count))
		if count != len(externals) {
			t.Fatalf("expected %d external rows, got %d", len(externals), count)
		}
	}
}
