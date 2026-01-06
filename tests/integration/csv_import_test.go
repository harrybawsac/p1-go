package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/harrybawsac/p1-go/src/services/csvloader"
)

// TestCSVImportIntegration tests the full CSV import workflow
func TestCSVImportIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Create temp directory for CSV files
	tmpDir := t.TempDir()

	// Create test power CSV with multiple days
	powerCSV := filepath.Join(tmpDir, "power-15m.csv")
	powerContent := `time,Import T1 kWh,Import T2 kWh,Export T1 kWh,Export T2 kWh,L1 max W,L2 max W,L3 max W
2025-06-02 20:30,8293.146,7210.113,1916.077,4181.422,173,1212,67
2025-06-02 20:45,8293.146,7210.236,1916.077,4181.422,127,48,77
2025-06-02 21:00,8293.147,7210.276,1916.077,4181.422,93,56,78
2025-06-03 08:00,8293.200,7210.300,1916.077,4181.422,100,60,80
2025-06-03 08:15,8293.250,7210.350,1916.077,4181.422,110,70,90`

	if err := os.WriteFile(powerCSV, []byte(powerContent), 0644); err != nil {
		t.Fatalf("write power file: %v", err)
	}

	// Create test gas CSV
	gasCSV := filepath.Join(tmpDir, "gas-15m.csv")
	gasContent := `time,Total gas used
2025-06-02 20:30,3488.524
2025-06-02 20:45,3488.524
2025-06-02 21:00,3488.530
2025-06-03 08:00,3488.600
2025-06-03 08:15,3488.650`

	if err := os.WriteFile(gasCSV, []byte(gasContent), 0644); err != nil {
		t.Fatalf("write gas file: %v", err)
	}

	// Load and merge CSV files
	loader := &csvloader.CSVLoader{DataDir: tmpDir}
	merged, err := loader.LoadAndMerge()
	if err != nil {
		t.Fatalf("LoadAndMerge failed: %v", err)
	}

	if len(merged) != 5 {
		t.Errorf("expected 5 merged records, got %d", len(merged))
	}

	// Group by day
	dayMap := csvloader.GroupByDay(merged)

	if len(dayMap) != 2 {
		t.Errorf("expected 2 days, got %d", len(dayMap))
	}

	// Verify day 1 has 3 records
	day1 := "2025-06-02"
	if len(dayMap[day1]) != 3 {
		t.Errorf("expected 3 records for day 1, got %d", len(dayMap[day1]))
	}

	// Verify day 2 has 2 records
	day2 := "2025-06-03"
	if len(dayMap[day2]) != 2 {
		t.Errorf("expected 2 records for day 2, got %d", len(dayMap[day2]))
	}

	// Verify data integrity for first record
	firstRecord := merged[0]
	expectedTime, _ := time.Parse("2006-01-02 15:04", "2025-06-02 20:30")
	if !firstRecord.Time.Equal(expectedTime) {
		t.Errorf("expected time %v, got %v", expectedTime, firstRecord.Time)
	}
	if firstRecord.ImportT1Kwh != 8293.146 {
		t.Errorf("expected ImportT1Kwh=8293.146, got %f", firstRecord.ImportT1Kwh)
	}
	if firstRecord.TotalGasM3 != 3488.524 {
		t.Errorf("expected TotalGasM3=3488.524, got %f", firstRecord.TotalGasM3)
	}
}

