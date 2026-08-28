package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	var myApp Config

	// 1. Carregar la configuració de l'usuari
	userCfg, err := LoadConfig()
	if err != nil {
		log.Println("No s'ha trobat config.json, usant configuració per defecte:", err)
		userCfg = UserConfig{
			MunicipioCodigo: "08001",
			Municipios:      []string{"08001"},
		}
	} else {
		log.Println("Configuració carregada amb èxit:", userCfg)
	}

	// 2. Crear l'aplicació Fyne
	fyneApp := app.NewWithID("cat.cibernarium.ecohortapp")
	myApp.App = fyneApp

	// 3. Crear archivo de logs en disco i configuració de loggers
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logFile = os.Stdout // Si falla la creación del archivo, usa la consola como alternativa
	}

	myApp.InfoLog = log.New(logFile, "INFO\t", log.Ldate|log.Ltime)
	myApp.ErrorLog = log.New(logFile, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// 4. Connexió i configuració de la base de dades
	sqlDB, err := myApp.connectSQL()
	if err != nil {
		log.Panic(err)
	}
	myApp.setupDB(sqlDB)

	// 5. Inicialitzar el client HTTP amb timeout
	myApp.HTTPClient = http.Client{
		Timeout: 15 * time.Second,
	}

	// 6. Assignar configuració i API Key
	myApp.UserConfig = userCfg
	myApp.municipi = userCfg.MunicipioCodigo
	myApp.apiKey = fyneApp.Preferences().StringWithFallback("aemet_api_key", os.Getenv("AEMET_API_KEY"))

	// 7. Càrrega silenciosa de municipis en segon pla (RAM / Cache)
	go func() {
		if m, err := myApp.ObtenirMunicipiosAEMET(); err == nil {
			myApp.Municipis = m
			myApp.InfoLog.Println("[INIT] Municipis carregats automàticament a la memòria RAM.")
		} else {
			myApp.ErrorLog.Println("[INIT] No s'han pogut carregar els municipis en segon pla:", err)
		}
	}()

	// 8. Crear la finestra principal
	myApp.MainWindow = fyneApp.NewWindow("Eco Hort App")
	myApp.MainWindow.Resize(fyne.NewSize(800, 700))
	myApp.MainWindow.SetFixedSize(true)
	myApp.MainWindow.SetMaster()

	myApp.makeUI()

	// 9. Executar l'aplicació
	myApp.MainWindow.ShowAndRun()
}
