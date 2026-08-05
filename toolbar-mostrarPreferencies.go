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
	vistaModePro, guardarPreferenciesPro := cfg.buildPestanyaModePro()

	pestanyes := container.NewAppTabs(
		container.NewTabItem("Per Municipi", vistaMunicipi),
		container.NewTabItem("Mode PRO", vistaModePro),
	)

	d := dialog.NewCustomConfirm("Preferències", "Guardar", "Cancel·lar", pestanyes, func(guardar bool) {
		if guardar {
			if guardarPreferenciesPro != nil {
				guardarPreferenciesPro()
			}
		}
	}, win)

	return d
}
