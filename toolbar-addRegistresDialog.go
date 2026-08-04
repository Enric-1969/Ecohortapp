package main

import (
	"ecohortapp/repository"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// addRegistresDialog mostra el diàleg per afegir registres
func (app *Config) addRegistresDialog() dialog.Dialog {
	dataRegistreEntrada := widget.NewEntry()
	precipitacioEntrada := widget.NewEntry()
	tempMaximaEntrada := widget.NewEntry()
	tempMinimaEntrada := widget.NewEntry()
	humitatEntrada := widget.NewEntry()

	app.AfegirRegistresDataRegistreEntrada = dataRegistreEntrada
	app.AfegirRegistresPrecipitacioEntrada = precipitacioEntrada
	app.AfegirRegistresTempMaximaEntrada = tempMaximaEntrada
	app.AfegirRegistresTempMinimaEntrada = tempMinimaEntrada
	app.AfegirRegistresHumitatEntrada = humitatEntrada

	validacioData := func(s string) error {
		if _, err := time.Parse("2006-01-02", s); err != nil {
			return err
		}
		return nil
	}
	dataRegistreEntrada.Validator = validacioData

	esIntValidador := func(s string) error {
		_, err := strconv.Atoi(s)
		return err
	}
	precipitacioEntrada.Validator = esIntValidador
	tempMaximaEntrada.Validator = esIntValidador
	tempMinimaEntrada.Validator = esIntValidador
	humitatEntrada.Validator = esIntValidador

	dataRegistreEntrada.PlaceHolder = "YYYY-MM-DD"

	addForm := dialog.NewForm(
		"Afegir Registre",
		"Afegir",
		"Cancelar",
		[]*widget.FormItem{
			{Text: "Data Registre", Widget: dataRegistreEntrada},
			{Text: "Probabilitat de precipitació", Widget: precipitacioEntrada},
			{Text: "Temperatura màxima", Widget: tempMaximaEntrada},
			{Text: "Temperatura mínima", Widget: tempMinimaEntrada},
			{Text: "Humitat", Widget: humitatEntrada},
		},
		func(valid bool) {
			if valid {
				dataRegistre, _ := time.Parse("2006-01-02", dataRegistreEntrada.Text)
				precipitacio, _ := strconv.Atoi(precipitacioEntrada.Text)
				tempMaxima, _ := strconv.Atoi(tempMaximaEntrada.Text)
				tempMinima, _ := strconv.Atoi(tempMinimaEntrada.Text)
				humitat, _ := strconv.Atoi(humitatEntrada.Text)

				if app.DB != nil {
					_, err := app.DB.InsertRegistre(repository.Registres{
						Data:         dataRegistre,
						Precipitacio: precipitacio,
						TempMaxima:   tempMaxima,
						TempMinima:   tempMinima,
						Humitat:      humitat,
					})
					if err != nil && app.ErrorLog != nil {
						app.ErrorLog.Println(err)
					}
				}
				app.actualitzarRegistresTable()
			}
		},
		app.MainWindow,
	)

	addForm.Resize(fyne.NewSize(400, 350))
	addForm.Show()

	return addForm
}
