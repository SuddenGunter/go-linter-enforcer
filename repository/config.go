package repository

import (
	"encoding/json"
	"os"
)

type config struct {
	Repositories []Repository `json:"repositories"`
}

func LoadListFromJSON(cfgFileName string) ([]Repository, error) {
	f, err := os.Open(cfgFileName)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(f)

	var cfg config

	if err = dec.Decode(&cfg); err != nil {
		return nil, err
	}

	return cfg.Repositories, nil
}
