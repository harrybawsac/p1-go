package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	tmp := filepath.Join(os.TempDir(), "p1-config-test.json")
	content := []byte(`{"meter_endpoint":"http://example","db_dsn":"postgres://u:p@127.0.0.1:5432/db?sslmode=disable"}`)
	if err := os.WriteFile(tmp, content, 0644); err != nil {
		t.Fatalf("write tmp config: %v", err)
	}
	defer os.Remove(tmp)

	cfg, err := Load(tmp)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cfg.MeterEndpoint != "http://example" {
		t.Fatalf("unexpected meter endpoint: %s", cfg.MeterEndpoint)
	}
	if cfg.DBDSN == "" {
		t.Fatalf("expected db dsn")
	}
}
