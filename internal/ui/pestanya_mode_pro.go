package ui

import (
	"ecohortapp/internal/domain"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// BuildPestanyaModePro retorna la vista del panel PRO y su callback de guardado
func BuildPestanyaModePro(cfg *domain.Config, win fyne.Window) (fyne.CanvasObject, func()) {
	apiKeyEntry := widget.NewEntry()
	apiKeyEntry.SetPlaceHolder("Enganxa la teva AEMET API Key aquí...")
	apiKeyEntry.SetText(cfg.ApiKey)

	lblResum := widget.NewLabel("")

	btnVeureDetall := widget.NewButtonWithIcon("Veure selecció", theme.InfoIcon(), func() {
		// Callback para diálogos extendidos
	})

	barraInferior := container.NewBorder(nil, nil, lblResum, btnVeureDetall)

	var treeContainer fyne.CanvasObject

	if len(cfg.Municipis) == 0 {
		treeContainer = container.NewCenter(
			widget.NewLabel("No hi ha municipis carregats en memòria.\nVés a la pestanya 'Per Municipi' i prem 'Carregar Municipis d'AEMET'."),
		)
	} else {
		// Contenedor estructural provisional hasta migrar los árboles jerárquicos
		treeContainer = container.NewVBox(widget.NewLabel("Arbre de municipis actiu"))
	}

	scrollTree := container.NewVScroll(treeContainer)
	scrollTree.SetMinSize(fyne.NewSize(400, 250))

	saveFunc := func() {
		cfg.ApiKey = apiKeyEntry.Text
	}

	topSection := container.NewVBox(
		widget.NewLabelWithStyle("AEMET API Key:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		apiKeyEntry,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Filtre Jeràrquic (Autonomies / Províncies / Municipis):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)

	return container.NewBorder(topSection, barraInferior, nil, nil, scrollTree), saveFunc
}
