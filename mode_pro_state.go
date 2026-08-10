package main

import (
	"encoding/json"
	"os"
)

// actualitzarMunicipisProvincia afegeix o elimina els codis INE d'una província en cfg.UserConfig.Municipios
func (cfg *Config) actualitzarMunicipisProvincia(prov string, afegir bool, provMap map[string][]string) {
	codes := provMap[prov]
	if len(codes) == 0 {
		return
	}

	set := make(map[string]bool)
	for _, c := range cfg.UserConfig.Municipios {
		set[c] = true
	}

	if afegir {
		for _, c := range codes {
			set[c] = true
		}
	} else {
		for _, c := range codes {
			delete(set, c)
		}
	}

	nousMunicipis := make([]string, 0, len(set))
	for c := range set {
		nousMunicipis = append(nousMunicipis, c)
	}
	cfg.UserConfig.Municipios = nousMunicipis
}

// isProvinciaSeleccionada comprova si tots els municipis d'una província estan en UserConfig.Municipios
func (cfg *Config) isProvinciaSeleccionada(prov string, provMap map[string][]string, selectedSet map[string]bool) bool {
	codes := provMap[prov]
	if len(codes) == 0 {
		return false
	}
	for _, c := range codes {
		if !selectedSet[c] {
			return false
		}
	}
	return true
}

// guardarUserConfigDisc desa la configuració a config.json
func (cfg *Config) guardarUserConfigDisc() {
	data, err := json.MarshalIndent(cfg.UserConfig, "", "  ")
	if err != nil {
		if cfg.ErrorLog != nil {
			cfg.ErrorLog.Println("Error serialitzant UserConfig:", err)
		}
		return
	}
	_ = os.WriteFile("config.json", data, 0644)
}
