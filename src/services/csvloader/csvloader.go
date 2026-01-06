package csvloader

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/harrybawsac/p1-go/src/models"
)

// CSVLoader reads and merges power and gas CSV files
type CSVLoader struct {
	DataDir string
}

// MergedReading represents a merged row from power and gas CSV files
type MergedReading struct {
	Time        time.Time
	ImportT1Kwh float64
	ImportT2Kwh float64
	ExportT1Kwh float64
	ExportT2Kwh float64
	L1MaxW      float64
	L2MaxW      float64
	L3MaxW      float64
	TotalGasM3  float64
}

// LoadAndMerge reads both CSV files and merges them in memory
func (l *CSVLoader) LoadAndMerge() ([]MergedReading, error) {
	powerFile := filepath.Join(l.DataDir, "power-15m.csv")
	gasFile := filepath.Join(l.DataDir, "gas-15m.csv")

	powerRecords, err := l.readPowerCSV(powerFile)
	if err != nil {
		return nil, fmt.Errorf("read power CSV: %w", err)
	}

	gasRecords, err := l.readGasCSV(gasFile)
	if err != nil {
		return nil, fmt.Errorf("read gas CSV: %w", err)
	}

	if len(powerRecords) != len(gasRecords) {
		return nil, fmt.Errorf("power and gas CSV files have different number of records: %d vs %d", len(powerRecords), len(gasRecords))
	}

	merged := make([]MergedReading, len(powerRecords))
	for i := 0; i < len(powerRecords); i++ {
		// Verify timestamps match
		if !powerRecords[i].Time.Equal(gasRecords[i].Time) {
			return nil, fmt.Errorf("timestamp mismatch at row %d: power=%s, gas=%s",
				i+1, powerRecords[i].Time, gasRecords[i].Time)
		}

		merged[i] = MergedReading{
			Time:        powerRecords[i].Time,
			ImportT1Kwh: powerRecords[i].ImportT1Kwh,
			ImportT2Kwh: powerRecords[i].ImportT2Kwh,
			ExportT1Kwh: powerRecords[i].ExportT1Kwh,
			ExportT2Kwh: powerRecords[i].ExportT2Kwh,
			L1MaxW:      powerRecords[i].L1MaxW,
			L2MaxW:      powerRecords[i].L2MaxW,
			L3MaxW:      powerRecords[i].L3MaxW,
			TotalGasM3:  gasRecords[i].TotalGasM3,
		}
	}

	return merged, nil
}

// powerRecord represents a single row from power-15m.csv
type powerRecord struct {
	Time        time.Time
	ImportT1Kwh float64
	ImportT2Kwh float64
	ExportT1Kwh float64
	ExportT2Kwh float64
	L1MaxW      float64
	L2MaxW      float64
	L3MaxW      float64
}

// gasRecord represents a single row from gas-15m.csv
type gasRecord struct {
	Time       time.Time
	TotalGasM3 float64
}

func (l *CSVLoader) readPowerCSV(path string) ([]powerRecord, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("power CSV is empty")
	}

	// Skip header
	rows = rows[1:]
	records := make([]powerRecord, 0, len(rows))

	for i, row := range rows {
		if len(row) != 8 {
			return nil, fmt.Errorf("power CSV row %d has %d columns, expected 8", i+2, len(row))
		}

		t, err := time.Parse("2006-01-02 15:04", row[0])
		if err != nil {
			return nil, fmt.Errorf("parse time at row %d: %w", i+2, err)
		}

		importT1, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			return nil, fmt.Errorf("parse ImportT1 at row %d: %w", i+2, err)
		}

		importT2, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("parse ImportT2 at row %d: %w", i+2, err)
		}

		exportT1, err := strconv.ParseFloat(row[3], 64)
		if err != nil {
			return nil, fmt.Errorf("parse ExportT1 at row %d: %w", i+2, err)
		}

		exportT2, err := strconv.ParseFloat(row[4], 64)
		if err != nil {
			return nil, fmt.Errorf("parse ExportT2 at row %d: %w", i+2, err)
		}

		l1Max, err := strconv.ParseFloat(row[5], 64)
		if err != nil {
			return nil, fmt.Errorf("parse L1Max at row %d: %w", i+2, err)
		}

		l2Max, err := strconv.ParseFloat(row[6], 64)
		if err != nil {
			return nil, fmt.Errorf("parse L2Max at row %d: %w", i+2, err)
		}

		l3Max, err := strconv.ParseFloat(row[7], 64)
		if err != nil {
			return nil, fmt.Errorf("parse L3Max at row %d: %w", i+2, err)
		}

		records = append(records, powerRecord{
			Time:        t,
			ImportT1Kwh: importT1,
			ImportT2Kwh: importT2,
			ExportT1Kwh: exportT1,
			ExportT2Kwh: exportT2,
			L1MaxW:      l1Max,
			L2MaxW:      l2Max,
			L3MaxW:      l3Max,
		})
	}

	return records, nil
}

func (l *CSVLoader) readGasCSV(path string) ([]gasRecord, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("gas CSV is empty")
	}

	// Skip header
	rows = rows[1:]
	records := make([]gasRecord, 0, len(rows))

	for i, row := range rows {
		if len(row) != 2 {
			return nil, fmt.Errorf("gas CSV row %d has %d columns, expected 2", i+2, len(row))
		}

		t, err := time.Parse("2006-01-02 15:04", row[0])
		if err != nil {
			return nil, fmt.Errorf("parse time at row %d: %w", i+2, err)
		}

		totalGas, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			return nil, fmt.Errorf("parse TotalGas at row %d: %w", i+2, err)
		}

		records = append(records, gasRecord{
			Time:       t,
			TotalGasM3: totalGas,
		})
	}

	return records, nil
}

// ToReading converts a MergedReading to models.Reading
func (m *MergedReading) ToReading() models.Reading {
	return models.Reading{
		CreatedAt:             m.Time,
		TotalPowerImportT1Kwh: m.ImportT1Kwh,
		TotalPowerImportT2Kwh: m.ImportT2Kwh,
		TotalPowerExportT1Kwh: m.ExportT1Kwh,
		TotalPowerExportT2Kwh: m.ExportT2Kwh,
		// Max values from CSV are in watts, store them as active power
		ActivePowerL1W: m.L1MaxW,
		ActivePowerL2W: m.L2MaxW,
		ActivePowerL3W: m.L3MaxW,
		TotalGasM3:     m.TotalGasM3,
	}
}

// GroupByDay groups merged readings by day
func GroupByDay(readings []MergedReading) map[string][]MergedReading {
	dayMap := make(map[string][]MergedReading)

	for _, r := range readings {
		day := r.Time.Format("2006-01-02")
		dayMap[day] = append(dayMap[day], r)
	}

	return dayMap
}
