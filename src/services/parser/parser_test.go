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
	r, err := ParseFullReading(data)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if r.ActiveTariff == 0 {
		t.Logf("note: active_tariff is 0 (may be expected)")
	}
	// quick sanity check types
	_ = models.Reading{}
}
