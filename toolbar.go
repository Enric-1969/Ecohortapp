package main

import (
	"fyne.io/fyne/v2"

	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// getToolBar construeix la barra d'eines superior
func (app *Config) getToolBar(_ fyne.Window) *widget.Toolbar {
	toolBar := widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			app.addRegistresDialog()
		}),
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			app.actualitzarClimaDadesContent()
		}),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			w := app.mostrarPreferencies()
			w.Resize(fyne.NewSize(350, 250))
			w.Show()
		}),
	)

	return toolBar
}
