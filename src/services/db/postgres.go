package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/harrybawsac/p1-go/src/models"
)

type PostgresAdapter struct {
	DB *sql.DB
}

func (p *PostgresAdapter) InsertReading(ctx context.Context, r models.Reading) error {
	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	cols := []string{
		"created_at", "active_tariff",
		"total_power_import_kwh", "total_power_import_t1_kwh", "total_power_import_t2_kwh",
		"total_power_export_kwh", "total_power_export_t1_kwh", "total_power_export_t2_kwh",
		"active_power_w", "active_power_l1_w", "active_power_l2_w", "active_power_l3_w",
		"active_voltage_l1_v", "active_voltage_l2_v", "active_voltage_l3_v",
		"active_current_a", "active_current_l1_a", "active_current_l2_a", "active_current_l3_a",
		"voltage_sag_l1_count", "voltage_sag_l2_count", "voltage_sag_l3_count",
		"voltage_swell_l1_count", "voltage_swell_l2_count", "voltage_swell_l3_count",
		"any_power_fail_count", "long_power_fail_count", "total_gas_m3", "gas_timestamp",
	}

	// ensure CreatedAt is set (DB default is now() but we include explicit value for reproducibility)
	if r.CreatedAt.IsZero() {
		r.CreatedAt = time.Now().UTC()
	}

	args := []interface{}{
		r.CreatedAt, r.ActiveTariff,
		r.TotalPowerImportKwh, r.TotalPowerImportT1Kwh, r.TotalPowerImportT2Kwh,
		r.TotalPowerExportKwh, r.TotalPowerExportT1Kwh, r.TotalPowerExportT2Kwh,
		r.ActivePowerW, r.ActivePowerL1W, r.ActivePowerL2W, r.ActivePowerL3W,
		r.ActiveVoltageL1V, r.ActiveVoltageL2V, r.ActiveVoltageL3V,
		r.ActiveCurrentA, r.ActiveCurrentL1A, r.ActiveCurrentL2A, r.ActiveCurrentL3A,
		r.VoltageSagL1Count, r.VoltageSagL2Count, r.VoltageSagL3Count,
		r.VoltageSwellL1Count, r.VoltageSwellL2Count, r.VoltageSwellL3Count,
		r.AnyPowerFailCount, r.LongPowerFailCount, r.TotalGasM3, r.GasTimestamp,
	}

	// build placeholders
	ph := make([]string, len(args))
	for i := range ph {
		ph[i] = fmt.Sprintf("$%d", i+1)
	}

	// Insert a new meter_readings row and return its id
	insert := fmt.Sprintf("INSERT INTO p1.meter_readings (%s) VALUES (%s) RETURNING id",
		strings.Join(cols, ", "), strings.Join(ph, ","))

	var readingID int64
	if err := tx.QueryRowContext(ctx, insert, args...).Scan(&readingID); err != nil {
		tx.Rollback()
		return fmt.Errorf("insert reading: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

// InsertReadingsBatch inserts multiple readings in a single transaction
func (p *PostgresAdapter) InsertReadingsBatch(ctx context.Context, readings []models.Reading) error {
	if len(readings) == 0 {
		return nil
	}

	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	cols := []string{
		"created_at", "active_tariff",
		"total_power_import_kwh", "total_power_import_t1_kwh", "total_power_import_t2_kwh",
		"total_power_export_kwh", "total_power_export_t1_kwh", "total_power_export_t2_kwh",
		"active_power_w", "active_power_l1_w", "active_power_l2_w", "active_power_l3_w",
		"active_voltage_l1_v", "active_voltage_l2_v", "active_voltage_l3_v",
		"active_current_a", "active_current_l1_a", "active_current_l2_a", "active_current_l3_a",
		"voltage_sag_l1_count", "voltage_sag_l2_count", "voltage_sag_l3_count",
		"voltage_swell_l1_count", "voltage_swell_l2_count", "voltage_swell_l3_count",
		"any_power_fail_count", "long_power_fail_count", "total_gas_m3", "gas_timestamp",
	}

	// Build multi-row insert statement
	var valueStrings []string
	var valueArgs []interface{}

	for i, r := range readings {
		if r.CreatedAt.IsZero() {
			r.CreatedAt = time.Now().UTC()
		}

		// Create placeholder string for this row
		rowPlaceholders := make([]string, len(cols))
		for j := range rowPlaceholders {
			rowPlaceholders[j] = fmt.Sprintf("$%d", i*len(cols)+j+1)
		}
		valueStrings = append(valueStrings, fmt.Sprintf("(%s)", strings.Join(rowPlaceholders, ",")))

		// Add values for this row
		valueArgs = append(valueArgs,
			r.CreatedAt, r.ActiveTariff,
			r.TotalPowerImportKwh, r.TotalPowerImportT1Kwh, r.TotalPowerImportT2Kwh,
			r.TotalPowerExportKwh, r.TotalPowerExportT1Kwh, r.TotalPowerExportT2Kwh,
			r.ActivePowerW, r.ActivePowerL1W, r.ActivePowerL2W, r.ActivePowerL3W,
			r.ActiveVoltageL1V, r.ActiveVoltageL2V, r.ActiveVoltageL3V,
			r.ActiveCurrentA, r.ActiveCurrentL1A, r.ActiveCurrentL2A, r.ActiveCurrentL3A,
			r.VoltageSagL1Count, r.VoltageSagL2Count, r.VoltageSagL3Count,
			r.VoltageSwellL1Count, r.VoltageSwellL2Count, r.VoltageSwellL3Count,
			r.AnyPowerFailCount, r.LongPowerFailCount, r.TotalGasM3, r.GasTimestamp,
		)
	}

	insert := fmt.Sprintf("INSERT INTO p1.meter_readings (%s) VALUES %s",
		strings.Join(cols, ", "), strings.Join(valueStrings, ","))

	if _, err := tx.ExecContext(ctx, insert, valueArgs...); err != nil {
		tx.Rollback()
		return fmt.Errorf("insert batch: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

// GenerateInsertSQL generates a single multi-row INSERT SQL statement for all readings
func (p *PostgresAdapter) GenerateInsertSQL(readings []models.Reading) string {
	if len(readings) == 0 {
		return ""
	}

	cols := []string{
		"created_at", "active_tariff",
		"total_power_import_kwh", "total_power_import_t1_kwh", "total_power_import_t2_kwh",
		"total_power_export_kwh", "total_power_export_t1_kwh", "total_power_export_t2_kwh",
		"active_power_w", "active_power_l1_w", "active_power_l2_w", "active_power_l3_w",
		"active_voltage_l1_v", "active_voltage_l2_v", "active_voltage_l3_v",
		"active_current_a", "active_current_l1_a", "active_current_l2_a", "active_current_l3_a",
		"voltage_sag_l1_count", "voltage_sag_l2_count", "voltage_sag_l3_count",
		"voltage_swell_l1_count", "voltage_swell_l2_count", "voltage_swell_l3_count",
		"any_power_fail_count", "long_power_fail_count", "total_gas_m3", "gas_timestamp",
	}

	var valueRows []string

	for _, r := range readings {
		if r.CreatedAt.IsZero() {
			r.CreatedAt = time.Now().UTC()
		}

		values := []string{
			fmt.Sprintf("'%s'", r.CreatedAt.Format("2006-01-02 15:04:05-07:00")),
			fmt.Sprintf("%d", r.ActiveTariff),
			fmt.Sprintf("%f", r.TotalPowerImportKwh),
			fmt.Sprintf("%f", r.TotalPowerImportT1Kwh),
			fmt.Sprintf("%f", r.TotalPowerImportT2Kwh),
			fmt.Sprintf("%f", r.TotalPowerExportKwh),
			fmt.Sprintf("%f", r.TotalPowerExportT1Kwh),
			fmt.Sprintf("%f", r.TotalPowerExportT2Kwh),
			fmt.Sprintf("%f", r.ActivePowerW),
			fmt.Sprintf("%f", r.ActivePowerL1W),
			fmt.Sprintf("%f", r.ActivePowerL2W),
			fmt.Sprintf("%f", r.ActivePowerL3W),
			fmt.Sprintf("%f", r.ActiveVoltageL1V),
			fmt.Sprintf("%f", r.ActiveVoltageL2V),
			fmt.Sprintf("%f", r.ActiveVoltageL3V),
			fmt.Sprintf("%f", r.ActiveCurrentA),
			fmt.Sprintf("%f", r.ActiveCurrentL1A),
			fmt.Sprintf("%f", r.ActiveCurrentL2A),
			fmt.Sprintf("%f", r.ActiveCurrentL3A),
			fmt.Sprintf("%d", r.VoltageSagL1Count),
			fmt.Sprintf("%d", r.VoltageSagL2Count),
			fmt.Sprintf("%d", r.VoltageSagL3Count),
			fmt.Sprintf("%d", r.VoltageSwellL1Count),
			fmt.Sprintf("%d", r.VoltageSwellL2Count),
			fmt.Sprintf("%d", r.VoltageSwellL3Count),
			fmt.Sprintf("%d", r.AnyPowerFailCount),
			fmt.Sprintf("%d", r.LongPowerFailCount),
			fmt.Sprintf("%f", r.TotalGasM3),
			fmt.Sprintf("%d", r.GasTimestamp),
		}

		valueRows = append(valueRows, fmt.Sprintf("(%s)", strings.Join(values, ", ")))
	}

	stmt := fmt.Sprintf("INSERT INTO p1.meter_readings (%s) VALUES\n  %s;",
		strings.Join(cols, ", "), strings.Join(valueRows, ",\n  "))

	return stmt
}
