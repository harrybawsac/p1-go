package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/harrybawsac/p1-go/src/models"
)

func TestParseFullReading_File(t *testing.T) {
	path := filepath.Join("..", "..", "..", "specs", "001-build-a-cli", "contracts", "meter_sample.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read sample json: %v", err)
	}
	r, ext, err := ParseFullReading(data)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if r.UniqueID == "" {
		t.Fatalf("expected unique_id")
	}
	if len(ext) == 0 {
		t.Fatalf("expected external readings")
	}
	// quick sanity check types
	_ = models.Reading{}
}
