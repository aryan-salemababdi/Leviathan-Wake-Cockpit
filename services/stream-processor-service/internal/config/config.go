package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	KeydbAddress           string `json:"keydb_address"`
	ArbitrumWsUrl          string `json:"arbitrum_ws_url"`
	ExecutionServerAddress int    `json:"execution_server_address"`
}

func Load(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = json.Unmarshal(file, &cfg)
	return &cfg, err
}