// TestCSVImportBatchByDay tests that readings are correctly batched by day
func TestCSVImportBatchByDay(t *testing.T) {
	tmpDir := t.TempDir()

	// Create CSV files spanning 3 days with varying record counts
	powerCSV := filepath.Join(tmpDir, "power-15m.csv")
	powerContent := `time,Import T1 kWh,Import T2 kWh,Export T1 kWh,Export T2 kWh,L1 max W,L2 max W,L3 max W
2025-06-01 23:45,100.0,200.0,50.0,75.0,100,200,150
2025-06-02 00:00,100.1,200.1,50.1,75.1,101,201,151
2025-06-02 00:15,100.2,200.2,50.2,75.2,102,202,152
2025-06-03 12:00,100.3,200.3,50.3,75.3,103,203,153`

	if err := os.WriteFile(powerCSV, []byte(powerContent), 0644); err != nil {
		t.Fatalf("write power file: %v", err)
	}

	gasCSV := filepath.Join(tmpDir, "gas-15m.csv")
	gasContent := `time,Total gas used
2025-06-01 23:45,1000.0
2025-06-02 00:00,1000.1
2025-06-02 00:15,1000.2
2025-06-03 12:00,1000.3`

	if err := os.WriteFile(gasCSV, []byte(gasContent), 0644); err != nil {
		t.Fatalf("write gas file: %v", err)
	}

	loader := &csvloader.CSVLoader{DataDir: tmpDir}
	merged, err := loader.LoadAndMerge()
	if err != nil {
		t.Fatalf("LoadAndMerge failed: %v", err)
	}

	dayMap := csvloader.GroupByDay(merged)

	// Should have 3 distinct days
	if len(dayMap) != 3 {
		t.Errorf("expected 3 days, got %d", len(dayMap))
	}

	// Verify each day's record count
	if len(dayMap["2025-06-01"]) != 1 {
		t.Errorf("expected 1 record for 2025-06-01, got %d", len(dayMap["2025-06-01"]))
	}
	if len(dayMap["2025-06-02"]) != 2 {
		t.Errorf("expected 2 records for 2025-06-02, got %d", len(dayMap["2025-06-02"]))
	}
	if len(dayMap["2025-06-03"]) != 1 {
		t.Errorf("expected 1 record for 2025-06-03, got %d", len(dayMap["2025-06-03"]))
	}
}

// TestCSVImportDataTypes tests that all data types are correctly converted
func TestCSVImportDataTypes(t *testing.T) {
	tmpDir := t.TempDir()

	powerCSV := filepath.Join(tmpDir, "power-15m.csv")
	powerContent := `time,Import T1 kWh,Import T2 kWh,Export T1 kWh,Export T2 kWh,L1 max W,L2 max W,L3 max W
2025-06-02 20:30,8293.146,7210.113,1916.077,4181.422,173.5,1212.8,67.2`

	if err := os.WriteFile(powerCSV, []byte(powerContent), 0644); err != nil {
		t.Fatalf("write power file: %v", err)
	}

	gasCSV := filepath.Join(tmpDir, "gas-15m.csv")
	gasContent := `time,Total gas used
2025-06-02 20:30,3488.524`

	if err := os.WriteFile(gasCSV, []byte(gasContent), 0644); err != nil {
		t.Fatalf("write gas file: %v", err)
	}

	loader := &csvloader.CSVLoader{DataDir: tmpDir}
	merged, err := loader.LoadAndMerge()
	if err != nil {
		t.Fatalf("LoadAndMerge failed: %v", err)
	}

	if len(merged) != 1 {
		t.Fatalf("expected 1 record, got %d", len(merged))
	}

	reading := merged[0].ToReading()

	// Verify all fields are correctly mapped
	if reading.TotalPowerImportT1Kwh != 8293.146 {
		t.Errorf("ImportT1Kwh: expected 8293.146, got %f", reading.TotalPowerImportT1Kwh)
	}
	if reading.TotalPowerImportT2Kwh != 7210.113 {
		t.Errorf("ImportT2Kwh: expected 7210.113, got %f", reading.TotalPowerImportT2Kwh)
	}
	if reading.TotalPowerExportT1Kwh != 1916.077 {
		t.Errorf("ExportT1Kwh: expected 1916.077, got %f", reading.TotalPowerExportT1Kwh)
	}
	if reading.TotalPowerExportT2Kwh != 4181.422 {
		t.Errorf("ExportT2Kwh: expected 4181.422, got %f", reading.TotalPowerExportT2Kwh)
	}
	if reading.ActivePowerL1W != 173.5 {
		t.Errorf("L1MaxW: expected 173.5, got %f", reading.ActivePowerL1W)
	}
	if reading.ActivePowerL2W != 1212.8 {
		t.Errorf("L2MaxW: expected 1212.8, got %f", reading.ActivePowerL2W)
	}
	if reading.ActivePowerL3W != 67.2 {
		t.Errorf("L3MaxW: expected 67.2, got %f", reading.ActivePowerL3W)
	}
	if reading.TotalGasM3 != 3488.524 {
		t.Errorf("TotalGasM3: expected 3488.524, got %f", reading.TotalGasM3)
	}
}

