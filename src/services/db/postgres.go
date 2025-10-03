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

func (p *PostgresAdapter) InsertReading(ctx context.Context, r models.Reading, externals []models.ExternalReading) error {
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
		"unique_id", "created_at", "wifi_ssid", "wifi_strength", "smr_version", "meter_model", "active_tariff",
		"total_power_import_kwh", "total_power_import_t1_kwh", "total_power_import_t2_kwh",
		"total_power_export_kwh", "total_power_export_t1_kwh", "total_power_export_t2_kwh",
		"active_power_w", "active_power_l1_w", "active_power_l2_w", "active_power_l3_w",
		"active_voltage_l1_v", "active_voltage_l2_v", "active_voltage_l3_v",
		"active_current_a", "active_current_l1_a", "active_current_l2_a", "active_current_l3_a",
		"voltage_sag_l1_count", "voltage_sag_l2_count", "voltage_sag_l3_count",
		"voltage_swell_l1_count", "voltage_swell_l2_count", "voltage_swell_l3_count",
		"any_power_fail_count", "long_power_fail_count", "total_gas_m3", "gas_timestamp", "gas_unique_id",
	}

	// ensure CreatedAt is set (DB default is now() but we include explicit value for reproducibility)
	if r.CreatedAt.IsZero() {
		r.CreatedAt = time.Now().UTC()
	}

	args := []interface{}{
		r.UniqueID, r.CreatedAt, r.WifiSSID, r.WifiStrength, r.SmrVersion, r.MeterModel, r.ActiveTariff,
		r.TotalPowerImportKwh, r.TotalPowerImportT1Kwh, r.TotalPowerImportT2Kwh,
		r.TotalPowerExportKwh, r.TotalPowerExportT1Kwh, r.TotalPowerExportT2Kwh,
		r.ActivePowerW, r.ActivePowerL1W, r.ActivePowerL2W, r.ActivePowerL3W,
		r.ActiveVoltageL1V, r.ActiveVoltageL2V, r.ActiveVoltageL3V,
		r.ActiveCurrentA, r.ActiveCurrentL1A, r.ActiveCurrentL2A, r.ActiveCurrentL3A,
		r.VoltageSagL1Count, r.VoltageSagL2Count, r.VoltageSagL3Count,
		r.VoltageSwellL1Count, r.VoltageSwellL2Count, r.VoltageSwellL3Count,
		r.AnyPowerFailCount, r.LongPowerFailCount, r.TotalGasM3, r.GasTimestamp, r.GasUniqueID,
	}

	// build placeholders
	ph := make([]string, len(args))
	for i := range ph {
		ph[i] = fmt.Sprintf("$%d", i+1)
	}

	// build update assignments for ON CONFLICT — do not overwrite created_at on update
	updates := make([]string, 0, len(cols)-2)
	for _, c := range cols {
		if c == "unique_id" || c == "created_at" {
			continue
		}
		updates = append(updates, fmt.Sprintf("%s = EXCLUDED.%s", c, c))
	}

	upsert := fmt.Sprintf("INSERT INTO p1.meter_readings (%s) VALUES (%s) ON CONFLICT (unique_id) DO UPDATE SET %s",
		strings.Join(cols, ", "), strings.Join(ph, ","), strings.Join(updates, ", "))

	_, err = tx.ExecContext(ctx, upsert, args...)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("upsert reading: %w", err)
	}

	// replace external readings for this meter_unique_id
	if len(externals) > 0 {
		if _, err := tx.ExecContext(ctx, "DELETE FROM p1.external_readings WHERE meter_reading_unique_id = $1", r.UniqueID); err != nil {
			tx.Rollback()
			return fmt.Errorf("delete externals: %w", err)
		}
		// created_at for external_readings will be set by the DB default (now())
		stmt := "INSERT INTO p1.external_readings (meter_reading_unique_id, unique_id, type, timestamp, value, unit) VALUES ($1,$2,$3,$4,$5,$6)"
		for _, e := range externals {
			if _, err := tx.ExecContext(ctx, stmt, r.UniqueID, e.UniqueID, e.Type, e.Timestamp, e.Value, e.Unit); err != nil {
				tx.Rollback()
				return fmt.Errorf("insert external: %w", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}
