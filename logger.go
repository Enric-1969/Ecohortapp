package main

import (
	"errors"
	"fmt"

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

// TratarErrorHTTP mapea los códigos de estado HTTP de AEMET
func (c *Config) TratarErrorHTTP(statusCode int, errTecnico error) {
	var msgUsuario string

	switch statusCode {
	case 401, 403:
		msgUsuario = "La clave API de AEMET no es válida. Revisa tus preferencias."
	case 429:
		msgUsuario = "Límite de peticiones alcanzado. Espera unos minutos."
	case 500, 502, 503:
		msgUsuario = "El servidor de AEMET está sufriendo problemas técnicos."
	default:
		msgUsuario = fmt.Sprintf("Error de red al conectar con AEMET (Código %d).", statusCode)
	}

	c.MostrarErrorCentralitzat(errTecnico, msgUsuario)
}
