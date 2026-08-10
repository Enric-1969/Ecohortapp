package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
)

// ============================================================================
// FUNCIÓ PRINCIPAL / DIÀLEG D'AJUSTAMENTS
// ============================================================================

func (cfg *Config) mostrarPreferencies(win fyne.Window) dialog.Dialog {
	vistaMunicipi := cfg.buildPestanyaMunicipi(win)
	vistaModePro, guardarPreferenciesPro := cfg.buildPestanyaModePro(win)
	pestanyes := container.NewAppTabs(
		container.NewTabItem("Per Municipi", vistaMunicipi),
		container.NewTabItem("Mode PRO", vistaModePro),
	)

	d := dialog.NewCustomConfirm("Preferències", "Guardar", "Cancel·lar", pestanyes, func(guardar bool) {
		if guardar {
			if guardarPreferenciesPro != nil {
				guardarPreferenciesPro()
			}

			// Opció A: Si l'usuari guarda des de la pestanya "Mode PRO", agafem el primer
			// municipi de la selecció (cfg.UserConfig.Municipios) per actualitzar la vista.
			if pestanyes.Selected() != nil && pestanyes.Selected().Text == "Mode PRO" {
				if len(cfg.UserConfig.Municipios) > 0 {
					primerMunicipi := cfg.UserConfig.Municipios[0]
					cfg.municipi = primerMunicipi
					cfg.UserConfig.MunicipioCodigo = primerMunicipi
				}
			}

			// 1. Guardar la configuració al fitxer config.json
			_ = SaveConfig(cfg.UserConfig)

			// 2. Refrescar la pantalla principal amb el nou municipi seleccionat
			go cfg.actualitzarClimaDadesContent()
		}
	}, win)

	return d
}
