package csvloader

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestLoadPowerCSV tests reading a power CSV file
func TestLoadPowerCSV(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()

	// Create test CSV
	powerCSV := filepath.Join(tmpDir, "power-15m.csv")
	content := `time,Import T1 kWh,Import T2 kWh,Export T1 kWh,Export T2 kWh,L1 max W,L2 max W,L3 max W
2025-06-02 20:30,8293.146,7210.113,1916.077,4181.422,173,1212,67
2025-06-02 20:45,8293.146,7210.236,1916.077,4181.422,127,48,77
2025-06-02 21:00,8293.147,7210.276,1916.077,4181.422,93,56,78`

	if err := os.WriteFile(powerCSV, []byte(content), 0644); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	loader := &CSVLoader{DataDir: tmpDir}
	records, err := loader.readPowerCSV(powerCSV)
	if err != nil {
		t.Fatalf("readPowerCSV failed: %v", err)
	}

	if len(records) != 3 {
		t.Errorf("expected 3 records, got %d", len(records))
	}

	// Check first record
	if records[0].ImportT1Kwh != 8293.146 {
		t.Errorf("expected ImportT1Kwh=8293.146, got %f", records[0].ImportT1Kwh)
	}
	if records[0].L1MaxW != 173 {
		t.Errorf("expected L1MaxW=173, got %f", records[0].L1MaxW)
	}

	expectedTime, _ := time.Parse("2006-01-02 15:04", "2025-06-02 20:30")
	if !records[0].Time.Equal(expectedTime) {
		t.Errorf("expected time %v, got %v", expectedTime, records[0].Time)
	}
}

// TestLoadGasCSV tests reading a gas CSV file
func TestLoadGasCSV(t *testing.T) {
	tmpDir := t.TempDir()

	gasCSV := filepath.Join(tmpDir, "gas-15m.csv")
	content := `time,Total gas used
2025-06-02 20:30,3488.524
2025-06-02 20:45,3488.524
2025-06-02 21:00,3488.600`

	if err := os.WriteFile(gasCSV, []byte(content), 0644); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	loader := &CSVLoader{DataDir: tmpDir}
	records, err := loader.readGasCSV(gasCSV)
	if err != nil {
		t.Fatalf("readGasCSV failed: %v", err)
	}

	if len(records) != 3 {
		t.Errorf("expected 3 records, got %d", len(records))
	}

	if records[0].TotalGasM3 != 3488.524 {
		t.Errorf("expected TotalGasM3=3488.524, got %f", records[0].TotalGasM3)
	}

	if records[2].TotalGasM3 != 3488.600 {
		t.Errorf("expected TotalGasM3=3488.600, got %f", records[2].TotalGasM3)
	}
}

// TestLoadAndMerge tests merging power and gas CSV files
func TestLoadAndMerge(t *testing.T) {
	tmpDir := t.TempDir()

	// Create power CSV
	powerCSV := filepath.Join(tmpDir, "power-15m.csv")
	powerContent := `time,Import T1 kWh,Import T2 kWh,Export T1 kWh,Export T2 kWh,L1 max W,L2 max W,L3 max W
2025-06-02 20:30,8293.146,7210.113,1916.077,4181.422,173,1212,67
2025-06-02 20:45,8293.146,7210.236,1916.077,4181.422,127,48,77`

	if err := os.WriteFile(powerCSV, []byte(powerContent), 0644); err != nil {
		t.Fatalf("write power file: %v", err)
	}

	// Create gas CSV
	gasCSV := filepath.Join(tmpDir, "gas-15m.csv")
	gasContent := `time,Total gas used
2025-06-02 20:30,3488.524
2025-06-02 20:45,3488.530`

	if err := os.WriteFile(gasCSV, []byte(gasContent), 0644); err != nil {
		t.Fatalf("write gas file: %v", err)
	}

	loader := &CSVLoader{DataDir: tmpDir}
	merged, err := loader.LoadAndMerge()
	if err != nil {
		t.Fatalf("LoadAndMerge failed: %v", err)
	}

	if len(merged) != 2 {
		t.Errorf("expected 2 merged records, got %d", len(merged))
	}

	// Check first merged record
	if merged[0].ImportT1Kwh != 8293.146 {
		t.Errorf("expected ImportT1Kwh=8293.146, got %f", merged[0].ImportT1Kwh)
	}
	if merged[0].TotalGasM3 != 3488.524 {
		t.Errorf("expected TotalGasM3=3488.524, got %f", merged[0].TotalGasM3)
	}

	// Check second merged record
	if merged[1].TotalGasM3 != 3488.530 {
		t.Errorf("expected TotalGasM3=3488.530, got %f", merged[1].TotalGasM3)
	}
}

