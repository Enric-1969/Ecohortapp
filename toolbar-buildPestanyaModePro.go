package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// buildPestanyaModePro retorna la vista i el callback de guardat per a toolbar-mostrarPreferencies.go
func (cfg *Config) buildPestanyaModePro() (fyne.CanvasObject, func()) {
	apiKeyEntry := widget.NewEntry()
	apiKeyEntry.SetPlaceHolder("Enganxa la teva AEMET API Key aquí...")
	apiKeyEntry.SetText(cfg.apiKey)

	var treeContainer fyne.CanvasObject

	// 1. Verificar si hi ha municipis en la memòria RAM (cfg.Municipis)
	if len(cfg.Municipis) == 0 {
		treeContainer = container.NewCenter(
			widget.NewLabel("No hi ha municipis carregats en memòria.\nVés a la pestanya 'Per Municipi' i prem 'Carregar Municipis d'AEMET'."),
		)
	} else {
		// Generar l'arbre utilitzant els municipis existents en RAM
		tree, ccaaKeys := BuildCCAATree(cfg.Municipis)
		treeContainer = cfg.buildHierarchyTreeUI(tree, ccaaKeys)
	}

	scrollTree := container.NewVScroll(treeContainer)
	scrollTree.SetMinSize(fyne.NewSize(400, 250))

	// 2. Funció de guardat devuelta a toolbar-mostrarPreferencies.go
	saveFunc := func() {
		cfg.apiKey = apiKeyEntry.Text
	}

	view := container.NewVBox(
		widget.NewLabelWithStyle("AEMET API Key:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		apiKeyEntry,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Filtre Jeràrquic (Autonomies / Províncies / Municipis):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		scrollTree,
	)

	return view, saveFunc
}

// buildHierarchyTreeUI genera les caselles de verificació a la interfície
func (cfg *Config) buildHierarchyTreeUI(tree CCAATree, ccaaKeys []string) fyne.CanvasObject {
	treeContainer := container.NewVBox()
	chkAll := widget.NewCheck("Seleccionar tota Espanya", nil)
	treeContainer.Add(chkAll)

	var allProvChecks []*widget.Check
	ccaaCheckboxes := make(map[string]*widget.Check)
	provCheckboxes := make(map[string][]*widget.Check)

	// Bandera de control (guard) per evitar bucles d'esdeveniments en cascada a Fyne
	var isUpdating bool

	for _, ccaa := range ccaaKeys {
		chkCCAA := widget.NewCheck(ccaa, nil)
		ccaaCheckboxes[ccaa] = chkCCAA
		treeContainer.Add(chkCCAA)

		for _, prov := range tree[ccaa] {
			chkProv := widget.NewCheck("   - "+prov+" (Tots els municipis)", nil)

			chkProv.OnChanged = func(val bool) {
				if isUpdating {
					return
				}
				if !val {
					isUpdating = true
					chkAll.SetChecked(false)
					isUpdating = false
				}
			}

			provCheckboxes[ccaa] = append(provCheckboxes[ccaa], chkProv)
			allProvChecks = append(allProvChecks, chkProv)
			treeContainer.Add(chkProv)
		}

		currentCCAA := ccaa
		chkCCAA.OnChanged = func(val bool) {
			if isUpdating {
				return
			}
			isUpdating = true
			for _, chkP := range provCheckboxes[currentCCAA] {
				chkP.SetChecked(val)
			}
			if !val {
				chkAll.SetChecked(false)
			}
			isUpdating = false
		}
	}

	chkAll.OnChanged = func(val bool) {
		if isUpdating {
			return
		}
		isUpdating = true
		for _, chkC := range ccaaCheckboxes {
			chkC.SetChecked(val)
		}
		for _, chkP := range allProvChecks {
			chkP.SetChecked(val)
		}
		isUpdating = false
	}

	return treeContainer
}
