package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

// Funció que retorna quatre elements de text amb Fyne i el possible error de la petició
func (app *Config) getClimaText() (*canvas.Text, *canvas.Text, *canvas.Text, *canvas.Text, error) {
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

	prediccio, err := g.GetPrediccions()

	// Si AEMET falla (ej. HTTP 429), retornem l'error i cap element visual
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Color per a la precipitació
	colorPrecipitacio := color.NRGBA{R: 0, G: 180, B: 0, A: 255}
	if prediccio.ProbPrecipitacio < 50 {
		colorPrecipitacio = color.NRGBA{R: 180, G: 0, B: 0, A: 255}
	}

	// Color de text segons el tema actual de Fyne
	textColor := theme.ForegroundColor()

	precipitacioTxt := fmt.Sprintf("Precipitació: %d%%", prediccio.ProbPrecipitacio)
	tempMaxTxt := fmt.Sprintf("Temp. Max: %d°C", prediccio.TemperaturaMax)
	tempMinTxt := fmt.Sprintf("Temp. Min: %d°C", prediccio.TemperaturaMin)
	humitatTxt := fmt.Sprintf("Humitat: %d%%", prediccio.HumitatRelativa)

	precipitacio := canvas.NewText(precipitacioTxt, colorPrecipitacio)
	tempMax := canvas.NewText(tempMaxTxt, textColor)
	tempMin := canvas.NewText(tempMinTxt, textColor)
	humitat := canvas.NewText(humitatTxt, textColor)

	// Alineació dels textos en la interfície
	precipitacio.Alignment = fyne.TextAlignLeading
	tempMax.Alignment = fyne.TextAlignCenter
	tempMin.Alignment = fyne.TextAlignCenter
	humitat.Alignment = fyne.TextAlignTrailing

	return precipitacio, tempMax, tempMin, humitat, nil
}