// TestLoadAndMergeMismatchedLength tests error handling for mismatched CSV lengths
func TestLoadAndMergeMismatchedLength(t *testing.T) {
	tmpDir := t.TempDir()

	// Create power CSV with 2 records
	powerCSV := filepath.Join(tmpDir, "power-15m.csv")
	powerContent := `time,Import T1 kWh,Import T2 kWh,Export T1 kWh,Export T2 kWh,L1 max W,L2 max W,L3 max W
2025-06-02 20:30,8293.146,7210.113,1916.077,4181.422,173,1212,67
2025-06-02 20:45,8293.146,7210.236,1916.077,4181.422,127,48,77`

	if err := os.WriteFile(powerCSV, []byte(powerContent), 0644); err != nil {
		t.Fatalf("write power file: %v", err)
	}

	// Create gas CSV with 3 records
	gasCSV := filepath.Join(tmpDir, "gas-15m.csv")
	gasContent := `time,Total gas used
2025-06-02 20:30,3488.524
2025-06-02 20:45,3488.530
2025-06-02 21:00,3488.540`

	if err := os.WriteFile(gasCSV, []byte(gasContent), 0644); err != nil {
		t.Fatalf("write gas file: %v", err)
	}

	loader := &CSVLoader{DataDir: tmpDir}
	_, err := loader.LoadAndMerge()
	if err == nil {
		t.Fatal("expected error for mismatched lengths, got nil")
	}
}

// TestLoadAndMergeMismatchedTimestamps tests error handling for mismatched timestamps
func TestLoadAndMergeMismatchedTimestamps(t *testing.T) {
	tmpDir := t.TempDir()

	// Create power CSV
	powerCSV := filepath.Join(tmpDir, "power-15m.csv")
	powerContent := `time,Import T1 kWh,Import T2 kWh,Export T1 kWh,Export T2 kWh,L1 max W,L2 max W,L3 max W
2025-06-02 20:30,8293.146,7210.113,1916.077,4181.422,173,1212,67
2025-06-02 20:45,8293.146,7210.236,1916.077,4181.422,127,48,77`

	if err := os.WriteFile(powerCSV, []byte(powerContent), 0644); err != nil {
		t.Fatalf("write power file: %v", err)
	}

	// Create gas CSV with different timestamp
	gasCSV := filepath.Join(tmpDir, "gas-15m.csv")
	gasContent := `time,Total gas used
2025-06-02 20:30,3488.524
2025-06-02 21:00,3488.530`

	if err := os.WriteFile(gasCSV, []byte(gasContent), 0644); err != nil {
		t.Fatalf("write gas file: %v", err)
	}

	loader := &CSVLoader{DataDir: tmpDir}
	_, err := loader.LoadAndMerge()
	if err == nil {
		t.Fatal("expected error for mismatched timestamps, got nil")
	}
}

// TestToReading tests converting MergedReading to models.Reading
func TestToReading(t *testing.T) {
	timeVal, _ := time.Parse("2006-01-02 15:04", "2025-06-02 20:30")
	merged := MergedReading{
		Time:        timeVal,
		ImportT1Kwh: 8293.146,
		ImportT2Kwh: 7210.113,
		ExportT1Kwh: 1916.077,
		ExportT2Kwh: 4181.422,
		L1MaxW:      173,
		L2MaxW:      1212,
		L3MaxW:      67,
		TotalGasM3:  3488.524,
	}

	reading := merged.ToReading()

	if !reading.CreatedAt.Equal(timeVal) {
		t.Errorf("expected CreatedAt=%v, got %v", timeVal, reading.CreatedAt)
	}
	if reading.TotalPowerImportT1Kwh != 8293.146 {
		t.Errorf("expected TotalPowerImportT1Kwh=8293.146, got %f", reading.TotalPowerImportT1Kwh)
	}
	if reading.ActivePowerL1W != 173 {
		t.Errorf("expected ActivePowerL1W=173, got %f", reading.ActivePowerL1W)
	}
	if reading.TotalGasM3 != 3488.524 {
		t.Errorf("expected TotalGasM3=3488.524, got %f", reading.TotalGasM3)
	}
}

