package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

func (app *Config) makeUI() {
	// Obtenir les dades de l'API (Precipitacions, Temp. Max/Min i Humitat) i el possible error
	precipitacio, tempMax, tempMin, humitat, err := app.getClimaText()
	if err != nil && app.ErrorLog != nil {
		app.ErrorLog.Printf("[AEMET WARNING] Error inicial en carregar clima: %v", err)
	}

	// Protecció en cas que fos nil en el primer arrencada sense connexió
	if precipitacio == nil {
		precipitacio = canvas.NewText("Precipitació: --", theme.ForegroundColor())
		tempMax = canvas.NewText("Temp. Max: --", theme.ForegroundColor())
		tempMin = canvas.NewText("Temp. Min: --", theme.ForegroundColor())
		humitat = canvas.NewText("Humitat: --", theme.ForegroundColor())
	}

	climaDadesContent := container.NewGridWithColumns(4,
		precipitacio,
		tempMax,
		tempMin,
		humitat,
	)

	app.ClimaDadesContainer = climaDadesContent

	// Obtenir la barra d'eines
	toolBar := app.getToolBar(app.MainWindow)

	pronosticTabContent := app.pronosticTab()
	registresTabContent := app.registresTab()

	// Obtenir les pestanyes de l'aplicació
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Pronòstic", theme.HomeIcon(), pronosticTabContent),
		container.NewTabItemWithIcon("Diari Meteorològic", theme.InfoIcon(), registresTabContent),
	)

	tabs.SetTabLocation(container.TabLocationTop)

	// Agrupem el clima i la toolbar a la part superior
	topContainer := container.NewVBox(climaDadesContent, toolBar)

	// Utilitzem NewBorder per a que les pestanyes ocupin tot l'espai central
	finalContent := container.NewBorder(topContainer, nil, nil, nil, tabs)

	app.MainWindow.SetContent(finalContent)

	// Goroutine en segon pla per actualitzar cada 15 minuts
	go func() {
		ticker := time.NewTicker(15 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			app.actualitzarClimaDadesContent()
		}
	}()
}

func (app *Config) actualitzarClimaDadesContent() {
	if app.InfoLog != nil {
		app.InfoLog.Print("actualitzar les dades meteorològiques")
	}

	// Capturem les etiquetes I l'error de la petició
	precipitacio, tempMax, tempMin, humitat, err := app.getClimaText()

	// SI HI HA UN ERROR (ex: HTTP 429 de límit d'AEMET), NO toquem el contenidor visual
	if err != nil {
		if app.ErrorLog != nil {
			app.ErrorLog.Printf("[AEMET WARNING] Error en actualitzar clima (%v). Mantenint dades anteriors.", err)
		}
		return // Sortim de la funció sense modificar ClimaDadesContainer
	}

	// Només si NO hi ha error, actualitzem la interfície amb les noves dades
	if app.ClimaDadesContainer != nil {
		app.ClimaDadesContainer.Objects = []fyne.CanvasObject{precipitacio, tempMax, tempMin, humitat}
		app.ClimaDadesContainer.Refresh()
	}

	if app.PronosticGraficContainer != nil {
		grafic := app.obtenirGrafic()
		app.PronosticGraficContainer.Objects = []fyne.CanvasObject{grafic}
		app.PronosticGraficContainer.Refresh()
	}
}

func (app *Config) actualitzarRegistresTable() {
	app.Registres = app.getRegistresSlice()
	if app.RegistresTable != nil {
		app.RegistresTable.Refresh()
	}
}
