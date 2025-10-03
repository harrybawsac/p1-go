package buffer

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestAppendAndDrain(t *testing.T) {
	tmp := filepath.Join(os.TempDir(), "p1-buffer-test.jsonl")
	os.Remove(tmp)
	defer os.Remove(tmp)

	b := New(tmp)
	sample := map[string]interface{}{"a": 1}
	if err := b.Append(sample); err != nil {
		t.Fatalf("append: %v", err)
	}

	called := 0
	persistFn := func(ctx context.Context, raw json.RawMessage) error {
		called++
		return nil
	}

	if err := b.Drain(context.Background(), persistFn); err != nil {
		t.Fatalf("drain: %v", err)
	}
	if called != 1 {
		t.Fatalf("expected called once, got %d", called)
	}
}
