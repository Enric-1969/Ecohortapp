package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// registresTab genera la pestanya de registres (coincideix amb la crida de ui.go)
func (app *Config) registresTab() *fyne.Container {
	tabla := app.getRegistresTable()

	dataEntrada := widget.NewEntry()
	dataEntrada.SetPlaceHolder("AAAA-MM-DD")

	precipitacioEntrada := widget.NewEntry()
	precipitacioEntrada.SetPlaceHolder("0")

	tempMaxEntrada := widget.NewEntry()
	tempMaxEntrada.SetPlaceHolder("0")

	tempMinEntrada := widget.NewEntry()
	tempMinEntrada.SetPlaceHolder("0")

	humitatEntrada := widget.NewEntry()
	humitatEntrada.SetPlaceHolder("0")

	app.AfegirRegistresDataRegistreEntrada = dataEntrada
	app.AfegirRegistresPrecipitacioEntrada = precipitacioEntrada
	app.AfegirRegistresTempMaximaEntrada = tempMaxEntrada
	app.AfegirRegistresTempMinimaEntrada = tempMinEntrada
	app.AfegirRegistresHumitatEntrada = humitatEntrada

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Data Registre", Widget: dataEntrada},
			{Text: "Precipitació (%)", Widget: precipitacioEntrada},
			{Text: "Temp. Màxima (°C)", Widget: tempMaxEntrada},
			{Text: "Temp. Mínima (°C)", Widget: tempMinEntrada},
			{Text: "Humitat (%)", Widget: humitatEntrada},
		},
		OnSubmit: func() {
			err := app.guardarRegistre()
			if err != nil {
				dialog.ShowError(err, app.MainWindow)
				return
			}
			dialog.ShowInformation("Èxit", "Registre afegit correctament", app.MainWindow)
			app.netejarFormulari()
			app.refreshRegistresTable()
		},
	}

	formContainer := container.NewVBox(
		widget.NewLabelWithStyle("Afegir Nou Registre Climatològic", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		form,
	)

	return container.NewBorder(formContainer, nil, nil, nil, tabla)
}

// getRegistresSlice retorna les dades en format matriu per a la taula
func (app *Config) getRegistresSlice() [][]interface{} {
	var slice [][]interface{}

	// Encapçalats
	slice = append(slice, []interface{}{"ID", "Data", "Precipitació", "T. Max", "T. Min", "Humitat"})

	// Afegir registres emmagatzemats si n'hi ha
	if len(app.Registres) > 0 {
		slice = append(slice, app.Registres...)
	}

	return slice
}

// getRegistresTable construeix el widget de taula de Fyne
func (app *Config) getRegistresTable() *widget.Table {
	dades := app.getRegistresSlice()

	tabla := widget.NewTable(
		func() (int, int) {
			return len(dades), len(dades[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Ample per defecte")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			label := o.(*widget.Label)
			label.SetText(fmt.Sprintf("%v", dades[i.Row][i.Col]))
		},
	)

	tabla.SetColumnWidth(0, 50)
	tabla.SetColumnWidth(1, 100)
	tabla.SetColumnWidth(2, 110)
	tabla.SetColumnWidth(3, 80)
	tabla.SetColumnWidth(4, 80)
	tabla.SetColumnWidth(5, 80)

	app.RegistresTable = tabla
	return tabla
}

// refreshRegistresTable actualitza la visualització de la taula
func (app *Config) refreshRegistresTable() {
	if app.RegistresTable != nil {
		app.RegistresTable.Refresh()
	}
}

// netejarFormulari buida els camps d'entrada
func (app *Config) netejarFormulari() {
	if app.AfegirRegistresDataRegistreEntrada != nil {
		app.AfegirRegistresDataRegistreEntrada.SetText("")
	}
	if app.AfegirRegistresPrecipitacioEntrada != nil {
		app.AfegirRegistresPrecipitacioEntrada.SetText("")
	}
	if app.AfegirRegistresTempMaximaEntrada != nil {
		app.AfegirRegistresTempMaximaEntrada.SetText("")
	}
	if app.AfegirRegistresTempMinimaEntrada != nil {
		app.AfegirRegistresTempMinimaEntrada.SetText("")
	}
	if app.AfegirRegistresHumitatEntrada != nil {
		app.AfegirRegistresHumitatEntrada.SetText("")
	}
}

// guardarRegistre gestiona la gravació de dades
func (app *Config) guardarRegistre() error {
	// Lògica de persistència
	return nil
}
