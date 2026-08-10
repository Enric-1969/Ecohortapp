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

// buildHierarchyTreeUI genera les caselles de verificació i enllaça els esdeveniments
func (cfg *Config) buildHierarchyTreeUI(tree CCAATree, ccaaKeys []string) fyne.CanvasObject {
	treeContainer := container.NewVBox()
	chkAll := widget.NewCheck("Seleccionar tota Espanya", nil)
	treeContainer.Add(chkAll)

	var allProvChecks []*widget.Check
	ccaaCheckboxes := make(map[string]*widget.Check)
	provCheckboxes := make(map[string][]*widget.Check)

	// Bandera de control (guard) per evitar bucles d'esdeveniments en cascada
	var isUpdating bool

	// 1. Construcció de la interfície i associació d'esdeveniments
	for _, ccaa := range ccaaKeys {
		chkCCAA := widget.NewCheck(ccaa, nil)
		ccaaCheckboxes[ccaa] = chkCCAA
		treeContainer.Add(chkCCAA)

		for _, prov := range tree[ccaa] {
			chkProv := widget.NewCheck("   - "+prov+" (Tots els municipis)", nil)
			setupProvinciaCheck(chkProv, chkAll, &isUpdating)

			provCheckboxes[ccaa] = append(provCheckboxes[ccaa], chkProv)
			allProvChecks = append(allProvChecks, chkProv)
			treeContainer.Add(chkProv)
		}

		setupCCAACheck(chkCCAA, provCheckboxes[ccaa], chkAll, &isUpdating)
	}

	// 2. Esdeveniment global per a la casella de tot el país
	setupNacionalCheck(chkAll, ccaaCheckboxes, allProvChecks, &isUpdating)

	return treeContainer
}

// --- FUNCIONS AUXILIARS DE GESTIÓ D'ESDEVENIMENTS (CLEAN CODE) ---

func setupProvinciaCheck(chkProv *widget.Check, chkAll *widget.Check, isUpdating *bool) {
	chkProv.OnChanged = func(val bool) {
		if *isUpdating {
			return
		}
		if !val {
			*isUpdating = true
			chkAll.SetChecked(false)
			*isUpdating = false
		}
	}
}

func setupCCAACheck(chkCCAA *widget.Check, provs []*widget.Check, chkAll *widget.Check, isUpdating *bool) {
	chkCCAA.OnChanged = func(val bool) {
		if *isUpdating {
			return
		}
		*isUpdating = true
		for _, chkP := range provs {
			chkP.SetChecked(val)
		}
		if !val {
			chkAll.SetChecked(false)
		}
		*isUpdating = false
	}
}

func setupNacionalCheck(chkAll *widget.Check, ccaaChecks map[string]*widget.Check, allProvs []*widget.Check, isUpdating *bool) {
	chkAll.OnChanged = func(val bool) {
		if *isUpdating {
			return
		}
		*isUpdating = true
		for _, chkC := range ccaaChecks {
			chkC.SetChecked(val)
		}
		for _, chkP := range allProvs {
			chkP.SetChecked(val)
		}
		*isUpdating = false
	}
}
