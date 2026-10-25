package ui

import (
	"ecohortapp/internal/domain"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// GetToolBar genera la barra superior de herramientas
func GetToolBar(app *domain.Config, win fyne.Window) fyne.CanvasObject {
	toolBar := widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			AddRegistresDialog(app)
		}),
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			// NOTA: Se conectará con el refresco global de la UI en pasos posteriores
		}),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			d := MostrarPreferencies(app, win)
			d.Resize(fyne.NewSize(350, 250))
			d.Show()
		}),
	)

	return container.NewHBox(layout.NewSpacer(), toolBar)
}
