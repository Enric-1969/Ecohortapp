package main

import (
	"encoding/json"
	"errors"
	"os"
)

const configFileName = "config.json"

// LoadConfig lee las preferencias guardadas desde config.json
func LoadConfig() (UserConfig, error) {
	var cfg UserConfig
	cfg.MunicipioCodigo = "08001" // Valor por defecto (Abrera)

	data, err := os.ReadFile(configFileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return cfg, err
	}

	err = json.Unmarshal(data, &cfg)
	return cfg, err
}

// SaveConfig guarda la configuración actual en config.json
func SaveConfig(cfg UserConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configFileName, data, 0644)
}
