package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// =============================================================================
// FUNCIÓ AUXILIAR: PESTANYA 1 (PER MUNICIPI)
// =============================================================================

func (app *Config) buildPestanyaMunicipi() (fyne.CanvasObject, func() string) {
	codigosANombres := map[string]string{
		"08001": "Abrera",
		"08019": "Barcelona",
		"08121": "Martorell",
	}

	nombresACodigos := map[string]string{
		"Abrera":    "08001",
		"Barcelona": "08019",
		"Martorell": "08121",
	}

	dadaMunicipi := widget.NewEntry()
	dadaMunicipi.Text = municipi

	opcionsNoms := []string{"Abrera", "Barcelona", "Martorell"}
	dadaNom := widget.NewSelect(opcionsNoms, nil)

	if nomInicial, existe := codigosANombres[dadaMunicipi.Text]; existe {
		dadaNom.SetSelected(nomInicial)
	}

	var sincronitzant bool = false

	dadaNom.OnChanged = func(nomSeleccionat string) {
		if sincronitzant {
			return
		}
		sincronitzant = true
		if codi, existe := nombresACodigos[nomSeleccionat]; existe {
			dadaMunicipi.SetText(codi)
		}
		sincronitzant = false
	}

	dadaMunicipi.OnChanged = func(codiEscrit string) {
		if sincronitzant {
			return
		}
		sincronitzant = true
		if nom, existe := codigosANombres[codiEscrit]; existe {
			dadaNom.SetSelected(nom)
		}
		sincronitzant = false
	}

	vista := container.NewVBox(
		widget.NewLabel("Selecciona el municipi per Nom o Codi INE/AEMET:"),
		widget.NewForm(
			widget.NewFormItem("Nom Municipi:", dadaNom),
			widget.NewFormItem("Codi Municipi:", dadaMunicipi),
		),
	)

	return vista, func() string {
		return dadaMunicipi.Text
	}
}
