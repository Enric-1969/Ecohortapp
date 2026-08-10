package main

import (
	"fmt"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (app *Config) buildPestanyaMunicipi(win fyne.Window) fyne.CanvasObject {
	labelInfo := widget.NewLabel("Selecciona la teva Comunitat Autònoma i Municipi:")

	selectCCAA := widget.NewSelect([]string{}, nil)
	selectCCAA.PlaceHolder = "Selecciona CCAA..."

	selectMunicipi := widget.NewSelect([]string{}, nil)
	selectMunicipi.PlaceHolder = "Selecciona Municipi..."
	selectMunicipi.Disable()

	municipisPerNom := make(map[string]Municipio)

	// Auxiliar per a poblar el desplegable de CCAA des de la memòria RAM
	poblarCCAA := func() {
		mapaCCAAUniques := make(map[string]bool)
		for _, m := range app.Municipis {
			if m.CCAA != "" {
				mapaCCAAUniques[m.CCAA] = true
			}
		}

		var llistaCCAA []string
		for ccaa := range mapaCCAAUniques {
			llistaCCAA = append(llistaCCAA, ccaa)
		}
		sort.Strings(llistaCCAA)

		selectCCAA.Options = llistaCCAA
		selectCCAA.Refresh()
	}

	selectCCAA.OnChanged = func(ccaaSeleccionada string) {
		if ccaaSeleccionada == "" {
			return
		}

		selectMunicipi.ClearSelected()
		municipisPerNom = make(map[string]Municipio)
		var opcionsMunicipis []string

		for _, m := range app.Municipis {
			if m.CCAA == ccaaSeleccionada {
				codiINE := m.CodigoINE()
				nomClau := fmt.Sprintf("%s (%s)", m.Nombre, codiINE)
				opcionsMunicipis = append(opcionsMunicipis, nomClau)
				municipisPerNom[nomClau] = m
			}
		}

		sort.Strings(opcionsMunicipis)
		selectMunicipi.Options = opcionsMunicipis
		selectMunicipi.Enable()
		selectMunicipi.Refresh()
	}

	selectMunicipi.OnChanged = func(municipiSeleccionat string) {
		if m, ok := municipisPerNom[municipiSeleccionat]; ok {
			codiINE := m.CodigoINE()
			app.UserConfig.MunicipioCodigo = codiINE
			app.UserConfig.MunicipioNombre = m.Nombre
			app.UserConfig.Latitud = parseAEMETCoord(m.Latitude)
			app.UserConfig.Longitud = parseAEMETCoord(m.Longitude)

			app.municipi = codiINE
		}
	}

	btnCarregar := widget.NewButton("Carregar Municipis d'AEMET", func() {
		progress := dialog.NewCustom("AEMET", "Cancel·lar", widget.NewLabel("Descarregant llista de municipis..."), win)
		progress.Show()

		go func() {
			var err error
			app.Municipis, err = app.ObtenirMunicipiosAEMET()
			progress.Hide()

			if err != nil {
				dialog.ShowError(fmt.Errorf("error descarregant municipis: %v", err), win)
				return
			}

			poblarCCAA()
			dialog.ShowInformation("Èxit", fmt.Sprintf("S'han carregat %d municipis.", len(app.Municipis)), win)
		}()
	})

	// SI JA HI HA MUNICIPIS EN RAM: Reutilitzar immediatament i seleccionar l'actual
	if len(app.Municipis) > 0 {
		poblarCCAA()

		// Seleccionar automàticament el municipi guardat actualment
		if app.UserConfig.MunicipioCodigo != "" {
			for _, m := range app.Municipis {
				if m.IDOld == app.UserConfig.MunicipioCodigo {
					selectCCAA.SetSelected(m.CCAA)
					nomClau := fmt.Sprintf("%s (%s)", m.Nombre, m.IDOld)
					selectMunicipi.SetSelected(nomClau)
					break
				}
			}
		}
	}

	return container.NewVBox(
		labelInfo,
		btnCarregar,
		widget.NewLabel("Comunitat Autònoma:"),
		selectCCAA,
		widget.NewLabel("Municipi:"),
		selectMunicipi,
	)
}
