package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

func (app *Config) makeUI() {
	// Obtenir les dades de l'API (Precipitacions, Temp. Max/Min i Humitat)
	precipitacio, tempMax, tempMin, humitat := app.getClimaText()

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

	precipitacio, tempMax, tempMin, humitat := app.getClimaText()

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
