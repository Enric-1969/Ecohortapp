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

	// Desplegables per CCAA i Municipi
	selectCCAA := widget.NewSelect([]string{}, nil)
	selectCCAA.PlaceHolder = "Selecciona CCAA..."

	selectMunicipi := widget.NewSelect([]string{}, nil)
	selectMunicipi.PlaceHolder = "Selecciona Municipi..."
	selectMunicipi.Disable() // Desactivat fins que es triï CCAA

	// Mapa intern per relacionar el nom del municipi amb el seu objecte/codi
	municipisPerNom := make(map[string]Municipio)
	var totsMunicipis []Municipio

	// Botó per carregar/actualitzar la llista des de l'API d'AEMET
	btnCarregar := widget.NewButton("Carregar Municipis d'AEMET", func() {
		progress := dialog.NewCustom("AEMET", "Cancel·lar", widget.NewLabel("Descarregant llista de municipis..."), win)
		progress.Show()

		go func() {
			var err error
			totsMunicipis, err = app.ObtenirMunicipiosAEMET()
			progress.Hide()

			if err != nil {
				dialog.ShowError(fmt.Errorf("error descarregant municipis: %v", err), win)
				return
			}

			// Extreure llista única de CCAA i ordenar-la
			mapaCCAAUniques := make(map[string]bool)
			for _, m := range totsMunicipis {
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
			dialog.ShowInformation("Èxit", fmt.Sprintf("S'han carregat %d municipis.", len(totsMunicipis)), win)
		}()
	})

	// Quan l'usuari tria una CCAA, filtrem els municipis d'aquesta regió
	selectCCAA.OnChanged = func(ccaaSeleccionada string) {
		selectMunicipi.ClearSelected()
		municipisPerNom = make(map[string]Municipio)
		var opcionsMunicipis []string

		for _, m := range totsMunicipis {
			if m.CCAA == ccaaSeleccionada {
				nomClau := fmt.Sprintf("%s (%s)", m.Nombre, m.IDOld)
				opcionsMunicipis = append(opcionsMunicipis, nomClau)
				municipisPerNom[nomClau] = m
			}
		}

		sort.Strings(opcionsMunicipis)
		selectMunicipi.Options = opcionsMunicipis
		selectMunicipi.Enable()
		selectMunicipi.Refresh()
	}

	// Quan l'usuari tria un municipi, guardem el codi INE a la configuració
	selectMunicipi.OnChanged = func(municipiSeleccionat string) {
		if m, ok := municipisPerNom[municipiSeleccionat]; ok {
			app.UserConfig.MunicipioCodigo = m.IDOld
			app.municipi = m.IDOld
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
