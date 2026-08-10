package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// obtenirResumSeleccioModePro calcula el total de municipis i les CCAA implicades
func (app *Config) obtenirResumSeleccioModePro() (int, int, map[string][]string) {
	mapaCCAA := make(map[string][]string)

	// Mapa per cercar municipis per codi INE en O(1)
	municipisMap := make(map[string]Municipio)
	for _, m := range app.Municipis {
		municipisMap[m.CodigoINE()] = m
	}

	for _, codi := range app.UserConfig.Municipios {
		if m, ok := municipisMap[codi]; ok {
			nomFormatat := fmt.Sprintf("%s (%s)", m.Nombre, codi)
			mapaCCAA[m.CCAA] = append(mapaCCAA[m.CCAA], nomFormatat)
		} else {
			mapaCCAA["Altres/Desconegut"] = append(mapaCCAA["Altres/Desconegut"], codi)
		}
	}

	return len(app.UserConfig.Municipios), len(mapaCCAA), mapaCCAA
}

// actualitzarLabelResum actualitza el text de la barra inferior
func (app *Config) actualitzarLabelResum(lbl *widget.Label) {
	totalM, totalC, _ := app.obtenirResumSeleccioModePro()
	if totalM == 0 {
		lbl.SetText("📊 Cap municipi seleccionat.")
	} else {
		lbl.SetText(fmt.Sprintf("📊 Seleccionats: %d municipis en %d CCAA", totalM, totalC))
	}
}

// obrirDialogDetallSeleccio mostra la finestra emergent amb el desglose per CCAA
func (app *Config) obrirDialogDetallSeleccio(win fyne.Window) {
	totalM, totalC, mapaCCAA := app.obtenirResumSeleccioModePro()

	if totalM == 0 {
		dialog.ShowInformation("Resum de la Selecció", "No hi ha cap municipi seleccionat en el Mode PRO.", win)
		return
	}

	content := container.NewVBox()
	header := widget.NewLabelWithStyle(
		fmt.Sprintf("Total seleccionat: %d municipis a %d Comunitats Autònomes", totalM, totalC),
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)
	content.Add(header)
	content.Add(widget.NewSeparator())

	for ccaa, llista := range mapaCCAA {
		titolCCAA := widget.NewLabelWithStyle(
			fmt.Sprintf("• %s (%d municipis)", ccaa, len(llista)),
			fyne.TextAlignLeading,
			fyne.TextStyle{Bold: true},
		)
		content.Add(titolCCAA)

		mostraMax := 8
		for i, mun := range llista {
			if i >= mostraMax {
				mésText := fmt.Sprintf("   ... i %d municipis més.", len(llista)-mostraMax)
				content.Add(widget.NewLabelWithStyle(mésText, fyne.TextAlignLeading, fyne.TextStyle{Italic: true}))
				break
			}
			content.Add(widget.NewLabel(fmt.Sprintf("   - %s", mun)))
		}
		content.Add(widget.NewSeparator())
	}

	scroll := container.NewScroll(content)
	scroll.SetMinSize(fyne.NewSize(450, 320))

	dialog.ShowCustom("Resum Detallat (Mode PRO)", "Tancar", scroll, win)
}
