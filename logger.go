package main

import (
	"errors"

	"fyne.io/fyne/v2/dialog"
)

// MostrarErrorCentralitzat registra el detalle técnico en app.log
// y muestra una alerta limpia al usuario final en la interfaz Fyne.
func (c *Config) MostrarErrorCentralitzat(errTecnico error, mensajeUsuario string) {
	// 1. Log técnico invisible para el usuario (se escribe en app.log)
	if errTecnico != nil {
		c.ErrorLog.Println(errTecnico)
	}

	// 2. Alerta amigable en la interfaz gráfica
	dialog.ShowError(errors.New(mensajeUsuario), c.MainWindow)
}
