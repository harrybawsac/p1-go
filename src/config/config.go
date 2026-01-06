package config

import (
	"encoding/json"
	"os"
)

// Config holds runtime configuration for the CLI
type Config struct {
	MeterEndpoint string `json:"meter_endpoint"`
	DBDSN         string `json:"db_dsn"`
	DataDir       string `json:"data_dir"`
}

// Load reads a JSON config file from path and unmarshals into Config
func Load(path string) (Config, error) {
	var cfg Config
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
