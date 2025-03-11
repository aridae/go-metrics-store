package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type jsonconf struct {
	CryptoKey                    *string `json:"crypto_key"`
	AddressOverride              *string `json:"address"`
	StoreIntervalSecondsOverride *int64  `json:"store_interval"`
	FileStoragePathOverride      *string `json:"store_file"`
	RestoreOverride              *bool   `json:"restore"`
	DatabaseDsnOverride          *string `json:"database_dsn"`
	TrustedSubnetOverride        *string `json:"trusted_subnet"`
}

func parseJSONFile(path string) (*jsonconf, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config json file: %w", err)
	}
	defer file.Close()

	cnf := jsonconf{}
	if err = json.NewDecoder(file).Decode(&cnf); err != nil {
		return nil, fmt.Errorf("failed to parse json file: %w", err)
	}

	return &cnf, nil
}

func (f jsonconf) override(cfg *Config) {
	if f.AddressOverride != nil {
		cfg.overrideAddressIfNotDefault(*f.AddressOverride, "json")
	}

	if f.StoreIntervalSecondsOverride != nil {
		storeInterval := time.Duration(*f.StoreIntervalSecondsOverride) * time.Second
		cfg.overrideStoreIntervalIfNotDefault(storeInterval, "json")
	}

	if f.FileStoragePathOverride != nil {
		cfg.overrideFileStoragePathIfNotDefault(*f.FileStoragePathOverride, "json")
	}

	if f.RestoreOverride != nil {
		cfg.overrideRestoreIfNotDefault(*f.RestoreOverride, "json")
	}

	if f.DatabaseDsnOverride != nil {
		cfg.overrideDatabaseDNSIfNotDefault(*f.DatabaseDsnOverride, "json")
	}

	if f.CryptoKey != nil {
		cfg.overrideCryptoKeyIfNotDefault(*f.CryptoKey, "json")
	}
}
