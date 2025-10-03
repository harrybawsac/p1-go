package parser

import (
	"encoding/json"
	"errors"

	"github.com/harrybawsac/p1-go/src/models"
)

type MeterPayload struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

func ParseMeterJSON(data []byte) (MeterPayload, error) {
	var p MeterPayload
	if err := json.Unmarshal(data, &p); err != nil {
		return MeterPayload{}, err
	}
	if p.Unit == "" {
		return MeterPayload{}, errors.New("missing unit")
	}
	return p, nil
}

// ParseFullReading parses the complete meter JSON payload into a Reading and external readings
func ParseFullReading(data []byte) (models.Reading, []models.ExternalReading, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return models.Reading{}, nil, err
	}

	r := models.Reading{}
	if v, ok := raw["unique_id"].(string); ok {
		r.UniqueID = v
	}
	if v, ok := raw["wifi_ssid"].(string); ok {
		r.WifiSSID = v
	}
	if v, ok := raw["wifi_strength"].(float64); ok {
		r.WifiStrength = int(v)
	}
	if v, ok := raw["smr_version"].(float64); ok {
		r.SmrVersion = int(v)
	}
	if v, ok := raw["meter_model"].(string); ok {
		r.MeterModel = v
	}
	if v, ok := raw["active_tariff"].(float64); ok {
		r.ActiveTariff = int(v)
	}
	// numeric fields
	mapFloat := func(key string) float64 {
		if vv, ok := raw[key].(float64); ok {
			return vv
		}
		return 0
	}
	r.TotalPowerImportKwh = mapFloat("total_power_import_kwh")
	r.TotalPowerImportT1Kwh = mapFloat("total_power_import_t1_kwh")
	r.TotalPowerImportT2Kwh = mapFloat("total_power_import_t2_kwh")
	r.TotalPowerExportKwh = mapFloat("total_power_export_kwh")
	r.TotalPowerExportT1Kwh = mapFloat("total_power_export_t1_kwh")
	r.TotalPowerExportT2Kwh = mapFloat("total_power_export_t2_kwh")
	r.ActivePowerW = mapFloat("active_power_w")
	r.ActivePowerL1W = mapFloat("active_power_l1_w")
	r.ActivePowerL2W = mapFloat("active_power_l2_w")
	r.ActivePowerL3W = mapFloat("active_power_l3_w")
	r.ActiveVoltageL1V = mapFloat("active_voltage_l1_v")
	r.ActiveVoltageL2V = mapFloat("active_voltage_l2_v")
	r.ActiveVoltageL3V = mapFloat("active_voltage_l3_v")
	r.ActiveCurrentA = mapFloat("active_current_a")
	r.ActiveCurrentL1A = mapFloat("active_current_l1_a")
	r.ActiveCurrentL2A = mapFloat("active_current_l2_a")
	r.ActiveCurrentL3A = mapFloat("active_current_l3_a")
	r.VoltageSwellL1Count = int(mapFloat("voltage_swell_l1_count"))
	r.VoltageSwellL2Count = int(mapFloat("voltage_swell_l2_count"))
	r.VoltageSwellL3Count = int(mapFloat("voltage_swell_l3_count"))
	r.VoltageSagL1Count = int(mapFloat("voltage_sag_l1_count"))
	r.VoltageSagL2Count = int(mapFloat("voltage_sag_l2_count"))
	r.VoltageSagL3Count = int(mapFloat("voltage_sag_l3_count"))
	r.AnyPowerFailCount = int(mapFloat("any_power_fail_count"))
	r.LongPowerFailCount = int(mapFloat("long_power_fail_count"))
	r.TotalGasM3 = mapFloat("total_gas_m3")
	if v, ok := raw["gas_timestamp"].(float64); ok {
		r.GasTimestamp = int64(v)
	}
	if v, ok := raw["gas_unique_id"].(string); ok {
		r.GasUniqueID = v
	}

	var externals []models.ExternalReading
	if arr, ok := raw["external"].([]interface{}); ok {
		for _, it := range arr {
			if m, ok := it.(map[string]interface{}); ok {
				er := models.ExternalReading{}
				if v, ok := m["unique_id"].(string); ok {
					er.UniqueID = v
				}
				if v, ok := m["type"].(string); ok {
					er.Type = v
				}
				if v, ok := m["timestamp"].(float64); ok {
					er.Timestamp = int64(v)
				}
				er.Value = 0
				if v, ok := m["value"].(float64); ok {
					er.Value = v
				}
				if v, ok := m["unit"].(string); ok {
					er.Unit = v
				}
				externals = append(externals, er)
			}
		}
	}

	return r, externals, nil
}
