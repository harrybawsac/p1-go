package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/harrybawsac/p1-go/src/buffer"
	"github.com/harrybawsac/p1-go/src/services/db"
	"github.com/harrybawsac/p1-go/src/services/parser"
)

// RunOnceWithDeps performs a single fetch -> parse -> persist cycle using injected dependencies.
func RunOnceWithDeps(ctx context.Context, adapter *db.PostgresAdapter, buf *buffer.Buffer) error {
	endpoint := os.Getenv("METER_ENDPOINT")
	if endpoint == "" {
		return fmt.Errorf("METER_ENDPOINT not set")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("meter endpoint status: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	r, externals, err := parser.ParseFullReading(body)
	if err != nil {
		return err
	}

	if err := adapter.InsertReading(ctx, r, externals); err != nil {
		// buffer raw payload for retry
		if berr := buf.Append(json.RawMessage(body)); berr != nil {
			return fmt.Errorf("insert failed: %v; buffer append failed: %v", err, berr)
		}
		return err
	}
	return nil
}
