package contract

import (
	"testing"

	"github.com/harrybawsac/p1-go/src/services/parser"
)

func TestParseFullReading_AcceptsExamplePayload(t *testing.T) {
	sample := []byte(`{
        "timestamp": "2025-10-03T12:00:00Z",
        "unique_id": "meter-001",
        "electricity": { "delivered": 123.45 },
        "externals": [{ "unique_id": "ext-1", "type": "sensor", "timestamp": 1696344000, "value": 12.3, "unit": "kWh" }]
    }`)

	_, _, err := parser.ParseFullReading(sample)
	if err != nil {
		t.Fatalf("expected ParseFullReading to accept sample payload, got error: %v", err)
	}
}
