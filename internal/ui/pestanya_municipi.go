package ui

import (
	"fmt"
	"sort"

	"ecohortapp/internal/domain"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// BuildPestanyaMunicipi construye la interfaz de selección geográfica básica
func BuildPestanyaMunicipi(app *domain.Config, win fyne.Window) fyne.CanvasObject {
	labelInfo := widget.NewLabel("Selecciona la teva Comunitat Autònoma i Municipi:")

	selectCCAA := widget.NewSelect([]string{}, nil)
	selectCCAA.PlaceHolder = "Selecciona CCAA..."

	selectMunicipi := widget.NewSelect([]string{}, nil)
	selectMunicipi.PlaceHolder = "Selecciona Municipi..."
	selectMunicipi.Disable()

	municipisPerNom := make(map[string]domain.Municipio)

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
		municipisPerNom = make(map[string]domain.Municipio)
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
			app.Municipi = codiINE
		}
	}

	btnCarregar := widget.NewButton("Carregar Municipis d'AEMET", func() {
		progress := dialog.NewCustom("AEMET", "Cancel·lar", widget.NewLabel("Descarregant llista de municipis..."), win)
		progress.Show()
	})

	if len(app.Municipis) > 0 {
		poblarCCAA()
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
