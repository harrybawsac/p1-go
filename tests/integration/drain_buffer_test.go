package integration

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/harrybawsac/p1-go/src/buffer"
)

func TestBufferDrain_PersistsLines(t *testing.T) {
	tmp := "/tmp/p1-buffer-test-integration.jsonl"
	os.Remove(tmp)
	f, err := os.Create(tmp)
	if err != nil {
		t.Fatalf("create tmp buffer: %v", err)
	}
	sample := map[string]interface{}{"foo": "bar"}
	if err := json.NewEncoder(f).Encode(sample); err != nil {
		t.Fatalf("write sample: %v", err)
	}
	f.Close()

	b := buffer.New(tmp)
	called := false
	persist := func(ctx context.Context, raw json.RawMessage) error {
		called = true
		return nil
	}

	if err := b.Drain(context.Background(), persist); err != nil {
		t.Fatalf("drain failed: %v", err)
	}
	if !called {
		t.Fatalf("persist function was not called")
	}
}
