package main

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
)

// =============================================================================
// FUNCIÓ PRINCIPAL / DIÀLEG D'AJUSTAMENTS
// =============================================================================

func (app *Config) mostrarPreferencies() dialog.Dialog {
	vistaMunicipi, getCodiMunicipi := app.buildPestanyaMunicipi()
	vistaModePro, guardarPreferenciesPro := app.buildPestanyaModePro()

	pestanyes := container.NewAppTabs(
		container.NewTabItem("Per Municipi", vistaMunicipi),
		container.NewTabItem("Mode PRO", vistaModePro),
	)

	return dialog.NewCustomConfirm(
		"Configurar ajustaments",
		"Guardar",
		"Cancelar",
		pestanyes,
		func(valid bool) {
			if valid {
				municipi = getCodiMunicipi()
				app.municipi = municipi
				app.App.Preferences().SetString("municipi", municipi)
				guardarPreferenciesPro()
				app.actualitzarClimaDadesContent()
			}
		},
		app.MainWindow,
	)
}
