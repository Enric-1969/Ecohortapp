package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// buildPestanyaModePro retorna la vista i el callback de guardat
func (cfg *Config) buildPestanyaModePro(win fyne.Window) (fyne.CanvasObject, func()) {
	apiKeyEntry := widget.NewEntry()
	apiKeyEntry.SetPlaceHolder("Enganxa la teva AEMET API Key aquí...")
	apiKeyEntry.SetText(cfg.apiKey)

	lblResum := widget.NewLabel("")
	cfg.actualitzarLabelResum(lblResum)

	btnVeureDetall := widget.NewButtonWithIcon("Veure selecció", theme.InfoIcon(), func() {
		cfg.obrirDialogDetallSeleccio(win)
	})

	barraInferior := container.NewBorder(nil, nil, lblResum, btnVeureDetall)

	var treeContainer fyne.CanvasObject

	if len(cfg.Municipis) == 0 {
		treeContainer = container.NewCenter(
			widget.NewLabel("No hi ha municipis carregats en memòria.\nVés a la pestanya 'Per Municipi' i prem 'Carregar Municipis d'AEMET'."),
		)
	} else {
		tree, ccaaKeys := BuildCCAATree(cfg.Municipis)
		treeContainer = cfg.buildHierarchyTreeUI(tree, ccaaKeys, func() {
			cfg.actualitzarLabelResum(lblResum)
		})
	}

	scrollTree := container.NewVScroll(treeContainer)
	scrollTree.SetMinSize(fyne.NewSize(400, 250))

	saveFunc := func() {
		cfg.apiKey = apiKeyEntry.Text
		cfg.guardarUserConfigDisc()
	}

	topSection := container.NewVBox(
		widget.NewLabelWithStyle("AEMET API Key:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		apiKeyEntry,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Filtre Jeràrquic (Autonomies / Províncies / Municipis):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)

	return container.NewBorder(topSection, barraInferior, nil, nil, scrollTree), saveFunc
}

// buildHierarchyTreeUI genera les caselles de verificació i enllaça els esdeveniments
func (cfg *Config) buildHierarchyTreeUI(tree CCAATree, ccaaKeys []string, onSelectionChanged func()) fyne.CanvasObject {
	treeContainer := container.NewVBox()

	provMap := make(map[string][]string)
	for _, m := range cfg.Municipis {
		provMap[m.Provincia] = append(provMap[m.Provincia], m.CodigoINE())
	}

	selectedSet := make(map[string]bool)
	for _, codi := range cfg.UserConfig.Municipios {
		selectedSet[codi] = true
	}

	chkAll := widget.NewCheck("Seleccionar tota Espanya", nil)
	treeContainer.Add(chkAll)

	var allProvChecks []*widget.Check
	ccaaCheckboxes := make(map[string]*widget.Check)
	provCheckboxes := make(map[string][]*widget.Check)
	var isUpdating bool

	allProvsCheckedInitial := true

	for _, ccaa := range ccaaKeys {
		chkCCAA := widget.NewCheck(ccaa, nil)
		ccaaCheckboxes[ccaa] = chkCCAA
		treeContainer.Add(chkCCAA)

		allProvsInCCAAChecked := true

		for _, prov := range tree[ccaa] {
			chkProv := widget.NewCheck("   - "+prov+" (Tots els municipis)", nil)
			provSelected := cfg.isProvinciaSeleccionada(prov, provMap, selectedSet)
			if !provSelected {
				allProvsInCCAAChecked = false
				allProvsCheckedInitial = false
			}
			chkProv.SetChecked(provSelected)

			provCheckboxes[ccaa] = append(provCheckboxes[ccaa], chkProv)
			allProvChecks = append(allProvChecks, chkProv)
			treeContainer.Add(chkProv)

			currentProv := prov
			chkProv.OnChanged = func(val bool) {
				if isUpdating {
					return
				}
				cfg.actualitzarMunicipisProvincia(currentProv, val, provMap)

				isUpdating = true
				if !val {
					chkCCAA.SetChecked(false)
					chkAll.SetChecked(false)
				}
				isUpdating = false

				if onSelectionChanged != nil {
					onSelectionChanged()
				}
			}
		}

		chkCCAA.SetChecked(allProvsInCCAAChecked)

		currentCCAA := ccaa
		chkCCAA.OnChanged = func(val bool) {
			if isUpdating {
				return
			}
			isUpdating = true
			for _, chkP := range provCheckboxes[currentCCAA] {
				chkP.SetChecked(val)
			}
			for _, provName := range tree[currentCCAA] {
				cfg.actualitzarMunicipisProvincia(provName, val, provMap)
			}
			if !val {
				chkAll.SetChecked(false)
			}
			isUpdating = false

			if onSelectionChanged != nil {
				onSelectionChanged()
			}
		}
	}

	chkAll.SetChecked(allProvsCheckedInitial && len(ccaaKeys) > 0)

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
		for provName := range provMap {
			cfg.actualitzarMunicipisProvincia(provName, val, provMap)
		}
		isUpdating = false

		if onSelectionChanged != nil {
			onSelectionChanged()
		}
	}

	return treeContainer
}
