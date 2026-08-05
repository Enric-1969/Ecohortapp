package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func (app *Config) getToolBar(win fyne.Window) fyne.CanvasObject {
	toolBar := widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			app.addRegistresDialog()
		}),
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			app.actualitzarClimaDadesContent()
		}),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			d := app.mostrarPreferencies(win)
			d.Resize(fyne.NewSize(350, 250))
			d.Show()
		}),
	)

	// El spacer empuja la barra de herramientas totalmente a la derecha
	return container.NewHBox(layout.NewSpacer(), toolBar)
}