// TestFullImportWorkflow simulates the complete import process
func TestFullImportWorkflow(t *testing.T) {
	tmpDir := t.TempDir()

	// Create realistic multi-day dataset
	powerCSV := filepath.Join(tmpDir, "power-15m.csv")
	var powerLines []string
	powerLines = append(powerLines, "time,Import T1 kWh,Import T2 kWh,Export T1 kWh,Export T2 kWh,L1 max W,L2 max W,L3 max W")

	gasCSV := filepath.Join(tmpDir, "gas-15m.csv")
	var gasLines []string
	gasLines = append(gasLines, "time,Total gas used")

	// Generate 3 days of data, 4 readings per day
	baseTime := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	for day := 0; day < 3; day++ {
		for reading := 0; reading < 4; reading++ {
			t := baseTime.Add(time.Duration(day*24+reading*6) * time.Hour)
			timeStr := t.Format("2006-01-02 15:04")

			powerLine := fmt.Sprintf("%s,%.3f,%.3f,%.3f,%.3f,%d,%d,%d",
				timeStr, 100.0+float64(day), 200.0+float64(day), 50.0, 75.0, 100, 200, 150)
			powerLines = append(powerLines, powerLine)

			gasLine := fmt.Sprintf("%s,%.3f", timeStr, 1000.0+float64(day)*10)
			gasLines = append(gasLines, gasLine)
		}
	}

	powerContent := ""
	for i, line := range powerLines {
		if i > 0 {
			powerContent += "\n"
		}
		powerContent += line
	}

	gasContent := ""
	for i, line := range gasLines {
		if i > 0 {
			gasContent += "\n"
		}
		gasContent += line
	}

	if err := os.WriteFile(powerCSV, []byte(powerContent), 0644); err != nil {
		t.Fatalf("write power file: %v", err)
	}

	if err := os.WriteFile(gasCSV, []byte(gasContent), 0644); err != nil {
		t.Fatalf("write gas file: %v", err)
	}

	// Execute the workflow
	loader := &csvloader.CSVLoader{DataDir: tmpDir}
	merged, err := loader.LoadAndMerge()
	if err != nil {
		t.Fatalf("LoadAndMerge failed: %v", err)
	}

	// Should have 12 total records (3 days * 4 readings)
	if len(merged) != 12 {
		t.Errorf("expected 12 merged records, got %d", len(merged))
	}

	// Group by day
	dayMap := csvloader.GroupByDay(merged)

	// Should have 3 days
	if len(dayMap) != 3 {
		t.Errorf("expected 3 days, got %d", len(dayMap))
	}

	// Each day should have 4 readings
	for day, readings := range dayMap {
		if len(readings) != 4 {
			t.Errorf("day %s: expected 4 readings, got %d", day, len(readings))
		}
	}

	// Simulate day-by-day processing
	processedDays := 0
	processedReadings := 0

	for _, dayReadings := range dayMap {
		// Convert to model readings
		for _, mr := range dayReadings {
			_ = mr.ToReading()
			processedReadings++
		}
		processedDays++
	}

	if processedDays != 3 {
		t.Errorf("expected to process 3 days, processed %d", processedDays)
	}

	if processedReadings != 12 {
		t.Errorf("expected to process 12 readings, processed %d", processedReadings)
	}
}
