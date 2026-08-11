package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// crearBannerAlerta genera un bàner d'alerta compost i elegant quan hi ha un avís actiu
func (app *Config) crearBannerAlerta(nivel string, missatge string) *fyne.Container {
	// Si no hi ha alerta, nivell verd o missatge buit, no mostrem res (interfície neta)
	if nivel == "" || nivel == "verde" || missatge == "" {
		return container.NewVBox()
	}

	var fonsColor color.Color
	var titolText string

	switch nivel {
	case "amarillo":
		fonsColor = color.RGBA{R: 255, G: 193, B: 7, A: 210} // Amarillo cálido
		titolText = "AVÍS GROC: RISC METEOROLÒGIC"
	case "naranja":
		fonsColor = color.RGBA{R: 255, G: 112, B: 67, A: 220} // Naranja alerta
		titolText = "AVÍS TARONJA: RISC IMPORTANT"
	case "rojo":
		fonsColor = color.RGBA{R: 229, G: 57, B: 53, A: 230} // Rojo peligro
		titolText = "AVÍS VERMELL: RISC EXTREM"
	default:
		fonsColor = color.RGBA{R: 33, G: 150, B: 243, A: 190} // Azul informativo
		titolText = "AVÍS METEOROLÒGIC"
	}

	// 1. Fondo de color para la tarjeta
	rect := canvas.NewRectangle(fonsColor)

	// 2. Textos formateados
	lblTitol := widget.NewLabel(titolText)
	lblTitol.TextStyle = fyne.TextStyle{Bold: true}
	lblTitol.Alignment = fyne.TextAlignCenter

	lblDesc := widget.NewLabel(missatge)
	lblDesc.Alignment = fyne.TextAlignCenter

	// 3. Empaquetado con padding (márgenes internos)
	contingut := container.NewVBox(lblTitol, lblDesc)
	targeta := container.NewMax(rect, container.NewPadded(contingut))

	// Retornamos la tarjeta envuelta en un margen externo para no pegar al gráfico
	return container.NewPadded(targeta)
}
