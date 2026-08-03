package main

import (
	"ecohortapp/repository"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Usamos '_' para indicar que el parámetro win no se usa explícitamente dentro de la función
func (app *Config) getToolBar(_ fyne.Window) *widget.Toolbar {
	toolBar := widget.NewToolbar(
		widget.NewToolbarSpacer(), //Crearem un espaciador que empenyi els diferents items cap a la dreta
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			app.addRegistresDialog()
		}), //Crearem una nova acció indicant quina icona i quina funcio estaràn involucrades
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			app.actualitzarClimaDadesContent()
		}),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			//Realitzem una crida a la funció que conté el dialeg per modificar les preferencies de l'app
			w := app.mostrarPreferencies()
			//Definim el tamany per aquesta nova finestra
			w.Resize(fyne.NewSize(300, 200))
			//Mostrem la finestra
			w.Show()

		}),
	)

	return toolBar
}

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

	// Retornem la vista UI i una funció getter per al codi final
	return vista, func() string {
		return dadaMunicipi.Text
	}
}

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

	// Càrrega d'estat inicial
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

	// Retornem la vista UI i una funció de desat per al callback de Guardar
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

// =============================================================================
// FUNCIÓ PRINCIPAL / DIÀLEG D'AJUSTAMENTS
// =============================================================================
func (app *Config) mostrarPreferencies() dialog.Dialog {
	vistaMunicipi, getCodiMunicipi := app.buildPestanyaMunicipi()
	vistaModePro, guardarPreferenciesPro := app.buildPestanyaModePro()

	pestanyes := container.NewAppTabs(
		container.NewTabItem("Per Municipi", vistaMunicipi),
		container.NewTabItem("Mode PRO", vistaModePro),
	)

	return dialog.NewCustomConfirm(
		"Configurar ajustaments",
		"Guardar",
		"Cancelar",
		pestanyes,
		func(valid bool) {
			if valid {
				municipi = getCodiMunicipi()
				app.App.Preferences().SetString("municipi", municipi)
				guardarPreferenciesPro()
				app.actualitzarClimaDadesContent()
			}
		},
		app.MainWindow,
	)
}

// Funció per afegir Registres a on referenciem el struct Config
func (app *Config) addRegistresDialog() dialog.Dialog {
	//Definim les variables a on guardarem el resultat del mètode d'entrada
	dataRegistreEntrada := widget.NewEntry()
	precipitacioEntrada := widget.NewEntry()
	tempMaximaEntrada := widget.NewEntry()
	tempMinimaEntrada := widget.NewEntry()
	humitatEntrada := widget.NewEntry()

	app.AfegirRegistresDataRegistreEntrada = dataRegistreEntrada
	app.AfegirRegistresPrecipitacioEntrada = precipitacioEntrada
	app.AfegirRegistresTempMaximaEntrada = tempMaximaEntrada
	// CORREGIDO: antes decía TempMaximaEntrada
	app.AfegirRegistresTempMinimaEntrada = tempMinimaEntrada
	app.AfegirRegistresHumitatEntrada = humitatEntrada
	app.AfegirRegistresHumitatEntrada = humitatEntrada

	validacioData := func(s string) error {
		//Apliquem el format de la data definit en el primer parametre al string indicat en el segon parametre i que correspon el aportat per l'usuari
		//S’ha de referenciar el format amb el layout standard 2006-01-02 i no pas amb cualsevol data
		if _, err := time.Parse("2006-01-02", s); err != nil {
			//Si es produeix algun error en el procés el retornem
			return err
		}
		return nil
	}
	dataRegistreEntrada.Validator = validacioData

	esIntValidador := func(s string) error {
		_, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		return nil
	}
	precipitacioEntrada.Validator = esIntValidador
	tempMaximaEntrada.Validator = esIntValidador
	tempMinimaEntrada.Validator = esIntValidador
	humitatEntrada.Validator = esIntValidador

	//Definim un placeholder i aixi facilitar l'usabilitat en el camp data de compra
	dataRegistreEntrada.PlaceHolder = "YYYY-MM-DD"

	//Crearem el dialeg creant un formulari
	addForm := dialog.NewForm(
		"Afegir Registre",
		"Afegir",
		"Cancelar",
		//Afegirem les etiquetes en forma de item per el formulari
		[]*widget.FormItem{
			{Text: "Data Registre", Widget: dataRegistreEntrada},
			{Text: "Probabilitat de precipitació", Widget: precipitacioEntrada},
			{Text: "Temperatura màxima", Widget: tempMaximaEntrada},
			{Text: "Temperatura minima", Widget: tempMinimaEntrada},
			{Text: "Humitat", Widget: humitatEntrada},
		},
		//A continuació realitzem la validació de les dades
		func(valid bool) {
			if valid {
				//Desenvolupar un filtratge, convertint les dades a els formats de la bd
				//S’ha de referenciar el format amb el layout standard 2006-01-02 i no pas amb cualsevol data
				dataRegistre, _ := time.Parse("2006-01-02", dataRegistreEntrada.Text)
				precipitacio, _ := strconv.Atoi(precipitacioEntrada.Text)
				tempMaxima, _ := strconv.Atoi(tempMaximaEntrada.Text)
				tempMinima, _ := strconv.Atoi(tempMinimaEntrada.Text)
				humitat, _ := strconv.Atoi(humitatEntrada.Text)

				//Invoquem el mètode de la base de dades per insertar registres i que poblarem amb les dades formatejades
				_, err := app.DB.InsertRegistre(repository.Registres{
					Data:         dataRegistre,
					Precipitacio: precipitacio,
					TempMaxima:   tempMaxima,
					TempMinima:   tempMinima,
					Humitat:      humitat,
				})
				//Controlem si és produeix algun error
				if err != nil {
					app.ErrorLog.Println(err)
				}
				//Invoquem el paremetre del struct Config per permetre que refresqui el widget de la taula amb el nou registre
				app.actualitzarRegistresTable()
			}
		},
		app.MainWindow)

	//Establim el tamany de la finestra i mostrem el dialeg
	addForm.Resize(fyne.Size{Width: 400})
	addForm.Show()

	return addForm
}
