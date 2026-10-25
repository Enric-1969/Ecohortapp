package repository

import (
	"log"
	"os"
)

// InicialitzarLoggers configura el fitxer físic app.log en el disc dur.
// Retorna el descriptor del fitxer i els dos loggers configurats.
func InicialitzarLoggers() (*os.File, *log.Logger, *log.Logger) {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logFile = os.Stdout // Fallback a la consola si falla el disc de Windows
	}

	infoLog := log.New(logFile, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(logFile, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return logFile, infoLog, errorLog
}