// TestGroupByDay tests grouping readings by day
func TestGroupByDay(t *testing.T) {
	day1, _ := time.Parse("2006-01-02 15:04", "2025-06-02 20:30")
	day1b, _ := time.Parse("2006-01-02 15:04", "2025-06-02 21:00")
	day2, _ := time.Parse("2006-01-02 15:04", "2025-06-03 10:00")

	readings := []MergedReading{
		{Time: day1, ImportT1Kwh: 100},
		{Time: day1b, ImportT1Kwh: 101},
		{Time: day2, ImportT1Kwh: 200},
	}

	grouped := GroupByDay(readings)

	if len(grouped) != 2 {
		t.Errorf("expected 2 days, got %d", len(grouped))
	}

	day1Key := "2025-06-02"
	if len(grouped[day1Key]) != 2 {
		t.Errorf("expected 2 readings for day1, got %d", len(grouped[day1Key]))
	}

	day2Key := "2025-06-03"
	if len(grouped[day2Key]) != 1 {
		t.Errorf("expected 1 reading for day2, got %d", len(grouped[day2Key]))
	}
}

// TestInvalidPowerCSVFormat tests error handling for invalid power CSV format
func TestInvalidPowerCSVFormat(t *testing.T) {
	tmpDir := t.TempDir()

	powerCSV := filepath.Join(tmpDir, "power-15m.csv")
	// Missing a column
	content := `time,Import T1 kWh,Import T2 kWh,Export T1 kWh,Export T2 kWh,L1 max W,L2 max W
2025-06-02 20:30,8293.146,7210.113,1916.077,4181.422,173,1212`

	if err := os.WriteFile(powerCSV, []byte(content), 0644); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	loader := &CSVLoader{DataDir: tmpDir}
	_, err := loader.readPowerCSV(powerCSV)
	if err == nil {
		t.Fatal("expected error for invalid format, got nil")
	}
}

// TestInvalidGasCSVFormat tests error handling for invalid gas CSV format
func TestInvalidGasCSVFormat(t *testing.T) {
	tmpDir := t.TempDir()

	gasCSV := filepath.Join(tmpDir, "gas-15m.csv")
	// Missing Total gas column
	content := `time
2025-06-02 20:30`

	if err := os.WriteFile(gasCSV, []byte(content), 0644); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	loader := &CSVLoader{DataDir: tmpDir}
	_, err := loader.readGasCSV(gasCSV)
	if err == nil {
		t.Fatal("expected error for invalid format, got nil")
	}
}

// TestInvalidNumberFormat tests error handling for invalid number parsing
func TestInvalidNumberFormat(t *testing.T) {
	tmpDir := t.TempDir()

	powerCSV := filepath.Join(tmpDir, "power-15m.csv")
	content := `time,Import T1 kWh,Import T2 kWh,Export T1 kWh,Export T2 kWh,L1 max W,L2 max W,L3 max W
2025-06-02 20:30,invalid,7210.113,1916.077,4181.422,173,1212,67`

	if err := os.WriteFile(powerCSV, []byte(content), 0644); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	loader := &CSVLoader{DataDir: tmpDir}
	_, err := loader.readPowerCSV(powerCSV)
	if err == nil {
		t.Fatal("expected error for invalid number format, got nil")
	}
}

// TestInvalidTimeFormat tests error handling for invalid time parsing
func TestInvalidTimeFormat(t *testing.T) {
	tmpDir := t.TempDir()

	gasCSV := filepath.Join(tmpDir, "gas-15m.csv")
	content := `time,Total gas used
2025/06/02 20:30,3488.524`

	if err := os.WriteFile(gasCSV, []byte(content), 0644); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	loader := &CSVLoader{DataDir: tmpDir}
	_, err := loader.readGasCSV(gasCSV)
	if err == nil {
		t.Fatal("expected error for invalid time format, got nil")
	}
}

// TestEmptyCSV tests error handling for empty CSV files
func TestEmptyCSV(t *testing.T) {
	tmpDir := t.TempDir()

	powerCSV := filepath.Join(tmpDir, "power-15m.csv")
	if err := os.WriteFile(powerCSV, []byte(""), 0644); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	loader := &CSVLoader{DataDir: tmpDir}
	_, err := loader.readPowerCSV(powerCSV)
	if err == nil {
		t.Fatal("expected error for empty CSV, got nil")
	}
}

// TestMissingFiles tests error handling when files don't exist
func TestMissingFiles(t *testing.T) {
	tmpDir := t.TempDir()

	loader := &CSVLoader{DataDir: tmpDir}
	_, err := loader.LoadAndMerge()
	if err == nil {
		t.Fatal("expected error for missing files, got nil")
	}
}
