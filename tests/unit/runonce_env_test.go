package unit

import (
	"context"
	"os"
	"testing"

	"github.com/harrybawsac/p1-go/src/app"
	"github.com/harrybawsac/p1-go/src/buffer"
	"github.com/harrybawsac/p1-go/src/services/db"
)

// This test ensures that RunOnceWithDeps returns an error when METER_ENDPOINT is missing.
func TestRunOnce_MissingEndpoint(t *testing.T) {
	// ensure env var unset
	os.Unsetenv("METER_ENDPOINT")

	// use nil adapter and buffer for this check — function should fail early
	adapter := &db.PostgresAdapter{DB: nil}
	buf := buffer.New("/tmp/p1-buffer-test.jsonl")

	if err := app.RunOnceWithDeps(context.Background(), adapter, buf, false); err == nil {
		t.Fatalf("expected error when METER_ENDPOINT missing, got nil")
	}
}
