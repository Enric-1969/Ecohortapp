package main

import (
	"encoding/json"
	"os"
)

// UserConfig define las preferencias del usuario para config.json
type UserConfig struct {
	MunicipioCodigo string   `json:"municipio_codigo"`
	Municipios      []string `json:"municipios"`
}

const configFileName = "config.json"

// LoadConfig lee las preferencias guardadas desde config.json
func LoadConfig() (UserConfig, error) {
	var cfg UserConfig
	data, err := os.ReadFile(configFileName)
	if err != nil {
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
