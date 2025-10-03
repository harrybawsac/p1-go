package parser

import (
	"encoding/json"
	"errors"
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
