package buffer

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// Buffer persists items to a JSON-lines file for later draining
type Buffer struct {
	path string
	mu   sync.Mutex
}

func New(path string) *Buffer {
	return &Buffer{path: path}
}

// Append writes an arbitrary JSON-serializable object as a line
func (b *Buffer) Append(v interface{}) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	f, err := os.OpenFile(b.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	if err := enc.Encode(v); err != nil {
		return err
	}
	return nil
}

// Drain reads all entries and attempts to persist them by calling persistFn.
// On success, the buffer file is truncated. If persistFn returns error for an
// entry, the entry is kept (not retried individually) and Drain returns error.
func (b *Buffer) Drain(ctx context.Context, persistFn func(context.Context, json.RawMessage) error) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	f, err := os.OpenFile(b.path, os.O_RDONLY, 0)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // nothing to do
		}
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var lines [][]byte
	for scanner.Scan() {
		line := append([]byte(nil), scanner.Bytes()...)
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	// attempt to persist all lines
	for _, l := range lines {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if err := persistFn(ctx, json.RawMessage(l)); err != nil {
			return fmt.Errorf("persist line: %w", err)
		}
	}

	// truncate file on success
	if err := os.Truncate(b.path, 0); err != nil {
		return err
	}
	return nil
}
