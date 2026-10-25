package repository

import (
	"encoding/json"
	"os"

	"ecohortapp/internal/domain"
)

// ActualitzarMunicipisProvincia afegeix o elimina els codis INE d'una província en cfg.UserConfig.Municipios
func ActualitzarMunicipisProvincia(cfg *domain.Config, prov string, afegir bool, provMap map[string][]string) {
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

// IsProvinciaSeleccionada comprova si tots els municipis d'una província estan en UserConfig.Municipios
func IsProvinciaSeleccionada(cfg *domain.Config, prov string, provMap map[string][]string, selectedSet map[string]bool) bool {
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

// GuardarUserConfigDisc desa la configuració a config.json
func GuardarUserConfigDisc(cfg *domain.Config) {
	data, err := json.MarshalIndent(cfg.UserConfig, "", "  ")
	if err != nil {
		if cfg.ErrorLog != nil {
			cfg.ErrorLog.Println("Error serialitzant UserConfig:", err)
		}
		return
	}
	_ = os.WriteFile("config.json", data, 0644)
}
