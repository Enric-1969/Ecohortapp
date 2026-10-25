package ui

import (
	"ecohortapp/internal/domain"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
)

// MostrarPreferencies gestiona las ventanas modales de pestañas de ajustes
func MostrarPreferencies(cfg *domain.Config, win fyne.Window) dialog.Dialog {
	vistaMunicipi := BuildPestanyaMunicipi(cfg, win)
	vistaModePro, guardarPreferenciesPro := BuildPestanyaModePro(cfg, win)

	pestanyes := container.NewAppTabs(
		container.NewTabItem("Per Municipi", vistaMunicipi),
		container.NewTabItem("Mode PRO", vistaModePro),
	)

	d := dialog.NewCustomConfirm("Preferències", "Guardar", "Cancel·lar", pestanyes, func(guardar bool) {
		if guardar {
			if guardarPreferenciesPro != nil {
				guardarPreferenciesPro()
			}

			if pestanyes.Selected() != nil && pestanyes.Selected().Text == "Mode PRO" {
				if len(cfg.UserConfig.Municipios) > 0 {
					primerMunicipi := cfg.UserConfig.Municipios[0]
					cfg.Municipi = primerMunicipi
					cfg.UserConfig.MunicipioCodigo = primerMunicipi
				}
			}
		}
	}, win)

	return d
}
