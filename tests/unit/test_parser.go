package unit

import (
    "testing"
    "github.com/harrybawsac/p1-go/src/services/parser"
)

func TestParseMeterJSON(t *testing.T) {
    data := []byte(`{"value": 12.34, "unit": "kWh"}`)
    p, err := parser.ParseMeterJSON(data)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if p.Unit != "kWh" {
        t.Fatalf("expected kWh, got %s", p.Unit)
    }
}

