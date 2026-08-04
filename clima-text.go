package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

// Funció que retorna quatre elements de text amb Fyne amb les dades climatològiques
func (app *Config) getClimaText() (*canvas.Text, *canvas.Text, *canvas.Text, *canvas.Text) {
	// 1. Obtenir el codi de municipi des de la memòria de l'aplicació
	codiMunicipi := app.municipi
	if codiMunicipi == "" && app.UserConfig.MunicipioCodigo != "" {
		codiMunicipi = app.UserConfig.MunicipioCodigo
	}
	if codiMunicipi == "" {
		codiMunicipi = "08001" // Valor per defecte (Abrera)
	}

	// 2. Assignem el codi a l'estructura Diaria
	var g Diaria
	g.CodiIne = codiMunicipi

	var precipitacio, tempMax, tempMin, humitat *canvas.Text

	prediccio, err := g.GetPrediccions()

	if err != nil {
		// Definim text en gris en cas d'error
		gris := color.NRGBA{R: 155, G: 155, B: 155, A: 255}
		precipitacio = canvas.NewText("Precipitació: No Definit", gris)
		tempMax = canvas.NewText("Temp. Max: No Definit", gris)
		tempMin = canvas.NewText("Temp. Min: No Definit", gris)
		humitat = canvas.NewText("Humitat: No Definit", gris)
	} else {
		// Color per a la precipitació
		colorPrecipitacio := color.NRGBA{R: 0, G: 180, B: 0, A: 255}
		if prediccio.ProbPrecipitacio < 50 {
			colorPrecipitacio = color.NRGBA{R: 180, G: 0, B: 0, A: 255}
		}

		// Utilitzem el color de text del tema actual (evita passar nil)
		textColor := theme.ForegroundColor()

		precipitacioTxt := fmt.Sprintf("Precipitació: %d%%", prediccio.ProbPrecipitacio)
		tempMaxTxt := fmt.Sprintf("Temp. Max: %d°C", prediccio.TemperaturaMax)
		tempMinTxt := fmt.Sprintf("Temp. Min: %d°C", prediccio.TemperaturaMin)
		humitatTxt := fmt.Sprintf("Humitat: %d%%", prediccio.HumitatRelativa)

		precipitacio = canvas.NewText(precipitacioTxt, colorPrecipitacio)
		tempMax = canvas.NewText(tempMaxTxt, textColor)
		tempMin = canvas.NewText(tempMinTxt, textColor)
		humitat = canvas.NewText(humitatTxt, textColor)
	}

	// Alineació dels textos en la interfície
	precipitacio.Alignment = fyne.TextAlignLeading
	tempMax.Alignment = fyne.TextAlignCenter
	tempMin.Alignment = fyne.TextAlignCenter
	humitat.Alignment = fyne.TextAlignTrailing

	return precipitacio, tempMax, tempMin, humitat
}
