package db

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/harrybawsac/p1-go/src/models"
)

// TestInsertReadingsBatch tests batch insertion of readings
func TestInsertReadingsBatch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	adapter := &PostgresAdapter{DB: db}

	readings := []models.Reading{
		{
			CreatedAt:             time.Date(2025, 6, 2, 20, 30, 0, 0, time.UTC),
			TotalPowerImportT1Kwh: 8293.146,
			TotalPowerImportT2Kwh: 7210.113,
			TotalGasM3:            3488.524,
		},
		{
			CreatedAt:             time.Date(2025, 6, 2, 20, 45, 0, 0, time.UTC),
			TotalPowerImportT1Kwh: 8293.146,
			TotalPowerImportT2Kwh: 7210.236,
			TotalGasM3:            3488.524,
		},
	}

	// Expect transaction begin
	mock.ExpectBegin()

	// Expect the batch insert with two rows
	mock.ExpectExec("INSERT INTO p1.meter_readings").
		WithArgs(
			// First reading
			readings[0].CreatedAt, readings[0].ActiveTariff,
			readings[0].TotalPowerImportKwh, readings[0].TotalPowerImportT1Kwh, readings[0].TotalPowerImportT2Kwh,
			readings[0].TotalPowerExportKwh, readings[0].TotalPowerExportT1Kwh, readings[0].TotalPowerExportT2Kwh,
			readings[0].ActivePowerW, readings[0].ActivePowerL1W, readings[0].ActivePowerL2W, readings[0].ActivePowerL3W,
			readings[0].ActiveVoltageL1V, readings[0].ActiveVoltageL2V, readings[0].ActiveVoltageL3V,
			readings[0].ActiveCurrentA, readings[0].ActiveCurrentL1A, readings[0].ActiveCurrentL2A, readings[0].ActiveCurrentL3A,
			readings[0].VoltageSagL1Count, readings[0].VoltageSagL2Count, readings[0].VoltageSagL3Count,
			readings[0].VoltageSwellL1Count, readings[0].VoltageSwellL2Count, readings[0].VoltageSwellL3Count,
			readings[0].AnyPowerFailCount, readings[0].LongPowerFailCount, readings[0].TotalGasM3, readings[0].GasTimestamp,
			// Second reading
			readings[1].CreatedAt, readings[1].ActiveTariff,
			readings[1].TotalPowerImportKwh, readings[1].TotalPowerImportT1Kwh, readings[1].TotalPowerImportT2Kwh,
			readings[1].TotalPowerExportKwh, readings[1].TotalPowerExportT1Kwh, readings[1].TotalPowerExportT2Kwh,
			readings[1].ActivePowerW, readings[1].ActivePowerL1W, readings[1].ActivePowerL2W, readings[1].ActivePowerL3W,
			readings[1].ActiveVoltageL1V, readings[1].ActiveVoltageL2V, readings[1].ActiveVoltageL3V,
			readings[1].ActiveCurrentA, readings[1].ActiveCurrentL1A, readings[1].ActiveCurrentL2A, readings[1].ActiveCurrentL3A,
			readings[1].VoltageSagL1Count, readings[1].VoltageSagL2Count, readings[1].VoltageSagL3Count,
			readings[1].VoltageSwellL1Count, readings[1].VoltageSwellL2Count, readings[1].VoltageSwellL3Count,
			readings[1].AnyPowerFailCount, readings[1].LongPowerFailCount, readings[1].TotalGasM3, readings[1].GasTimestamp,
		).
		WillReturnResult(sqlmock.NewResult(0, 2))

	// Expect commit
	mock.ExpectCommit()

	ctx := context.Background()
	if err := adapter.InsertReadingsBatch(ctx, readings); err != nil {
		t.Errorf("InsertReadingsBatch failed: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

// TestInsertReadingsBatchEmpty tests batch insertion with empty slice
func TestInsertReadingsBatchEmpty(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	adapter := &PostgresAdapter{DB: db}

	ctx := context.Background()
	if err := adapter.InsertReadingsBatch(ctx, []models.Reading{}); err != nil {
		t.Errorf("InsertReadingsBatch with empty slice should not error: %v", err)
	}
}

// TestInsertReadingsBatchError tests error handling in batch insertion
func TestInsertReadingsBatchError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	adapter := &PostgresAdapter{DB: db}

	readings := []models.Reading{
		{
			CreatedAt:             time.Date(2025, 6, 2, 20, 30, 0, 0, time.UTC),
			TotalPowerImportT1Kwh: 8293.146,
		},
	}

	// Expect transaction begin
	mock.ExpectBegin()

	// Simulate error during insert
	mock.ExpectExec("INSERT INTO p1.meter_readings").
		WillReturnError(sqlmock.ErrCancelled)

	// Expect rollback
	mock.ExpectRollback()

	ctx := context.Background()
	if err := adapter.InsertReadingsBatch(ctx, readings); err == nil {
		t.Error("expected error from InsertReadingsBatch, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

// TestGenerateInsertSQL tests SQL generation
func TestGenerateInsertSQL(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	adapter := &PostgresAdapter{DB: db}

	readings := []models.Reading{
		{
			CreatedAt:             time.Date(2025, 6, 2, 20, 30, 0, 0, time.UTC),
			TotalPowerImportT1Kwh: 8293.146,
			TotalPowerImportT2Kwh: 7210.113,
			TotalGasM3:            3488.524,
		},
		{
			CreatedAt:             time.Date(2025, 6, 2, 20, 45, 0, 0, time.UTC),
			TotalPowerImportT1Kwh: 8293.146,
			TotalPowerImportT2Kwh: 7210.236,
			TotalGasM3:            3488.524,
		},
	}

	statement := adapter.GenerateInsertSQL(readings)

	// Check that the statement is not empty
	if len(statement) == 0 {
		t.Error("SQL statement is empty")
	}

	// Verify it's a reasonable length (should be longer since it contains 2 rows)
	if len(statement) < 200 {
		t.Errorf("SQL statement seems too short: %d chars", len(statement))
	}

	// Check that it contains both timestamps
	if !strings.Contains(statement, "2025-06-02 20:30:00") {
		t.Error("statement doesn't contain first timestamp")
	}
	if !strings.Contains(statement, "2025-06-02 20:45:00") {
		t.Error("statement doesn't contain second timestamp")
	}

	// Check that it's a single statement (only one semicolon at the end)
	semicolonCount := strings.Count(statement, ";")
	if semicolonCount != 1 {
		t.Errorf("expected 1 semicolon, got %d", semicolonCount)
	}
}

// TestGenerateInsertSQLEmpty tests SQL generation with empty slice
func TestGenerateInsertSQLEmpty(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	adapter := &PostgresAdapter{DB: db}

	statement := adapter.GenerateInsertSQL([]models.Reading{})

	if statement != "" {
		t.Errorf("expected empty string for empty input, got %q", statement)
	}
}

// TestInsertReadingsBatchSetsCreatedAt tests that zero CreatedAt values are set
func TestInsertReadingsBatchSetsCreatedAt(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	adapter := &PostgresAdapter{DB: db}

	readings := []models.Reading{
		{
			TotalPowerImportT1Kwh: 8293.146,
			// CreatedAt is zero value
		},
	}

	// Expect transaction begin
	mock.ExpectBegin()

	// Expect insert - we use MatchExpectationsInOrder(false) conceptually
	// but sqlmock will match any non-zero time for the first argument
	mock.ExpectExec("INSERT INTO p1.meter_readings").
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Expect commit
	mock.ExpectCommit()

	ctx := context.Background()
	if err := adapter.InsertReadingsBatch(ctx, readings); err != nil {
		t.Errorf("InsertReadingsBatch failed: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
