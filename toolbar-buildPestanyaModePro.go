package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// =============================================================================
// FUNCIÓ AUXILIAR: PESTANYA 2 (MODE PRO)
// =============================================================================

func (app *Config) buildPestanyaModePro() (fyne.CanvasObject, func()) {
	var checkTodaEspana, checkCatalunya, checkBarcelona, checkGirona, checkMadridCCAA, checkMadridProv *widget.Check
	var actualitzantCascada bool = false

	checkBarcelona = widget.NewCheck("    - Barcelona (Tots els municipis)", func(marcat bool) {
		if actualitzantCascada {
			return
		}
		actualitzantCascada = true
		if !marcat {
			checkCatalunya.SetChecked(false)
			checkTodaEspana.SetChecked(false)
		}
		actualitzantCascada = false
	})

	checkGirona = widget.NewCheck("    - Girona (Tots els municipis)", func(marcat bool) {
		if actualitzantCascada {
			return
		}
		actualitzantCascada = true
		if !marcat {
			checkCatalunya.SetChecked(false)
			checkTodaEspana.SetChecked(false)
		}
		actualitzantCascada = false
	})

	checkCatalunya = widget.NewCheck("  Catalunya", func(marcat bool) {
		if actualitzantCascada {
			return
		}
		actualitzantCascada = true
		checkBarcelona.SetChecked(marcat)
		checkGirona.SetChecked(marcat)
		if !marcat {
			checkTodaEspana.SetChecked(false)
		}
		actualitzantCascada = false
	})

	checkMadridProv = widget.NewCheck("    - Madrid (Tots els municipis)", func(marcat bool) {
		if actualitzantCascada {
			return
		}
		actualitzantCascada = true
		if !marcat {
			checkMadridCCAA.SetChecked(false)
			checkTodaEspana.SetChecked(false)
		}
		actualitzantCascada = false
	})

	checkMadridCCAA = widget.NewCheck("  Comunidad de Madrid", func(marcat bool) {
		if actualitzantCascada {
			return
		}
		actualitzantCascada = true
		checkMadridProv.SetChecked(marcat)
		if !marcat {
			checkTodaEspana.SetChecked(false)
		}
		actualitzantCascada = false
	})

	checkTodaEspana = widget.NewCheck("Seleccionar tota Espanya", func(marcat bool) {
		if actualitzantCascada {
			return
		}
		actualitzantCascada = true
		checkCatalunya.SetChecked(marcat)
		checkBarcelona.SetChecked(marcat)
		checkGirona.SetChecked(marcat)
		checkMadridCCAA.SetChecked(marcat)
		checkMadridProv.SetChecked(marcat)
		actualitzantCascada = false
	})

	actualitzantCascada = true
	checkTodaEspana.SetChecked(app.App.Preferences().BoolWithFallback("pro_tota_espanya", false))
	checkCatalunya.SetChecked(app.App.Preferences().BoolWithFallback("pro_catalunya", false))
	checkBarcelona.SetChecked(app.App.Preferences().BoolWithFallback("pro_barcelona", false))
	checkGirona.SetChecked(app.App.Preferences().BoolWithFallback("pro_girona", false))
	checkMadridCCAA.SetChecked(app.App.Preferences().BoolWithFallback("pro_madrid_ccaa", false))
	checkMadridProv.SetChecked(app.App.Preferences().BoolWithFallback("pro_madrid_prov", false))
	actualitzantCascada = false

	arbreContengut := container.NewVBox(
		checkTodaEspana,
		widget.NewSeparator(),
		checkCatalunya,
		checkBarcelona,
		checkGirona,
		widget.NewSeparator(),
		checkMadridCCAA,
		checkMadridProv,
	)

	scrollArbre := container.NewVScroll(arbreContengut)
	scrollArbre.SetMinSize(fyne.NewSize(320, 140))

	vista := container.NewVBox(
		widget.NewLabel("Filtre Jeràrquic (Autonomies / Províncies / Municipis):"),
		scrollArbre,
	)

	guardarPreferenciesPro := func() {
		app.App.Preferences().SetBool("pro_tota_espanya", checkTodaEspana.Checked)
		app.App.Preferences().SetBool("pro_catalunya", checkCatalunya.Checked)
		app.App.Preferences().SetBool("pro_barcelona", checkBarcelona.Checked)
		app.App.Preferences().SetBool("pro_girona", checkGirona.Checked)
		app.App.Preferences().SetBool("pro_madrid_ccaa", checkMadridCCAA.Checked)
		app.App.Preferences().SetBool("pro_madrid_prov", checkMadridProv.Checked)
	}

	return vista, guardarPreferenciesPro
}
