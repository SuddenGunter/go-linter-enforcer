package repository

import (
	"encoding/json"
	"os"
)

type Config struct {
	Repositories []Repository `json:"repositories"`
}

func ConfigFromJSON(cfgFileName string) (*Config, error) {
	f, err := os.Open(cfgFileName)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(f)

	var cfg Config

	if err = dec.Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
